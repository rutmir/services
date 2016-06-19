package v2

import (
	"fmt"

	"github.com/golang/protobuf/proto"
)

var _ = proto.Marshal

// InternalMessage
type InternalMessage struct {
	Header Header `protobuf:"varint,1,opt,name=header" json:"header,omitempty"`
	Body   []byte `protobuf:"bytes,2,opt,name=body" json:"body,omitempty"`
}

func (im *InternalMessage) ToString() string {
	return fmt.Sprintf("InternalMessage@ Header: %v, Body: %v, Meta: %v.", im.Header, im.Body)
}

/* Required for proto.Message */
// Reset required for proto.Message
func (im *InternalMessage) Reset() {
	*im = InternalMessage{}
}

// String required for proto.Message
func (im *InternalMessage) String() string {
	return proto.CompactTextString(im)
}

// ProtoMessage required for proto.Message
func (*InternalMessage) ProtoMessage() {}

func init() {
}
