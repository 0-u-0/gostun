package libs

import (
	"fmt"
	"net"
)

func stunMessageHandle(message *Message,raddr *net.UDPAddr,tcp bool) (response []byte) {
	switch message.MessageType {
	case TypeBindingRequest:
		//fmt.Printf("binding request : %s \n",msg)

		//todo : handle with origin
		respMsg := new(Message)
		respMsg.MessageType = TypeBindingResponse
		respMsg.TransID = message.TransID
		respMsg.Attributes = make([]*Attribute,0)

		respMsg.addAttribute(newAttrXORMappedAddress(raddr))

		var err error
		response, err = Marshal(respMsg)

		if err != nil {
			fmt.Println(err)
			return
		}
	}

	return
}

