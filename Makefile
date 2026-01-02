.PHONY: all build-http build-cli run-http http

all: build-http build-cli

tidy:
	cd backend && go mod tidy

build-http:
	go build -C backend/http -tags netgo -ldflags '-s -w' -o ../../build/http

build-cli:
	go build -C backend/cli -tags netgo -ldflags '-s -w' -o ../../build/bridge-tab

run-http:
	./build/http

http: build-http run-http
