package libs

import (
	"net"
	"encoding/binary"
	"fmt"
	"crypto/md5"
	"time"
	"encoding/hex"
)

type Attribute struct{
	AttrType uint16
	Length uint16
	Value []byte
}

var (
	AttrSoftware = newAttr(AttributeSoftware,[]byte("Example server, version 1.17"))
	AttrRealm = newAttr(AttributeRealm,[]byte("realm"))
	AttrError401 = newAttrError401()
	AttrLifetime = newAttrLifetime()
	AttrDummyMessageIntegrity = newAttrNoValue(AttributeMessageIntegrity)
)

func xorAddress(port uint16, addr []byte) []byte {

	xport := make([]byte, 2)
	xaddr := make([]byte, 4)
	binary.BigEndian.PutUint16(xport, port^uint16(magicCookie>>16))
	binary.BigEndian.PutUint32(xaddr, binary.BigEndian.Uint32(addr)^magicCookie)
	return append(xport, xaddr...)

}

func unXorAddress(xorAddress []byte) (port uint16,addr []byte)  {
	addr = make([]byte,4)
	port = binary.BigEndian.Uint16(xorAddress[2:4])^uint16(magicCookie>>16)
	binary.BigEndian.PutUint32(addr,binary.BigEndian.Uint32(xorAddress[4:8])^magicCookie)
	return
}

func padding(bytes []byte) []byte {
	length := uint16(len(bytes))
	return append(bytes, make([]byte, align(length)-length)...)
}

func align(n uint16) uint16 {
	return (n + 3) & 0xfffc
}

// https://tools.ietf.org/html/rfc5389#page-32
func newAttr(attrType uint16,value []byte) *Attribute {
	att := new(Attribute)
	att.AttrType = attrType
	att.Length = uint16(len(value))
	att.Value = value
	return att
}

func newAttrNoValue(attrType uint16) *Attribute {
	att := new(Attribute)
	att.AttrType = attrType
	att.Value = make([]byte,20)
	att.Length = uint16(len(att.Value))
	return att
}

func newAttrMappedAddress(remoteAddress *net.UDPAddr) *Attribute  {
	port := make([]byte, 2)
	binary.BigEndian.PutUint16(port, uint16(remoteAddress.Port))
	reflexiveAddress := remoteAddress.IP.To4()
	value := append([]byte{0, attributeFamilyIPv4}, append(port, reflexiveAddress...)...)
	return newAttr(AttributeMappedAddress,value)
}

func newAttrXORMappedAddress(remoteAddress *net.UDPAddr) *Attribute  {
	port := uint16(remoteAddress.Port)
	reflexiveAddress := remoteAddress.IP.To4()
	//reflexiveAddress := net.ParseIP("11.11.11.11").To4()
	xorBytes := xorAddress(port, reflexiveAddress)

	value := append([]byte{0, attributeFamilyIPv4}, xorBytes...)
	return newAttr(AttributeXorMappedAddress,value)
}

func newAttrNonce() *Attribute{
	//fixme : 20 min expire
	timestampBytes := make([]byte, 4)
	timestamp := time.Now().Unix() + 20*60
	binary.BigEndian.PutUint32(timestampBytes, uint32(timestamp^magicCookie))
	nonce := hex.EncodeToString(timestampBytes)
	return newAttr(AttributeNonce,[]byte(nonce))
}

func validNonce(nonce []byte) bool{
	step1,err  :=  hex.DecodeString(string(nonce))

	if err != nil {
		return  false
	}
	timestamp := binary.BigEndian.Uint32(step1)^magicCookie

	if timestamp > uint32(time.Now().Unix()){
		return true
	}else{
		return false
	}
}

func newAttrError401() *Attribute{
	reason := "Unauthorized"
	error401 := make([]byte,4 + len([]byte(reason)))
	error401[0] = 0;
	error401[1] = 0;
	error401[2] = uint8(401 / 100)
	error401[3] = uint8(401 % 100)
	error401 = append(error401[:4],[]byte(reason)...)
	return newAttr(AttributeErrorCode,error401)
}

func getRelayAddress() (raddr string) {
	if external_ip != nil{
		if IsValidIPv4(*external_ip) {
			raddr = *external_ip
			return
		}
	}

	ipAddress , err := HostIP()
	if err != nil {
		raddr = "110.110.110.110"
		Log.Error("Can not find relay address")
	}else{
		raddr = ipAddress
	}
	return
}

func newAttrXORRelayedAddress(relay string ,rport int) *Attribute{
	//relayedAddress := net.ParseIP("22.22.22.22").To4()
	relayedAddress := net.ParseIP(relay).To4()
	port := uint16(rport)
	xorBytes := xorAddress(port, relayedAddress)
	value := append([]byte{0, attributeFamilyIPv4}, xorBytes...)
	return newAttr(AttributeXorRelayedAddress,value)
}

func newAttrXORPeerAddress(peer string, pport int) *Attribute {
	relayedAddress := net.ParseIP(peer).To4()
	port := uint16(pport)
	xorBytes := xorAddress(port, relayedAddress)
	value := append([]byte{0, attributeFamilyIPv4}, xorBytes...)
	return newAttr(AttributeXorPeerAddress,value)
}


