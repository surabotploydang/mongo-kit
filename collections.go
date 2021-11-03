package mk

import (
	"github.com/fatih/structs"
	"github.com/gertd/go-pluralize"
	changecase "github.com/ku/go-change-case"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)


func (m mongoKit) InitCollection(i interface{}, replace ...string) MongoKit {
	var collectionName string

	if len(replace) != 0 {
		collectionName = replace[0]
	} else {
		plural := pluralize.NewClient()
		collectionName = changecase.Snake(plural.Plural(structs.Name(i)))
	}
	m.collection = m.db.Collection(collectionName)
	return m
}

func (m mongoKit) Collection() *mongo.Collection {
	return m.collection
}

func SchemaCreated() Schema {
	return Schema{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func SchemaUpdated() Schema {
	return Schema{
		UpdatedAt: time.Now(),
	}
}