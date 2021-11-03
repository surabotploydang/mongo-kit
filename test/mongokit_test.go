package test

import (
	"context"
	"fmt"
	"git.onespace.co.th/osgolib/go-mongo-kit"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

type MockingSigle struct {
	mk.Schema `bson:"inline"`
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Foo string `json:"foo" bson:"foo"`
}

func TestMongoKit(t *testing.T) {
	conn := mk.InitMongoKit().Connection(options.Client().ApplyURI(`mongodb://root:example@203.151.50.101:27017/?authSource=admin`)).Ping()
	ctx := conn.InitDB(`mongokit`).InitCollection(MockingSigle{}, `name`)
	ctxColl := ctx.Collection()

	// insert
	mockCreated, _ := ctxColl.InsertOne(ctx.CtxTB(), &MockingSigle{
		Foo: `bar`,
		Schema: mk.SchemaCreated(),
	})
	fmt.Printf("%v", mockCreated.InsertedID)

	// find
	var mockFind MockingSigle
	err := ctxColl.FindOne(context.TODO(), bson.D{{"_id", mockCreated.InsertedID}}).Decode(&mockFind)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", mockFind)

	// update
	valUpdate, err := mk.BsonD(MockingSigle{
		Foo: `bar2`,
		Schema: mk.SchemaUpdated(),
	})
	opts := options.Update().SetUpsert(true)
	update := bson.D{{"$set", valUpdate}}

	result, err := ctxColl.UpdateByID(context.TODO(), mockFind.ID, update, opts)
	if err != nil {
		panic(err)
	}

	if result.MatchedCount != 0 {
		fmt.Println("matched and replaced an existing document")
	}
	if result.UpsertedCount != 0 {
		fmt.Printf("inserted a new document with ID %v\n", result.UpsertedID)
	}

	// delete
	resDelete, err := ctxColl.DeleteOne(context.TODO(), bson.D{{`_id`,mockFind.ID}})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", resDelete)

}
