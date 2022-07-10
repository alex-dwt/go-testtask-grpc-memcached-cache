package storage

import (
	"context"
	"sync"

	"github.com/alex-dwt/go-testtask-grpc-memcached-cache/internal/server"
)

type Storage struct {
	mu   sync.Mutex
	data map[string]string
}

func New() *Storage {
	return &Storage{
		data: make(map[string]string),
	}
}

func (s *Storage) Get(ctx context.Context, key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if val, found := s.data[key]; found {
		return val, nil
	}

	return "", server.ErrKeyNotFound
}

func (s *Storage) Set(ctx context.Context, key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
	return nil
}

func (s *Storage) Delete(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
	return nil
}
