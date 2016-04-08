package controllers

import (
	"fmt"
	"net/http"
)

// Alive http handler for check service status request
func Alive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "alive")
}
