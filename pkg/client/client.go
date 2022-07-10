package client

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/alex-dwt/go-testtask-grpc-memcached-cache/internal/server"
)

type Client struct {
	conn net.Conn
}

func New(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	return &Client{
		conn: conn,
	}, err
}

func (c *Client) Set(ctx context.Context, key, value string) error {
	command := fmt.Sprintf("%s %s %s\n", server.SetCommand, key, value)

	if _, err := c.conn.Write([]byte(command)); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	data, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	data = strings.TrimSpace(data)
	if data != server.OkResp {
		return fmt.Errorf("wrong response: %s", data)
	}

	return nil
}

func (c *Client) Get(ctx context.Context, key string) (string, bool, error) {
	command := fmt.Sprintf("%s %s\n", server.GetCommand, key)

	if _, err := c.conn.Write([]byte(command)); err != nil {
		return "", false, fmt.Errorf("write: %w", err)
	}

	data, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil {
		return "", false, fmt.Errorf("read: %w", err)
	}

	parts := strings.Split(strings.TrimSpace(data), " ")
	if len(parts) != 2 || parts[0] != server.DataResp {
		return "", false, fmt.Errorf("wrong response: %s", data)
	}

	if parts[1] == server.KeyNotFoundResp {
		return "", false, nil
	}

	return parts[1], true, nil
}

func (c *Client) Delete(ctx context.Context, key string) error {
	command := fmt.Sprintf("%s %s\n", server.DeleteCommand, key)

	if _, err := c.conn.Write([]byte(command)); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	data, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	data = strings.TrimSpace(data)
	if data != server.OkResp {
		return fmt.Errorf("wrong response: %s", data)
	}

	return nil
}
