package v2

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// FileLink representation of FileLink entity
type FileLink struct {
	ID          bson.ObjectId `json:"id,omitempty" bson:"_id"`
	FileID      bson.ObjectId `json:"fileID,required" bson:"fileID"`
	OwnerID     bson.ObjectId `json:"ownerID" bson:"ownerID"`
	Context     int16         `json:"context,required" bson:"context"`
	ContextID   bson.ObjectId `json:"contextID,required" bson:"contextID"`
	Permissions int16         `json:"permissions,required" bson:"permissions"`
	Visibility  int16         `json:"visibility,required" bson:"visibility"`
	CreatedDate time.Time     `json:"createdDate,required" bson:"createdDate"`
}
