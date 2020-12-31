package main

import (
	"errors"
	"fmt"
	"log"
	"os"

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
							fmt.Print("Nothing to see here.\nRun `add 'task'` to add a task")
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
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
