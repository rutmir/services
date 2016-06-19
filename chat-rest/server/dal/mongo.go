package dal

import (
	"os"

	"github.com/rutmir/services/core/log"
	"gopkg.in/mgo.v2"
)

var Session *mgo.Session

func init() {
	mongoUrl := os.Getenv("MONGO_URL")
	if len(mongoUrl) == 0 {
		log.Fatal("MONGO error: Required to set 'MONGO_URL' environment")
		return
	}

	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		log.Fatal("NO DB Connection to: %s", mongoUrl)

	}
	Session = session
}
