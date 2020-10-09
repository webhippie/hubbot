package main

import (
	"os"
	"log"
	"flag"
	"net/http"
	"github.com/pkg/errors"
	// "github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"github.com/drone/drone-go/drone"
	"github.com/webhippie/hubbot/pkg/webhookHandler"
)

type configuration struct {
	hub_webhook_secret *string
	drone_server *string
	drone_token *string
}

func (c *configuration) parseEnv() (err error) {
	c.hub_webhook_secret = flag.String("hub_webhook_secret", "", "Github webhook secret")
	c.drone_server = flag.String("drone_server", "https://cloud.drone.io/", "Drone server")
	c.drone_token = flag.String("drone_token", "", "Drone token")
	flag.Parse()

	if os.Getenv("HUB_WEBHOOK_SECRET") != "" {
		*c.hub_webhook_secret = os.Getenv("HUB_WEBHOOK_SECRET")
		log.Println("Using env HUB_WEBHOOK_SECRET")
	}
	if os.Getenv("DRONE_SERVER") != "" {
		*c.drone_server = os.Getenv("DRONE_SERVER")
		log.Println("Using env DRONE_SERVER")
	}
	if os.Getenv("DRONE_TOKEN") != "" {
		*c.drone_token = os.Getenv("DRONE_TOKEN")
		log.Println("Using env DRONE_TOKEN")
	}
	if len(*c.drone_token) == 0 {
		err = errors.New("Error: you must provide your Drone access token.")
	}
	return err
}

func main() {
	c := new(configuration)
	c.parseEnv()

	config := new(oauth2.Config)
	auther := config.Client(
		oauth2.NoContext,
		&oauth2.Token{
			AccessToken: *c.drone_token,
		},
	)
	droneClient := drone.NewClient(*c.drone_server, auther)
	// log.Println(droneClient.Self())

	// githubClient := github.NewClient(nil)
	
	wh := new(webhookHandler.WebhookHandler)
	wh.WebhookSecret = *c.hub_webhook_secret
	wh.DroneClient = &droneClient

	log.Println("server started")
	http.Handle("/webhook", wh)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
