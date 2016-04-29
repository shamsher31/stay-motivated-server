package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
)

type Quote struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

func main() {

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		qoute := Quote{"Education is the key to success", "Sam"}

		return c.JSON(http.StatusOK, qoute)
	})
	e.Run(fasthttp.New(":4000"))
}
