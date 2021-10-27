package mongostorage

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"miniurl/generator"
	storage2 "miniurl/storage"
	"time"
)

const dbName = "shortUrls"
const collName = "urls"

type storage struct {
	client *mongo.Client
	urls   *mongo.Collection
}

func NewStorage(mongoURL string) *storage {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		panic(err)
	}

	collection := client.Database(dbName).Collection(collName)

	return &storage{
		client: client,
		urls:   collection,
	}
}

func (s *storage) PutURL(ctx context.Context, url storage2.ShortedURL) (storage2.URLKey, error) {
	for attempt := 0; attempt < 5; attempt++ {
		key := storage2.URLKey(generator.GetRandomKey())
		item := urlItem{
			Key: key,
			URL: url,
		}

		_, err := s.urls.InsertOne(ctx, item)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				continue
			}
			return "", fmt.Errorf("something went wrong - %w", storage2.StorageError)
		}

		return key, nil
	}
	return "", fmt.Errorf("too much attempts during inserting - %w", storage2.ErrCollision)
}

func (s *storage) GetURL(ctx context.Context, key storage2.URLKey) (storage2.ShortedURL, error) {
	var result urlItem
	err := s.urls.FindOne(ctx, bson.M{"_id": key}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", fmt.Errorf("no document with key %v - %w", key, storage2.ErrNotFound)
		}
		return "", fmt.Errorf("somehting went wroing - %w", storage2.StorageError)
	}
	return result.URL, nil
}

func (s *storage) EnsureIndices(ctx context.Context) error {
	// ensure primary index
	indexModels := mongo.IndexModel{
		Keys: bsonx.Doc{{Key: "_id", Value: bsonx.Int32(1)}},
	}
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	_, err := s.urls.Indexes().CreateOne(ctx, indexModels, opts)
	if err != nil {
		return err
	}

	// ensure collection is sharded over primary index
	if err := s.client.Database("admin").RunCommand(ctx, bson.D{
		{"shardCollection", fmt.Sprintf("%s.%s", dbName, collName)},
		{"key", bson.D{{"_id", 1}}}, // range-based sharding
		{"unique", true},
		//"options": bson.M{"locale": "simple"},
	}).Err(); err != nil {
		return err
	}

	return nil
}

type urlItem struct {
	Key storage2.URLKey     `bson:"_id"`
	URL storage2.ShortedURL `bson:"url"`
}
