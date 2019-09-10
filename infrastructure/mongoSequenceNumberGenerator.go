package infrastructure

import (
	"context"
	"github.com/stainour/test10/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

var sequenceDocId, _ = primitive.ObjectIDFromHex("5aff52eb72d9993560b2d3cc")
var idFilter = bson.D{{"_id", sequenceDocId}}

const sequenceCollectionName = "sequence"
const valueField = "value"

type sequence struct {
	Id    primitive.ObjectID
	Value int64
}

type MongoSequenceNumberGenerator struct {
	collection *mongo.Collection
}

func (generator *MongoSequenceNumberGenerator) NextValue(context context.Context) (int64, error) {
	collection := generator.collection

	update := bson.D{
		{"$inc", bson.D{
			{valueField, 1},
		}},
	}
	var result sequence
	returnDocument := options.After
	err := collection.FindOneAndUpdate(context, idFilter, update, &options.FindOneAndUpdateOptions{
		ReturnDocument: &returnDocument,
	}).Decode(&result)

	return result.Value, err
}

func NewMongoSequenceGenerator(connection MongoConnectionSetting) (domain.SequenceNumberGenerator, error) {
	clientOptions := options.Client().ApplyURI(connection.connectionUri)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		return nil, err
	}

	database := client.Database(connection.database)
	collection := database.Collection(sequenceCollectionName, &options.CollectionOptions{
		WriteConcern: writeconcern.New(writeconcern.J(true)),
	})

	update := bson.D{
		{"$setOnInsert", bson.D{
			{valueField, 0},
		}},
	}

	updateOptions := options.Update().SetUpsert(true)
	_, err = collection.UpdateOne(context.TODO(), idFilter, update, updateOptions)

	if err != nil {
		return nil, err
	}

	return &MongoSequenceNumberGenerator{collection: collection}, nil
}
