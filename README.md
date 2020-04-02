# Docker registry v2 cleaner application

### How to use it:

- Clone the repo
- Don`t forget to prepare & edit the config:

```bash
mv config.example.toml config.toml && vim config.toml
```

- Build the application

```bash
go build ./... # build on your platform
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./... #crossplatform build for linux
```

- Run builded binary to clean the registry

```bash
./docker-registry-cleaner
```

- Profit!
