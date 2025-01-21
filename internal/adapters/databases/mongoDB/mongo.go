package mongoDB

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

func NewMongoClient(mongoURI string, mongoTimeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewMongoRepository(mongoURI, mongoDB string, mongoTimeout int) (*MongoRepository, error) {
	repo := &MongoRepository{
		timeout:  time.Duration(mongoTimeout) * time.Second,
		database: mongoDB,
	}

	client, err := NewMongoClient(mongoURI, mongoTimeout)
	if err != nil {
		return nil, err
	}
	repo.client = client
	return repo, nil
}

func (mr *MongoRepository) GetCollection(collectionName string) *mongo.Collection {
	return mr.client.Database(mr.database).Collection(collectionName)
}
