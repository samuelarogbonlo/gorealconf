package gorealconf

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hashicorp/consul/api"
)

type ConsulSource[T any] struct {
	client *api.Client
	key    string
}

func NewConsulSource[T any](address string, key string) (*ConsulSource[T], error) {
	config := api.DefaultConfig()
	config.Address = address

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &ConsulSource[T]{
		client: client,
		key:    key,
	}, nil
}

func (s *ConsulSource[T]) Load(ctx context.Context) (T, error) {
	pair, _, err := s.client.KV().Get(s.key, nil)
	if err != nil {
		var zero T
		return zero, err
	}

	if pair == nil {
		var zero T
		return zero, nil
	}

	var value T
	if err := json.Unmarshal(pair.Value, &value); err != nil {
		var zero T
		return zero, err
	}

	return value, nil
}

func (s *ConsulSource[T]) Watch(ctx context.Context) (<-chan T, error) {
	ch := make(chan T, 1)

	go func() {
		defer close(ch)

		var index uint64
		for {
			select {
			case <-ctx.Done():
				return
			default:
				pair, meta, err := s.client.KV().Get(s.key, &api.QueryOptions{
					WaitIndex: index,
					WaitTime:  5 * time.Minute,
				})
				if err != nil {
					time.Sleep(time.Second)
					continue
				}

				if meta.LastIndex <= index {
					continue
				}

				index = meta.LastIndex

				if pair == nil {
					continue
				}

				var value T
				if err := json.Unmarshal(pair.Value, &value); err != nil {
					continue
				}

				select {
				case ch <- value:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return ch, nil
}
