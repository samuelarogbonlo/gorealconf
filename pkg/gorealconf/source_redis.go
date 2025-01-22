package gorealconf

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisSource[T any] struct {
	client  *redis.Client
	key     string
	channel string
}

func NewRedisSource[T any](addr, password, key, channel string) (*RedisSource[T], error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisSource[T]{
		client:  client,
		key:     key,
		channel: channel,
	}, nil
}

func (s *RedisSource[T]) Load(ctx context.Context) (T, error) {
	data, err := s.client.Get(ctx, s.key).Bytes()
	if err != nil {
		var zero T
		return zero, err
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		var zero T
		return zero, err
	}

	return value, nil
}

func (s *RedisSource[T]) Watch(ctx context.Context) (<-chan T, error) {
	ch := make(chan T, 1)
	pubsub := s.client.Subscribe(ctx, s.channel)

	go func() {
		defer close(ch)
		defer pubsub.Close()

		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-pubsub.Channel():
				var newValue T
				if err := json.Unmarshal([]byte(msg.Payload), &newValue); err == nil {
					select {
					case ch <- newValue:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	return ch, nil
}
