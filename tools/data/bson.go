package data

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AsBsonDocument(v interface{}) (primitive.M, error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return nil, err
	}
	doc := primitive.M{}

	err = bson.Unmarshal(data, &doc)
	return doc, err
}
