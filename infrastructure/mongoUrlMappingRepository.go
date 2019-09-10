package infrastructure

import (
	"context"
	"github.com/stainour/test10/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"net/url"
)

const batchSize = 1000
const mappingCollectionName = "urlMapping"
const hitCountName = "hitCount"
const shortenedKeyName = "shortenedKey"
const idName = "_id"

type mongoUrlMapping struct {
	HitCount     int64  `bson:"hitCount"`
	ShortenedKey string `bson:"shortenedKey"`
	Url          string `bson:"_id"`
}

type MongoUrlMappingRepository struct {
	collection *mongo.Collection
}

func (m MongoUrlMappingRepository) AddIfNotExists(context context.Context, uriMapping *domain.UrlMapping) (domain.AddResult, error) {
	update := bson.D{
		{"$setOnInsert", toBson(uriMapping)},
	}

	result, err := m.collection.UpdateOne(context,
		bson.D{{idName, uriMapping.Uri()}},
		update, options.Update().SetUpsert(true))

	if err != nil {
		return domain.Fail, nil
	}

	if result.MatchedCount == 1 {
		return domain.AlreadyExists, nil
	}

	return domain.OK, nil
}

func (m *MongoUrlMappingRepository) FindById(context context.Context, id *url.URL) (*domain.UrlMapping, error) {
	var mapping mongoUrlMapping
	err := m.collection.FindOne(context, bson.D{{idName, id.String()}}).Decode(&mapping)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return toUrlMapping(&mapping), nil
}

func (m *MongoUrlMappingRepository) IncrementHitCount(context context.Context, id *url.URL) error {
	filter := bson.D{{idName, id.String()}}
	update := bson.D{
		{"$inc", bson.D{
			{hitCountName, 1},
		}},
	}
	_, err := m.collection.UpdateOne(context, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoUrlMappingRepository) GetAll(context context.Context) (values <-chan *domain.UrlMapping, errors <-chan error) {
	errorsChan := make(chan error, 1)

	collection := m.collection
	find := options.Find()
	find.SetNoCursorTimeout(true)
	find.SetBatchSize(batchSize)
	cursor, err := collection.Find(context, bson.D{{}}, find)

	if err != nil {
		defer close(errorsChan)
		errorsChan <- err
		return nil, errorsChan
	}

	mappingsChan := make(chan *domain.UrlMapping, batchSize)

	go func() {
		defer close(mappingsChan)
		defer close(errorsChan)
		defer cursor.Close(context)

		for cursor.Next(context) {
			var mapping mongoUrlMapping
			err := cursor.Decode(&mapping)
			if err != nil {
				errorsChan <- err
				return
			}
			mappingsChan <- toUrlMapping(&mapping)
		}
		if err := cursor.Err(); err != nil {
			errorsChan <- err
		}
	}()

	return mappingsChan, errorsChan
}

func (m *MongoUrlMappingRepository) FindByShortenedKey(context context.Context, shortenedKey string) (*domain.UrlMapping, error) {
	var mapping mongoUrlMapping
	err := m.collection.FindOne(context, bson.D{{shortenedKeyName, shortenedKey}}).Decode(&mapping)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return toUrlMapping(&mapping), nil
}

func NewMongoUrlMappingRepository(connection MongoConnectionSetting) (domain.UrlMappingRepository, error) {
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
	collection := database.Collection(mappingCollectionName, &options.CollectionOptions{
		WriteConcern: writeconcern.New(writeconcern.J(true)),
	})

	_, err = collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.D{{shortenedKeyName, 1}},
		Options: options.Index().SetUnique(true).SetBackground(true),
	})
	if err != nil {
		return nil, err
	}

	return &MongoUrlMappingRepository{
		collection: collection,
	}, nil
}
func toBson(url *domain.UrlMapping) bson.D {
	return bson.D{
		{idName, url.Uri()},
		{hitCountName, url.HitCount()},
		{shortenedKeyName, url.ShortenedKey()},
	}
}
func toUrlMapping(mapping *mongoUrlMapping) *domain.UrlMapping {
	parse, _ := url.Parse(mapping.Url)
	return domain.NewFullUrlMapping(parse, mapping.ShortenedKey, mapping.HitCount)
}
