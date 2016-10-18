package main

import (
	"git.learning-tech.cn/realtimecat/gostun/libs"
)

func main() {


	server := libs.NewRelayServer(33333)
	server.Serve()
	libs.Init()

}
