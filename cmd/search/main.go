package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"golang.org/x/sync/errgroup"
)

const (
	msgFmt = "==== %s ====\n"
)

func main() {
	uri := flag.String("uri", "http://localhost:19530", "milvus addr")
	user := flag.String("user", "", "milvus user")
	pass := flag.String("pass", "", "milvus password")
	conncurrency := flag.Int("concurrency", 10, "number of concurrent connections")
	requests := flag.Int("reqs", 1000, "number of requests to test")
	dim := flag.Int("dim", 768, "dimension of vector")
	vectorField := flag.String("vector", "title_vector", "vector field name")
	collectionName := flag.String("collection", "medium_articles", "collection name")
	flag.Parse()
	ctx := context.Background()
	fmt.Printf(msgFmt, "start connecting to Milvus")
	c, err := client.NewDefaultGrpcClientWithURI(ctx, *uri, *user, *pass)
	if err != nil {
		log.Fatalf("failed to connect to milvus, err: %v", err)
	}
	fmt.Printf(msgFmt, "connected")
	defer c.Close()

	hasCollection, err := c.HasCollection(ctx, *collectionName)
	if err != nil {
		log.Fatalf("failed to check collection, err: %v", err)
	}
	if !hasCollection {
		log.Fatalf("collection %s not found", *collectionName)
	}

	eg := new(errgroup.Group)

	vector := make(entity.FloatVector, *dim)
	searchVectors := []entity.Vector{
		vector,
	}
	var sendRequests int32
	var totalDuration time.Duration
	searchParam, err := entity.NewIndexAUTOINDEXSearchParam(1)
	if err != nil {
		log.Fatalf("failed to create search param, err: %v", err)
	}
	fmt.Printf("sending search requests to Milvus %d times with concurrency %d\n", *requests, *conncurrency)
	for i := 0; i < *conncurrency; i++ {
		eg.Go(func() error {
			for {
				newSendRequests := atomic.AddInt32(&sendRequests, 1)
				if newSendRequests > int32(*requests) {
					return nil
				}
				startTime := time.Now()
				_, err := c.Search(ctx, *collectionName, []string{}, "", []string{}, searchVectors, *vectorField, entity.L2, 10, searchParam, client.WithSearchQueryConsistencyLevel(entity.ClEventually))
				if err != nil {
					return err
				}
				timeCost := time.Since(startTime)
				atomic.AddInt64((*int64)(&totalDuration), int64(timeCost))
			}
		})
	}
	err = eg.Wait()
	if err != nil {
		log.Fatalf("failed to search, err: %v", err)
	}
	fmt.Printf("avg search latency: %v\n", totalDuration/time.Duration(*requests))
}
