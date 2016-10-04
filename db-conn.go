package main

import (
	"fmt"
	"os"

	"gopkg.in/mgo.v2"
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
