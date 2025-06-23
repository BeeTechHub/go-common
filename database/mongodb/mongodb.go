package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClientWrapper struct {
	Client *mongo.Client
}

type MongoDatabaseWrapper struct {
	Database *mongo.Database
}

type MongoCollectionWrapper struct {
	Collection *mongo.Collection
}

func ConnectDB(mongoUri string, timeOutConnection time.Duration) (MongoClientWrapper, error) {
	/*cmdMonitor := &event.CommandMonitor{
		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			log.Print(evt.Command)
		},
	}*/

	fmt.Println("Try to connect to MongoDB:" + mongoUri)
	//mongodbOption := options.Client().ApplyURI(EnvMongoURI()).SetTimeout(constants.TIME_OUT_CONNECTION * time.Second).SetMonitor(cmdMonitor)
	mongodbOption := options.Client().ApplyURI(mongoUri).SetTimeout(timeOutConnection * time.Second)
	client, err := mongo.NewClient(mongodbOption)
	if err != nil {
		fmt.Println(err)
		return MongoClientWrapper{client}, err
	}

	ctx, _ := context.WithTimeout(context.Background(), timeOutConnection*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		fmt.Println(err)
		return MongoClientWrapper{client}, err
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println(err)
		return MongoClientWrapper{client}, err
	}
	fmt.Println("Connected to MongoDB")
	return MongoClientWrapper{client}, nil
}

func (client MongoClientWrapper) GetDatabase(databaseName string) MongoDatabaseWrapper {
	return MongoDatabaseWrapper{client.Client.Database(databaseName)}
}

func (database MongoDatabaseWrapper) GetCollection(collectionName string) MongoCollectionWrapper {
	return MongoCollectionWrapper{database.Database.Collection(collectionName)}
}

func (client MongoClientWrapper) StartSession(opts ...*options.SessionOptions) (mongo.Session, error) {
	return client.Client.StartSession(opts...)
}
