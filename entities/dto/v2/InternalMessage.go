package v2

import (
	"fmt"

	"github.com/golang/protobuf/proto"
)

var _ = proto.Marshal

// InternalMessage
type InternalMessage struct {
	Header *Header `protobuf:"bytes,1,req,name=header" json:"header,omitempty"`
	Body   []byte `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
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

// GetRequiredHeader required for proto.Message
func (im *InternalMessage) GetHeader() *Header {
	if im != nil {
		return im.Header
	}
	return nil
}

func init() {
	proto.RegisterType((*InternalMessage)(nil), "InternalMessage")
}
