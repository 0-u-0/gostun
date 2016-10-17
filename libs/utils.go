package libs

import (
	"log"
	"crypto/rand"
)

var counter uint32


func RandBytes(length int) (r []byte) {
	if length < 64 {
		r = make([]byte, length)
		_, err := rand.Read(r)

		if err != nil {
			log.Panicln(err)
		}
	}else {
		log.Panicf("the max length of randbyte is 64 , %d not supported \n",length)
	}
	return
}



func PrintModuleLoaded(moduleName string)  {
	log.Printf("< %s > module loads successfully",moduleName)
}

func PrintModuleRelease(moduleName string)  {
	log.Printf("< %s > module releases successfully",moduleName)
}

