package libs

import (
	"net"
	"time"
)


const (
	MinTimeRefresh = 300
	MaxTimeRefresh = 600
)

const (
	ClientProtocolUDP = 17
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
	ClientProtocol uint8
	ClientAddress *net.UDPAddr
	RelayAddress *net.UDPAddr
	PeerRelayAddress *net.UDPAddr
	ByteSend int
	ByteRecv int
	Channels []Channel
	Permissions []Permission
}

//todo: init channel and permission
func NewAllocate(username ,nonce string, protocol uint8, lifetime int,client,relay *net.UDPAddr) *Allocate  {
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



	return allocate
}

func (a *Allocate)TimerRun()  {
	go func() {
		for range a.ExpiresTicker.C {
			a.ConnectTime++;
			a.ExpiresTime--;
			if a.ExpiresTime <= 0 {
				a.IsExpired = true
				break
			}
		}
	}()
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


