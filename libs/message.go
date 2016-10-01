package libs

import (
	"encoding/binary"
	"errors"
	"net"
	"math/big"
	"fmt"
)

//ErrX are the errors to be expected during message handling
var (
	ErrInvalidRequest = errors.New("Invalid STUN request")
	ErrRFC3489 = errors.New("no magic cookie , RFC3489")
)

//Message holds the information about a STUN Message
type Message struct {
	MessageType   uint16
	MessageLength uint16
	TransID       *big.Int
	Attributes    map[uint16][]byte
}

//UnMarshal creates a Message object from data received by the STUN server
func UnMarshal(data []byte) (*Message, error) {
	length := len(data)
	if length < 20 {
		return nil, ErrInvalidRequest
	}

	pkgType := binary.BigEndian.Uint16(data[0:2])
	// check 00
	if pkgType >  (1 << 15 - 1 ) {
		return nil, ErrInvalidRequest
	}

	//check magic cookie
	magicCookieCheck := binary.BigEndian.Uint32(data[4:8]);
	if(magicCookie != magicCookieCheck){
		return nil, ErrRFC3489
	}

	msg := new(Message)

	//parse the header
	msg.MessageType = pkgType
	msg.MessageLength = binary.BigEndian.Uint16(data[2:4])


	tid := new(big.Int)
	tid.SetBytes(data[4:20])
	msg.TransID = tid

	//if we have leftover data, parse as attributes
	if length > 20 {
		msg.Attributes = make(map[uint16][]byte)
		i := 20
		for i < length {
			attrType := binary.BigEndian.Uint16(data[i : i+2])
			attrLength := binary.BigEndian.Uint16(data[i+2 : i+4])
			i += 4 + int(attrLength)
			msg.Attributes[attrType] = data[i-int(attrLength) : i]
			if pad := int(attrLength) % 4; pad > 0 {
				i += 4 - pad
			}
		}
		//recover here to catch any index errors
		if recover() != nil {
			return nil, ErrInvalidRequest
		}
	}

	return msg, nil
}

//Marshal transforms a message into a byte array
func Marshal(m *Message) ([]byte, error) {
	result := make([]byte, 36)
	//first do the header
	binary.BigEndian.PutUint16(result[:2], m.MessageType)
	result = append(result[:4], m.TransID.Bytes()...)

	//now we do the attributes
	if m.Attributes != nil {
		i := 20
		for t, v := range m.Attributes {
			length := len(v)
			binary.BigEndian.PutUint16(result[i:i+2], t)
			binary.BigEndian.PutUint16(result[i+2:i+4], uint16(len(v)))
			result = append(result[:i+4], v...)
			i += 4 + length
			//if we need to pad, do so
			if pad := length % 4; pad > 0 {
				result = append(result, make([]byte, 4-pad)...)
				i += 4 - pad
			}
		}

		//add length
		binary.BigEndian.PutUint16(result[2:4], uint16(i-20))
	}
	return result, nil
}

func addMappedAddress(m *Message, raddr *net.UDPAddr) {
	port := make([]byte, 2)
	binary.BigEndian.PutUint16(port, uint16(raddr.Port))
	addr := raddr.IP.To4()
	m.Attributes[AttributeMappedAddress] = append([]byte{0, attributeFamilyIPv4}, append(port, addr...)...)
}

func addXORMappedAddress(m *Message, raddr *net.UDPAddr) {

	//addr := raddr.IP.To4()
	addr := net.ParseIP("11.11.11.11").To4()
	port := uint16(raddr.Port)
	xbytes := xorAddress(port, addr)
	m.Attributes[AttributeXorMappedAddress] = append([]byte{0, attributeFamilyIPv4}, xbytes...)

}

func xorAddress(port uint16, addr []byte) []byte {

	xport := make([]byte, 2)
	xaddr := make([]byte, 4)
	binary.BigEndian.PutUint16(xport, port^uint16(magicCookie>>16))
	binary.BigEndian.PutUint32(xaddr, binary.BigEndian.Uint32(addr)^magicCookie)
	return append(xport, xaddr...)

}

func (m Message) TypeToString() (typeString string)  {
	switch m.MessageType {
	case TypeBindingRequest:
		typeString = "BindRequest"
	case TypeAllocate:
		typeString = "Allocate"
	case TypeBindingResponse:
		typeString = "BindingResponse"

	}
	return
}

func (m Message) String() string {
	attrString:= ""
	if len(m.Attributes) > 0 {
		attrString = "\n attributes : \n"

		for k,v := range m.Attributes {
			attrString += fmt.Sprintf(`	attr: type -> %s , length -> %d , value -> %s`,
				attrTypeToString(k), len(v), v)
		}
	}

	return fmt.Sprintf(`packet : type -> %s , length -> %d , tid -> %d , length of the attr -> %d	%s
			 `,
		m.TypeToString(),m.MessageLength,m.TransID,len(m.Attributes),attrString)
}

func attrTypeToString(attrType uint16) (typeString string)  {
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
	default:
		typeString = "fuck??"
	}

	return
}