package libs

import (
	"net"
)

func stunMessageHandle(message *Message,raddr *net.UDPAddr,tcp bool) (response []byte,err error) {
	Log.Verbosef("stun request : %s ",message)

	switch message.MessageType {
	case TypeBindingRequest:

		//todo : handle with origin
		respMsg := new(Message)
		respMsg.MessageType = TypeBindingResponse
		respMsg.TransID = message.TransID
		respMsg.Attributes = make([]*Attribute,0)

		respMsg.addAttribute(newAttrXORMappedAddress(raddr))

		response, err = Marshal(respMsg)

		if err != nil {
			return
		}
	}

	return
}

