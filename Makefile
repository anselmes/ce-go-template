NAME := cecli
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null | sed 's/-\([0-9][0-9]*\)-g/+\1.g/')

# Go build flags
LDFLAGS := -X 'github.com/anselmes/ce-go-template/cmd.Name=$(NAME)' \
           -X 'github.com/anselmes/ce-go-template/cmd.Version=$(VERSION)'

.PHONY: all proto build clean config rootca ca cert webhook listen send send-test help
all: build

help:
	@echo "Available targets:"
	@echo "  build         - Build the $(NAME) binary"
	@echo "  proto         - Generate protobuf files"
	@echo "  clean         - Clean build artifacts and generated files"
	@echo "  config        - Generate OpenSSL config from cert.yaml"
	@echo "  rootca        - Generate root CA certificate"
	@echo "  ca            - Generate intermediate CA certificate (depends on rootca)"
	@echo "  cert          - Generate TLS certificates (depends on ca)"
	@echo "  webhook       - Start webhook server"
	@echo "  listen        - Start event listener"
	@echo "  send          - Send event (requires DATA variable)"
	@echo "  send-test     - Send test event"
	@echo "  help          - Show this help message"

proto:
	buf generate

build: proto
	mkdir -p .build
	go build -ldflags "$(LDFLAGS)" -o .build/$(NAME) .
	source .env

clean:
	go clean .
	rm -rf .build/
	rm -f api/*.pb.go
	rm -f *.pem *.csr *.json

config:
	yq '.config' cert.yaml -o json >openssl.json

rootca: config
	yq '.ca' cert.yaml -o json >ca.json
	cfssl genkey -config openssl.json -profile ca -initca ca.json | cfssljson -bare ca

ca: rootca
	yq '.intermediate' cert.yaml -o json >intermediate.json
	cfssl gencert \
		-config openssl.json \
		-profile ca \
		-ca ca.pem \
		-ca-key ca-key.pem intermediate.json \
		| cfssljson -bare intermediate
	cat intermediate.pem ca.pem >ca-bundle.pem

cert: ca
	yq '.tls' cert.yaml -o json >tls.json
	cfssl gencert \
		-config openssl.json \
		-profile tls \
		-ca intermediate.pem \
		-ca-key intermediate-key.pem tls.json \
		| cfssljson -bare tls
	cat tls.pem ca-bundle.pem >tls-bundle.pem

webhook:
	.build/$(NAME) event webhook \
		--cert tls-bundle.pem \
		--key tls-key.pem \
		--port 8080

listen:
	.build/$(NAME) event listen \
		--cert tls-bundle.pem \
		--key tls-key.pem \
		--address localhost \
		--port 8443

send:
	@if [ "x$(DATA)" = "x" ]; then \
		echo "Usage: make send DATA='your-data-here'"; \
		echo "Example: make send DATA='Hello World'"; \
		echo "Example: make send DATA='{\"message\": \"Hello World\"}'"; \
		echo "Or use: make send-test for a quick test"; \
		exit 1; \
	fi
	.build/$(NAME) event send \
		--cert tls-bundle.pem \
		--key tls-key.pem \
		--address localhost \
		--port 8443 \
		--data '$(DATA)'

send-test:
	.build/$(NAME) event send \
		--cert tls-bundle.pem \
		--key tls-key.pem \
		--address localhost \
		--port 8443 \
		--data '{"message": "Hello from CloudEvent!!!", "timestamp": "'$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")'"}'
