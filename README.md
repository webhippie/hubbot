[![Build Status](https://cloud.drone.io/api/badges/webhippie/hubbot/status.svg)](https://cloud.drone.io/webhippie/hubbot)
# hubbot
bot for github in golang with fun

## Local dev

```
go run main.go --help 
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

# configure
export DRONE_TOKEN=12345
export HUB_WEBHOOK_SECRET=67890

# run
docker run --rm -ti -e DRONE_TOKEN -e HUB_WEBHOOK_SECRET webhippie/hubbot:dev
```

## Build

```
# build
go build -o bin/hubbot main.go

# run
./bin/hubbot --help
```

## Configs

See `hubbot --help``

- Github webhook secret `--hub_webhook_secret | $HUB_WEBHOOK_SECRET`
- Drone access token `--drone_token | $DRONE_TOKEN`
