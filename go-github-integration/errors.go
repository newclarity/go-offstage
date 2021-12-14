package gghi

import "github.com/newclarity/rep-go-offstage/wraperr"

var (
	ErrorBranchNameCannotBeEmpty         = wraperr.New("Branch name cannot be empty for '%s'")
	ErrorContextNotInitialized           = wraperr.New("context not initialized")
	ErrorFetchingPullRequests            = wraperr.New("unable to fetch pull requests for '%s'")
	ErrorFetchingReference               = wraperr.New("unable to fetch reference '%s' for %s")
	ErrorFetchingRepositoriesGetContents = wraperr.New("failed calling github.Repositories.GetContents() for %s")
	ErrorFetchingRepository              = wraperr.New("unable to fetch repository at %s")
	ErrorFetchingRepositoryContent       = wraperr.New("unable to fetch github.RepositoryContent for %s")
	ErrorFetchingRepositoryContentList   = wraperr.New("unable to fetch RepositoryContent instances for %s")
	ErrorFilepathCannotBeEmpty           = wraperr.New("filepath cannot be empty for '%s'")
	ErrorGettingFileContents             = wraperr.New("unable to get file contents for %s")
	ErrorGettingGithubClient             = wraperr.New("failed to get GitHub client for %s")
	ErrorReferenceAlreadyExists          = wraperr.New("reference '%s' already exists")
	ErrorRepositoryFileNotFound          = wraperr.New("unable to find repository file %s")
	ErrorUnexpectedReference             = wraperr.New("unexpected reference %s returned; expected %s for %s")
)
