# BUILD VARIABLES
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VERSION := dev
ifdef COMMIT
	COMMIT := $(COMMIT)
else
	COMMIT := $(shell git rev-parse --short=12 HEAD)
endif

default: backend

backend:
	@CGO_ENABLED=0 go build -ldflags="-X 'main.version=$(VERSION)' -X 'main.commit=$(COMMIT)' -X 'main.date=$(DATE)' -s -w" -o bin/todo-server ./backend/cmd/todo-server

clean:
	@rm bin/todo-server

image:
	@docker build -f backend/Dockerfile -t todo-server:latest --build-arg VERSION=$(VERSION) --build-arg DATE=$(DATE) --build-arg COMMIT=$(COMMIT) .

test:
	go test ./... -cover

itest:
	go test ./... -cover -tags integration

.PHONY: backend clean image test itest