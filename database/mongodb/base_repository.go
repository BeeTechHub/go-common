package mongodb

import (
	"context"
	"time"

	"github.com/BeeTechHub/go-common/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (collection MongoCollectionWrapper) FindMany(sessContext mongo.SessionContext, records any, filter bson.M, opts ...*options.FindOptions) error {
	ctx, cancel := getOrCreateContext(sessContext, collection.Timeout)
	defer cancel()

	cursor, err := collection.Collection.Find(ctx, filter, opts...)
	if err != nil {
		return err
	}
	defer cursor.Close(context.TODO())

	return cursor.All(context.TODO(), records)
}

func (collection MongoCollectionWrapper) FindManyWithAggregation(sessContext mongo.SessionContext, records any, pipeline []bson.M, opts ...*options.AggregateOptions) error {
	ctx, cancel := getOrCreateContext(sessContext, collection.Timeout)
	defer cancel()

	cursor, err := collection.Collection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return err
	}
	defer cursor.Close(context.TODO())

	return cursor.All(context.TODO(), records)
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

func (collection MongoCollectionWrapper) FindOneById(sessContext mongo.SessionContext, record any, id primitive.ObjectID) error {
	ctx, cancel := getOrCreateContext(sessContext, collection.Timeout)
	defer cancel()

	return collection.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(record)
}

func (collection MongoCollectionWrapper) DeleteOneById(sessContext mongo.SessionContext, id primitive.ObjectID) error {
	ctx, cancel := getOrCreateContext(sessContext, collection.Timeout)
	defer cancel()

	_, err := collection.Collection.DeleteMany(ctx, bson.M{"_id": id})
	return err
}

func (collection MongoCollectionWrapper) FindManyByIds(sessContext mongo.SessionContext, records any, ids []primitive.ObjectID, opts ...*options.FindOptions) error {
	ctx, cancel := getOrCreateContext(sessContext, collection.Timeout)
	defer cancel()

	cursor, err := collection.Collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}}, opts...)
	if err != nil {
		return err
	}
	defer cursor.Close(context.TODO())

	return cursor.All(context.TODO(), records)
}

func (collection MongoCollectionWrapper) DeleteManyByIds(sessContext mongo.SessionContext, ids []primitive.ObjectID) error {
	ctx, cancel := getOrCreateContext(sessContext, collection.Timeout)
	defer cancel()

	_, err := collection.Collection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
	return err
}

func (collection MongoCollectionWrapper) FindPaginated(sessContext mongo.SessionContext, records any, page int64, size int64, filter bson.M, sortParams bson.D, pipe ...bson.M) (int64, error) {
	ctx, cancel := getOrCreateContext(sessContext, collection.Timeout)
	defer cancel()

	matchFilter := bson.M{
		"$match": filter,
	}

	if len(sortParams) == 0 {
		sortParams = bson.D{
			{"_id", 1},
		}
	}

	sort := bson.M{
		"$sort": sortParams,
	}

	limit := bson.M{
		"$limit": size,
	}

	skip := bson.M{
		"$skip": utils.CalculatePaginatedSkip(page, size),
	}

	pipeline := []bson.M{matchFilter}

	if len(pipe) > 0 {
		pipeline = append(pipeline, pipe...)
	}

	pipeline = append(pipeline, sort, limit, skip)

	count, err := collection.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	} else if count == 0 {
		return 0, nil
	}

	cursor, err := collection.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(context.TODO())

	err = cursor.All(context.TODO(), records)
	if err != nil {
		return 0, err
	}

	return count, err
}
