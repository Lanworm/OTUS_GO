BIN := "./bin"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN)/calendar -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN)/scheduler -ldflags "$(LDFLAGS)" ./cmd/scheduler
	go build -v -o $(BIN)/sender -ldflags "$(LDFLAGS)" ./cmd/sender

run: build
	$(BIN) -config ./build/config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.55.2

deps-up:
	COMPOSE_PROJECT_NAME="calendar" docker compose -f ./build/local/docker-compose.yml up -d --build

migrate-create:
	if [[ -z "$(NAME)" ]]; then \
		echo "No received param 'NAME'"; \
	else \
		migrate create -ext sql -dir build/app/migrations $(NAME); \
	fi


migrate-up:
	migrate -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" -path build/app/migrations up

migrate-down:
	migrate -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" -path build/app/migrations down

lint: install-lint-deps
	golangci-lint run ./...

generate:
	protoc ./docs/protobuf/calendar.proto --go_out=. --go-grpc_out=.

.PHONY: build run build-img run-img version test lint proto
