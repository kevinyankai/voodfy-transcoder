BUILD_FLAGS=CGO_ENABLED=0

build: 
	$(BUILD_FLAGS) go install ./...
.PHONY: build

build-cli:
	$(BUILD_FLAGS) go install ./voodfycli/
.PHONY: build-cli

test:
	go test -short -parallel 6 -race -timeout 30m ./... 
.PHONY: test

run: 
	docker-compose up --build -V

.PHONY: localnet
