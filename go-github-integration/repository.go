package gghi

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/mikeschinkel/go-only"
)

type RepositoryArgs struct {
	Context      ContextPropertyGetter
	Organization string
	Repository   string
	URL          string
}

func (args *RepositoryArgs) getRepoPath() string {
	return fmt.Sprintf("%s/%s",
		args.Organization,
		args.Repository)
}

var repositoryMap = make(map[string]*github.Repository, 0)

// FetchRepository returns a pointer to a github.Repository instance
func FetchRepository(args *RepositoryArgs) (ghr *github.Repository, err error) {
	for range only.Once {

		var ok bool
		path := args.getRepoPath()
		ghr, ok = repositoryMap[path]
		if ok {
			break
		}

		ctx := args.Context
		if !ctx.IsInitialized() {
			err = ErrorContextNotInitialized
			break
		}

		client, err2 := ctx.GetGithubClient()
		if err2 != nil {
			err = err2
			break
		}

		ghr, _, err = client.Repositories.Get(
			ctx.Unwrap(),
			args.Organization,
			args.Repository)
		if err != nil {
			break
		}

		repositoryMap[path] = ghr

	}
	if err != nil {
		err = ErrorFetchingRepository.Wrap(err,
			args.URL)
	}
	return ghr, err
}
