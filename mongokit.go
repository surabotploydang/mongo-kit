package mk

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

type Schema struct {
	CreatedAt       time.Time `json:"created_at,omitempty" bson:"created_at,omitempty" faker:"-"`
	UpdatedAt       time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty" faker:"-"`
}

type (
	MongoKit interface {
		Error() error

		Client() *mongo.Client

		Connection(clientOptions *options.ClientOptions) MongoKit
		InitDB(db string, option ...*options.DatabaseOptions) MongoKit
		DB() *mongo.Database
		Ping() MongoKit

		InitCollection(interface{}, ...string) MongoKit
		Collection() *mongo.Collection

		CtxTB(...time.Duration) context.Context
	}

	CTXTimeoutFunc func(...time.Duration) context.Context
	mongoKit struct {
		err error
		client *mongo.Client
		collection *mongo.Collection
		db *mongo.Database
		ctxTimeout CTXTimeoutFunc
	}
)

func ObjectID (id string) primitive.ObjectID {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}
	return objectID
}

func BsonD(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}

func (m mongoKit) CtxTB(duration ...time.Duration) context.Context {
	var timeSet time.Duration
	timeSet = 5
	if len(duration) != 0 {
		timeSet = duration[0]
	}
	ctx, _ := context.WithTimeout(context.Background(), timeSet*time.Second)
	//defer cancel()/
	return ctx
}

func (m mongoKit) DB() *mongo.Database {
	return m.db
}

func (m mongoKit) InitDB(db string, option ...*options.DatabaseOptions) MongoKit {
	m.db = m.client.Database(db, option...)
	return m
}

func (m mongoKit) Client() *mongo.Client {
	return m.client
}

func (m mongoKit) Ping() MongoKit {
	err := m.client.Ping(m.ctxTimeout(2), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf(`✅ MongoDB: connecting`)
	}
	return m
}

func BasicURI() *options.ClientOptions {
	return options.Client().ApplyURI(
		"mongodb://MONGO_USERNAME:MONGO_PASSWORD@MONGO_HOST:MONGO_PORT/?authSource=admin",
		)
}

func (m mongoKit) Connection(uri *options.ClientOptions) MongoKit {
	client, err := mongo.NewClient(uri)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(m.ctxTimeout(10))
	if err != nil {
		log.Fatal(err)
	}
	m.client = client
	return m
}

func (m mongoKit) Error() error {
	return m.err
}

func InitMongoKit() MongoKit {
	return &mongoKit{
		err: nil,
		ctxTimeout: func(second ...time.Duration) context.Context {
			var timeSet time.Duration
			timeSet = 5
			if len(second) != 0 {
				timeSet = second[0]
			}
			ctx, _ := context.WithTimeout(context.Background(), timeSet*time.Second)
			//defer cancel()/
			return ctx
		},
	}
}

