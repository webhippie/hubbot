package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	// "github.com/google/go-github/github"
	"github.com/drone/drone-go/drone"
	"github.com/webhippie/hubbot/pkg/config"
	"github.com/webhippie/hubbot/pkg/webhookHandler"
	"golang.org/x/oauth2"
)

type configuration struct {
	hub_webhook_secret *string
	drone_server       *string
	drone_token        *string
	drone_debug        *bool
}

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	c := config.New()

	cfg := new(oauth2.Config)
	auther := cfg.Client(
		oauth2.NoContext,
		&oauth2.Token{
			AccessToken: c.Drone.Token,
		},
	)
	droneClient := drone.NewClient(c.Drone.Server, auther)

	wh := new(webhookHandler.WebhookHandler)
	wh.WebhookSecret = c.GitHub.WebhookSecret
	wh.DroneClient = &droneClient

	log.Println("server started")
	http.Handle("/webhook", wh)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
