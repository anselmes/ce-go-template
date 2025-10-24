# CloudEvents Go Template [![CodeQL](https://github.com/anselmes/ce-go-template/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/anselmes/ce-go-template/actions/workflows/github-code-scanning/codeql)

---

## Build

```shell
make build
```

## TLS

```shell
make config
make ca
make tls
```

## Usage

- install [krew](https://krew.sigs.k8s.io)
- install [kubectl](https://kubernetes.io/docs/tasks/tools)
- install [kubernetes](https://kubernetes.io)

  ```shell
  # set environment variables
  cp -f env.example .env
  source .env
  ```

### Event Webhook

```shell
cecli event webhook

# using cli
cecli event send -d '{"key": "value"}'

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
cecli event listen
```

### Send Event

```shell
cecli event send -d '{"hello": "world"}'
```

## Cleanup

```shell
make clean
```

---

Copyright (c) 2025 Schubert Anselme <schubert@anselm.es>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.
