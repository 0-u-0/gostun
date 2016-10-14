package libs

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

type Entry struct {
	Port       int
	udpConn *net.UDPConn
}

func (s *Entry) serve() {
	for {
		var buf = make([]byte, 1024)
		size, remoteAddr, err := s.udpConn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		go s.handleData(remoteAddr, buf[:size])
	}
}

func (s *Entry) Serve() {
	laddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(s.Port))
	if err != nil {
		log.Fatal(err)
	}
	s.udpConn, err = net.ListenUDP("udp", laddr)
	if err != nil {
		log.Fatal(err)
	}
	s.serve()
}

func NewServer(port int) *Entry {
	ret := new(Entry)
	ret.Port = port

	return ret
}

func (s *Entry) handleData(raddr *net.UDPAddr, data []byte) {
	msg, err := UnMarshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch msg.MessageType {
	case TypeBindingRequest:
		//todo : add connection
		//stunMessageHandle()
	case TypeAllocate:

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
			_, err = s.udpConn.WriteToUDP(response, raddr)
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
			_, err = s.udpConn.WriteToUDP(response, raddr)
			if err != nil {
				fmt.Println(err)
			}
			//fmt.Printf("binding response : %s \n",respMsg)
		}


	}


}
