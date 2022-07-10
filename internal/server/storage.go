package server

import (
	"context"
	"errors"
)

var ErrKeyNotFound = errors.New("key not found")

type Storage interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
	Delete(ctx context.Context, key string) error
}
