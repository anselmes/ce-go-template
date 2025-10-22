NAME := cecli
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null | sed 's/-\([0-9][0-9]*\)-g/+\1.g/')

# Go build flags
LDFLAGS := -X 'github.com/anselmes/ce-go-template/cmd.Name=$(NAME)' \
           -X 'github.com/anselmes/ce-go-template/cmd.Version=$(VERSION)'

.PHONY: all build clean
all: build

build:
	mkdir -p .build
	go build -ldflags "$(LDFLAGS)" -o .build/$(NAME) .

clean:
	go clean .
	rm -rf .build/
