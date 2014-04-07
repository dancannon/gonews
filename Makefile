all: build

build:
	go build -o bin/gonews
	
deploy:
	go get -u
	go install

build_docker:
	docker build -t dancannon/gonews .
