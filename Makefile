all: build

build:
	go build -o bin/gonews
	
deploy:
	go get -u
	go install
	restart gonews

build_docker:
	docker build -t dancannon/gonews .
