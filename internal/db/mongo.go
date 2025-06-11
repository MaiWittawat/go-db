package db

import (
	"context"
	"errors"
	"fmt"
	appcore_config "go-rebuild/cmd/go-rebuild/config"
	"reflect"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoRepo struct {
	client *mongo.Client
	dbName string
}

func InitMongoDB(ctx context.Context) (*mongo.Client, error) {
	url := appcore_config.Config.MongoConnString
	opts := options.Client().ApplyURI(url).SetMaxPoolSize(100).SetRetryWrites(true)
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

// ------------------------ Constructor ------------------------
func NewMongoRepo(client *mongo.Client, dbName string) DB {
	return &mongoRepo{client: client, dbName: dbName}
}

// ------------------------ Method ------------------------
func (m *mongoRepo) setCollection(name string) *mongo.Collection {
	return m.client.Database(m.dbName).Collection(name)
}

func (m *mongoRepo) modelToBSONDoc(model any) (bson.M, error) {
    // Marshal เป็น BSON แล้ว Unmarshal เป็น bson.M
    data, err := bson.Marshal(model)
    if err != nil {
        return nil, err
    }
    
    var doc bson.M
    if err := bson.Unmarshal(data, &doc); err != nil {
        return nil, err
    }
    
    // จัดการ ID field
    if idValue, exists := doc["id"]; exists {
        if idStr, ok := idValue.(string); ok && idStr != "" {
            if primitive.IsValidObjectID(idStr) {
                objID, _ := primitive.ObjectIDFromHex(idStr)
                doc["_id"] = objID
            } else {
                // ถ้าไม่ใช่ ObjectID format ให้สร้างใหม่
                doc["_id"] = primitive.NewObjectID()
            }
            delete(doc, "id") // ลบ id field เดิม
        } else {
            // ถ้า id เป็น empty หรือไม่ใช่ string
            doc["_id"] = primitive.NewObjectID()
            delete(doc, "id")
        }
    } else {
        // ถ้าไม่มี id field
        doc["_id"] = primitive.NewObjectID()
    }
    
    return doc, nil
}

// ------------------------ Method Basic CUD ------------------------
func (m *mongoRepo) Create(ctx context.Context, coll string, model any) error {
    doc, err := m.modelToBSONDoc(model)
    if err != nil {
        return err
    }
    
    _, err = m.setCollection(coll).InsertOne(ctx, doc)
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

func (m *mongoRepo) Delete(ctx context.Context, coll string, model any, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = m.setCollection(coll).DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

// ------------------------ Method Basic Query ------------------------
func (m *mongoRepo) GetAll(ctx context.Context, coll string, results any) error {
	cursor, err := m.setCollection(coll).Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	// ตรวจสอบว่า results เป็น pointer ไปยัง slice จริง ๆ
	slicePtr := reflect.ValueOf(results)
	if slicePtr.Kind() != reflect.Ptr {
		return errors.New("results must be a pointer to a slice")
	}

	sliceVal := slicePtr.Elem()
	elemType := sliceVal.Type().Elem()

	for cursor.Next(ctx) {
		elemPtr := reflect.New(elemType) // เช่น *Product

		if err := cursor.Decode(elemPtr.Interface()); err != nil {
			continue
		}

		// append *elemPtr ลง slice
		sliceVal.Set(reflect.Append(sliceVal, elemPtr.Elem()))
	}

	return nil
}


func (m *mongoRepo) GetByID(ctx context.Context, coll string, id string, result any) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	if err = m.setCollection(coll).FindOne(ctx, bson.M{"_id": objID}).Decode(result); err != nil {
		return err
	}

	return nil
}


func (m *mongoRepo) GetByField(ctx context.Context, coll string, field string, value any, result any) error {
	filter := bson.M{field: value}

	if err := m.setCollection(coll).FindOne(ctx, filter).Decode(result); err != nil {
		return err
	}

	return nil

}
