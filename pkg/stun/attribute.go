package stun

type AttributeType uint16

const (
	MappedAddress     AttributeType = 0x0001
	Username                        = 0x0006
	MessageIntegrity                = 0x0008
	ErrorCode                       = 0x0009
	UnknownAttributes               = 0x000A
	Realm                           = 0x0014
	Nonce                           = 0x0015
	XORMappedAddress                = 0x0020
	Software                        = 0x8022
	AlternateServer                 = 0x8023
	Fingerprint                     = 0x8028
)

type Attribute struct {
	typ    AttributeType
	length uint16
	value  []byte
}

func NewAttribute(typ AttributeType, value []byte) *Attribute {
	return &Attribute{typ, uint16(len(value)), value}
}
