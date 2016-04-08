// +build !appengine

// appcfg.py -A gab update app.yaml
package main

import (
	"net/http"
	"os"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/rutmir/services/core/log"
	"github.com/rutmir/services/web-rest/server"
)

func main() {

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	n := negroni.New()

	// Init middleware
	n.Use(negroni.NewRecovery())
	//n.Use(newLogger())

	// static router
	n.Use(negroni.NewStatic(http.Dir(dir + "/static")))

	// doc router
	docSMiddleware := negroni.NewStatic(http.Dir(dir + "/swagger-doc"))
	docSMiddleware.Prefix = "/doc"
	n.Use(docSMiddleware)

	n.UseHandler(server.NewRouter("/api"))

	n.Run(":8070")

	//log.Fatal(http.ListenAndServe(":8080", nil))
}

// Logger is a middleware handler that logs the request as it goes in and the response as it goes out.
type Logger struct {
}

func newLogger() *Logger {
	return &Logger{}
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	next(rw, r)
	res := rw.(negroni.ResponseWriter)
	//latency := int64(time.Since(start).Seconds() / 1000)
	/*if l.framework.Stats.Connected {
		Debug("Sending Stats to statsd...")
		logIfErr(l.framework.Stats.Client.Inc(r.URL.Path, 1, 1.0))
		logIfErr(l.framework.Stats.Client.Inc(strconv.Itoa(res.Status()), 1, 1.0))
		logIfErr(l.framework.Stats.Client.Inc(strconv.Itoa(res.Status()) + " - " + r.URL.String(), 1, 1.0))
		logIfErr(l.framework.Stats.Client.Gauge(l.framework.Id + "-avg-response-time", latency, 1.0))
		logIfErr(l.framework.Stats.Client.Gauge(l.framework.Id + "/" + r.URL.Path + "-avg-response-time", latency, 1.0))
	}*/
	log.Debug(res.Status(), http.StatusText(res.Status()), time.Since(start))
}
