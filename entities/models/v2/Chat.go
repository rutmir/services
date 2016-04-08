package v2

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Chat representation of chat entity
type Chat struct {
	ID          bson.ObjectId   `json:"id,omitempty" bson"_id"`
	OwnerID     bson.ObjectId   `json:"ownerID,required" bson"ownerID"`
	Members     []bson.ObjectId `json:"members" bson"members"`
	CreatedDate time.Time       `json:"createdDate,required" bson"createdDate"`
	UpdatedDate time.Time       `json:"updatedDate,required" bson"updatedDate"`
	ChangedDate time.Time       `json:"changedDate,required" bson"changedDate"`
	Status      int16           `json:"status,required" bson"status"`
	Image       string          `json:"image" bson"image"`
	Title       string          `json:"title" bson"title"`
}
