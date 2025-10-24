NAME := cecli
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null | sed 's/-\([0-9][0-9]*\)-g/+\1.g/')

# Go build flags
LDFLAGS := -X 'github.com/anselmes/ce-go-template/cmd.Name=$(NAME)' \
           -X 'github.com/anselmes/ce-go-template/cmd.Version=$(VERSION)'

.PHONY: all build clean config ca intermediateca cert amqp amqpcert tls
all: build

build:
	mkdir -p .build
	go build -ldflags "$(LDFLAGS)" -o .build/$(NAME) .

clean:
	go clean .
	rm -rf .build/

config:
	yq '.config' cert.yaml -o json >openssl.json

ca: rootca intermediateca
cert: amqp tls

rootca:
	yq '.ca' cert.yaml -o json >ca.json
	cfssl genkey -config openssl.json -profile ca -initca ca.json | cfssljson -bare ca

intermediateca:
	yq '.intermediate' cert.yaml -o json >intermediate.json
	cfssl gencert -config openssl.json -profile ca -ca ca.pem -ca-key ca-key.pem intermediate.json | cfssljson -bare intermediate
	cat intermediate.pem ca.pem >ca-bundle.pem

tls:
	yq '.tls' cert.yaml -o json >tls.json
	cfssl gencert -config openssl.json -profile tls -ca intermediate.pem -ca-key intermediate-key.pem tls.json | cfssljson -bare tls
	cat tls.pem ca-bundle.pem >tls-bundle.pem
