package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/skyquest/server/internal/models"
)

const (
	flightsCacheKey = "flights:all"
	flightCacheTTL  = 5 * time.Minute
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(addr, password string) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	return &RedisClient{client: client}
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

func (r *RedisClient) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// Flight caching methods

func (r *RedisClient) CacheFlights(ctx context.Context, flights []models.Flight) error {
	data, err := json.Marshal(flights)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, flightsCacheKey, data, flightCacheTTL).Err()
}

func (r *RedisClient) GetCachedFlights(ctx context.Context) ([]models.Flight, error) {
	data, err := r.client.Get(ctx, flightsCacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var flights []models.Flight
	if err := json.Unmarshal(data, &flights); err != nil {
		return nil, err
	}
	return flights, nil
}

func (r *RedisClient) InvalidateFlightsCache(ctx context.Context) error {
	return r.client.Del(ctx, flightsCacheKey).Err()
}

// Session caching for active games

func (r *RedisClient) CacheSession(ctx context.Context, sessionID string, session *models.GameSession, ttl time.Duration) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, "session:"+sessionID, data, ttl).Err()
}

func (r *RedisClient) GetCachedSession(ctx context.Context, sessionID string) (*models.GameSession, error) {
	data, err := r.client.Get(ctx, "session:"+sessionID).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var session models.GameSession
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *RedisClient) DeleteSession(ctx context.Context, sessionID string) error {
	return r.client.Del(ctx, "session:"+sessionID).Err()
}

// Rate limiting

func (r *RedisClient) CheckRateLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	current, err := r.client.Incr(ctx, "ratelimit:"+key).Result()
	if err != nil {
		return false, err
	}

	if current == 1 {
		r.client.Expire(ctx, "ratelimit:"+key, window)
	}

	return current <= int64(limit), nil
}

