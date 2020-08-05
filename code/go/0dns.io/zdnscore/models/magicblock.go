package models

import (
	"context"

	"0dns.io/core/datastore"
	. "0dns.io/core/logging"

	"github.com/0chain/gosdk/core/block"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2/bson"
)

type MagicBlock struct {
	block.MagicBlock
}

func (MagicBlock) GetCollection() *mongo.Collection {
	return datastore.GetStore().GetDB().Collection("magicblocks")
}

func InsertMagicBlock(ctx context.Context, magicBlock *block.MagicBlock) error {
	magicBlocksCollection := (&MagicBlock{}).GetCollection()
	_, err := magicBlocksCollection.InsertOne(ctx, magicBlock)
	if err != nil {
		Logger.Error("Failed to insert magic block data", zap.Error(err))
		return err
	}
	return nil
}

func CheckMagicBlockPresentInDB(ctx context.Context, magicBlockNumber int64) bool {
	var magicBlock block.MagicBlock
	filter := bson.M{"magicblocknumber": magicBlockNumber}
	if err := (&MagicBlock{}).GetCollection().FindOne(ctx, filter).Decode(&magicBlock); err != nil {
		return false
	}
	return true
}

func GetLatestMagicBlockInDB(ctx context.Context) (*MagicBlock, error) {
	var magicBlock block.MagicBlock
	opts := options.FindOne()
	opts.SetSort(bson.M{"magicblocknumber": -1})
	if err := (&MagicBlock{}).GetCollection().FindOne(ctx, bson.M{}, opts).Decode(&magicBlock); err != nil {
		Logger.Error("Failed to get latest magic block from DB", zap.Error(err))
		return nil, err
	}
	return &MagicBlock{magicBlock}, nil
}
