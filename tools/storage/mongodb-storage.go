package storage

import (
	"context"
	"github.com/dohr-michael/roll-and-paper-bot/tools/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDB struct {
	database   *mongo.Database
	collection string
}

func NewMongoDBStorage(database *mongo.Database, collection string) Storage {
	return &mongoDB{
		database:   database,
		collection: collection,
	}
}

func (m *mongoDB) FindAll(result interface{}, context context.Context) error {
	cursor, err := m.database.Collection(m.collection).Find(context, bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return err
	}
	if err := cursor.All(context, result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return err
	}
	return nil
}

func (m *mongoDB) FindOne(id string, result interface{}, context context.Context) error {
	cursor := m.database.Collection(m.collection).FindOne(context, bson.M{"_id": id})
	if cursor.Err() != nil {
		if cursor.Err() == mongo.ErrNoDocuments {
			return nil
		}
		return cursor.Err()
	}
	if err := cursor.Decode(result); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return err
	}
	return nil
}

func (m *mongoDB) FindOneOrCreate(id string, result interface{}, context context.Context) error {
	col := m.database.Collection(m.collection)
	cursor := col.FindOne(context, bson.M{"_id": id})
	if cursor.Err() != nil {
		if cursor.Err() == mongo.ErrNoDocuments {
			d, err := data.AsBsonDocument(result)
			if err != nil {
				return err
			}
			d["_id"] = id
			_, err = col.InsertOne(context, d)
			if err != nil {
				return err
			}
			return nil
		}
		return cursor.Err()
	}
	if err := cursor.Decode(result); err != nil {

		return err
	}
	return nil
}

func (m *mongoDB) Remove(id string, context context.Context) error {
	_, err := m.database.Collection(m.collection).DeleteOne(context, bson.M{"_id": id})
	return err
}

func (m *mongoDB) Save(id string, updated interface{}, result interface{}, context context.Context) error {
	toUpdate, err := data.AsBsonDocument(updated)
	if err != nil {
		return err
	}
	delete(toUpdate, "_id")

	cursor := m.database.Collection(m.collection).FindOneAndUpdate(
		context,
		bson.M{"_id": id},
		bson.M{"$set": toUpdate, "$setOnInsert": bson.M{"_id": id}},
		options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true),
	)
	if cursor.Err() != nil {
		return cursor.Err()
	}
	if err := cursor.Decode(result); err != nil {
		return err
	}
	return nil
}
