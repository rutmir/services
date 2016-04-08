package v2

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

// ClientApplication representation of client application
type ClientApplication struct {
	ID           bson.ObjectId `json:"id,omitempty" bson:"_id"`
	Name         string        `json:"name,required" bson:"name"` // unique
	ClientID     string        `json:"clientID,required" bson:"clientID"` // unique
	ClientSecret string        `json:"clientSecret,required" bson:"clientSecret"`
}

// ToString stringify ClientApplication object
func (ca *ClientApplication) ToString() string {
	return fmt.Sprintf("ClientApplication@ ID: %v, Name: %s, ClientId: %s, ClientSecret: %s.", ca.ID, ca.Name, ca.ClientID, ca.ClientSecret)
}
