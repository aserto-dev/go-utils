package testutil

import (
	"context"

	"github.com/google/go-github/v43/github"
	"golang.org/x/oauth2"
)

const tokenType = "token"

type GithubMopper struct {
	Client *github.Client
}

func NewGithubMopper(accessToken string) *GithubMopper {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: accessToken,
			TokenType:   tokenType,
		},
	)
	clientWithToken := oauth2.NewClient(context.Background(), tokenSource)

	githubClient := github.NewClient(clientWithToken)

	return &GithubMopper{
		Client: githubClient,
	}
}

func (gm *GithubMopper) DeleteRepoFromOrg(org, repo string) error {
	_, err := gm.Client.Repositories.Delete(context.Background(), org, repo)
	return err
}
