package v2

import (
	"crypto/sha1"
	"encoding/hex"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// UserAccount representation of user account entity
type UserAccount struct {
	ID               bson.ObjectId `json:"id,omitempty" bson:"_id"`
	NameIdentity     string        `json:"nameIdentity,required" bson:"nameIdentity"`
	IdentityProvider string        `json:"identityProvider,required" bson:"identityProvider"`
	Password         string        `json:"-" bson:"password"`
	ProfileID        bson.ObjectId `json:"profileID,required" bson:"profileID"`
	Token            string        `json:"-" bson:"token"`
	Email            string        `json:"email" bson:"email"`
	LastLoginDate    time.Time     `json:"lastLoginDate" bson:"lastLoginDate"`
	LoginAttempts    int16         `json:"loginAttempts" bson:"loginAttempts"`
	LockUntil        time.Time     `json:"lockUntil" bson:"lockUntil"`
	IsLockedOut      bool          `json:"isLockedOut" bson:"isLockedOut"`
	IsApproved       bool          `json:"isApproved" bson:"isApproved"`
	CreatedBy        string        `json:"createdBy" bson:"createdBy"`
	CreatedDate      time.Time     `json:"createdDate,required" bson:"createdDate"`
	UpdatedBy        string        `json:"updatedBy" bson:"updatedBy"`
	UpdatedDate      time.Time     `json:"updatedDate,required" bson:"updatedDate"`
}

// IsLocked check is current user locked
func (ua *UserAccount) IsLocked() bool {
	return ua.IsLockedOut || ua.LockUntil.After(time.Now())
}

// SetPassword hash and set password to UserAccount instance
func (ua *UserAccount) SetPassword(pass string) {
	ua.Password = HashPassword(pass)
}

// ComparePassword compare is candidatePassword equal to existing password
func (ua *UserAccount) ComparePassword(candidatePassword string) bool {
	return ua.Password == HashPassword(candidatePassword)
}

// HashPassword returns given string as SHA1 hash
func HashPassword(pass string) string {
	h := sha1.New()
	h.Write([]byte(pass))
	sha1Hash := hex.EncodeToString(h.Sum(nil))
	return sha1Hash
}
