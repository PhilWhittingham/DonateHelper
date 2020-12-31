package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Charity struct {
	ID        primitive.ObjectID `bson:"_id"`
	CharityID string             `bson:"charityID"`
	CompanyID string             `bson:"companyID"`
	Name      string             `bson:"name"`
	Website   string             `bson:"website"`
}
