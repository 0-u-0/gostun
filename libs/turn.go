package libs

import (
	"net"

	"strings"
	"strconv"
	"encoding/hex"
	"bytes"
	"sync"
)

var Allocates map[string]*Allocate
var Mutex *sync.Mutex

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

func messageIntegrityCalculate(username string,responseMessage *Message) (response []byte, err error) {
	var m_i_response []byte
	m_i_response, err = Marshal(responseMessage,false)

	if err != nil {
		return nil,err
	}

	password := hex.EncodeToString(HmacSha1([]byte(username),[]byte("passwordkey")))

	key := generateKey(username,password,"realm")

	hmacValue := MessageIntegrityHmac(m_i_response[:len(m_i_response)-24],key)

	response = append(m_i_response[:len(m_i_response)-20],hmacValue...)

	return response,nil
}


func turnMessageHandle(requestMessage *Message,raddr *net.UDPAddr,tcp bool) ([]byte, *net.UDPAddr,error)  {
	//Log.Verbosef("turn request : %s",message)


	switch requestMessage.MessageType {
	case TypeAllocate:

		// check message is auth
		ok := requestMessage.hasAttribute(AttributeRealm)

		Log.Infof("request : %s",requestMessage)
		// mi = message integrity
		if ok {
			nonce := requestMessage.getAttribute(AttributeNonce)

			if(validNonce(nonce.Value)){
				err := messageIntegrityCheck(requestMessage)

				if err != nil {
					//todo 400 error
				}else{
					originUsername := requestMessage.getAttribute(AttributeUsername)
					strUsername := string(originUsername.Value)
					
					originUsernameArray := strings.Split(strUsername,":")
					if len(originUsernameArray) > 1 {
						accountExpireTime, err := strconv.Atoi(originUsernameArray[0])

						if err == nil {
							//todo : check accountExpireTime
							if accountExpireTime > 0 {
								realUsername := originUsernameArray[1]
								protocol := requestMessage.getAttribute(AttributeRequestedTransport).Value

								// create Allocate
								//respMsg.MessageType = TypeAllocateResponse

								port := RelayPortPool.RandSelectPort()

								server := NewRelayServer(port)
								server.Serve()
								//todo : check relay address is available
								relayAddress := getRelayAddress()

								relay := new(net.UDPAddr)
								relay.Port = port
								relay.IP = net.ParseIP(relayAddress)

								allocate := NewAllocate(realUsername,protocol[0],MaxTimeRefresh,raddr,relay)

								clientAddressString := raddr.String()

								Mutex.Lock()
								Allocates[clientAddressString] = allocate
								Mutex.Unlock()

								respMsg := NewResponse(TypeAllocateResponse,requestMessage.TransID,
									newAttrXORRelayedAddress(relayAddress,port),
									newAttrXORMappedAddress(raddr),
									AttrLifetime,
									AttrSoftware,
									AttrDummyMessageIntegrity,
								)


								response,err := messageIntegrityCalculate(strUsername,respMsg)


								return response,nil,err
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


		}else{
			respMsg := NewResponse(TypeAllocateErrorResponse,requestMessage.TransID,
				newAttrNonce(),
				AttrRealm,
				AttrError401,
				AttrSoftware,
			)

			response, err := Marshal(respMsg,false)

			if err != nil {
				return nil,nil,err
			}
			return response,nil,nil
		}
	case TypeCreatePermisiion:
		Log.Info(" permission \n")


		//todo : check
		respMsg := NewResponse(TypeCreatePermisiionResponse,requestMessage.TransID,
			AttrSoftware,
			AttrDummyMessageIntegrity,
		)

		originUsername := requestMessage.getAttribute(AttributeUsername)
		strUsername := string(originUsername.Value)


		response ,err := messageIntegrityCalculate(strUsername,respMsg)
		return response,nil,err
	case TypeSendIndication:
		Log.Info(" indication \n")

		//peerAddress := requestMessage.getAttribute(AttributeXorPeerAddress)
		//
		//if peerAddress != nil{
		//	port,address := unXorAddress(peerAddress.Value)
		//
		//
		//	saddr := fmt.Sprintf("%s:%d",net.IP(address),port)
		//	caddr := strings.Split(getClientAddress(RelayMap,saddr),":")
		//	Log.Infof("client : %s ,relay local : %s",saddr,port,caddr)
		//	responseAddress = new(net.UDPAddr)
		//	responseAddress.IP =  net.ParseIP(caddr[0]).To4()
		//	responseAddress.Port,_ = strconv.Atoi(caddr[1])
		//
		//
		//	relayAddress := RelayMap[raddr.String()]
		//
		//	peerPortStr := strings.Split(relayAddress,":")[1]
		//	pAddress := strings.Split(relayAddress,":")[0]
		//	peerPort , _ := strconv.Atoi(peerPortStr)
		//	respMsg.MessageType = TypeDataIndication
		//	respMsg.addAttribute(requestMessage.getAttribute(AttributeData))
		//	respMsg.addAttribute(newAttrXORPeerAddress(pAddress,peerPort))
		//	respMsg.addAttribute(AttrSoftware)
		//
		//	response, err = Marshal(respMsg,false)
		//
		//	if err != nil {
		//		return
		//	}
		//}else{
		//	//todo : error
		//}


	case TypeChannelBinding:
		//respMsg.MessageType = TypeChannelBindingResponse
		//respMsg.addAttribute(AttrSoftware)
		//respMsg.addAttribute(newAttrDummyMessageIntegrity())
		//
		//var m_i_response []byte
		//m_i_response, err = Marshal(respMsg,false)
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
	case TypeRefreshRequest:
		//todo : ....

	}

	return nil,nil,nil
}