func newAttrLifetime() *Attribute {
	time := make([]byte,4)
	binary.BigEndian.PutUint32(time,600)
	return newAttr(AttributeLifetime, time)
}

func newAttrMessageIntegrity(value []byte) *Attribute {
	return newAttr(
		AttributeMessageIntegrity,value)
}



// message-integrity
func generateKey(username,password,realm string) []byte  {
	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprintf("%s:%s:%s",username,realm,password)))
	key := hasher.Sum(nil)
	return key
}

func MessageIntegrityHmac(value,key []byte) []byte {
	return HmacSha1(value,key)
}



func  AttrTypeToString(attrType uint16) (typeString string)  {
	switch attrType {
	case AttributeMappedAddress:
		typeString = "MappedAddress"
	case AttributeResponseAddress:
		typeString = "ResponseAddress"
	case AttributeChangeRequest:
		typeString = "ChangeRequest"
	case AttributeSourceAddress:
		typeString = "SourceAddress"
	case AttributeChangedAddress:
		typeString = "ChangedAddress"
	case  AttributeUsername:
		typeString = "Username"
	case  AttributePassword:
		typeString = "Password"
	case AttributeMessageIntegrity:
		typeString = "MessageIntegrity"
	case AttributeErrorCode:
		typeString = "ErrorCode"
	case AttributeUnknownAttributes:
		typeString = "UnknownAttributes"
	case AttributeReflectedFrom:
		typeString = "ReflectedFrom"
	case AttributeChannelNumber:
		typeString = "ChannelNumber"
	case AttributeLifetime:
		typeString = "Lifetime"
	case AttributeBandwidth:
		typeString = "Bandwidth"
	case AttributeXorPeerAddress:
		typeString = "XorPeerAddress"
	case AttributeData:
		typeString = "Data"
	case AttributeRealm:
		typeString = "Realm"
	case AttributeNonce:
		typeString = "Nonce"
	case AttributeXorRelayedAddress:
		typeString = "XorRelayedAddress"
	case AttributeRequestedAddressFamily:
		typeString = "RequestedAddressFamily"
	case AttributeEvenPort:
		typeString = "EvenPort"
	case AttributeRequestedTransport:
		typeString = "RequestedTransport"
	case AttributeDontFragment:
		typeString = "DontFragment"
	case AttributeXorMappedAddress:
		typeString = "XorMappedAddress"
	case AttributeTimerVal:
		typeString = "TimerVal"
	case AttributeReservationToken:
		typeString = "ReservationToken"
	case AttributePriority:
		typeString = "Priority"
	case AttributeUseCandidate:
		typeString = "UseCandidate"
	case AttributePadding:
		typeString = "Padding"
	case AttributeResponsePort:
		typeString = "ResponsePort"
	case AttributeConnectionID:
		typeString = "ConnectionID"
	case AttributeXorMappedAddressExp:
		typeString = "XorMappedAddressExp"
	case AttributeSoftware:
		typeString = "Software"
	case AttributeAlternateServer:
		typeString = "AlternateServer"
	case AttributeCacheTimeout:
		typeString = "CacheTimeout"
	case AttributeFingerprint:
		typeString = "Fingerprint"
	case AttributeIceControlled:
		typeString = "IceControlled"
	case AttributeIceControlling:
		typeString = "IceControlling"
	case AttributeResponseOrigin:
		typeString = "ResponseOrigin"
	case AttributeOtherAddress:
		typeString = "OtherAddress"
	case AttributeEcnCheckStun:
		typeString = "EcnCheckStun"
	case AttributeCiscoFlowdata:
		typeString = "CiscoFlowdata"
	case AttributeOrigin:
		typeString = "Origin"
	case AttributeNetworkInfo :
		typeString = "NetworkInfo"
	default:
		typeString = "fuck??"
	}

	return
}

func (a Attribute) String() string {
	attrString := ""
	switch a.AttrType {
	case AttributeRequestedTransport,AttributePriority,
		AttributeIceControlled,AttributeIceControlling:
		attrString = fmt.Sprintf("	attr: type -> %s , length -> %d , value -> %d \n",
			AttrTypeToString(a.AttrType), a.Length,  uint8(a.Value[0]) )
	case AttributeLifetime:
		attrString = fmt.Sprintf("	attr: type -> %s , length -> %d , value -> %d \n",
			AttrTypeToString(a.AttrType), a.Length,  binary.BigEndian.Uint32(a.Value) )
	case AttributeMessageIntegrity,AttributeFingerprint:
		attrString = fmt.Sprintf("	attr: type -> %s , length -> %d , value -> %x \n",
			AttrTypeToString(a.AttrType), a.Length,  a.Value )
	case AttributeXorPeerAddress:
		port ,addr :=  unXorAddress(a.Value)
		attrString = fmt.Sprintf("	attr: type -> %s , length -> %d , value -> %s:%d \n",
			AttrTypeToString(a.AttrType), a.Length,  net.IP(addr),port )
	default:
		attrString = fmt.Sprintf("	attr: type -> %s , length -> %d , value -> %s \n",
			AttrTypeToString(a.AttrType), a.Length, a.Value)
	}

	return attrString
}


