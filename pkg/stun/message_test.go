package stun

import "testing"
import "bytes"

func TestDecodeEmptyMessage(t *testing.T) {
	data := make([]byte, 0)
	dec := NewDecoder(bytes.NewReader(data))
	message := dec.Decode()
	//TODO: should return an error

	if message == nil {
		t.Errorf("message is nil\n")
	}
}

func TestDecodeZeroAttributes(t *testing.T) {
	data := []byte{
		0x00, 0x01, 0x00, 0x00,
		0x21, 0x12, 0xA4, 0x42,
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06,
		0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C,
	}
	dec := NewDecoder(bytes.NewReader(data))
	message := dec.Decode()

	if message.length != 0 {
		t.Errorf("Message Length is %d, expected 0\n", message.length)
	}

	if len(message.attributes) != 0 {
		t.Errorf(
			"Number of attributes is %d, expected %d\n",
			len(message.attributes),
			0,
		)
	}
}

func TestDecodeMethod(t *testing.T) {
	data := []byte{
		0x00, 0x01, 0x00, 0x00,
		0x21, 0x12, 0xA4, 0x42,
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06,
		0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C,
	}
	dec := NewDecoder(bytes.NewReader(data))
	message := dec.Decode()

	// TODO: test all types of methods
	if message.typ.method != Binding {
		t.Errorf("Message Method is %v, expected %v\n", message.typ.method, Binding)
	}
}

func TestDecodeClass(t *testing.T) {
	data := []byte{
		0x00, 0x01, 0x00, 0x00,
		0x21, 0x12, 0xA4, 0x42,
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06,
		0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C,
	}
	dec := NewDecoder(bytes.NewReader(data))
	message := dec.Decode()

	//TODO: test all types of classes
	if message.typ.class != Request {
		t.Errorf("Message Class is %v, expected %v\n", message.typ.class, Request)
	}
}

func TestDecodeTransactionID(t *testing.T) {
	data := []byte{
		0x00, 0x01, 0x00, 0x00,
		0x21, 0x12, 0xA4, 0x42,
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06,
		0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C,
	}
	dec := NewDecoder(bytes.NewReader(data))
	message := dec.Decode()

	if len(message.transactionID) != 12 {
		t.Errorf(
			"Length of Transaction ID is %d, expected %d\n",
			message.transactionID,
			12,
		)
	}

	if message.transactionID[0] != 0x01 || message.transactionID[11] != 0x0C {
		t.Errorf(
			"Transaction ID is %X, expected %X\n",
			message.transactionID,
			[]byte{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06,
				0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C,
			},
		)
	}
}
