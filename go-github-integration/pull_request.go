package gghi

import (
	"github.com/google/go-github/github"
	"github.com/mikeschinkel/go-only"
)

type PullRequestState string

const (
	PRIsOpen   PullRequestState = "open"
	PRIsClosed PullRequestState = "closed"
)

// CreatePullRequest returns a pointer to a newly instantiated PullRequest
func CreatePullRequest(ctx ContextPropertyGetter, newpr *github.NewPullRequest) (pr *github.PullRequest, err error) {

	for range only.Once {
		client, err2 := ctx.GetGithubClient()
		if err2 != nil {
			err = err2
			break
		}

		pr, _, err = client.PullRequests.Create(
			ctx.Unwrap(),
			ctx.GetOrganizationName(),
			ctx.GetRepositoryName(),
			newpr,
		)

	}
	return pr, err
}
