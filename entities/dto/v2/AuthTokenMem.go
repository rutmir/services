package v2

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"
	"github.com/golang/protobuf/proto"
)

var _ = proto.Marshal

// AuthTokenMem
type AuthTokenMem struct {
	AccessToken  string        `protobuf:"bytes,1,opt,name=accessToken" json:"access_token,omitempty"`
	RefreshToken string        `protobuf:"bytes,2,opt,name=refreshToken" json:"refresh_token,omitempty"`
	ClientID     string        `protobuf:"bytes,3,rep,name=clientID" json:"clientID,omitempty"`
	AccountID    bson.ObjectId `protobuf:"bytes,4,rep,name=accountID" json:"accountID,omitempty"`
	ProfileID    bson.ObjectId `protobuf:"bytes,5,rep,name=profileID" json:"profileID,omitempty"`
}

func (atm *AuthTokenMem) ToString() string {
	return fmt.Sprintf("AuthTokenMem@ AccessToken: %v, RefreshToken: %s, ClientID: %s, AccountID: %v, ProfileID: %v.", atm.AccessToken, atm.RefreshToken, atm.ClientID, atm.AccountID, atm.ProfileID)
}

/* Required for proto.Message */
// Reset required for proto.Message
func (atm *AuthTokenMem) Reset() {
	*atm = AuthTokenMem{}
}

// String required for proto.Message
func (atm *AuthTokenMem) String() string {
	return proto.CompactTextString(atm)
}

// ProtoMessage required for proto.Message
func (*AuthTokenMem) ProtoMessage() {}

func init() {
}
