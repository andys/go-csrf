GOCMD=go
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

all: format test
test: 
	 $(GOTEST) -v ./...
deps:
	 $(GOGET) golang.org/x/tools/cmd/goimports
format:
	 goimports -w=true $(shell find . -type f -name '*.go' -not -path './vendor/*' )
