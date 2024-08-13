# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=text2sql
BINARY_UNIX=$(BINARY_NAME)_unix

# Main package location
MAIN_PACKAGE=./cmd/text2sql

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PACKAGE)

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PACKAGE)
	./$(BINARY_NAME)

deps:
	$(GOGET) github.com/spf13/cobra
	$(GOGET) github.com/lib/pq
	$(GOGET) github.com/sashabaranov/go-openai
	$(GOGET) github.com/anthropic-ai/anthropic-sdk-golang

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v $(MAIN_PACKAGE)

docker-build:
	docker build -t $(BINARY_NAME):latest .

tidy:
	$(GOMOD) tidy

.PHONY: all build test clean run deps build-linux docker-build tidy

