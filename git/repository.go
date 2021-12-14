package git

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/mikeschinkel/go-only"
)

const (
	APIDomain          = "api.github.com"
	GithubDomain       = "github.com"
	RepositoryIDFormat = "%s/%s/%s"
	PullHTMLURLFormat  = "https://%s/%s/%s/pull/%d"
	PullsAPIURLFormat  = "https://%s/repos/%s/%s/pulls"
)

// Repositories is a named type for a slice of pointers to Repository instances
type Repositories []*Repository

type referenceMap map[string]*github.Reference

// Repository struct for representing a hosted Git repository
type Repository struct {
	Name         string `json:"repository"`
	Organization string `json:"organization"`
	HostDomain   string `json:"host_domain"`
	APIDomain    string `json:"api_domain"`
	repository   *github.Repository
	referenceMap referenceMap
}

// NewRepository returns a pointer to an instance of Repository
func NewRepository(name, org string) *Repository {
	return &Repository{
		Name:         name,
		Organization: org,
		HostDomain:   GithubDomain,
		APIDomain:    APIDomain,
	}
}

// NewRepositoryFromGoGithubRepository returns a pointer to an instance of Repository
// given a *github.Repository
func NewRepositoryFromGoGithubRepository(r *github.Repository) *Repository {
	return NewRepository(r.GetName(),
		r.GetOrganization().GetName())
}

// GetID returns the ID for the repository,
// i.e. the repository URL without prefixing protocol
func (r *Repository) GetID() string {
	return fmt.Sprintf(RepositoryIDFormat,
		r.HostDomain,
		r.Organization,
		r.Name)
}

// GetURL returns the URL for the specified repository
func (r *Repository) GetURL() string {
	return fmt.Sprintf("https://%s", r.GetID())
}

// GetPullHTMLURL returns the HTML URL for a given pull request
func (r *Repository) GetPullHTMLURL(pullno int) string {
	return fmt.Sprintf(PullHTMLURLFormat,
		r.HostDomain,
		r.Organization,
		r.Name,
		pullno)
}

// GetPullsAPIURL returns the API URL for list of pull requests
func (r *Repository) GetPullsAPIURL() string {
	return fmt.Sprintf(PullsAPIURLFormat,
		r.APIDomain,
		r.Organization,
		r.Name)
}

// FetchOrCreateBranch creates a new branch up on GitHub for the current repository
func (r *Repository) FetchOrCreateBranch(ctx ContextPropertyGetter, base Branch, branch string) (b *Branch, existed bool, err error) {
	for range only.Once {
		ghrepo, err1 := r.FetchGoGithubRepositoryForDefaultRepo(ctx)
		if err1 != nil {
			err = err1
			break
		}

		var ghref *github.Reference
		ghref, existed, err = FetchOrCreateGoGithubRepository(ctx,
			*base.GetGoGithubReference(),
			branch)
		if err != nil {
			break
		}

		b = NewBranchFromGoGithubRepositoryAndReference(*ghrepo, *ghref)
	}
	if err != nil {
		err = ErrorCreatingBranch.Wrap(err,
			branch,
			base.Name,
			r.GetURL())
	}
	return b, existed, err
}

// FetchDefaultBranchName returns a pointer to an instance for the context repository's default branch
func (r *Repository) FetchDefaultBranchName(ctx ContextPropertyGetter) (name string, err error) {
	for range only.Once {
		var ghrepo *github.Repository
		ghrepo, err = r.FetchGoGithubRepositoryForDefaultRepo(ctx)
		if err != nil {
			err = ErrorFetchingDefaultBranchName.Wrap(err,
				r.GetURL())
			break
		}
		name = *ghrepo.DefaultBranch
	}
	return name, err
}

// FetchDefaultBranch returns a pointer to a Branch instance for the context repository's default branch
func (r *Repository) FetchDefaultBranch(ctx ContextPropertyGetter) (b *Branch, err error) {
	for range only.Once {
		ghrepo, err1 := r.FetchGoGithubRepositoryForDefaultRepo(ctx)
		if err1 != nil {
			err = err1
			break
		}
		ghref, err2 := FetchGoGithubReference(ctx, *ghrepo.DefaultBranch)
		if err2 != nil {
			// TODO Wrap error
			err = err2
			break
		}
		b = NewBranchFromGoGithubRepositoryAndReference(*ghrepo, *ghref)
	}
	return b, err
}

// FetchGoGithubRepository returns pointer to github.Repository instance for this Repository instance
func (r *Repository) FetchGoGithubRepository(ctx ContextPropertyGetter) (ghr *github.Repository, err error) {
	return FetchGoGithubRepository(&RepositoryArgs{
		Context:      ctx,
		Organization: r.Organization,
		Repository:   r.Name,
		URL:          r.GetURL(),
	})
}

// FetchOrCreateGoGithubRepository retrieves or creates and returns a pointer to a go-github Reference
func FetchOrCreateGoGithubRepository(ctx ContextPropertyGetter, base github.Reference, branch string) (ref *github.Reference, existed bool, err error) {
	for range only.Once {
		existed = false
		ref, err = CreateGoGithubReference(ctx, base, branch)
		if err == nil {
			break
		}
		existed = true
		if err != ErrorReferenceAlreadyExists {
			err = ErrorCreatingGoGitHubReference.Wrap(err,
				branch,
				ctx.GetRepoURL())
			break
		}
		ref, err = FetchGoGithubReference(ctx, branch)
		if err != nil {
			err = ErrorFetchingReference.Wrap(err,
				branch,
				ctx.GetRepoURL())
			break
		}
	}
	return ref, existed, err
}

// FetchGoGithubReference retrieves a pointer to a github.Reference instance for this branch
func (r *Repository) FetchGoGithubReference(ctx ContextPropertyGetter, branch string) (ref *github.Reference, err error) {
	for range only.Once {

		var ok bool
		ref, ok = r.referenceMap[branch]
		if ok {
			break
		}
		if ctx.GetOrganizationName() != r.Organization {
			err = ErrorMismatchedOrganization.Errorf(
				ctx.GetOrganizationName(),
				r.Organization)
			break
		}

		if ctx.GetRepositoryName() != r.Name {
			err = ErrorMismatchedRepository.Errorf(
				ctx.GetRepositoryName(),
				r.Name)
			break
		}

		ref, err = FetchGoGithubReference(ctx, branch)
		if err != nil {
			break
		}

		r.referenceMap[branch] = ref

	}
	if err != nil {
		err = ErrorFetchingGoGitHubReference.Wrap(err,
			branch,
			ctx.GetRepoURL(),
		)
		ref = nil
	}
	return ref, err
}

// FetchGoGithubRepositoryForDefaultRepo returns pointer to a
// github.Repository instance for this Repository instance, and it
// ensures that the repository has a default branch.
func (r *Repository) FetchGoGithubRepositoryForDefaultRepo(ctx ContextPropertyGetter) (ghr *github.Repository, err error) {
	for range only.Once {
		ghr, err = r.FetchGoGithubRepository(ctx)
		if err != nil {
			err = ErrorFetchingRepository.Wrap(err, r.GetURL())
			break
		}
		if *ghr.DefaultBranch == "" {
			err = ErrorNoDefaultBranch.Errorf(r.GetURL())
			break
		}
	}
	return ghr, err
}
