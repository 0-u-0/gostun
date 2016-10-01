package main

import (
	"./libs"
)
func main() {
	server := libs.Server{Port:3478}
	server.Serve()

}
