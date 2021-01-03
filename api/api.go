package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/PhilWhittingham/DonateHelper/db"
)

func Initialise() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/all", func(c echo.Context) error {
		charities, err := db.GetAll()
		if err != nil {
			if err == mongo.ErrNoDocuments {
				fmt.Print("No charities are present")
				return nil
			}
			return err
		}
		return c.JSON(http.StatusOK, charities)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
