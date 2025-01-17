package db

import (
	"artion-api-graphql/internal/types"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// initUserCollection initializes collection with indexes and additional parameters.
func (db *MongoDbBridge) initUserCollection(col *mongo.Collection) {
	// prepare index models
	ix := make([]mongo.IndexModel, 0)

	// create indexes
	if _, err := col.Indexes().CreateMany(context.Background(), ix); err != nil {
		db.log.Panicf("can not create indexes for transaction collection; %s", err.Error())
	}

	// log we are done that
	db.log.Debugf("transactions collection initialized")
}

func (db *MongoDbBridge) GetUser(address common.Address) (user *types.User, err error) {
	col := db.client.Database(db.dbName).Collection(types.CoUsers)

	filter := bson.D{{ Key: fieldId, Value: address.String() }}
	result := col.FindOne(context.Background(), filter)

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	var row types.User
	if err = result.Decode(&row); err != nil {
		db.log.Errorf("can not decode user; %s", err.Error())
		return nil, err
	}

	return &row, err
}

func (db *MongoDbBridge) UpsertUser(User *types.User) error {
	if User == nil {
		return fmt.Errorf("no value to store")
	}
	col := db.client.Database(db.dbName).Collection(types.CoUsers)

	filter := bson.D{{ Key: fieldId, Value: User.Address.String() }}
	if _, err := col.ReplaceOne(context.Background(), filter, User, options.Replace().SetUpsert(true)); err != nil {
		db.log.Errorf("can not update User; %s", err)
		return err
	}
	return nil
}
