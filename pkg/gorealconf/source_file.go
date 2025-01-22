package gorealconf

import (
	"context"
	"encoding/json"
	"os"

	"github.com/fsnotify/fsnotify"
)

type FileSource[T any] struct {
	path    string
	watcher *fsnotify.Watcher
}

func NewFileSource[T any](path string) (*FileSource[T], error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &FileSource[T]{
		path:    path,
		watcher: watcher,
	}, nil
}

func (s *FileSource[T]) Load(ctx context.Context) (T, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		var zero T
		return zero, err
	}

	var config T
	if err := json.Unmarshal(data, &config); err != nil {
		var zero T
		return zero, err
	}

	return config, nil
}

func (s *FileSource[T]) Watch(ctx context.Context) (<-chan T, error) {
	if err := s.watcher.Add(s.path); err != nil {
		return nil, err
	}

	ch := make(chan T, 1)
	go func() {
		defer close(ch)
		defer s.watcher.Close()

		for {
			select {
			case <-ctx.Done():
				return
			case event := <-s.watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					if config, err := s.Load(ctx); err == nil {
						ch <- config
					}
				}
			}
		}
	}()

	return ch, nil
}
