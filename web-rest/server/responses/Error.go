package responses

import "fmt"

type Error struct {
	Error       string `json:"error"`
	Description string `json:"error_description"`
}

func (this *Error) ToString() string {
	return fmt.Sprintf("Error@ Error: %s, Description: %s.", this.Error, this.Description)
}
