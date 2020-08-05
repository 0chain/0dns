package datastore

import (
	"context"

	"0dns.io/core/common"
	"0dns.io/core/config"
	. "0dns.io/core/logging"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
)

var (
	//DBOpenError - Error opening the db
	DBOpenError = common.NewError("db_open_error", "Error opening the DB connection")
)

const CONNECTION_CONTEXT_KEY = "connection"

type Store struct {
	db *mongo.Database
}

var store Store

func GetStore() *Store {
	return &store
}

func (store *Store) Open(ctx context.Context) error {
	db, err := mongo.Connect(ctx,
		options.Client().
			ApplyURI(config.Configuration.MongoURL).
			SetMaxPoolSize(uint64(config.Configuration.MongoPoolSize)))
	if err != nil {
		Logger.Error("Failed to connect to DB", zap.Error(err))
		return DBOpenError
	}

	err = db.Ping(ctx, readpref.Primary())
	if err != nil {
		Logger.Error("Unable to connect to mongoDB", zap.Error(err))
		return err
	}
	store.db = db.Database(config.Configuration.DBName)
	return nil
}

func (store *Store) GetDB() *mongo.Database {
	return store.db
}
