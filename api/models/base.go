package models

import (
	"log"
	"os"

	"gopkg.in/mgo.v2"
)

var (
	collection = "blacklists"
	database   = "urls"
	mongo      = os.Getenv("MONGO_DIAL")
)

// MongoCollection wraps a mgo.Collection to embed methods in models.
type MongoCollection struct {
	*mgo.Collection
}

// Collection is an interface to access to the collection struct.
type Collection interface {
	Find(query interface{}) *mgo.Query
}

// MongoDatabase wraps a mgo.Database to embed methods in models.
type MongoDatabase struct {
	*mgo.Database
}

// C shadows *mgo.DB to return a DataLayer interface instead of *mgo.Database.
func (d MongoDatabase) C(name string) Collection {
	return &MongoCollection{Collection: d.Database.C(name)}
}

// DataLayer is an interface to access the database struct
type DataLayer interface {
	C(name string) Collection
	GetBlacklists(d string) (URL, error)
}

// Session is an interface to access to the Session struct.
type Session interface {
	DB(name string) DataLayer
	Close()
	Copy() Session
}

// MongoSession is a Mongo session.
type MongoSession struct {
	*mgo.Session
}

// DB shadows *mgo.DB to return a DataLayer interface instead of *mgo.Database.
func (s MongoSession) DB(name string) DataLayer {
	return &MongoDatabase{Database: s.Session.DB(name)}
}

// Copy mocks mgo.Session.Copy()
func (s MongoSession) Copy() Session {
	return MongoSession{s.Session.Copy()}
}

// NewSession returns a new Mongo Session.
func NewSession() Session {
	mgoSession, err := mgo.Dial(mongo)
	if err != nil {
		log.Fatal("cannot dial mongo: ", err)
	}
	log.Printf("Connected to %s", mongo)
	session := MongoSession{mgoSession}
	return session
}
