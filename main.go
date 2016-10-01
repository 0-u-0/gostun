package main

import (
	"git.learning-tech.cn/realtimecat/gostun/libs"
)

func main() {
	server := libs.Server{Port:3478}
	server.Serve()

}
