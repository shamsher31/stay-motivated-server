package qoutes

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/shamsher31/stay-motivated-server/db"
	"github.com/shamsher31/stay-motivated-server/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type StatusType int

const (
	PENDING StatusType = iota
	APPROVED
	CANCELED
)

// Quote defines structure of quote
type Quote struct {
	ID        bson.ObjectId `json:"-" bson:"_id,omitempty"`
	Title     string        `json:"title" bson:"title"`
	Author    string        `json:"author" bson:"author,omitempty"`
	Status    StatusType    `json:"status" bson:"status"`
	Tag       []string      `json:"tag" bson:"tag"`
	Timestamp time.Time     `json:"timestamp" bson:"timestamp"`
}

// Server provides common mongo session
type Server struct {
	session *mgo.Session
}

// Global variable to hold mongo session
var gServer Server

func CreateQoute(c echo.Context) error {

	session := gServer.session.Copy()
	defer session.Close()

	title := c.FormValue("title")
	author := c.FormValue("author")

	id := db.GenerateID()
	qoutes := db.GetCollection(session, "qoutes")
	err := qoutes.Insert(&Quote{
		ID:        id,
		Title:     title,
		Author:    author,
		Status:    PENDING,
		Timestamp: time.Now(),
	})

	utils.CheckError(err)

	return c.JSON(http.StatusOK, qoutes)
}

func GetAllQoutes(c echo.Context) error {

	var results []Quote

	session := gServer.session.Copy()
	defer session.Close()

	qoutes := db.GetCollection(session, "qoutes")
	err := qoutes.Find(bson.M{}).One(&results)

	utils.CheckError(err)
	return c.JSON(http.StatusOK, results)
}

func GetQoute(c echo.Context) error {

	session := gServer.session.Copy()
	defer session.Close()

	result := Quote{}
	qoutes := db.GetCollection(session, "qoutes")
	err := qoutes.Find(bson.M{}).One(&result)

	utils.CheckError(err)

	return c.JSON(http.StatusOK, result)
}

func UpdateQoute(c echo.Context) error {
	session := gServer.session.Copy()
	defer session.Close()

	id := db.GetHexID(c.Param("id"))
	author := c.FormValue("author")

	qoutes := db.GetCollection(session, "qoutes")
	err := qoutes.Update(bson.M{"_id": id}, bson.M{"$set": bson.M{"author": author}})

	utils.CheckError(err)

	return c.JSON(http.StatusOK, qoutes)
}

func DeleteQoute(c echo.Context) error {

	session := gServer.session.Copy()
	defer session.Close()

	id := db.GetHexID(c.Param("id"))

	qoutes := db.GetCollection(session, "qoutes")
	err := qoutes.Remove(bson.M{"_id": id})

	utils.CheckError(err)

	return c.JSON(http.StatusOK, qoutes)

}
