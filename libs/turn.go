package libs

import (
	"net"
	"fmt"
)

func turnMessageHandle(message *Message,raddr *net.UDPAddr,tcp bool) (response []byte,err error)  {
	//Log.Verbosef("turn request : %s",message)

	respMsg := new(Message)
	respMsg.TransID = message.TransID
	respMsg.Attributes = make([]*Attribute,0)

	switch message.MessageType {
	case TypeAllocate:

		ok := message.hasAttribute(AttributeRealm)

		if ok {

			respMsg.MessageType = TypeAllocateResponse

			respMsg.addAttribute(newAttrXORRelayedAddress())
			respMsg.addAttribute(newAttrXORMappedAddress(raddr))
			respMsg.addAttribute(newAttrLifetime())
			respMsg.addAttribute(newAttrSoftware())
			respMsg.addAttribute(newAttrDummyMessageIntegrity())

			var m_i_response []byte
			m_i_response, err = Marshal(respMsg)

			if err != nil {
				return
			}

			key := generateKey("user","pass","realm")

			hmacValue := MessageIntegrityHmac(m_i_response[:len(m_i_response)-24],key)

			response = append(m_i_response[:len(m_i_response)-20],hmacValue...)



		}else{
			respMsg.MessageType = TypeAllocateErrorResponse

			respMsg.addAttribute(newAttrNonce())
			respMsg.addAttribute(newAttrRealm())
			respMsg.addAttribute(newAttrError401())
			respMsg.addAttribute(newAttrSoftware())

			// addMappedAddress(respMsg, raddr)


			fmt.Printf("allocate response : %s \n",respMsg)

			response, err = Marshal(respMsg)

			if err != nil {
				return
			}

		}
	case TypeCreatePermisiion:
		Log.Info("create permission \n")

		respMsg.MessageType = TypeCreatePermisiionResponse
		respMsg.addAttribute(newAttrDummyMessageIntegrity())

		var m_i_response []byte
		m_i_response, err = Marshal(respMsg)

		if err != nil {
			return
		}

		key := generateKey("user","pass","realm")

		hmacValue := MessageIntegrityHmac(m_i_response[:len(m_i_response)-24],key)

		response = append(m_i_response[:len(m_i_response)-20],hmacValue...)
	}


	return
}