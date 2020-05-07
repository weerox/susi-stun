package stun

import (
	"crypto/rand"
	"encoding/binary"
	"io"
)

type Method uint16
type Class byte

const (
	Binding Method = 0b000000000001
)

// TODO: check if the Class identifier can be removed for everyone except the first
const (
	Request    Class = 0b00
	Indication Class = 0b01
	Success    Class = 0b10
	Error      Class = 0b11
)

const magicCookie = 0x2112A442

type MessageType struct {
	method Method
	class  Class
}

func NewMessageType(method Method, class Class) *MessageType {
	return &MessageType{method, class}
}

type Message struct {
	typ           MessageType
	length        uint16
	magicCookie   uint32
	transactionID [12]byte
	attributes    []Attribute
}

func NewMessage(typ MessageType) *Message {
	message := Message{typ: typ, magicCookie: magicCookie}
	_, err := rand.Read(message.transactionID[:])
	if err != nil {
		// TODO: handle error
	}
	message.attributes = nil
	message.length = 0
	return &message
}

func (m *Message) AddAttribute(attr Attribute) {
	m.attributes = append(m.attributes, attr)
}

// The net.Conn interface has both Read and Write methods which both takes
// a byte array as a single argument.

// The net.Conn interface can be used as an io.Reader interface.

// NewMessageFromConn:
// Create a new message by calling `message := stun.NewMessageFromConn(conn)`

// Decoder:
// Create a Decoder struct and a NewDecoder function which will take
// an io.Reader as the argument. Then a message can be constructed from
// a net.Conn by `message := stun.NewDecoder(conn).Decode()`

type Decoder struct {
	r io.Reader
}

type Encoder struct {
	w io.Writer
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r}
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}

func (d *Decoder) Decode() *Message {
	var message = &Message{}
	var header [20]byte

	_, err := d.r.Read(header[:])
	// TODO: also check that we actually read 20 bytes
	if err != nil {
		// TODO: handle error
	}

	// TODO: check that this is actually a valid STUN message
	// and handle magicCookie

	// Parse the message type
	// See Figure 3 in RFC 5389
	var mt MessageType
	var t uint16 = binary.BigEndian.Uint16(header[0:2])

	var m [3]uint16 // will hold the three different part of the method
	var c [2]byte   // will hold the two different parts of the class

	m[0] = t & 0b0000000000001111
	m[1] = (t >> 1) & 0b0000000001110000
	m[2] = (t >> 2) & 0b0000111110000000

	mt.method = Method(m[0] + m[1] + m[2])

	c[0] = byte((t >> 4) & 0b01)
	c[1] = byte((t >> 7) & 0b10)

	mt.class = Class(c[0] + c[1])

	message.typ = mt

	// Store the message length
	message.length = binary.BigEndian.Uint16(header[2:4])

	// Store the transaction ID
	copy(message.transactionID[:], header[8:])

	// Read message.length number of bytes
	var body []byte = make([]byte, message.length)

	_, err = d.r.Read(body)
	// TODO: check that we actually read message.length number of bytes
	if err != nil {
		// TODO: handle error
	}

	// TODO: Decode attributes

	// NOTE: Unknown comprehension-optional attributes MUST be ignored,
	// and unknown comprehension-required attributes MUST result in
	// an error response, see last two paragraphs in section 7.3
	// and first paragraph of section 7.3.1 in RFC 5389.
	// Known-but-unexpected attributes SHOULD be ignored.
	// This means that we have to perform some sort of check when decoding
	// the attributes.

	var pos uint16 = 0

	for pos < message.length {
		// TODO: should the NewAttribute function be used instead?
		var attr Attribute
		var t, l uint16
		t = binary.BigEndian.Uint16(body[pos : pos+2])
		pos += 2
		l = binary.BigEndian.Uint16(body[pos : pos+2])
		pos += 2

		attr.typ = AttributeType(t)
		attr.length = l
		attr.value = make([]byte, l)
		copy(attr.value, body)
		//TODO: check that the correct length is copied

		// Consume the padding bytes
		padding := l % 4
		pos += padding

		message.attributes = append(message.attributes, attr)
	}

	return message
}
