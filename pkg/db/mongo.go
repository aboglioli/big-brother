package db

import (
	"context"

	"github.com/aboglioli/big-brother/pkg/config"
	"github.com/aboglioli/big-brother/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrMongoConnect = errors.Internal.New("mongo.connect")
)

func ConnectMongo(url, database, username, password string) (*mongo.Database, error) {
	config := config.Get()
	ctx := context.Background()

	options := options.Client().ApplyURI(url).SetAuth(
		options.Credential{
			AuthSource: config.MongoAuthSource,
			Username:   username,
			Password:   password,
		},
	)

	client, err := mongo.Connect(ctx, options)
	if err != nil {
		return nil, ErrMongoConnect.M("failed to connect to Mongo in %s", url).
			C("mongoUrl", url).
			C("mongoUsername", username).
			C("mongoPassword", password).
			Wrap(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, ErrMongoConnect.M("failed to ping Mongo").Wrap(err)
	}

	return client.Database(database), nil
}
