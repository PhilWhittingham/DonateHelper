package main

import (
   "context"
   "errors"
   "fmt"
   "log"
   //"time"
   "os"

   "github.com/urfave/cli/v2"
   "go.mongodb.org/mongo-driver/bson"
   
   "go.mongodb.org/mongo-driver/bson/primitive"
   "go.mongodb.org/mongo-driver/mongo"
   "go.mongodb.org/mongo-driver/mongo/options"
   "gopkg.in/gookit/color.v1"
   //"go.mongodb.org/mongo-driver/mongo/readpref"
)

var collection *mongo.Collection
var ctx = context.TODO()

func init() {
    clientOptions := options.Client().ApplyURI("connect string")
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
    ID        primitive.ObjectID `bson:"_id"`
    CharityID string             `bson:"charityID"`
    CompanyID string             `bson:"companyID"`
    Name      string             `bson:"name"`
    Website   string             `bson:"website"`
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
                            fmt.Print("Nothing to see here.\nRun `add 'task'` to add a task")
                            return nil                            
                        }

                        return err
                    }

                    printCharities(charities)
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