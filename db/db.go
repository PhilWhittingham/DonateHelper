package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/PhilWhittingham/DonateHelper/types"
)

var collection *mongo.Collection
var ctx = context.TODO()

func init() {
	clientOptions := options.Client().ApplyURI("str")
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

func CreateCharity(charity *types.Charity) error {
	_, err := collection.InsertOne(ctx, charity)
	return err
}

func GetAll() ([]*types.Charity, error) {
	// passing bson.D{{}} matches all documents in the collection
	filter := bson.D{{}}
	return FilterCharities(filter)
}

func FilterCharities(filter interface{}) ([]*types.Charity, error) {
	// A slice of charities for storing the decoded documents
	var charities []*types.Charity

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return charities, err
	}

	for cur.Next(ctx) {
		var t types.Charity
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
