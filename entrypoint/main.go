package main

import (
	"git.learning-tech.cn/realtimecat/gostun/libs"
)

func main() {

	entry := libs.NewEntry(3478)
	entry.Serve()

}
