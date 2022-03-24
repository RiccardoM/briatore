VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT  := $(shell git log -1 --format='%H')

export GO111MODULE = on

###############################################################################
###                                   All                                   ###
###############################################################################

all: lint test-unit install

###############################################################################
###                                  Build                                  ###
###############################################################################

build: go.sum
ifeq ($(OS),Windows_NT)
	@echo "building briatore binary..."
	@go build -mod=readonly -o build/briatore.exe ./cmd/briatore
else
	@echo "building briatore binary..."
	@go build -mod=readonly -o build/briatore ./cmd/briatore
endif
.PHONY: build

###############################################################################
###                                 Install                                 ###
###############################################################################

install: go.sum
	@echo "installing briatore binary..."
	@go install -mod=readonly ./cmd/briatore
.PHONY: install

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

coverage:
	@echo "viewing test coverage..."
	@go tool cover --html=coverage.out
.PHONY: coverage

test-unit:
	@echo "Executing unit tests..."
	@go test -mod=readonly -v -coverprofile coverage.txt ./...
.PHONY: test-unit

lint:
	golangci-lint run --out-format=tab
.PHONY: lint

lint-fix:
	golangci-lint run --fix --out-format=tab --issues-exit-code=0
.PHONY: lint-fix

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -name '*.pb.go' | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -name '*.pb.go' | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -name '*.pb.go' | xargs goimports -w -local github.com/riccardom/briatore
.PHONY: format

clean:
	rm -f tools-stamp ./build/**
.PHONY: clean
