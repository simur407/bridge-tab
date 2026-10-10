.PHONY: all build-http build-cli build-admin run-http run-admin http admin test test-e2e

all: build-http build-cli build-admin

tidy:
	cd backend && go mod tidy

build-http:
	go build -C backend/http -tags netgo -ldflags '-s -w' -o ../../build/http

build-cli:
	go build -C backend/cli -tags netgo -ldflags '-s -w' -o ../../build/bridge-tab

build-admin:
	go build -C backend/admin -tags netgo -ldflags '-s -w' -o ../../build/admin

run-http:
	./build/http

run-admin:
	./build/admin

http: build-http run-http

admin: build-admin run-admin

test:
	cd backend && go test ./...

# Requires TEST_DATABASE_STRING pointing at a Postgres test database.
test-e2e:
	cd backend && go test ./e2e/ -v -count=1
