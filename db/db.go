package db

import (
	"fmt"
	"os"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// DBName define name of database
var DBName = "motivational"

// ConnectDB connects to mongodb
func ConnectDB() (session *mgo.Session) {

	session, err := mgo.Dial(os.Getenv("DB_URL"))
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	session.SetSafe(&mgo.Safe{})
	return session
}

// GenerateID gives new mongo id
func GenerateID() (id bson.ObjectId) {
	return bson.NewObjectId()
}

// GetCollection gives collection name you want to select
func GetCollection(s *mgo.Session, colName string) (qoutes *mgo.Collection) {
	return s.DB(DBName).C(colName)
}
