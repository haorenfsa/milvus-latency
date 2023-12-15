.PHONY: all

all: search-x86-64.tar.gz latency-x86-64.tar.gz

search-x86-64.tar.gz: search
	tar -czf search-x86-64.tar.gz search

latency-x86-64.tar.gz: latency
	tar -czf latency-x86-64.tar.gz latency

latency:
	GOOS=linux go build ./cmd/latency

search:
	GOOS=linux go build ./cmd/search
