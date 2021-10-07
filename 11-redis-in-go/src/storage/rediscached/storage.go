package rediscached

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"miniurl/storage"
	"time"
)

func NewStorage(persistentStorage storage.Storage, client *redis.Client) *Storage {
	return &Storage{
		client:            client,
		persistentStorage: persistentStorage,
	}
}

type Storage struct {
	client            *redis.Client
	persistentStorage storage.Storage
}

var _ storage.Storage = (*Storage)(nil)

func (s *Storage) PutURL(ctx context.Context, url storage.ShortedURL) (storage.URLKey, error) {
	key, err := s.persistentStorage.PutURL(ctx, url)
	if err != nil {
		return "", err
	}

	if err := s.storeURL(ctx, key, url); err != nil {
		log.Printf("Failed to insert key %s into cache due to an error: %s\n", key, err)
	}

	return key, nil
}

func (s *Storage) GetURL(ctx context.Context, key storage.URLKey) (storage.ShortedURL, error) {
	result := s.client.Get(ctx, s.fullKey(key))
	switch rawURL, err := result.Result(); {
	case err == redis.Nil:
	// continue execution
	case err != nil:
		return "", fmt.Errorf("%w: failed to get value from redis due to error %s", storage.StorageError, err)
	default:
		log.Printf("Successfully obtained url from cache for key %s", key)
		return storage.ShortedURL(rawURL), nil
	}

	log.Printf("Loading url by key %s from persistent storage", key)
	url, err := s.persistentStorage.GetURL(ctx, key)
	if err != nil {
		return "", err
	}
	if err := s.storeURL(ctx, key, url); err != nil {
		log.Printf("Failed to insert key %s into cache due to an error: %s\n", key, err)
	}

	return url, nil
}

func (s *Storage) storeURL(ctx context.Context, key storage.URLKey, url storage.ShortedURL) error {
	fullKey := s.fullKey(key)
	return s.client.Set(ctx, fullKey, string(url), time.Hour).Err()
}

func (s *Storage) fullKey(key storage.URLKey) string {
	return "surl:" + string(key)
}
