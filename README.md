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
make config
make ca
make tls
```

## Using

- install [krew](https://krew.sigs.k8s.io)
- install [kubectl](https://kubernetes.io/docs/tasks/tools)
- install [kubernetes](https://kubernetes.io)

### RabbitMQ

```shell
kubectl krew install rabbitmq
kubectl rabbitmq install-cluster-operator

make amqp
```

### Event Webhook

```shell
.build/cecli event webhook

# using cli
.build/cecli event send -d '{"key": "value"}'

# using curl
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -H "Ce-Specversion: 1.0" \
  -H "Ce-Type: com.example.string" \
  -H "Ce-Source: example/source" \
  -H "Ce-Id: 1234" \
  -d '{"key": "value"}' \
  --max-time 5
```

### Event Listener

```shell
.build/cecli event listen
```

### Send Event

```shell
.build/cecli event send -d '{"hello": "world"}'
```
