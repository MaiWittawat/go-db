package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func InitMongoDB(ctx context.Context) (*mongo.Client, error) {
	uri := os.Getenv("MONGO_URI")
	log.Println("uri: ", uri)

	opts := options.Client().ApplyURI(uri).SetMaxPoolSize(100).SetRetryWrites(true)
	client, err := mongo.Connect(ctx, opts)

	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("fail to connect mongo: %v", err)
	}

	return client, nil
}

func DisconnectMongoDB(mgDB *mongo.Database, ctx context.Context) error {
	if err := mgDB.Client().Disconnect(ctx); err != nil {
		return err
	}
	fmt.Println("disconnected mongodb successfully")
	return nil
}