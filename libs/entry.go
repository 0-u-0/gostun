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
		fmt.Println(err)
		return
	}

	var response []byte
	switch msg.MessageType {
	case TypeBindingRequest:
		response = stunMessageHandle(msg,raddr,false)
	case TypeAllocate:
		response = turnMessageHandle(msg,raddr,false)
	}

	if !tcp {
		if response != nil {
			_, err := entry.udpConn.WriteToUDP(response, raddr)
			if err != nil {
				fmt.Println(err)
			}
		}else {
			fmt.Println("no response.")
		}

	}else {
		//todo : add tcp
	}
}
