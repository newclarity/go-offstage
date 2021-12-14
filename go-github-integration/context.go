package gghi

import (
	"context"
	"github.com/google/go-github/github"
)

type ContextPropertyGetter interface {
	IsInitialized() bool
	GetGithubClient() (*github.Client, error)
	GetOrganizationName() string
	GetRepositoryName() string
	GetRepoURL() string
	Unwrap() context.Context
}
