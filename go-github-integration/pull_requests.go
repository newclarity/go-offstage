package gghi

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/mikeschinkel/go-only"
)

var defaultState PullRequestState = PRIsOpen

type PullRequestQueryArgs struct {
	BranchName   string
	Organization string
	Repository   string
	State        PullRequestState
}

// FetchPullRequests returns a pointer to an slice of github.PullRequest instances
func FetchPullRequests(ctx ContextPropertyGetter, args PullRequestQueryArgs) (prs []*github.PullRequest, err error) {
	for range only.Once {

		var client *github.Client
		client, err = ctx.GetGithubClient()
		if err != nil {
			break
		}

		if args.State == "" {
			args.State = defaultState
		}

		opts := &github.PullRequestListOptions{
			State: string(args.State),
		}

		if args.Organization == "" {
			args.Organization = ctx.GetOrganizationName()
		}
		if args.Repository == "" {
			args.Repository = ctx.GetRepositoryName()
		}

		if args.BranchName != "" {
			opts.Head = fmt.Sprintf("%s:%s",
				args.Organization,
				args.BranchName)
		}

		prs, _, err = client.PullRequests.List(
			ctx.Unwrap(),
			args.Organization,
			args.Repository,
			opts)

		if len(prs) == 0 {
			// Not an error, we just don't have any PRs
			break
		}

		if err != nil {
			err = ErrorFetchingPullRequests.Wrap(err, ctx.GetRepoURL())
			break
		}
	}
	return prs, err
}
