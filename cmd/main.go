package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	// "crypto/tls"
	// "github.com/jackspirou/syscerts"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"github.com/drone/drone-go/drone"
)

var (
	hub_webhook_secret = flag.String("hub_webhook_secret", "", "Github webhook secret")
	drone_server = flag.String("drone_server", "cloud.drone.io", "Drone server")
	drone_token = flag.String("drone_token", "", "Drone token")
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
	case *github.IssueCommentEvent:
		log.Println("Issue Comment Event Action", *e.Action)
		log.Println("User:", *e.Sender.Login)
		log.Println("Comment:", *e.Comment.Body)

	case *github.CommitCommentEvent:
		log.Println("Commit Comment Event")
	case *github.PushEvent:
		log.Println("Push Event")
	case *github.PullRequestEvent:
		log.Println("Pull Request Event")
	case *github.WatchEvent:
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

	if os.Getenv("HUB_WEBHOOK_SECRET") != "" {
		*hub_webhook_secret = os.Getenv("HUB_WEBHOOK_SECRET")
	}
	if os.Getenv("DRONE_SERVER") != "" {
		*drone_server = os.Getenv("DRONE_SERVER")
	}
	if os.Getenv("DRONE_TOKEN") != "" {
		*drone_token = os.Getenv("DRONE_TOKEN")
	}
	if len(*drone_token) == 0 {
		log.Fatal("Error: you must provide your Drone access token.")
	}

	log.Println("Webhook secret: ", *hub_webhook_secret)
	
	*drone_server = strings.TrimRight(*drone_server, "/")
	config := new(oauth2.Config)
	auther := config.Client(
		oauth2.NoContext,
		&oauth2.Token{
			AccessToken: *drone_token,
		},
	)
	// certs := syscerts.SystemRootsPool()
	// tlsConfig := &tls.Config{
	// 	RootCAs:            certs,
	// 	InsecureSkipVerify: true,
	// }
	// trans, _ := auther.Transport.(*oauth2.Transport)
	// trans.Base = &http.Transport{
	// 	TLSClientConfig: tlsConfig,
	// 	Proxy:           http.ProxyFromEnvironment,
	// }
	client := drone.NewClient(*drone_server, auther)

	owner := "webhippie"
	name := "hubbot"
	build, err := client.BuildLast(owner, name, "")
	if err != nil {
		log.Fatal(err)
	}
	number := int(build.Number)
	
	log.Println("Last Build", number)
	log.Println(owner, name, build)


	log.Println("server started")
	http.HandleFunc("/webhook", handleWebhook)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
