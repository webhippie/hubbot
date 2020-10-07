package main

import (
	"os"
	"log"
	"flag"
	"reflect"
	"net/http"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"github.com/drone/drone-go/drone"
)

var droneClient = drone.New("")

var (
	hub_webhook_secret = flag.String("hub_webhook_secret", "", "Github webhook secret")
	drone_server = flag.String("drone_server", "https://cloud.drone.io/", "Drone server")
	drone_token = flag.String("drone_token", "", "Drone token")
)

func handleGithubWebhook(w http.ResponseWriter, r *http.Request) {
	log.Println("Webhook triggered")
	payload, err := github.ValidatePayload(r, []byte(*hub_webhook_secret))
	if err != nil {
		log.Printf("error validating request body: err=%s\n", err)
		return
	}
	defer r.Body.Close()

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Printf("could not parse webhook: err=%s\n", err)
		return
	}
	
	log.Println("Event", reflect.TypeOf(event))
	switch e := event.(type) {
	case *github.IssueCommentEvent:
		handleIssueCommitEvent(e)
	default:
		log.Printf("unknown event type %s\n", github.WebHookType(r))
		return
	}
}

func handleIssueCommitEvent(e *github.IssueCommentEvent) {
	log.Println("Issue Comment Event Action", *e.Action)
	log.Println("User:", *e.Sender.Login)
	log.Println("Comment:", *e.Comment.Body)

	owner := "webhippie"
	name := "hubbot"
	build, err := droneClient.BuildLast(owner, name, "")
	if err != nil {
		log.Fatal(err)
	}
	number := int(build.Number)
	
	log.Println("Last Build", number)
	log.Println(owner, name, build)
}

func main() {
	flag.Parse()

	if os.Getenv("HUB_WEBHOOK_SECRET") != "" {
		*hub_webhook_secret = os.Getenv("HUB_WEBHOOK_SECRET")
		log.Println("Using env HUB_WEBHOOK_SECRET")
	}
	if os.Getenv("DRONE_SERVER") != "" {
		*drone_server = os.Getenv("DRONE_SERVER")
		log.Println("Using env DRONE_SERVER")
	}
	if os.Getenv("DRONE_TOKEN") != "" {
		*drone_token = os.Getenv("DRONE_TOKEN")
		log.Println("Using env DRONE_TOKEN")
	}
	if len(*drone_token) == 0 {
		log.Fatal("Error: you must provide your Drone access token.")
	}

	config := new(oauth2.Config)
	auther := config.Client(
		oauth2.NoContext,
		&oauth2.Token{
			AccessToken: *drone_token,
		},
	)
	droneClient = drone.NewClient(*drone_server, auther)
	// log.Println(droneClient.Self())

	log.Println("server started")
	http.HandleFunc("/webhook", handleGithubWebhook)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
