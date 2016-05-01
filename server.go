package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Quote struct {
	ID        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Title     string        `json:"title" bson:"title"`
	Author    string        `json:"author" bson:"author"`
	Timestamp time.Time     `json:"timestamp" bson:"timestamp"`
}

func dbConnection() session {
  session, err := mgo.Dial("mongodb://127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
  
  return session
}

func generateId() id {
  id := bson.NewObjectId()
  return id
}

func main() {

	e := echo.New()
  
  session := dbConnection()
  
  e.POST("/qoutes", createQoute)
  e.GET("/qoutes/", getAllQoutes)
  e.GET("/qoutes/:id", getQoute)
  e.PUT("/qoutes/:id", updateQoute)
  e.DELETE("/qoutes/:id", deleteQoute)

  e.Run(fasthttp.New(":4000"))
  
  func createQoute(c echo.Context) error {
      
      title := c.FormValue("title")
      author := c.FormValue("author")
      
      c := session.DB("motivational").C("qoutes")
      id := generateId()
      err = c.Insert(&Quote{id, "Education is the key to success", "Sam", time.Now()})
      if err != nil {
        log.Fatal(err)
      }
      
      return c.JSON(http.StatusOK)
  }
  
  func getAllQoutes(c echo.Context) {
    var results []Quote
	  err = c.Find(bson.M{"name": "Sam"}).Sort("-timestamp").All(&results)

	 if err != nil {
		 panic(err)
	 }
   
   return c.JSON(http.StatusOK, results)
  }
  
  func getQoute(c echo.Context)  {
    result := Quote{}
	  err = c.Find(bson.M{"name": "Sam"}).One(&result)
	  if err != nil {
		  panic(err)
	  }
    
    return c.JSON(http.StatusOK, result)
  }
}
