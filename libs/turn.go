package libs

import (
	"net"
	"fmt"
	"strings"
	"strconv"
	"encoding/hex"
	"bytes"
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

func messageIntegrityCheck(requestMessage *Message) (err error) {
	miAttr := requestMessage.getAttribute(AttributeMessageIntegrity)

	if miAttr != nil {
		userAttr := requestMessage.getAttribute(AttributeUsername)
		if userAttr != nil {
			username := string(userAttr.Value)
			password := hex.EncodeToString(HmacSha1(userAttr.Value,[]byte("passwordkey")))
			//Log.Infof("password %s",password)
			key := generateKey(username,password,"realm")
			requestValue , err :=  Marshal(requestMessage,true)

			if err == nil {
				calculateMi := MessageIntegrityHmac(requestValue,key)

				if(!bytes.Equal(calculateMi,miAttr.Value)){
					//todo : not equal
				}
			}else{
				//todo
			}
		}else{
			//todo
		}
	}else{
		//todo : error response
	}

	return
}

func turnMessageHandle(requestMessage *Message,raddr *net.UDPAddr,tcp bool) (response []byte,responseAddress *net.UDPAddr,err error)  {
	//Log.Verbosef("turn request : %s",message)

	respMsg := new(Message)
	respMsg.TransID = requestMessage.TransID
	respMsg.Attributes = make([]*Attribute,0)

	switch requestMessage.MessageType {
	case TypeAllocate:

		// check message is auth
		ok := requestMessage.hasAttribute(AttributeRealm)

		// mi = message integrity
		if ok {
			nonce := requestMessage.getAttribute(AttributeNonce)

			if(validNonce(nonce.Value)){
				err := messageIntegrityCheck(requestMessage)

				if err != nil {
					//todo 400 error
				}else{
					originUsername := requestMessage.getAttribute(AttributeUsername)

					Log.Info(originUsername.String())
					originUsernameArray := strings.Split(string(originUsername.Value),":")
					if len(originUsernameArray) > 1 {
						accountExpireTime, err := strconv.Atoi(originUsernameArray[0])

						if err == nil {
							//todo : check accountExpireTime
							if accountExpireTime > 0 {
								realUsername := originUsernameArray[1]
								transport := requestMessage.getAttribute(AttributeRequestedTransport).Value

								Log.Infof("username : %s , len : %d , t %d",realUsername,len(realUsername),transport[0] )
							}
						}else{
							//todo : not available username
						}
					}else{
						//todo : not available username
					}
				}
			}else{
				//todo : error stale nonce
			}



			//
			//// create Allocate
			//respMsg.MessageType = TypeAllocateResponse
			//
			//port := RelayPortPool.RandSelectPort()
			//Log.Infof("random port %d",port)
			//server := NewRelayServer(port)
			//server.Serve()
			////todo : check relay address is available
			//relayAddress := getRelayAddress()
			//
			//rKey := fmt.Sprintf("%s:%d",relayAddress,port)
			//RelayMap[raddr.String()] = rKey
			//
			//
			//
			//
			//
			//
			//respMsg.addAttribute(newAttrXORRelayedAddress(relayAddress,port))
			//respMsg.addAttribute(newAttrXORMappedAddress(raddr))
			//respMsg.addAttribute(newAttrLifetime())
			//respMsg.addAttribute(newAttrSoftware())
			//respMsg.addAttribute(newAttrDummyMessageIntegrity())
			//
			//var m_i_response []byte
			//m_i_response, err = Marshal(respMsg)
			//
			//if err != nil {
			//	return
			//}
			//
			//key := generateKey("user","pass","realm")
			//
			//hmacValue := MessageIntegrityHmac(m_i_response[:len(m_i_response)-24],key)
			//
			//response = append(m_i_response[:len(m_i_response)-20],hmacValue...)



		}else{
			respMsg.MessageType = TypeAllocateErrorResponse

			respMsg.addAttribute(newAttrNonce())
			respMsg.addAttribute(AttrRealm)
			respMsg.addAttribute(newAttrError401())
			respMsg.addAttribute(AttrSoftware)

			response, err = Marshal(respMsg,false)

			if err != nil {
				return
			}
		}
	case TypeCreatePermisiion:
		Log.Info(" permission \n")

		respMsg.MessageType = TypeCreatePermisiionResponse
		respMsg.addAttribute(AttrSoftware)
		respMsg.addAttribute(newAttrDummyMessageIntegrity())

		var m_i_response []byte
		m_i_response, err = Marshal(respMsg,false)

		if err != nil {
			return
		}

		key := generateKey("user","pass","realm")

		hmacValue := MessageIntegrityHmac(m_i_response[:len(m_i_response)-24],key)

		response = append(m_i_response[:len(m_i_response)-20],hmacValue...)
	case TypeSendIndication:
		Log.Info(" indication \n")

		peerAddress := requestMessage.getAttribute(AttributeXorPeerAddress)

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
			respMsg.addAttribute(requestMessage.getAttribute(AttributeData))
			respMsg.addAttribute(newAttrXORPeerAddress(pAddress,peerPort))
			respMsg.addAttribute(AttrSoftware)

			response, err = Marshal(respMsg,false)

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
		respMsg.addAttribute(AttrSoftware)
		respMsg.addAttribute(newAttrDummyMessageIntegrity())

		var m_i_response []byte
		m_i_response, err = Marshal(respMsg,false)

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