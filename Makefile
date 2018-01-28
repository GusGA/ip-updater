# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=ip_updater
BINARY_ARM=$(BINARY_NAME)_arm
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build
build:
	mkdir build && $(GOBUILD) -o build/$(BINARY_NAME) -v
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -rf build/

run:
	mkdir build && $(GOBUILD) -o build/$(BINARY_NAME) -v ./...
	./build/$(BINARY_NAME)
	#deps:
	#$(GOGET) github.com/<account>/<project>


# Cross compilation
build-arm:
	# Compilation for Raspberry Pi model b
	mkdir build && CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 $(GOBUILD) -o build/$(BINARY_ARM) -v
build-linux:
	mkdir build && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o build/$(BINARY_UNIX) -v
