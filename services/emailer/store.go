package emailer

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storer interface {
	SaveEmail(EmailRecord) (string, error)

	Close(context.Context)
}

type MongoStore struct {
	client *mongo.Client
}

func NewMongoStore(uri string) *MongoStore {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri).SetAuth(options.Credential{
		Username: "root",
		Password: "example",
	}))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("connected to db: %s", uri)

	return &MongoStore{
		client: client,
	}
}

type EmailRecord struct {
	Subject string `bson:"subject"`

	ViewURL string `bson:"view_url"`
}

func (ms MongoStore) SaveEmail(doc EmailRecord) (string, error) {
	coll := ms.client.Database("mailer").Collection("emails")

	insert := bson.M{
		"subject":  doc.Subject,
		"view_url": doc.ViewURL,
	}

	res, err := coll.InsertOne(context.Background(), insert)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(string), nil
}

func (ms MongoStore) Close(ctx context.Context) {
	if err := ms.client.Disconnect(ctx); err != nil {
		log.Panic(err)
	}
}
