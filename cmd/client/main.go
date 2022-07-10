package main

import (
	"context"
	"log"

	"github.com/alex-dwt/go-testtask-grpc-memcached-cache/pkg/client"
)

func main() {
	c, err := client.New(":7779")
	if err != nil {
		log.Fatalf("client: %s", err)
	}

	ctx := context.Background()

	if err := c.Set(ctx, "some-key", "some-value"); err != nil {
		log.Fatalf("Set: %s", err)
	}

	res, found, err := c.Get(ctx, "some-key")
	if err != nil {
		log.Fatalf("Get: %s", err)
	}
	log.Printf("get 1 - res: %v, found: %v\n", res, found)

	if err := c.Delete(ctx, "some-key"); err != nil {
		log.Fatalf("Delete: %s", err)
	}

	res, found, err = c.Get(ctx, "some-key")
	if err != nil {
		log.Fatalf("Get: %s", err)
	}
	log.Printf("get 2 - res: %v, found: %v\n", res, found)
}
