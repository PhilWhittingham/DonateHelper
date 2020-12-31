package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/urfave/cli/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/PhilWhittingham/DonateHelper/db"
	"github.com/PhilWhittingham/DonateHelper/types"
)

func main() {
	app := &cli.App{
		Name:  "Donate Helper API",
		Usage: "The API service for the donate helper system",
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add a charity to the list",
				Action: func(c *cli.Context) error {
					name := c.Args().First()
					if name == "" {
						return errors.New("Cannot add charity with no name")
					}

					charity := &types.Charity{
						ID:        primitive.NewObjectID(),
						CharityID: "2",
						CompanyID: "2",
						Name:      name,
						Website:   "https://en.wikipedia.org/wiki/Main_Page",
					}

					return db.CreateCharity(charity)
				},
			},
			{
				Name:    "all",
				Aliases: []string{"l"},
				Usage:   "list all charities",
				Action: func(c *cli.Context) error {
					charities, err := db.GetAll()
					if err != nil {
						if err == mongo.ErrNoDocuments {
							fmt.Print("No charities are present")
							return nil
						}

						return err
					}

					for i, v := range charities {
						fmt.Printf("%d: %s\n", i+1, v.Name)
					}
					return nil
				},
			},
			{
				Name:    "csv",
				Aliases: []string{"c"},
				Usage:   "add charities from a csv",
				Action: func(c *cli.Context) error {
					filepath := c.Args().First()
					if filepath == "" {
						return errors.New("Must contain a file path")
					}
					errList := loadCharitiesFromFile(filepath)
					for err := range errList {
						fmt.Println(err)
					}

					return nil
				},
			},
			{
				Name:    "api",
				Aliases: []string{"r"},
				Usage:   "start the rest api",
				Action: func(c *cli.Context) error {
					e := echo.New()
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
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func loadCharitiesFromFile(filepath string) []error {
	inFile, err := os.Open(filepath)

	if err != nil {
		return []error{err}
	}

	defer inFile.Close()
	errList := []error{}
	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ",")
		charity := &types.Charity{
			ID:        primitive.NewObjectID(),
			CharityID: s[0],
			CompanyID: s[1],
			Name:      s[2],
			Website:   s[3],
		}
		if err := db.CreateCharity(charity); err != nil {
			errList = append(errList, err)
		}
	}

	return errList
}
