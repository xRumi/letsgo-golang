# Let's Go Golang
Learning Go Backend development from the Book **Let's Go by Alex Edwards** and Docker from the internet.

# Prerequisites
1. Go
2. Docker (docker, docker-buildx, docker-compose)

# Building & Running
## With docker
```
$ cp example.env .env      # config
$ docker compose up        # start
$ docker compose down      # stop
```
## Without docker
```
$ cp example.env .env                   # config
$ set -a && . .env && set +a            # apply config
$ go mod download
$ mkdir build
$ go build -o ./build/ ./cmd/web/       # building
$ ./build/web                           # run
$ go run ./cmd/web/                     # build & run
```
### Available Flags
|Flag       |Description|Default|
|-----------|-----------|-------|
|-dsn       |MySQL data source name|*root:very-hard-password@/mydb?parseTime=true*|
|-addr      |HTTP network address to bind|:4000|
|-devMode   |Disables template caching|false|
|||