.PHONY: run build dist clean

COMMIT := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u +"%Y%m%d-%H%M%S")
VERSION := $(COMMIT)-$(BUILD_DATE)

LDFLAGS := -s -w \
	-X 'main.Version=$(VERSION)' \
	-X 'main.Commit=$(COMMIT)' \
	-X 'main.BuildDate=$(BUILD_DATE)'

run: build
	./bin/lattice

build:
	mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o bin/lattice cmd/lattice/main.go

dist: build
	rm -rf dist
	mkdir -p dist/bin
	cp bin/* dist/bin/

clean:
	rm -rf bin
	rm -rf dist

# end
