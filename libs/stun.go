package libs

import (
	"fmt"
	"net"
)

func stunMessageHandle(entry *Entry,message *Message,raddr *net.UDPAddr)  {
	switch message.MessageType {
	case TypeBindingRequest:
		//fmt.Printf("binding request : %s \n",msg)

		//todo : handle with origin
		respMsg := new(Message)
		respMsg.MessageType = TypeBindingResponse
		respMsg.TransID = message.TransID
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
		_, err = entry.udpConn.WriteToUDP(response, raddr)
		if err != nil {
			fmt.Println(err)
		}
	}
}

