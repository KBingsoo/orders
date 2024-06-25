package database

import (
	"context"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mux = &sync.Mutex{}

type MongoDB struct{}

var singleMongoDB *mongo.Client

func GetInstance() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	mux.Lock()
	defer mux.Unlock()

	var err error
	if singleMongoDB == nil {
		mongoUri := os.Getenv("MONGO_URI")
		singleMongoDB, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
		if err != nil {
			return nil, err
		}

	}

	return singleMongoDB, nil
}

func getDatabase(name string) (*mongo.Database, error) {
	c, err := GetInstance()
	if err != nil {
		return nil, err
	}
	return c.Database(name), nil
}

func GetCollection(name string) (*mongo.Collection, error) {
	mongoDatabase := os.Getenv("MONGO_DATABASE")
	d, err := getDatabase(mongoDatabase)
	if err != nil {
		return nil, err
	}

	return d.Collection(name), nil
}

func Close() {
	if singleMongoDB != nil {
		singleMongoDB.Disconnect(context.Background())
	}
}
