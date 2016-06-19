package v2

import (
	"fmt"

	"github.com/golang/protobuf/proto"
)

var _ = proto.Marshal

// Header
type Header struct {
	Timestamp int64  `protobuf:"varint,1,opt,name=timestamp" json:"timestamp,omitempty"`
	Action    Action `protobuf:"bytes,2,opt,name=action,json=action,enum=v2.Action" json:"action,omitempty"`
	Meta      string `protobuf:"bytes,3,opt,name=meta" json:"meta,omitempty"`
}

func (h *Header) ToString() string {
	return fmt.Sprintf("Header@ Timestamp: %v, Action: %v, Meta: %v.", h.Timestamp, h.Action, h.Meta)
}

/* Required for proto.Message */
// Reset required for proto.Message
func (h *Header) Reset() {
	*h = Header{}
}

// String required for proto.Message
func (h *Header) String() string {
	return proto.CompactTextString(h)
}

// ProtoMessage required for proto.Message
func (*Header) ProtoMessage() {}

func init() {
}
