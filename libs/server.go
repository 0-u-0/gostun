package libs

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

//Server is the main struct that contains the connection and registry information
type Server struct {
	Port       int
	Registry   *Registry
	connection *net.UDPConn
}

func (s *Server) serve() {
	for {
		var buf = make([]byte, 1024)
		size, remoteAddr, err := s.connection.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		go s.handleData(remoteAddr, buf[:size])
	}
}

//Serve initiates a UDP connection that listens on any port for incoming data
func (s *Server) Serve() {
	laddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(s.Port))
	if err != nil {
		log.Fatal(err)
	}
	s.connection, err = net.ListenUDP("udp", laddr)
	if err != nil {
		log.Fatal(err)
	}
	s.serve()
}

//NewServer conveniently creates a new server from the given port
func NewServer(port int) *Server {
	ret := new(Server)
	ret.Port = port
	ret.Registry = new(Registry)
	ret.Registry.mappings = make(map[string]*Client)
	return ret
}

func (s *Server) handleData(raddr *net.UDPAddr, data []byte) {
	msg, err := UnMarshal(data)
	if err != nil {
		return
	}

	switch msg.MessageType {
	case TypeBindingRequest:
		//fmt.Printf("binding request : %s \n",msg)

		//todo : handle with origin
		respMsg := new(Message)
		respMsg.MessageType = TypeBindingResponse
		respMsg.TransID = msg.TransID
		respMsg.Attributes = make([]*Attribute,0)

		respMsg.addAttribute(newAttrXORMappedAddress(raddr))
		// addMappedAddress(respMsg, raddr)


		//fmt.Printf("binding response : %s \n",respMsg)

		response, err := Marshal(respMsg)
		if err != nil {
			fmt.Println(err)
			return
		}
		//send response
		_, err = s.connection.WriteToUDP(response, raddr)
		if err != nil {
			fmt.Println(err)
		}
	case TypeAllocate:
		//fmt.Printf("allocate request : %s \n",msg)

		ok := msg.hasAttribute(AttributeRealm)

		if ok {

			fakeMsg := new(Message)
			fakeMsg.MessageType = TypeAllocateResponse
			fakeMsg.TransID = msg.TransID
			fakeMsg.Attributes = make([]*Attribute,0)

			fakeMsg.addAttribute(newAttrXORRelayedAddress())
			fakeMsg.addAttribute(newAttrXORMappedAddress(raddr))
			fakeMsg.addAttribute(newAttrLifetime())
			fakeMsg.addAttribute(newAttrSoftware())
			//fakeMsg.addAttribute(newAttrFakeMessageIntegrity())
			fakeMsg.MessageLength += 24

			fmt.Printf("fake response : %s \n",fakeMsg)


			fakeResponse, err := Marshal(fakeMsg)

			//fmt.Printf("fake length : %d \n", binary.BigEndian.Uint16(fakeResponse[2:4]))

			key := generateKey("user","pass","realm")
			hmacValue := messageIntegrityHmac(fakeResponse,key)


			respMsg := new(Message)
			respMsg.MessageType = TypeAllocateResponse
			respMsg.TransID = msg.TransID
			respMsg.Attributes = make([]*Attribute,0)

			respMsg.addAttribute(newAttrXORRelayedAddress())
			respMsg.addAttribute(newAttrXORMappedAddress(raddr))
			respMsg.addAttribute(newAttrLifetime())
			respMsg.addAttribute(newAttrSoftware())
			respMsg.addAttribute(newAttrMessageIntegrity (hmacValue))

			/*
			 Transaction-Id=0xC271E932AD7446A32C234492     |             |
		|    SOFTWARE="Example server, version 1.17"       |             |
		|    LIFETIME=1200 (20 minutes)      |             |             |
		|    XOR-RELAYED-ADDRESS=192.0.2.15:50000          |             |
		|    XOR-MAPPED-ADDRESS=192.0.2.1:7000             |             |
        |    MESSAGE-INTEGRITY=...
			*/
			fmt.Printf("allocate response : %s \n",respMsg)

			response, err := Marshal(respMsg)

			//fmt.Printf("real length : %d \n", binary.BigEndian.Uint16(response[2:4]))

			if err != nil {
				fmt.Println(err)
				return
			}
			//send response
			_, err = s.connection.WriteToUDP(response, raddr)
			if err != nil {
				fmt.Println(err)
			}
			//fmt.Printf("binding response : %s \n",respMsg)


		}else{
			respMsg := new(Message)
			respMsg.MessageType = TypeAllocateErrorResponse
			respMsg.TransID = msg.TransID
			respMsg.Attributes = make([]*Attribute,0)

			respMsg.addAttribute(newAttrNonce())
			respMsg.addAttribute(newAttrRealm())
			respMsg.addAttribute(newAttrError401())
			respMsg.addAttribute(newAttrSoftware())

			// addMappedAddress(respMsg, raddr)


			//fmt.Printf("allocate response : %s \n",respMsg)

			response, err := Marshal(respMsg)
			if err != nil {
				fmt.Println(err)
				return
			}
			//send response
			_, err = s.connection.WriteToUDP(response, raddr)
			if err != nil {
				fmt.Println(err)
			}
			//fmt.Printf("binding response : %s \n",respMsg)
		}


	}


}
