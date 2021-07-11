.PHONY: all
all: build run

.PHONY: build
build:
	go build -o dist/secretary ./cmd/secretary

.PHONY: run
run:
	./dist/secretary

.PHONY: docker
docker: docker-build docker-tag docker-push

.PHONY: docker-build
docker-build:
	docker build --no-cache -t secretary .

.PHONY: docker-tag
docker-tag:
	docker tag secretary:latest sammobach/secretary:latest

.PHONY: docker-push
docker-push:
	docker push sammobach/secretary:latest