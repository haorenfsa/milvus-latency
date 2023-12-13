package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

const (
	collectionName = "hello_milvus"
	msgFmt         = "==== %s ====\n"
)

func main() {
	uri := flag.String("uri", "http://localhost:19530", "milvus addr")
	user := flag.String("user", "", "milvus user")
	pass := flag.String("pass", "", "milvus password")
	requests := flag.Int("reqs", 1000, "number of requests to test")
	flag.Parse()
	ctx := context.Background()
	fmt.Printf(msgFmt, "start connecting to Milvus")
	c, err := client.NewDefaultGrpcClientWithURI(ctx, *uri, *user, *pass)
	if err != nil {
		log.Fatalf("failed to connect to milvus, err: %v", err)
	}
	fmt.Printf(msgFmt, "connected")
	defer c.Close()

	// check whether collection if exists
	fmt.Printf("sending simple requests to Milvus %d times\n", *requests)
	t := time.Now()
	for i := 0; i < *requests; i++ {
		_, err = c.HasCollection(ctx, collectionName)
		if err != nil {
			break
		}
	}
	if err != nil {
		log.Fatalf("failed to check collection, err: %v", err)
	}
	fmt.Printf("avg simple request latency: %v\n", time.Since(t)/time.Duration(*requests))
}
