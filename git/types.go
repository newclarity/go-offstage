package git

import (
	"github.com/google/go-github/github"
	gghi "github.com/newclarity/go-offstage/go-github-integration"
)

type ContextPropertyGetter = gghi.ContextPropertyGetter
type RepositoryArgs = gghi.RepositoryArgs
type PullRequestQueryArgs = gghi.PullRequestQueryArgs

type PullRequestArgs = github.NewPullRequest
