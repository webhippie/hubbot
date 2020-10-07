[![Build Status](https://cloud.drone.io/api/badges/webhippie/hubbot/status.svg)](https://cloud.drone.io/webhippie/hubbot)
# hubbot
bot for github in golang with fun

## Local dev

```
go run cmd/main.go --help 
```

## Setup github webhook with ngrok domain for local server
Use ngrok to provide a temporary domain, allowing github to reach your local webhook. 

```
ngrok http 8080
```

You can inspect and replay messages on http://127.0.0.1:4040/inspect/http

## Local dev with docker

```
# build
docker build -t webhippie/hubbot:dev .

# run
docker run --rm -ti webhippie/hubbot:dev --help
```

## Build

```
# build
go build -o bin/hubbot ./cmd

# run
./bin/hubbot --help
```

## Configs

See `hubbot --help``

- Github webhook secret `--hub_webhook_secret | $HUB_WEBHOOK_SECRET`
- Drone access token `--drone_token | $DRONE_TOKEN`
