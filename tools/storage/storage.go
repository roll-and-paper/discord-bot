package storage

import "context"

type Storage interface {
	FindAll(result interface{}, context context.Context) error
	FindOne(id string, result interface{}, context context.Context) error
	FindOneOrCreate(id string, data interface{}, context context.Context) error
	Remove(id string, context context.Context) error
	Save(id string, update interface{}, result interface{}, context context.Context) error
}
