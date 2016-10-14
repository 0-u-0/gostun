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
		fmt.Println(err)
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

			respMsg := new(Message)
			respMsg.MessageType = TypeAllocateResponse
			respMsg.TransID = msg.TransID
			respMsg.Attributes = make([]*Attribute,0)

			respMsg.addAttribute(newAttrXORRelayedAddress())
			respMsg.addAttribute(newAttrXORMappedAddress(raddr))
			respMsg.addAttribute(newAttrLifetime())
			respMsg.addAttribute(newAttrSoftware())
			respMsg.addAttribute(newAttrDummyMessageIntegrity())

			fmt.Printf("m-i response : %s \n",respMsg)

			mi, err := Marshal(respMsg)


			fmt.Printf("response hex : %x \n",mi)

			if err != nil {
				fmt.Println(err)
				return
			}
			//fmt.Printf("fake length : %d \n", binary.BigEndian.Uint16(mi[2:4]))

			key := generateKey("user","pass","realm")

			hmacValue := MessageIntegrityHmac(mi[0:len(mi)-24],key)

			fmt.Printf("hmac2 length %d , hmac2 %x \n",len(hmacValue),hmacValue)


			response := append(mi[:len(mi)-20],hmacValue...)

			testResponse, err := UnMarshal(response)

			fmt.Printf("test response : %s \n",testResponse)
			fmt.Printf("test response hex : %x \n",response)


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
