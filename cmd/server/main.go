package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"go.uber.org/zap"

	"github.com/alex-dwt/go-testtask-grpc-memcached-cache/internal/server"
	"github.com/alex-dwt/go-testtask-grpc-memcached-cache/internal/storage"
)

func main() {
	//logger, err := zap.NewProduction()
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Sync()

	srv := server.New(storage.New(), logger)

	l, err := net.Listen("tcp", ":7779")
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}
	defer l.Close()

	logger.Debug("started listening")

	conn, err := l.Accept()
	if err != nil {
		logger.Fatal("failed to accept connection", zap.Error(err))
	}

	logger.Debug("client accepted")

	for {
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			logger.Fatal("failed to read", zap.Error(err))
		}

		res, err := srv.RunCommand(context.TODO(), strings.TrimSpace(data))
		if err != nil {
			logger.Fatal("failed to run command", zap.Error(err))
		}

		var response string
		if res != "" {
			response = fmt.Sprintf("%s %s\n", server.DataResp, res)
		} else {
			response = server.OkResp + "\n"
		}

		logger.Debug("send response to client", zap.String("response", response))

		if _, err := conn.Write([]byte(response)); err != nil {
			logger.Fatal("failed to write", zap.Error(err))
		}
	}
}
