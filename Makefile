NAME := cecli
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null | sed 's/-\([0-9][0-9]*\)-g/+\1.g/')

PREFIX ?= /usr/local/bin

# Go build flags
LDFLAGS := -X 'github.com/anselmes/ce-go-template/cmd.Name=$(NAME)' \
           -X 'github.com/anselmes/ce-go-template/cmd.Version=$(VERSION)'

.PHONY: all install uninstall proto build clean config rootca ca cert webhook listen send send-test help
all: build cert
	@$(MAKE) proto || echo "⚠️ 'proto' target failed (possibly rate limited), continuing..."
	@echo "🎯 All targets completed successfully!"

help:
	@echo "Available targets:"
	@echo "  build         - Build the $(NAME) binary"
	@echo "  install       - Install $(NAME) to $(PREFIX) (requires elevated privileges)"
	@echo "  uninstall     - Uninstall $(NAME) from $(PREFIX) (requires elevated privileges)"
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

install:
	@if [ ! -d .build ]; then \
		echo "Please run 'make build' first."; \
		exit 1; \
	fi
	sudo -E install -m 0755 .build/$(NAME) $(DESTDIR)$(PREFIX)/$(NAME)
	@echo "✅ Installed $(NAME) to $(DESTDIR)$(PREFIX)/$(NAME)"

uninstall:
	@if [ -f $(DESTDIR)$(PREFIX)/$(NAME) ]; then \
		sudo -E rm -f $(DESTDIR)$(PREFIX)/$(NAME); \
		echo "🗑️ Uninstalled $(NAME) from $(DESTDIR)$(PREFIX)/$(NAME)"; \
	else \
		echo "$(DESTDIR)$(PREFIX)/$(NAME) not found."; \
	fi

proto:
	buf generate
	@echo "🔄 Protocol buffer files generated successfully!"

build:
	@if [ ! -d .build ]; then mkdir -p .build; fi
	go build -ldflags "$(LDFLAGS)" -o .build/$(NAME) .
	@echo "✅ Build complete!"
	@echo "📝 To add cecli to PATH and enable completion, run:"
	@echo "   source .env"

clean:
	go clean .
	rm -rf .build/
	rm -f *.pem *.csr *.json
	@echo "🧹 Clean complete!"

config:
	yq '.config' cert.yaml -o json >openssl.json
	@echo "⚙️ OpenSSL configuration generated from cert.yaml!"

rootca: config
	yq '.ca' cert.yaml -o json >ca.json
	cfssl genkey -config openssl.json -profile ca -initca ca.json | cfssljson -bare ca
	@echo "🔐 Root CA certificate generated successfully!"

ca: rootca
	yq '.intermediate' cert.yaml -o json >intermediate.json
	cfssl gencert \
		-config openssl.json \
		-profile ca \
		-ca ca.pem \
		-ca-key ca-key.pem intermediate.json \
		| cfssljson -bare intermediate
	cat intermediate.pem ca.pem >ca-bundle.pem
	@echo "🔗 Intermediate CA certificate and bundle generated successfully!"

cert: ca
	yq '.tls' cert.yaml -o json >tls.json
	cfssl gencert \
		-config openssl.json \
		-profile tls \
		-ca intermediate.pem \
		-ca-key intermediate-key.pem tls.json \
		| cfssljson -bare tls
	cat tls.pem ca-bundle.pem >tls-bundle.pem
	@echo "🔒 TLS certificates and bundle generated successfully!"

webhook:
	.build/$(NAME) event webhook \
		--cert tls-bundle.pem \
		--key tls-key.pem \
		--port 8080
	@echo "🕸️ Webhook server started on port 8080!"

listen:
	.build/$(NAME) event listen \
		--cert tls-bundle.pem \
		--key tls-key.pem \
		--address localhost \
		--port 8443
	@echo "👂 Event listener started on localhost:8443!"

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
	@echo "📤 Event sent successfully with data: $(DATA)"

send-test:
	.build/$(NAME) event send \
		--cert tls-bundle.pem \
		--key tls-key.pem \
		--address localhost \
		--port 8443 \
		--data '{"message": "Hello from CloudEvent!!!", "timestamp": "'$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")'"}'
	@echo "🧪 Test event sent successfully!"
