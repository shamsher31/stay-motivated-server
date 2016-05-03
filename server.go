package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
	"github.com/labstack/echo/middleware"

	"github.com/shamsher31/stay-motivated-server/db"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Quote defines structure of quote
type Quote struct {
	ID        bson.ObjectId `json:"-" bson:"_id,omitempty"`
	Title     string        `json:"title" bson:"title"`
	Author    string        `json:"author" bson:"author,omitempty"`
	Timestamp time.Time     `json:"timestamp" bson:"timestamp"`
}

func main() {

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

func getCollection() (qoutes *mgo.Collection) {
	session := db.GetSession()
	defer session.Close()

	qoutes = session.DB(db.DBName).C("qoutes")
	return qoutes
}

func createQoute(c echo.Context) error {

	title := c.FormValue("title")
	author := c.FormValue("author")

	qoutes := getCollection()
	id := db.GenerateID()
	err := qoutes.Insert(&Quote{id, title, author, time.Now()})
	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, qoutes)
}

func getAllQoutes(c echo.Context) error {

	qoutes := getCollection()

	var results []Quote
	err := qoutes.Find(bson.M{"name": "Sam"}).Sort("-timestamp").All(&results)

	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, results)
}

func getQoute(c echo.Context) error {
	qoutes := getCollection()

	result := Quote{}
	err := qoutes.Find(bson.M{"name": "Sam"}).One(&result)
	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, result)
}
