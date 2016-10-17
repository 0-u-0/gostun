package libs

import (
	"log"
	"net"
	"strconv"
)

type Entry struct {
	Port       int
	udpConn *net.UDPConn
}

func LoadEntryModule()  {
	PrintModuleLoaded("Entry")
	entry := NewEntry(3478)
	entry.Serve()
}
func (s *Entry) serveUDP() {
	laddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(s.Port))
	if err != nil {
		log.Fatal(err)
	}
	s.udpConn, err = net.ListenUDP("udp", laddr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		var buf = make([]byte, 1024)
		size, remoteAddr, err := s.udpConn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		go s.handleData(remoteAddr, buf[:size],false)
	}
}

func (s *Entry) Serve() {
	serverTCP, serverTLS := false,false
	serverUDP := true

	if serverTCP{

	}

	if serverTLS {

	}

	if serverUDP {
		s.serveUDP()
	}

}

func NewEntry(port int) *Entry {
	ret := new(Entry)
	ret.Port = port
	return ret
}

func (entry *Entry) handleData(raddr *net.UDPAddr, data []byte,tcp bool) {
	msg, err := UnMarshal(data)
	if err != nil {
		Log.Warning(err)
		return
	}

	Log.Verbosef("msg : %s \n",msg)

	var response []byte
	var response_err error
	switch msg.MessageType {
	case TypeBindingRequest:
		response,response_err = stunMessageHandle(msg,raddr,false)
	case TypeAllocate , TypeCreatePermisiion:
		response,response_err = turnMessageHandle(msg,raddr,false)
	}

	if response_err == nil{
		if !tcp {
			if response != nil {
				_, err := entry.udpConn.WriteToUDP(response, raddr)
				if err != nil {
					Log.Warning(err)
				}
			}else {
				//todo add message type check
				//Log.Warning("no response.")
			}

		}else {
			//todo : add tcp
		}
	}else{
		Log.Warningf("response error : %s",response_err)
	}

}
