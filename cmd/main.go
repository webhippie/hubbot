package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/github"
)

var (
	hub_webhook_secret = flag.String("hub_webhook_secret", "", "Github webhook secret")
)

func handleWebhook(w http.ResponseWriter, r *http.Request) {
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

	switch e := event.(type) {
	case *github.PushEvent:
		log.Println("Push Event")
	case *github.IssueCommentEvent:
		log.Println("Issue Comment Event")
	case *github.CommitCommentEvent:
		log.Println("Commit Comment Event")
	case *github.PullRequestEvent:
		log.Println("Pull Request Event")
	case *github.WatchEvent:
		// https://developer.github.com/v3/activity/events/types/#watchevent
		// someone starred our repository
		if e.Action != nil && *e.Action == "starred" {
			fmt.Printf("%s starred repository %s\n",
				*e.Sender.Login, *e.Repo.FullName)
		}
	default:
		log.Printf("unknown event type %s\n", github.WebHookType(r))
		return
	}
}

func main() {
	flag.Parse()

	if *hub_webhook_secret == "" {
		*hub_webhook_secret = os.Getenv("HUB_WEBHOOK_SECRET")
	}
	log.Println("Webhook secret: ", *hub_webhook_secret)

	log.Println("server started")
	http.HandleFunc("/webhook", handleWebhook)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
