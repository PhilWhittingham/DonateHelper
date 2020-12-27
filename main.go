package main

import (
   "context"
   "errors"
   "fmt"
   "log"
   "os"
   "bufio"
   "strings"
   "net/http"

   "github.com/labstack/echo/v4"
   "github.com/urfave/cli/v2"
   "go.mongodb.org/mongo-driver/bson"
   "go.mongodb.org/mongo-driver/bson/primitive"
   "go.mongodb.org/mongo-driver/mongo"
   "go.mongodb.org/mongo-driver/mongo/options"
   "gopkg.in/gookit/color.v1"
)

var collection *mongo.Collection
var ctx = context.TODO()

func init() {
    clientOptions := options.Client().ApplyURI("connection string")
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    err = client.Ping(ctx, nil)
    if err != nil {
      log.Fatal(err)
    }
  
    collection = client.Database("donate_helper").Collection("charities")
}

type Charity struct {
    ID        primitive.ObjectID `json:"_id" bson:"_id"`
    CharityID string             `json:"charityID" bson:"charityID"`
    CompanyID string             `json:"companyID" bson:"companyID"`
    Name      string             `json:"name" bson:"name"`
    Website   string             `json:"website" bson:"website"`
}

func main() {
    app := &cli.App{
        Name:     "Donate Helper API",
        Usage:    "The API service for the donate helper system",
        Commands: []*cli.Command{
            {
                Name: "add",
                Aliases: []string{"a"},
                Usage: "add a charity to the list",
                Action: func(c *cli.Context) error {
                    name := c.Args().First()
                    if name == "" {
                        return errors.New("Cannot add charity with no name")
                    }

                    charity := &Charity{
                        ID: primitive.NewObjectID(),
                        CharityID: "2",
                        CompanyID: "2",
                        Name: name,
                        Website: "https://en.wikipedia.org/wiki/Main_Page",
                    }

                    return createCharity(charity)
                },
            },
            {
                Name: "all",
                Aliases: []string{"l"},
                Usage: "list all charities",
                Action: func(c *cli.Context) error {
                    charities, err := getAll()
                    if err != nil {
                        if err == mongo.ErrNoDocuments {
                            fmt.Print("No charities are present")
                            return nil                            
                        }

                        return err
                    }

                    printCharities(charities)
                    return nil
                },
            },
            {
                Name: "csv",
                Aliases: []string{"c"},
                Usage: "add charities from a csv",
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
                Name: "api",
                Aliases: []string{"r"},
                Usage: "start the rest api",
                Action: func(c *cli.Context) error {
                    e := echo.New()
                    e.GET("/", func(c echo.Context) error {
                        return c.String(http.StatusOK, "Hello, World!")
                    })
                    e.GET("/all", func(c echo.Context) error {
                        charities, err := getAll()
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

func printCharities(charities []*Charity) {
    for i, v := range charities {
        color.Green.Printf("%d: %s\n", i+1, v.Name)
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
        charity := &Charity{
            ID: primitive.NewObjectID(),
            CharityID: s[0],
            CompanyID: s[1],
            Name: s[2],
            Website: s[3],
        }
        if err := createCharity(charity); err != nil {
            errList = append(errList, err)
        }
    }
    
    return errList
}

func createCharity(charity *Charity) error {
    _, err := collection.InsertOne(ctx, charity)
    return err
}

func getAll() ([]*Charity, error) {
    // passing bson.D{{}} matches all documents in the collection
    filter := bson.D{{}}
    return filterCharities(filter)
}
  
func filterCharities(filter interface{}) ([]*Charity, error) {
    // A slice of charities for storing the decoded documents
    var charities []*Charity
  
    cur, err := collection.Find(ctx, filter)
    if err != nil {
        return charities, err
    }
  
    for cur.Next(ctx) {
        var t Charity
        err := cur.Decode(&t)
        if err != nil {
            return charities, err
        }
  
        charities = append(charities, &t)
    }
  
    if err := cur.Err(); err != nil {
        return charities, err
    }
  
    // once exhausted, close the cursor
    cur.Close(ctx)
  
    if len(charities) == 0 {
        return charities, mongo.ErrNoDocuments
    }
  
    return charities, nil
}