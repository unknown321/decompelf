PLATFORMS=darwin linux windows
ARCHITECTURES=386 amd64 arm arm64
BINARY=decompelf

build:
	@go build .

build_all:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build -o $(BINARY)-$(GOOS)-$(GOARCH))))

lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.56.2 golangci-lint run -v

clean:
	-@rm -v decompelf \
		src/tinyelf/main \
		decompelf-*

tinyelf:
	cd src/tinyelf/ && go build main.go