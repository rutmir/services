package v2

import "fmt"

// NameBlock described person name
type NameBlock struct {
	LastName   string `json:"lastName" bson:"lastName"`
	FirstName  string `json:"firstName" bson:"firstName"`
	MiddleName string `json:"middleName" bson:"middleName"`
}

// ToString stringify object
func (nb *NameBlock) ToString() string {
	return fmt.Sprintf("NameBlock@ LastName: %s, FirstName: %s, MiddleName: %s.", nb.LastName, nb.FirstName, nb.MiddleName)
}
