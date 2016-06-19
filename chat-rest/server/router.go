package server
import (
	"net/http"
	"github.com/gorilla/mux"
	"time"
	"github.com/rutmir/services/core/log"
)


func NewRouter(pathPrefix string) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	for _, element := range routes {
		var handler http.Handler
		route := Route(element)
		handler = route.HandlerFunc
		handler = logger(handler, route.Name)

		router.Methods(route.Method).Path(pathPrefix + route.Pattern).Name(route.Name).Handler(handler)
	}
	return router
}


func logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		log.Info("%s\t%s\t%s\t%s", r.Method, r.RequestURI, name, time.Since(start))
	})
}