# docs/configuration-sources.md
# Configuration Sources

gorealconf supports multiple configuration sources:

## File Source

```go
source, err := gorealconf.NewFileSource[AppConfig]("config.json")
cfg := gorealconf.New[AppConfig](gorealconf.WithSource[AppConfig](source))
```

## Etcd Source

```go
source, err := gorealconf.NewEtcdSource[AppConfig]([]string{"localhost:2379"}, "/app/config")
cfg := gorealconf.New[AppConfig](gorealconf.WithSource[AppConfig](source))
```

## Consul Source

```go
source, err := gorealconf.NewConsulSource[AppConfig]("localhost:8500", "app/config")
cfg := gorealconf.New[AppConfig](gorealconf.WithSource[AppConfig](source))
```
