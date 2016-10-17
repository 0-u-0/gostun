package libs

import (
	"net"
	"strconv"
	"log"
	"fmt"
)

type RelayServer struct {
	Port       int
	connection *net.UDPConn
}

func (s *RelayServer) serve() {
	for {
		var buf = make([]byte, 1024)
		size, remoteAddr, err := s.connection.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		go s.handleData(remoteAddr, buf[:size])
	}
}

//Serve initiates a UDP connection that listens on any port for incoming data
func (s *RelayServer) Serve() {
	laddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(s.Port))
	if err != nil {
		log.Fatal(err)
	}
	s.connection, err = net.ListenUDP("udp", laddr)
	if err != nil {
		log.Fatal(err)
	}
	go s.serve()
}

//NewServer conveniently creates a new server from the given port
func NewRelayServer(port int) *RelayServer {
	ret := new(RelayServer)
	ret.Port = port

	return ret
}

func (s *RelayServer) handleData(raddr *net.UDPAddr, data []byte) {
	msg , err := UnMarshal(data)
	if err != nil {
		return
	}

	fmt.Printf("ffffff request : %s \n",msg)

}

