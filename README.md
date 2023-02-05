### Build with lwip stack
`go build -v -ldflags="-s -w" -trimpath -tags="with_lwip"`

### Supported net stack
Build with tag `gvisor`,`lwip`,`system`

### Enable debug mode
Build with tag `debug`
