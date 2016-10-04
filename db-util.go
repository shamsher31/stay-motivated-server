package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// GenerateID gives new mongo id
func GenerateID() (id bson.ObjectId) {
	return bson.NewObjectId()
}

// GetCollection gives collection name you want to select
func GetCollection(s *mgo.Session, colName string) (qoutes *mgo.Collection) {
	return s.DB(DBName).C(colName)
}

// GetHexID returns mongo id hex so you can find docs
func GetHexID(id string) bson.ObjectId {
	return bson.ObjectIdHex(id)
}
