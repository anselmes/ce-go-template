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
cat intermediate.pem ca.pem >ca-bundle.pem

# generate cert
yq '.amqp' cert.yaml -o json >amqp.json
cfssl gencert -config openssl.json -profile tls -ca intermediate.pem -ca-key intermediate-key.pem amqp.json | cfssljson -bare amqp

yq '.tls' cert.yaml -o json >tls.json
cfssl gencert -config openssl.json -profile tls -ca intermediate.pem -ca-key intermediate-key.pem tls.json | cfssljson -bare tls
cat tls.pem ca-bundle.pem >tls-bundle.pem
```

## Prerequesite

- install [krew](https://krew.sigs.k8s.io)
- install [kubectl](https://kubernetes.io/docs/tasks/tools)
- install [kubernetes](https://kubernetes.io)

## RabbitMQ

```shell
# install operator
kubectl krew install rabbitmq
kubectl rabbitmq install-cluster-operator

# tls
kubectl create secret tls amqp-tls-secret --cert=amqp.pem --key=amqp-key.pem --dry-run=client -o yaml | kubectl apply -f -

# amqp
kubectl apply -f rabbit.yaml
```

## Event Listener

```shell
.build/cecli event listen
```

## Send Event

```shell
.build/cecli event send -d '{"hello": "world"}'
```
