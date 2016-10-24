package libs

import (
	"net"
	"fmt"
	"strings"
	"strconv"
)

// client : relay
var RelayMap map[string]string

func getClientAddress(m map[string]string,relay string) string {
	for k,v := range m{
		if v == relay  {
			return k
		}
	}
	return ""
}

func turnMessageHandle(message *Message,raddr *net.UDPAddr,tcp bool) (response []byte,responseAddress *net.UDPAddr,err error)  {
	//Log.Verbosef("turn request : %s",message)

	respMsg := new(Message)
	respMsg.TransID = message.TransID
	respMsg.Attributes = make([]*Attribute,0)

	switch message.MessageType {
	case TypeAllocate:

		ok := message.hasAttribute(AttributeRealm)

		if ok {

			respMsg.MessageType = TypeAllocateResponse

			port := RelayPortPool.RandSelectPort()
			Log.Infof("random port %d",port)
			server := NewRelayServer(port)
			server.Serve()
			//todo :
			relayAddress := getRelayAddress()

			rKey := fmt.Sprintf("%s:%d",relayAddress,port)
			RelayMap[raddr.String()] = rKey

			respMsg.addAttribute(newAttrXORRelayedAddress(relayAddress,port))
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
		Log.Info(" permission \n")

		respMsg.MessageType = TypeCreatePermisiionResponse
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
	case TypeSendIndication:
		Log.Info(" indication \n")

		peerAddress := message.getAttribute(AttributeXorPeerAddress)

		if peerAddress != nil{
			port,address := unXorAddress(peerAddress.Value)


			saddr := fmt.Sprintf("%s:%d",net.IP(address),port)
			caddr := strings.Split(getClientAddress(RelayMap,saddr),":")
			Log.Infof("client : %s ,relay local : %s",saddr,port,caddr)
			responseAddress = new(net.UDPAddr)
			responseAddress.IP =  net.ParseIP(caddr[0]).To4()
			responseAddress.Port,_ = strconv.Atoi(caddr[1])


			relayAddress := RelayMap[raddr.String()]

			peerPortStr := strings.Split(relayAddress,":")[1]
			pAddress := strings.Split(relayAddress,":")[0]
			peerPort , _ := strconv.Atoi(peerPortStr)
			respMsg.MessageType = TypeDataIndication
			respMsg.addAttribute(message.getAttribute(AttributeData))
			respMsg.addAttribute(newAttrXORPeerAddress(pAddress,peerPort))
			respMsg.addAttribute(newAttrSoftware())

			response, err = Marshal(respMsg)

			if err != nil {
				return
			}
		}else{
			//todo : error
		}
		//respMsg.MessageType = TypeDataIndication
		//respMsg.addAttribute()
	case TypeChannelBinding:
		respMsg.MessageType = TypeChannelBindingResponse
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
	case TypeRefreshRequest:
		//todo : ....
	}


	return
}