package git

import (
	"github.com/newclarity/rep-go-offstage/go-github-integration"
	"github.com/newclarity/rep-go-offstage/wraperr"
)

var (
	ErrorReferenceAlreadyExists    = wraperr.New("reference '%s' already exists")
	ErrorAuthTokenNotSetForContext = wraperr.New("authorization token not set for Context")
	ErrorContentNotInitialized     = wraperr.New("content handle not initialized; must call namespace.FetchGoGithubRepositoryContent() first")
	ErrorCreatingBranch            = wraperr.New("unable to create branch '%s' from '%s' for %s")
	ErrorCreatingGoGitHubReference = wraperr.New("unable to create github.Reference '%s' for %s")
	ErrorFetchingDefaultBranchName = wraperr.New("unable to fetch default branch name for %s")
	ErrorFetchingGoGitHubReference = wraperr.New("unable to fetch github.Reference '%s' for %s")
	ErrorGettingGithubUser         = wraperr.New("failed to get GitHub user")
	ErrorMismatchedOrganization    = wraperr.New("mismatch organization '%s' vs. '%s'")
	ErrorMismatchedRepository      = wraperr.New("mismatch repository '%s' vs. '%s'")
	ErrorNoDefaultBranch           = wraperr.New("no default branch for %s")
)

//goland:noinspection ALL
var (
	ErrorBranchNameCannotBeEmpty         = gghi.ErrorBranchNameCannotBeEmpty
	ErrorContextNotInitialized           = gghi.ErrorContextNotInitialized
	ErrorFetchingPullRequests            = gghi.ErrorFetchingPullRequests
	ErrorFetchingReference               = gghi.ErrorFetchingReference
	ErrorFetchingRepositoriesGetContents = gghi.ErrorFetchingRepositoriesGetContents
	ErrorFetchingRepository              = gghi.ErrorFetchingRepository
	ErrorFetchingRepositoryContent       = gghi.ErrorFetchingRepositoryContent
	ErrorFetchingRepositoryContentList   = gghi.ErrorFetchingRepositoryContentList
	ErrorFilepathCannotBeEmpty           = gghi.ErrorFilepathCannotBeEmpty
	ErrorGettingFileContents             = gghi.ErrorGettingFileContents
	ErrorGettingGithubClient             = gghi.ErrorGettingGithubClient
	ErrorRepositoryFileNotFound          = gghi.ErrorRepositoryFileNotFound
	ErrorUnexpectedReference             = gghi.ErrorUnexpectedReference
)
