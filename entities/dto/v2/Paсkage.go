package v2

import (
	"fmt"

	"github.com/golang/protobuf/proto"
)

var _ = proto.Marshal

// Package
type Package struct {
	MessageID string `protobuf:"bytes,1,opt,name=messageID" json:"messageID,omitempty"`
	PackageNo int32  `protobuf:"varint,2,opt,name=packageNo" json:"packageNo,omitempty"`
	Closed    bool   `protobuf:"varint,3,opt,name=closed" json:"closed,omitempty"`
	Data      []byte `protobuf:"bytes,4,opt,name=data" json:"data,omitempty"`
}

func (p *Package) ToString() string {
	return fmt.Sprintf("Package@ MessageID: %v, PackageNo: %v, Closed: %v.", p.MessageID, p.PackageNo, p.Closed)
}

func (p *Package) ToFullString() string {
	return fmt.Sprintf("Package@ MessageID: %v, PackageNo: %v, Closed: %v, Data: %v.", p.MessageID, p.PackageNo, p.Closed, p.Data)
}

/* Required for proto.Message */
// Reset required for proto.Message
func (p *Package) Reset() {
	*p = Package{}
}

// String required for proto.Message
func (p *Package) String() string {
	return proto.CompactTextString(p)
}

// ProtoMessage required for proto.Message
func (*Package) ProtoMessage() {}

func init() {
}
