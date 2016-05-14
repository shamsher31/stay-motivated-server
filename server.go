package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
	"github.com/labstack/echo/middleware"

	"github.com/shamsher31/stay-motivated-server/db"
	"github.com/shamsher31/stay-motivated-server/utils"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Status int

const (
	PENDING Status = iota
	APPROVED
	CANCELED
)

// Quote defines structure of quote
type Quote struct {
	ID        bson.ObjectId `json:"-" bson:"_id,omitempty"`
	Title     string        `json:"title" bson:"title"`
	Author    string        `json:"author" bson:"author,omitempty"`
	Status    string        `json:"status" bson:"status"`
	Tag       []string      `json:"tag" bson:"tag"`
	Timestamp time.Time     `json:"timestamp" bson:"timestamp"`
}

// Server provides common mongo session
type Server struct {
	session *mgo.Session
}

// Global variable to hold mongo session
var gServer Server

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {

	session := db.ConnectDB()
	defer session.Close()

	gServer = Server{session}

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/qoutes/", createQoute)
	e.GET("/qoutes/", getAllQoutes)
	e.GET("/qoutes/:id", getQoute)
	// e.PUT("/qoutes/:id", updateQoute)
	e.DELETE("/qoutes/:id", deleteQoute)

	e.Run(fasthttp.New(os.Getenv("PORT")))

}

func createQoute(c echo.Context) error {

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
		Timestamp: time.Now(),
	})

	utils.CheckError(err)

	return c.JSON(http.StatusOK, qoutes)
}

func getAllQoutes(c echo.Context) error {

	var results []Quote

	session := gServer.session.Copy()
	defer session.Close()

	qoutes := db.GetCollection(session, "qoutes")
	err := qoutes.Find(bson.M{}).One(&results)

	utils.CheckError(err)
	return c.JSON(http.StatusOK, results)
}

func getQoute(c echo.Context) error {

	session := gServer.session.Copy()
	defer session.Close()

	result := Quote{}
	qoutes := db.GetCollection(session, "qoutes")
	err := qoutes.Find(bson.M{}).One(&result)

	utils.CheckError(err)

	return c.JSON(http.StatusOK, result)
}

func deleteQoute(c echo.Context) error {

	session := gServer.session.Copy()
	defer session.Close()

	id := db.GetHexID(c.Param("id"))

	qoutes := db.GetCollection(session, "qoutes")
	err := qoutes.Remove(bson.M{"_id": id})

	utils.CheckError(err)

	return c.JSON(http.StatusOK, qoutes)

}
