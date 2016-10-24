package libs

import (
	"net"
	"time"
	"fmt"
)


const (
	MinTimeRefresh = 300
	MaxTimeRefresh = 600
)

const (
	ClientProtocolUDP = 0
	ClientProtocolTCP
	ClientProtocolTCPTLS
)


type Allocate struct {
	UserName string
	ExpiresTime int
	IsExpired bool
	ConnectTime int
	ExpiresTicker *time.Ticker
	Nonce string
	ClientProtocol int
	ClientAddress *net.UDPAddr
	RelayAddress *net.UDPAddr
	PeerRelayAddress *net.UDPAddr
	ByteSend int
	ByteRecv int
}

func NewAllocate(username ,nonce string, protocol, lifetime int,client,relay *net.UDPAddr) *Allocate  {
	allocate := &Allocate{
		UserName:username,
		IsExpired:false,
		ConnectTime:0,
		ExpiresTime:lifetime,
		ExpiresTicker:time.NewTicker(1 * time.Second),
		Nonce:nonce,
		ClientProtocol:protocol,
		ClientAddress:client,
		RelayAddress:relay,
		ByteRecv:0,
		ByteSend:0,
	}

	go func() {
		for range allocate.ExpiresTicker.C {
			allocate.ConnectTime++;
			allocate.ExpiresTime--;
			if allocate.ExpiresTime <= 0 {

				break
			}

		}
	}()

	return allocate
}

func (a *Allocate)Refresh(time int)  {
	if time > MinTimeRefresh && time < MaxTimeRefresh {
		a.ExpiresTime = time
	}else if time < MinTimeRefresh {
		a.ExpiresTime = MinTimeRefresh
	}else if time > MaxTimeRefresh {
		a.ExpiresTime = MaxTimeRefresh
	}
}


