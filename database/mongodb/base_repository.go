package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getOrCreateContext(ctx mongo.SessionContext, timeout time.Duration) (context.Context, context.CancelFunc) {
	if ctx != nil {
		return ctx, func() {}
	}

	return context.WithTimeout(context.Background(), timeout)
}

func (collection MongoCollectionWrapper) Count(sessContext mongo.SessionContext, filter bson.M, opts ...*options.CountOptions) (int64, error) {
	ctx, cancel := getOrCreateContext(sessContext, collection.Timeout)
	defer cancel()

	return collection.Collection.CountDocuments(ctx, filter, opts...)
}

func (collection MongoCollectionWrapper) FindOne(sessContext mongo.SessionContext, filter bson.M, opts ...*options.FindOneOptions) *mongo.SingleResult {
	ctx, cancel := getOrCreateContext(sessContext, collection.Timeout)
	defer cancel()

	return collection.Collection.FindOne(ctx, filter, opts...)
}

func (collection MongoCollectionWrapper) FindMany(sessContext mongo.SessionContext, filter bson.M, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	ctx, cancel := getOrCreateContext(sessContext, collection.Timeout)
	defer cancel()

	return collection.Collection.Find(ctx, filter, opts...)
}

func (collection MongoCollectionWrapper) FindManyWithAggregation(sessContext mongo.SessionContext, pipeline []bson.M, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	ctx, cancel := getOrCreateContext(sessContext, collection.Timeout)
	defer cancel()

	return collection.Collection.Aggregate(ctx, pipeline, opts...)
}

func (collection MongoCollectionWrapper) UpdateOne(sessContext mongo.SessionContext, filter bson.M, update bson.M, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	ctx, cancel := getOrCreateContext(sessContext, collection.Timeout)
	defer cancel()

	return collection.Collection.UpdateOne(ctx, filter, update, opts...)
}

func (collection MongoCollectionWrapper) UpdateMany(sessContext mongo.SessionContext, filter bson.M, update bson.M, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	ctx, cancel := getOrCreateContext(sessContext, collection.Timeout)
	defer cancel()

	return collection.Collection.UpdateMany(ctx, filter, update, opts...)
}

func (collection MongoCollectionWrapper) InsertOne(sessContext mongo.SessionContext, newRecord interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	ctx, cancel := getOrCreateContext(sessContext, collection.Timeout)
	defer cancel()

	return collection.Collection.InsertOne(ctx, newRecord, opts...)
}

func (collection MongoCollectionWrapper) InsertMany(sessContext mongo.SessionContext, newRecords []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	ctx, cancel := getOrCreateContext(sessContext, collection.Timeout)
	defer cancel()

	return collection.Collection.InsertMany(ctx, newRecords, opts...)
}

func (collection MongoCollectionWrapper) DeleteMany(sessContext mongo.SessionContext, filter bson.M, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	ctx, cancel := getOrCreateContext(sessContext, collection.Timeout)
	defer cancel()

	return collection.Collection.DeleteMany(ctx, filter, opts...)
}

func (collection MongoCollectionWrapper) BulkWrite(sessContext mongo.SessionContext, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	ctx, cancel := getOrCreateContext(sessContext, collection.Timeout)
	defer cancel()

	return collection.Collection.BulkWrite(ctx, models, opts...)
}
