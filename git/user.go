package git

import (
	"github.com/google/go-github/github"
	"github.com/mikeschinkel/go-only"
)

var user *github.User

// GetUser returns a pointer to a github.User for the API Object's BranchName at HEAD
// Also returns an error object.
func GetUser(ctx ContextPropertyGetter) (u *github.User, err error) {
	for range only.Once {
		if user != nil {
			u = user
			break
		}
		var client *github.Client
		client, err = ctx.GetGithubClient()
		if err != nil {
			break
		}
		u, _, err = client.Users.Get(ctx.Unwrap(), "")
		if err != nil {
			err = ErrorGettingGithubUser.Wrap(err)
			break
		}
		user = u
	}
	return u, err
}
