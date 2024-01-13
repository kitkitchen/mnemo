PHONY: build

build: build_client build_bin

build_bin: 
	go build -o bin/store store.go
build_client:
	cd client && \
	npm run build-web && \
	cd ..

run:
	bin/store

dev:
	go run main.go
