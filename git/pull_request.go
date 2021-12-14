package git

import (
	"github.com/google/go-github/github"
	"github.com/mikeschinkel/go-only"
)

type PullRequest struct {
	*github.PullRequest
	Number     int     `json:"number"`
	ID         int64   `json:"id"`
	Title      string  `json:"title"`
	URL        string  `json:"url"`
	HTMLURL    string  `json:"html_url"`
	BranchName string  `json:"branch_name"`
	Head       *Branch `json:"head"`
	Base       *Branch `json:"base"`
}

func NewPullRequestFromGoGithubPullRequest(pr *github.PullRequest) *PullRequest {
	return &PullRequest{
		Number:      *pr.Number,
		ID:          *pr.ID,
		Title:       *pr.Title,
		BranchName:  pr.Head.GetRef(),
		URL:         *pr.URL,
		HTMLURL:     *pr.HTMLURL,
		Head:        NewBranchFromGoGithubPulLRequestBranch(*pr.Head),
		Base:        NewBranchFromGoGithubPulLRequestBranch(*pr.Base),
		PullRequest: pr,
	}
}

// FetchNewPullRequest returns a pointer to a newly instantiated github.PullRequest
func FetchNewPullRequest(ctx ContextPropertyGetter, newpr PullRequestArgs) (pr *PullRequest, err error) {
	for range only.Once {

		gghpr, err1 := CreateGoGithubPullRequest(ctx, &newpr)
		if err1 != nil {
			err = err1
			break
		}
		pr = NewPullRequestFromGoGithubPullRequest(gghpr)
	}
	return pr, err
}

// FetchPullRequestListForContext returns list of open pull requests
// for the Context's Github organization and repository
func FetchPullRequestListForContext(ctx ContextPropertyGetter) (prlist PullRequests, err error) {
	args := PullRequestQueryArgs{
		Organization: ctx.GetOrganizationName(),
		Repository:   ctx.GetRepositoryName(),
	}
	return FetchPullRequestList(ctx, args)
}

// FetchPullRequestList returns
func FetchPullRequestList(ctx ContextPropertyGetter, args PullRequestQueryArgs) (prlist PullRequests, err error) {
	for range only.Once {
		prs, err1 := FetchGoGithubPullRequests(ctx, args)
		if err1 != nil {
			err = err1
			break
		}
		prlist = make(PullRequests, len(prs))
		for i, pr := range prs {
			prlist[i] = NewPullRequestFromGoGithubPullRequest(pr)
		}
	}
	return prlist, err
}
