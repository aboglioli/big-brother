package db

import (
	"context"

	"github.com/aboglioli/big-brother/pkg/config"
	"github.com/aboglioli/big-brother/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrConnect = errors.Internal.New("mongo.connect")
	ErrPing    = errors.Internal.New("mongo.ping")
)

func Connect() (*mongo.Client, error) {
	config := config.Get()
	ctx := context.Background()

	options := options.Client().ApplyURI(config.MongoURL).SetAuth(
		options.Credential{
			AuthSource: config.MongoAuthSource,
			Username:   config.MongoUsername,
			Password:   config.MongoPassword,
		},
	)

	client, err := mongo.Connect(ctx, options)
	if err != nil {
		return nil, ErrConnect.M("failed to connect to Mongo in %s", config.MongoURL).
			C("mongoUrl", config.MongoURL).
			C("mongoUsername", config.MongoUsername).
			C("mongoPassword", config.MongoPassword).
			Wrap(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, ErrPing.M("failed to ping Mongo").Wrap(err)
	}

	return client, nil
}
