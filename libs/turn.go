package libs

import (
	"fmt"
	"net"
)

func turnMessageHandle(message *Message,raddr *net.UDPAddr,tcp bool) (response []byte)  {
	switch message.MessageType {
	case TypeAllocate:

		ok := message.hasAttribute(AttributeRealm)

		if ok {

			respMsg := new(Message)
			respMsg.MessageType = TypeAllocateResponse
			respMsg.TransID = message.TransID
			respMsg.Attributes = make([]*Attribute,0)

			respMsg.addAttribute(newAttrXORRelayedAddress())
			respMsg.addAttribute(newAttrXORMappedAddress(raddr))
			respMsg.addAttribute(newAttrLifetime())
			respMsg.addAttribute(newAttrSoftware())
			respMsg.addAttribute(newAttrDummyMessageIntegrity())


			m_i_response, err := Marshal(respMsg)

			if err != nil {
				fmt.Println(err)
				return
			}

			key := generateKey("user","pass","realm")

			hmacValue := MessageIntegrityHmac(m_i_response[:len(m_i_response)-24],key)

			response = append(m_i_response[:len(m_i_response)-20],hmacValue...)

			//fmt.Printf("binding response : %s \n",respMsg)


		}else{
			respMsg := new(Message)
			respMsg.MessageType = TypeAllocateErrorResponse
			respMsg.TransID = message.TransID
			respMsg.Attributes = make([]*Attribute,0)

			respMsg.addAttribute(newAttrNonce())
			respMsg.addAttribute(newAttrRealm())
			respMsg.addAttribute(newAttrError401())
			respMsg.addAttribute(newAttrSoftware())

			// addMappedAddress(respMsg, raddr)


			//fmt.Printf("allocate response : %s \n",respMsg)

			var err error
			response, err = Marshal(respMsg)
			if err != nil {
				fmt.Println(err)
				return
			}

		}
	}

	return
}