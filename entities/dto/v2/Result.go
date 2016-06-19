package v2

import (
	"fmt"

	"github.com/golang/protobuf/proto"
)

var _ = proto.Marshal

// Header
type Result struct {
	Code    int32  `protobuf:"varint,1,opt,name=code" json:"code,omitempty"`
	Result  string `protobuf:"bytes,2,opt,name=result" json:"result,omitempty"`
	Message string `protobuf:"bytes,3,opt,name=message" json:"message,omitempty"`
}

func (r *Result) ToString() string {
	return fmt.Sprintf("Result@ Code: %v, Result: %v, Message: %v.", r.Code, r.Result, r.Message)
}

/* Required for proto.Message */
// Reset required for proto.Message
func (r *Result) Reset() {
	*r = Result{}
}

// String required for proto.Message
func (r *Result) String() string {
	return proto.CompactTextString(r)
}

// ProtoMessage required for proto.Message
func (*Result) ProtoMessage() {}

func init() {
}
