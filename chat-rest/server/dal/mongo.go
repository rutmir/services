package dal

import (
	"gopkg.in/mgo.v2"
	"github.com/rutmir/services/core/log"
)

var Session   *mgo.Session

func init() {
	eUrl := "192.168.2.177:27017"
	session, err := mgo.Dial(eUrl)
	if err != nil {
		log.Fatal("NO DB Connection to: %s", eUrl)

	}
	Session = session
}
