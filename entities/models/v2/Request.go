package v2

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Request representation of message entity
type Request struct {
	ID          bson.ObjectId `json:"id,omitempty" bson:"_id"`
	SenderID    bson.ObjectId `json:"senderID,required" bson:"senderID"`
	TargetID    bson.ObjectId `json:"targetID,required" bson:"targetID"`
	Type        string        `json:"type,required" bson:"type"`
	Text        string        `json:"text" bson:"text"`
	Status      int16         `json:"status,required" bson:"status"`
	CreatedDate time.Time     `json:"createdDate,required" bson:"createdDate"`
	UpdatedDate time.Time     `json:"updatedDate,required" bson:"updatedDate"`
}
