package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
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
	e.Use(middleware.CORS())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*", "http://192.168.43.108:8100/"},
	}))

	e.POST("/qoutes/", CreateQoute)
	e.GET("/qoutes/", GetAllQoutes)
	e.GET("/qoutes/by/:status", GetQoutesByStatus)
	e.GET("/qoutes/:id", GetQoute)
	e.PUT("/qoutes/:id", UpdateQoute)
	e.PUT("/qoutes/:id/:status", UpdateStatus)
	e.DELETE("/qoutes/:id", DeleteQoute)
	log.Print("Server listing on port", os.Getenv("PORT"))
	e.Run(standard.New(os.Getenv("PORT")))

}
