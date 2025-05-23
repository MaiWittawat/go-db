package repository

import (
	"context"
	"go-rebuild/model"
	"go-rebuild/module/port"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoUserRepo struct {
	collection *mongo.Collection
}

// ------------------------ Constructor ------------------------
func NewMongoUserRepo(client *mongo.Client) port.UserDB {
	dbName := "miniproject"
	coll := client.Database(dbName).Collection("users")
	return &mongoUserRepo{collection: coll}
}

// ------------------------ Method ------------------------

func (m *mongoUserRepo) Create(ctx context.Context, u *model.User) error {
	_, err := m.collection.InsertOne(ctx, u)
	return err
}

func (m *mongoUserRepo) Update(ctx context.Context, u *model.User, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = m.collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": u})
	return err
}

func (m *mongoUserRepo) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = m.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func (m *mongoUserRepo) FindByID(ctx context.Context, id string) (*model.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user model.User
	err = m.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	return &user, err
}