package main

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const database = "plotsdb"
const collection = "links"

// Links Record
type LinksRecord struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	Links     []string  `json:"links" bson:"links"`
	Submitted time.Time `json:"submitted,omitempty" bson:"submitted,omitempty"`
}

func process(unit func(*mongo.Collection, *context.Context) error) error {
	url := "mongodb://" + mongoRouterHost + ":27017"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))

	if err != nil {
		log.Printf("Could not establish a connection to %s.\n", url)
		log.Fatal(err)
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Could not properly disconnect from %s.\n", url)
			log.Fatal(err)
		}
	}()

	databaseHandler := client.Database(database)
	collectionHandler := databaseHandler.Collection(collection)

	err = unit(collectionHandler, &ctx)

	if err != nil {
		log.Printf("Could not process the mongodb operation.\n")
		log.Printf("got the following %v", err)
	}

	return err
}

func getLinksRecord(id uuid.UUID) (*LinksRecord, error) {
	var record *LinksRecord
	log.Println("querying the database for known channels")

	querier := func(collection *mongo.Collection, ctx *context.Context) error {
		err := collection.FindOne(*ctx, bson.M{"_id": id.String()}).Decode(&record)

		return err
	}

	err := process(querier)

	if err != nil {
		log.Printf("querying returned an error: %s\n", err.Error())
		return nil, err
	}

	log.Printf("extracted record: %v\n", record)

	return record, nil
}

func storeLinks(links []string) (string, error) {
	var id string

	inserter := func(collection *mongo.Collection, ctx *context.Context) error {
		id = uuid.NewString()
		channelRecord := LinksRecord{
			ID:        id,
			Links:     links,
			Submitted: time.Now(),
		}

		log.Printf("storing record: %v\n", channelRecord)

		_, err := collection.InsertOne(*ctx, channelRecord)

		if err != nil {
			log.Printf("storing returned an error: %s\n", err.Error())
		}

		return err
	}

	err := process(inserter)

	return id, err
}
