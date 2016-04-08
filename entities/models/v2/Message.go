package v2

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Message representation of message entity
type Message struct {
	ID          bson.ObjectId `json:"id,omitempty" bson:"_id"`
	SenderID    bson.ObjectId `json:"senderID,required" bson:"senderID"`
	ChatID      bson.ObjectId `json:"chatID,required" bson:"chatID"`
	Type        string        `json:"type,required" bson:"type"`
	Body        []byte        `json:"body" bson:"body"`
	Status      int16         `json:"status,required" bson:"status"`
	CreatedDate time.Time     `json:"createdDate,required" bson:"createdDate"`
	UpdatedDate time.Time     `json:"updatedDate,required" bson:"updatedDate"`
}
