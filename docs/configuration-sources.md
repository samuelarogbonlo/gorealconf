# docs/configuration-sources.md
# Configuration Sources

DynConf supports multiple configuration sources:

## File Source

```go
source, err := dynconf.NewFileSource[AppConfig]("config.json")
cfg := dynconf.New[AppConfig](dynconf.WithSource[AppConfig](source))
```

## Etcd Source

```go
source, err := dynconf.NewEtcdSource[AppConfig]([]string{"localhost:2379"}, "/app/config")
cfg := dynconf.New[AppConfig](dynconf.WithSource[AppConfig](source))
```

## Consul Source

```go
source, err := dynconf.NewConsulSource[AppConfig]("localhost:8500", "app/config")
cfg := dynconf.New[AppConfig](dynconf.WithSource[AppConfig](source))
```

