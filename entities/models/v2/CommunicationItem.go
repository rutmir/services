package v2

import "fmt"

// CommunicationItem representation of personal contact information
type CommunicationItem struct {
	Value string `json:"value,required" bson:"value"`
	Type  string `json:"type" bson:"type"`
}

// ToString stringify CommunicationItem object
func (ci *CommunicationItem) ToString() string {
	return fmt.Sprintf("CommunicationItem@ Type: %s, Value: %s.", ci.Type, ci.Value)
}
