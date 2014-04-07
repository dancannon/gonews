all: build

build:
	go build -o bin/gonews
	
deploy:
	git pull origin master
	go install
	restart gonews

build_docker:
	docker build -t dancannon/gonews .
