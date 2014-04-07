all: build

build:
	go build -o bin/gonews

build_docker:
	docker build -t dancannon/gonews .
