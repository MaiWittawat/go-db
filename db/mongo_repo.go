package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepo struct {
	client *mongo.Client
	dbName string
}

// ------------------------ Constructor ------------------------
func NewMongoRepo(client *mongo.Client, dbName string) DB {
	return &mongoRepo{client: client, dbName: dbName}
}

// ------------------------ Method ------------------------

func (m *mongoRepo) setCollection(name string) *mongo.Collection {
	return m.client.Database(m.dbName).Collection(name)
}

// *********************************************************

func (m *mongoRepo) Create(ctx context.Context, coll string, model any) error {
	_, err := m.setCollection(coll).InsertOne(ctx, model)
	return err
}

func (m *mongoRepo) Update(ctx context.Context, coll string, model any, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = m.setCollection(coll).UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": model})
	return err
}

func (m *mongoRepo) Delete(ctx context.Context, coll string, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = m.setCollection(coll).DeleteOne(ctx, bson.M{"_id": objID})
	return err
}