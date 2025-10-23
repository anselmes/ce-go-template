# CloudEvents Go Template

## Build

```shell
make build
```

### Cleanup

```shell
make clean
```

## TLS

```shell
# generate config
yq '.config' cert.yaml -o json >openssl.json

# generate ca
yq '.ca' cert.yaml -o json >ca.json
cfssl genkey -config openssl.json -profile ca -initca ca.json | cfssljson -bare ca

# generate intermediate ca
yq '.intermediate' cert.yaml -o json >intermediate.json
cfssl gencert -config openssl.json -profile ca -ca ca.pem -ca-key ca-key.pem intermediate.json | cfssljson -bare intermediate

# generate cert
yq '.tls' cert.yaml -o json >tls.json
cfssl gencert -config openssl.json -profile tls -ca intermediate.pem -ca-key intermediate-key.pem tls.json | cfssljson -bare tls
cat tls.pem intermediate.pem ca.pem >tls-bundle.pem
```

## Event Listener

```shell
.build/cecli event listen
```

## Send Event

```shell
.build/cecli event send -d '{"hello": "world"}'
```
