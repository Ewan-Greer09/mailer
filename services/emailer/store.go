package emailer

import (
	"context"
	"log"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storer interface {
	SaveEmail(EmailRecord) (string, error)
	GetEmail(uid string) (EmailRecord, error)
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
	CommType string `bson:"email_type"`

	ViewURL string `bson:"view_url"`
}

func (ms MongoStore) SaveEmail(doc EmailRecord) (string, error) {
	coll := ms.client.Database("mailer").Collection("emails")
	uid := uuid.NewString()
	insert := bson.M{
		"email_type": doc.CommType,
		"uuid":       uid,
		"view_url":   doc.ViewURL,
	}

	_, err := coll.InsertOne(context.Background(), insert)
	if err != nil {
		return "", err
	}

	return uid, nil
}

func (ms MongoStore) GetEmail(uid string) (EmailRecord, error) {
	coll := ms.client.Database("mailer").Collection("emails")

	res := coll.FindOne(context.TODO(), bson.M{
		"uuid": uid,
	})
	if res.Err() != nil {
		return EmailRecord{}, res.Err()
	}

	var record EmailRecord
	err := res.Decode(&record)
	if err != nil {
		return EmailRecord{}, err
	}

	return record, nil
}

func (ms MongoStore) Close(ctx context.Context) {
	if err := ms.client.Disconnect(ctx); err != nil {
		log.Panic(err)
	}
}
