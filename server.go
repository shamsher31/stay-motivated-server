package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
	"github.com/labstack/echo/middleware"
	"github.com/shamsher31/stay-motivated-server/controllers"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/qoutes/", qoutes.CreateQoute)
	e.GET("/qoutes/", qoutes.GetAllQoutes)
	e.GET("/qoutes/by/:status", qoutes.GetQoutesByStatus)
	e.GET("/qoutes/:id", qoutes.GetQoute)
	e.PUT("/qoutes/:id", qoutes.UpdateQoute)
	e.PUT("/qoutes/:id/:status", qoutes.UpdateStatus)
	e.DELETE("/qoutes/:id", qoutes.DeleteQoute)
	log.Print("Server listing on port", os.Getenv("PORT"))
	e.Run(fasthttp.New(os.Getenv("PORT")))

}
