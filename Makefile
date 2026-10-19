.PHONY: all build-http build-cli build-admin build-tablet run-http run-admin run-tablet http admin tablet test test-e2e

all: build-http build-cli build-admin build-tablet

tidy:
	cd backend && go mod tidy

build-http:
	go build -C backend/http -tags netgo -ldflags '-s -w' -o ../../build/http

build-cli:
	go build -C backend/cli -tags netgo -ldflags '-s -w' -o ../../build/bridge-tab

build-admin:
	go build -C backend/admin -tags netgo -ldflags '-s -w' -o ../../build/admin

build-tablet:
	go build -C backend/tablet -tags netgo -ldflags '-s -w' -o ../../build/tablet

run-http:
	./build/http

run-admin:
	./build/admin

run-tablet:
	./build/tablet

http: build-http run-http

admin: build-admin run-admin

tablet: build-tablet run-tablet

test:
	cd backend && go test ./...

# Requires TEST_DATABASE_STRING pointing at a Postgres test database.
test-e2e:
	cd backend && go test ./e2e/ -v -count=1
