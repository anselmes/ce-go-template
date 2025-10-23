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

## Prerequesite

- install [krew](https://krew.sigs.k8s.io)
- install [kubectl](https://kubernetes.io/docs/tasks/tools)
- install [kubernetes](https://kubernetes.io)

## RabbitMQ

```shell
kubectl krew install rabbitmq
kubectl rabbitmq install-cluster-operator

make amqp
```

## Event Listener

```shell
.build/cecli event listen
```

## Send Event

```shell
.build/cecli event send -d '{"hello": "world"}'
```
