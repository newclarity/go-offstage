# Go Github Integration package: `gghi`

The `gghi` package provides a set of functions that simplify calling the functions and methods of the [Google Go-Github](https://gihubt.com/google/go-github) package.

Specfically these functions all accept a `Context` parameter that is this interface allowing the commonly used parameter for these function to be grouped into as passed as a single parameter:

```go 
type Context interface {
    IsInitialized() bool
    GetGithubClient() (*github.Client,error)
    GetOrganizationName() string
    GetRepositoryName() string
    GetRepoURL() string
    Unwrap() context.Context
}
```

This package also defines named errors to simplify detection from calling code: 

- `gghi.ErrorContextNotInitialized`
- `gghi.ErrorCallingRepositoriesGetContents`
- `gghi.ErrorReferenceAlreadyExists`
- `gghi.ErrorRetrievingRepository`
- `gghi.ErrorRetrievingRepositoryContentList`
- `gghi.ErrorRetrievingReference`
- `gghi.ErrorUnexpectedReference`
- `gghi.ErrorRetrievingPullRequests`
- `gghi.ErrorGettingFileContents`
- `gghi.ErrorFilepathCannotBeEmpty`
- `gghi.ErrorBranchNameCannotBeEmpty`
- `gghi.ErrorGettingGithubClient`
- `gghi.ErrorRetrievingRepositoryContent`