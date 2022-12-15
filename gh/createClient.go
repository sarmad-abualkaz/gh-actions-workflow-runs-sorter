package gh

import (
	"context"
	"os"
	
	"github.com/google/go-github/v47/github"
	"golang.org/x/oauth2"

	log "github.com/sirupsen/logrus"
)

func CreateClient() (context.Context, *github.Client) {

	log.WithFields(log.Fields{
	}).Info("Initializing Github client ...")

	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GH_TOKEN")},
	)

	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return ctx, client
}
