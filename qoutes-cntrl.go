package main

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type statusType int

const (
	Pending statusType = iota
	Approved
	Rejected
)

var ErrInvalidStatusType = errors.New("Invalid statusType")

// Quote defines structure of quote
type Quote struct {
	ID        bson.ObjectId `json:"-" bson:"_id,omitempty"`
	Title     string        `json:"title" bson:"title"`
	Author    string        `json:"author" bson:"author,omitempty"`
	Status    statusType    `json:"status" bson:"status"`
	Tag       []string      `json:"tag" bson:"tag"`
	Timestamp time.Time     `json:"timestamp" bson:"timestamp"`
}

func CreateQoute(c echo.Context) error {

	session := ConnectDB()
	defer session.Close()

	title := c.FormValue("title")
	author := c.FormValue("author")

	id := GenerateID()
	qoutes := GetCollection(session, "qoutes")
	err := qoutes.Insert(&Quote{
		ID:        id,
		Title:     title,
		Author:    author,
		Status:    Pending,
		Timestamp: time.Now(),
	})

	CheckError(err)

	return c.JSON(http.StatusOK, id)
}

func GetAllQoutes(c echo.Context) error {

	var results []Quote

	session := ConnectDB()
	defer session.Close()

	qoutes := GetCollection(session, "qoutes")
	err := qoutes.Find(bson.M{}).All(&results)

	CheckError(err)
	return c.JSON(http.StatusOK, results)
}

func GetQoutesByStatus(c echo.Context) error {

	var results []Quote

	session := ConnectDB()
	defer session.Close()

	status, err := ValidateStatusType(c.Param("status"))

	CheckError(err)

	qoutes := GetCollection(session, "qoutes")
	err = qoutes.Find(bson.M{"status": status}).All(&results)

	CheckError(err)
	return c.JSON(http.StatusOK, results)
}

func GetQoute(c echo.Context) error {

	session := ConnectDB()
	defer session.Close()

	result := Quote{}
	qoutes := GetCollection(session, "qoutes")
	err := qoutes.Find(bson.M{}).One(&result)

	CheckError(err)

	return c.JSON(http.StatusOK, result)
}

func UpdateQoute(c echo.Context) error {
	session := ConnectDB()
	defer session.Close()

	id := GetHexID(c.Param("id"))

	var qoute Quote

	qoutes := GetCollection(session, "qoutes")

	change := mgo.Change{
		Update: bson.M{"$set": bson.M{
			"title":  c.FormValue("title"),
			"author": c.FormValue("author"),
		}},
		ReturnNew: true,
	}

	info, err := qoutes.Find(bson.M{"_id": id}).Apply(change, &qoute)

	CheckError(err)

	return c.JSON(http.StatusOK, info)
}

func DeleteQoute(c echo.Context) error {

	session := ConnectDB()
	defer session.Close()

	id := GetHexID(c.Param("id"))

	qoutes := GetCollection(session, "qoutes")
	err := qoutes.Remove(bson.M{"_id": id})

	CheckError(err)

	return c.JSON(http.StatusOK, qoutes)

}

func UpdateStatus(c echo.Context) error {

	session := ConnectDB()
	defer session.Close()

	id := GetHexID(c.Param("id"))
	status, err := ValidateStatusType(c.Param("status"))

	CheckError(err)

	qoutes := GetCollection(session, "qoutes")
	err = qoutes.Update(
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"status": status}})

	CheckError(err)

	return c.JSON(http.StatusOK, id)
}

func ValidateStatusType(status string) (statusType, error) {
	i, err := strconv.Atoi(string(status))

	CheckError(err)

	v := statusType(i)

	if v < Pending || v > Rejected {
		return 0, ErrInvalidStatusType
	}

	return v, nil
}
