package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

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

// type Server struct {
// 	db *mgo.Session
// }

func main() {

	// db, err := mgo.Dial(os.Getenv("DB_URL"))
	// utils.CheckError(err)
	// defer db.Close()

	// server := &Server{db: db}

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/qoutes", createQoute)
	e.GET("/qoutes/", getAllQoutes)
	e.GET("/qoutes/:id", getQoute)
	// e.PUT("/qoutes/:id", updateQoute)
	// e.DELETE("/qoutes/:id", deleteQoute)

	e.Run(fasthttp.New(os.Getenv("PORT")))

}

func createQoute(c echo.Context) error {

	title := c.FormValue("title")
	author := c.FormValue("author")

	session, err := db.GetSession()
	utils.CheckError(err)
	defer session.Close()

	qoutes := db.GetCollection(session, "qoutes")
	id := db.GenerateID()
	err = qoutes.Insert(&Quote{id, title, author, time.Now()})

	utils.CheckError(err)

	return c.JSON(http.StatusOK, qoutes)
}

func getAllQoutes(c echo.Context) error {

	var results []Quote

	session, err := mgo.Dial(os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}

	qoutes := session.DB(db.DBName).C("qoutes")
	err = qoutes.Find(bson.M{"author": "Sam"}).One(&results)

	if err != nil {
		fmt.Println(err)
	}
	return c.JSON(http.StatusOK, results)
}

func getQoute(c echo.Context) error {

	session, err := db.GetSession()
	utils.CheckError(err)
	defer session.Close()

	qoutes := db.GetCollection(session, "qoutes")

	result := Quote{}
	err = qoutes.Find(bson.M{"name": "Sam"}).One(&result)

	utils.CheckError(err)

	return c.JSON(http.StatusOK, result)
}
