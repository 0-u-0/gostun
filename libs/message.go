package libs

import (
	"encoding/binary"
	"errors"
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
	Attributes    []*Attribute
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
		msg.Attributes = make([]*Attribute,0)
		i := 20
		for i < length {
			attrType := binary.BigEndian.Uint16(data[i : i+2])
			attrLength := binary.BigEndian.Uint16(data[i+2 : i+4])
			i += 4 + int(attrLength)

			msg.Attributes = append(msg.Attributes,newAttr(attrType,data[i-int(attrLength) : i]))

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
	result := make([]byte, 60)
	//first do the header
	binary.BigEndian.PutUint16(result[:2], m.MessageType)
	result = append(result[:4], m.TransID.Bytes()...)

	//now we do the attributes
	if m.Attributes != nil {
		i := 20
		for _ , attr := range m.Attributes {
			binary.BigEndian.PutUint16(result[i:i+2], attr.AttrType)
			binary.BigEndian.PutUint16(result[i+2:i+4], attr.Length)
			result = append(result[:i+4], attr.Value...)
			i += 4 + int(attr.Length)
			//if we need to pad, do so
			if pad := int(attr.Length % 4); pad > 0 {
				result = append(result, make([]byte, 4-pad)...)
				i += 4 - pad
			}
		}

		//add length
		binary.BigEndian.PutUint16(result[2:4], uint16(i-20))
	}
	return result, nil
}

func (m Message) hasAttribute(attrType uint16) bool  {
	for _, a := range m.Attributes {
		if a.AttrType == attrType {
			return true
		}
	}
	return false
}

func (m Message) getAttribute(attrType uint16) *Attribute  {
	for _, a := range m.Attributes {
		if a.AttrType == attrType {
			return a
		}
	}
	return nil
}


func (m Message) TypeToString() (typeString string)  {
	switch m.MessageType {
	case TypeBindingRequest:
		typeString = "BindRequest"
	case TypeAllocate:
		typeString = "Allocate"
	case TypeBindingResponse:
		typeString = "BindingResponse"
	case TypeAllocateErrorResponse:
		typeString = "AllocateErrorResponse"
	case TypeAllocateResponse:
		typeString = "AllocateResponse"
	}
	return
}

func (m Message) String() string {

	attrString := ""
	if len(m.Attributes) > 0{
		attrString = "\n Attributes : \n"

		for _ , attr := range m.Attributes{
			attrString += attr.String()
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