package v2

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// UserProfile representation user profile entity
type UserProfile struct {
	ID             bson.ObjectId       `json:"id,omitempty" bson:"_id"`
	DisplayName    string              `json:"displayName,required" bson:"displayName"`
	Name           NameBlock           `json:"name" bson:"name"`
	Communications []CommunicationItem `json:"communications" bson:"communications"`
	Photos         []string            `json:"photos" bson:"photos"`
	CreatedBy      string              `json:"createdBy" bson:"createdBy"`
	CreatedDate    time.Time           `json:"createdDate,required" bson:"createdDate"`
	UpdatedBy      string              `json:"updatedBy" bson:"updatedBy"`
	UpdatedDate    time.Time           `json:"updatedDate,required" bson:"updatedDate"`
	OnlineStatus   int16               `json:"onlineStatus" bson:"onlineStatus"`
	Contacts       []ContactItem       `json:"contacts" bson:"contacts"`
	BlackList      []ContactItem       `json:"blackList" bson:"blackList"`
	Roles          []string            `json:"roles" bson:"roles"`
}
