package gghi

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/mikeschinkel/go-only"
	"net/http"
	"strings"
)

var repositoryContentMap = make(map[string]*github.RepositoryContent, 0)

// FetchRepositoryContent retrieves namespace services from Github
// TestCoverage: YES
func FetchRepositoryContent(ctx ContextPropertyGetter, filepath string, branch string) (rc *github.RepositoryContent, err error) {
	var url string
	for range only.Once {
		var ok bool
		key := fmt.Sprintf("%s:%s", branch, filepath)
		if rc, ok = repositoryContentMap[key]; ok {
			break
		}

		url = fmt.Sprintf("%s/%s", ctx.GetRepoURL(), filepath)

		if filepath == "" {
			err = ErrorFilepathCannotBeEmpty.Errorf(url)
			break
		}

		if branch == "" {
			err = ErrorBranchNameCannotBeEmpty.Errorf(url)
			break
		}

		url = fmt.Sprintf("%s?ref=%s", url, branch)

		var client *github.Client
		client, err = ctx.GetGithubClient()
		if err != nil {
			err = ErrorGettingGithubClient.Wrap(err)
			break
		}

		rc, _, _, err = client.Repositories.GetContents(
			ctx.Unwrap(),
			ctx.GetOrganizationName(),
			ctx.GetRepositoryName(),
			filepath,
			&github.RepositoryContentGetOptions{
				Ref: branch,
			})

		if err != nil {
			maybe404, ok := err.(*github.ErrorResponse)
			if ok && is404NotFound(maybe404) {
				err = ErrorRepositoryFileNotFound.Wrap(err, url)
				break
			}
			err = ErrorFetchingRepositoriesGetContents.Wrap(err, url)
			break
		}
		repositoryContentMap[key] = rc

	}
	if err != nil {
		err = ErrorFetchingRepositoryContent.Wrap(err, url)
	}
	return rc, err
}

func is404NotFound(err *github.ErrorResponse) bool {
	return err.Response != nil &&
		err.Response.StatusCode == http.StatusNotFound
}

// FetchRepositoryContentList retrieves github.RepositoryContent instances from the GitHub API
func FetchRepositoryContentList(ctx ContextPropertyGetter, dir string) (rcs []*github.RepositoryContent, err error) {
	for range only.Once {
		var client *github.Client
		client, err = ctx.GetGithubClient()
		if err != nil {
			break
		}
		_, rcs, _, err = client.Repositories.GetContents(
			ctx.Unwrap(),
			ctx.GetOrganizationName(),
			ctx.GetRepositoryName(),
			strings.TrimLeft(dir, "/"),
			&github.RepositoryContentGetOptions{})
		if err != nil {
			err = ErrorGettingFileContents.Wrap(err,
				ctx.GetRepoURL())
			break
		}
	}
	if err != nil {
		err = ErrorFetchingRepositoryContentList.Wrap(err,
			ctx.GetRepoURL())
	}
	return rcs, err
}
