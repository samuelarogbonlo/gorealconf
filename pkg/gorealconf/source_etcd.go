package gorealconf

import (
	"context"
	"encoding/json"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdSource[T any] struct {
	client *clientv3.Client
	key    string
}

func NewEtcdSource[T any](endpoints []string, key string) (*EtcdSource[T], error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: endpoints,
	})
	if err != nil {
		return nil, err
	}

	return &EtcdSource[T]{
		client: client,
		key:    key,
	}, nil
}

func (s *EtcdSource[T]) Load(ctx context.Context) (T, error) {
	resp, err := s.client.Get(ctx, s.key)
	if err != nil {
		var zero T
		return zero, err
	}

	if len(resp.Kvs) == 0 {
		var zero T
		return zero, nil
	}

	var config T
	if err := json.Unmarshal(resp.Kvs[0].Value, &config); err != nil {
		var zero T
		return zero, err
	}

	return config, nil
}

func (s *EtcdSource[T]) Watch(ctx context.Context) (<-chan T, error) {
	ch := make(chan T, 1)
	watcher := s.client.Watch(ctx, s.key)

	go func() {
		defer close(ch)
		for resp := range watcher {
			for _, ev := range resp.Events {
				if ev.Type == clientv3.EventTypePut {
					var config T
					if err := json.Unmarshal(ev.Kv.Value, &config); err == nil {
						ch <- config
					}
				}
			}
		}
	}()

	return ch, nil
}
