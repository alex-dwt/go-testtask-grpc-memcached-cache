package server

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

const (
	KeyNotFoundResp = "KEY_NOT_FOUND"
	OkResp          = "OK"
	DataResp        = "DATA"
)

const (
	GetCommand    = "GET"
	SetCommand    = "SET"
	DeleteCommand = "DELETE"
)

var (
	ErrWrongCommandSignature = errors.New("wrong command signature")
	ErrCommandNotFound       = errors.New("command is not found")
)

type Server struct {
	storage Storage
	logger  *zap.Logger
}

func New(storage Storage, logger *zap.Logger) *Server {
	return &Server{
		storage: storage,
		logger:  logger.Named("server"),
	}
}

func (s *Server) RunCommand(ctx context.Context, command string) (string, error) {
	parts := strings.Split(command, " ")
	if len(parts) == 0 {
		return "", ErrCommandNotFound
	}

	switch parts[0] {
	case GetCommand:
		if len(parts) != 2 {
			return "", ErrWrongCommandSignature
		}

		s.logger.Debug("received GetCommand", zap.String("command", command))

		data, err := s.storage.Get(ctx, parts[1])
		if err != nil {
			if errors.Is(err, ErrKeyNotFound) {
				return KeyNotFoundResp, nil
			}
			return "", fmt.Errorf("storage Get: %w", err)
		}
		return data, nil
	case SetCommand:
		if len(parts) != 3 {
			return "", ErrWrongCommandSignature
		}

		s.logger.Debug("received SetCommand", zap.String("command", command))

		if err := s.storage.Set(ctx, parts[1], parts[2]); err != nil {
			return "", fmt.Errorf("storage Set: %w", err)
		}
		return "", nil
	case DeleteCommand:
		if len(parts) != 2 {
			return "", ErrWrongCommandSignature
		}

		s.logger.Debug("received DeleteCommand", zap.String("command", command))

		if err := s.storage.Delete(ctx, parts[1]); err != nil {
			return "", fmt.Errorf("storage Delete: %w", err)
		}
		return "", nil
	}

	return "", ErrCommandNotFound
}
