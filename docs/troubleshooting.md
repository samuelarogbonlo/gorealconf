# Troubleshooting

## Common Issues

### Configuration Not Updating
- Check file permissions
- Verify etcd/consul connectivity
- Ensure correct path/key configuration

### Validation Errors
- Check validation function logic
- Verify configuration structure
- Enable debug logging

### Performance Issues
- Reduce update frequency
- Optimize validation functions
- Check network connectivity for remote sources

## Debug Mode
Enable debug logging for more information:
```go
cfg := gorealconf.New[Config](
    gorealconf.WithDebug(true),
)