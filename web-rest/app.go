package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rutmir/services/web-rest/server"
)

var router *mux.Router

func init() {
	router = server.NewRouter("")
	http.Handle("/api/", router)
}
