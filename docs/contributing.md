

# docs/contributing.md
# Contributing to DynConf

## Development Setup

1. Fork the repository
2. Clone your fork:
```bash
git clone https://github.com/samuelarogbonlo/dynconf.git
```

3. Install dependencies:
```bash
go mod download
```

4. Run tests:
```bash
go test -v ./...
```

## Pull Request Process

1. Create a new branch for your feature
2. Add tests for new functionality
3. Update documentation as needed
4. Submit PR with description of changes

## Code Style

- Follow standard Go conventions
- Use gofmt for formatting
- Add comments for exported symbols