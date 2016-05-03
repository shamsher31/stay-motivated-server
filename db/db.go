package db

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// DBName define name of database
var (
	mgoSession *mgo.Session
	DBName     = os.Getenv("DB_NAME")
)

// GetSession gives new mongo session
func GetSession() *mgo.Session {
	if mgoSession == nil {
		mgoSession, err := mgo.Dial(os.Getenv("DB_URL"))
		if err != nil {
			panic(err)
		}
		// Optional. Switch the session to a monotonic behavior.
		mgoSession.SetMode(mgo.Monotonic, true)
	}
	return mgoSession.Clone()
}

// GenerateID gives new mongo id
func GenerateID() (id bson.ObjectId) {
	id = bson.NewObjectId()
	return id
}
