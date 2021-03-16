package discord

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Lock struct {
	database *mongo.Database
}

type lockDocument struct {
	Id    string `bson:"_id"`
	Count int    `bson:"count"`
}

func NewLock(database *mongo.Database) *Lock {
	return &Lock{database: database}
}

func (l *Lock) Received(id string, ctx context.Context) bool {
	res := l.database.
		Collection("lock").
		FindOneAndUpdate(ctx,
			bson.M{"_id": id},
			bson.M{"$inc": bson.M{"count": 1}},
			options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
		)
	doc := &lockDocument{}
	if err := res.Decode(doc); err != nil {
		println(err.Error())
		return false
	}
	return doc.Count == 1
}

func (l *Lock) End(id string, ctx context.Context) error {
	return nil
}
