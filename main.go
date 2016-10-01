package main

import (
	"github.com/zombiecong/golang-stun/libs"
)
func main() {
	server := libs.Server{Port:3478}
	server.Serve()

}
