package v2

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

// ContactItem representation of link to another object
type ContactItem struct {
	ProfileID bson.ObjectId `json:"profileID,required" bson:"profileID"`
	Status    int16         `json:"status,required" bson:"status"`
}

// ToString stringify ContactItem object
func (ci *ContactItem) ToString() string {
	return fmt.Sprintf("ContactItem@ ProfileID: %v, Status: %v.", ci.ProfileID, ci.Status)
}
