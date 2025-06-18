package mongo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Migration001Categories(ctx context.Context, db interface{}) error {

	mongoDB, ok := db.(*mongo.Database)
	if !ok {
		return errors.New("invalid db type")
	}

	err := mongoDB.CreateCollection(ctx, "categories")

	if err != nil {
		var cmdErr mongo.CommandError
		if !(errors.As(err, &cmdErr) && cmdErr.Code == 48) {
			return err
		}
	}

	_, err = mongoDB.Collection("categories").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]int32{"name": 1},
		Options: options.Index().SetUnique(true),
	})

	if err != nil {
		return err
	}

	return err
}
