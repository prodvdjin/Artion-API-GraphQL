package db

import (
	"artion-api-graphql/internal/types"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const fieldId = "_id"

// getTotalCount obtains count of documents in MongoDB collection corresponding the filter.
func (db *MongoDbBridge) getTotalCount(col *mongo.Collection, filter *bson.D) (int64, error) {
	ctx := context.Background()
	totalCount, err := col.CountDocuments(ctx, *filter)
	if err != nil {
		log.Errorf("can not get total count")
	}
	return totalCount, err
}

// findPaginated obtains one page of filtered results from given collection of MongoDB.
func (db *MongoDbBridge) findPaginated(col *mongo.Collection, inputFilter *bson.D, cursor types.Cursor, count int, backward bool) (mc *mongo.Cursor, err error) {
	ctx := context.Background()
	filter := *inputFilter

	if cursor != "" {
		cur, err := cursor.ToObjectId()
		if err != nil {
			log.Errorf("unable to decode cursor %s to ObjectId; %s", cursor, err)
			return nil, err
		}

		if backward {
			filter = append(filter, bson.E{Key: fieldId, Value: bson.D{{Key: "$gt", Value: cur }}})
		} else {
			filter = append(filter, bson.E{Key: fieldId, Value: bson.D{{Key: "$lt", Value: cur }}})
		}
	}

	opt := options.Find()
	if backward {
		opt.SetSort(bson.D{{Key: fieldId, Value: 1 }})
	} else {
		opt.SetSort(bson.D{{Key: fieldId, Value: -1 }})
	}
	opt.SetLimit(int64(count + 1))

	mc, err = col.Find(ctx, filter, opt)
	if err != nil {
		log.Errorf("error loading list; %s", err.Error())
	}
	return mc, err
}
