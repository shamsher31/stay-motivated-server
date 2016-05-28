package qoutes

import (
	"errors"
	"net/http"
	"strconv"
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
	REJECTED
)

var ErrInvalidStatusType = errors.New("Invalid StatusType")

// Quote defines structure of quote
type Quote struct {
	ID        bson.ObjectId `json:"-" bson:"_id,omitempty"`
	Title     string        `json:"title" bson:"title"`
	Author    string        `json:"author" bson:"author,omitempty"`
	Status    StatusType    `json:"status" bson:"status"`
	Tag       []string      `json:"tag" bson:"tag"`
	Timestamp time.Time     `json:"timestamp" bson:"timestamp"`
}

func CreateQoute(c echo.Context) error {

	session := db.ConnectDB()
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

	return c.JSON(http.StatusOK, id)
}

func GetAllQoutes(c echo.Context) error {

	var results []Quote

	session := db.ConnectDB()
	defer session.Close()

	qoutes := db.GetCollection(session, "qoutes")
	err := qoutes.Find(bson.M{}).One(&results)

	utils.CheckError(err)
	return c.JSON(http.StatusOK, results)
}

func GetQoute(c echo.Context) error {

	session := db.ConnectDB()
	defer session.Close()

	result := Quote{}
	qoutes := db.GetCollection(session, "qoutes")
	err := qoutes.Find(bson.M{}).One(&result)

	utils.CheckError(err)

	return c.JSON(http.StatusOK, result)
}

func UpdateQoute(c echo.Context) error {
	session := db.ConnectDB()
	defer session.Close()

	id := db.GetHexID(c.Param("id"))

	var qoute Quote

	qoutes := db.GetCollection(session, "qoutes")

	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"title":  c.FormValue("title"),
			"author": c.FormValue("author"),
		}},
		ReturnNew: true,
	}

	info, err := qoutes.Find(bson.M{"_id": id}).Apply(change, &qoute)

	utils.CheckError(err)

	return c.JSON(http.StatusOK, info)
}

func DeleteQoute(c echo.Context) error {

	session := db.ConnectDB()
	defer session.Close()

	id := db.GetHexID(c.Param("id"))

	qoutes := db.GetCollection(session, "qoutes")
	err := qoutes.Remove(bson.M{"_id": id})

	utils.CheckError(err)

	return c.JSON(http.StatusOK, qoutes)

}

func UpdateStatus(c echo.Context) error {

	session := db.ConnectDB()
	defer session.Close()

	id := db.GetHexID(c.Param("id"))
	status, err := ValidateStatusType(c.Param("status"))

	utils.CheckError(err)

	qoutes := db.GetCollection(session, "qoutes")
	err = qoutes.Update(
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"status": status}})

	utils.CheckError(err)

	return c.JSON(http.StatusOK, id)
}

func ValidateStatusType(status string) (StatusType, error) {
	i, err := strconv.Atoi(string(status))

	utils.CheckError(err)

	v := StatusType(i)

	if v < PENDING || v > REJECTED {
		return 0, ErrInvalidStatusType
	}

	return v, nil
}
