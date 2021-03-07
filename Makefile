# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=metatable
BINARY_LINUX=$(BINARY_NAME)_linux

all: build
build: 
	$(GOBUILD) -o $(BINARY_NAME) -v
test: 
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_LINUX)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
deps:
	$(GOGET) -v ./...
linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_LINUX) -v
