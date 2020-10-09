package webhookHandler

import (
	"log"
	"net/http"
	"reflect"
	"github.com/pkg/errors"
	"github.com/drone/drone-go/drone"
	"github.com/google/go-github/github"
)

type WebhookHandler struct {
	WebhookSecret string
	DroneClient *drone.Client
}

func (h *WebhookHandler) parseMessage(r *http.Request) (event interface {}, err error) {
	payload, err := github.ValidatePayload(r, []byte(h.WebhookSecret))
	if err != nil {
		return event, errors.Wrap(err, "error validating request body")
	}
	defer r.Body.Close()

	event, err = github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		return event, errors.Wrap(err, "could not parse webhook")
	}
	return event, err
}

func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Webhook triggered")
	event, err := h.parseMessage(r)
	if err != nil {
		log.Println("Error while parsing webhook message", err)
	}
	log.Println("Event", reflect.TypeOf(event))
	
	switch e := event.(type) {
	case *github.IssueCommentEvent:
		h.handleIssueCommitEvent(e)
		//	if parseMessage(e) == cancelPR
		// 		pr = getGithubPR(e)
		// 		drone.getBuild(pr)
		// 		drone.cancelBuild
		
	default:
		log.Printf("unknown event type %s\n", github.WebHookType(r))
		return
	}
}

func (h *WebhookHandler) handleIssueCommitEvent(e *github.IssueCommentEvent) {
	log.Println("Issue Comment Event Action", *e.Action)
	log.Println("User:", *e.Sender.Login)
	log.Println("Comment:", *e.Comment.Body)

	// if *e.Issue.PullRequestLinks {
	// 	log.Println(*e.Issue.PullRequestLinks.URL)
	// }
	// ctx := context.Background()
	// commits, err := githubClient.PullRequestsService.ListCommits(ctx, "webhippie", "hubbot", 2, nil)
	owner := "webhippie"
	name := "hubbot"
	
	drone := *h.DroneClient
	build, err := drone.BuildLast(owner, name, "")
	if err != nil {
		log.Fatal(err)
	}
	number := int(build.Number)
	
	log.Println("Last Drone Build", number)
	log.Println("Drone Build", build)

	if *e.Comment.Body == "/drone cancel" {
		log.Println("Drone Build canceled")
		drone.BuildCancel(owner, name, number)
	}
	if *e.Comment.Body == "/drone restart" {
		log.Println("Drone Build restarted")
		drone.BuildRestart(owner, name, number, nil)
	}
	//  PullRequestsService.ListCommits() https://github.com/google/go-github/blob/master/github/pulls.go#L359
	// 
}
