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
	mkdir -p bin
	@CGO_ENABLED=0 go build -ldflags="-X 'main.version=$(VERSION)' -X 'main.commit=$(COMMIT)' -X 'main.date=$(DATE)' -s -w" -o bin/todo-server ./backend/cmd/todo-server

frontend:
	mkdir -p bin
	cd frontend && npm run build && cp -R build ../bin

clean:
	@rm bin/todo-server

image:
	@docker build -t todo-app:latest --build-arg VERSION=$(VERSION) --build-arg DATE=$(DATE) --build-arg COMMIT=$(COMMIT) .

test:
	go test ./... -cover

itest:
	go test ./... -cover -tags integration

.PHONY: backend frontend clean image test itest