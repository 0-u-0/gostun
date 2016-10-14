package main

import (
	"git.learning-tech.cn/realtimecat/gostun/libs"
)

func main() {

	relay := libs.RelayServer{Port:33333}
	relay.Serve()

	server := libs.Entry{Port:3478}
	server.Serve()

}
