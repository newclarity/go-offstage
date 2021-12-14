package git

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/mikeschinkel/go-only"
	"golang.org/x/oauth2"
	"io/ioutil"
	"os"
	"strings"
)

const (
	TokenFilepathInGitInfo = ".git/info/GITHUB_TOKEN"
)

var _ ContextPropertyGetter = (*Context)(nil)

type Context struct {
	context.Context
	AuthToken     string
	CommitterName string
	repository    *Repository
	initialized   bool
	// clientMap is a global map to allow reuse of *github.Client instances.
	// It is global because
	clientMap map[string]*github.Client
}

// ContextArgs allows arguments to be passed to NewContext
type ContextArgs struct {
	Context       context.Context
	AuthToken     string
	CommitterName string
	Organization  string
	Repository    string
}

// NewContext returns a pointer to an instance of Context
func NewContext(args ContextArgs) *Context {
	if args.Context == nil {
		args.Context = context.Background()
	}
	if args.AuthToken == "" {
		args.AuthToken = GetAuthToken()
	}
	return &Context{
		Context:       args.Context,
		AuthToken:     args.AuthToken,
		CommitterName: args.CommitterName,
		repository:    NewRepository(args.Repository, args.Organization),
		initialized:   true,
		clientMap:     make(map[string]*github.Client, 0),
	}
}

// GetGithubClient returns a GitHub clientMap
func (ctx *Context) GetGithubClient() (client *github.Client, err error) {
	for range only.Once {
		if ctx.AuthToken == "" {
			err = ErrorAuthTokenNotSetForContext.Wrap(err)
			break
		}
		var ok bool
		if client, ok = ctx.clientMap[ctx.AuthToken]; ok {
			break
		}
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: ctx.AuthToken},
		)
		tc := oauth2.NewClient(ctx, ts)

		client = github.NewClient(tc)
		ctx.clientMap[ctx.AuthToken] = client
	}
	if err != nil {
		err = ErrorGettingGithubClient.Wrap(err)
	}
	return client, err
}

// SetAuthToken allows setting the Github token
func (ctx *Context) SetAuthToken(token string) {
	ctx.AuthToken = token
}

// Unwrap returns an unwrapped context.Context
func (ctx *Context) Unwrap() context.Context {
	return ctx.Context
}

// GetRepoURL returns URL for target repository on GitHub
func (ctx *Context) GetRepoURL() string {
	return ctx.repository.GetURL()
}

// IsInitialized returns true if Context was initialized,
// which happens when context is created via NewContext()
func (ctx *Context) IsInitialized() bool {
	return ctx.initialized
}

// GetRepository creates and returns an instance of a Repository
func (ctx *Context) GetRepository() *Repository {
	return ctx.repository
}

// GetCommitterName returns name of the Committer
func (ctx *Context) GetCommitterName() string {
	return ctx.CommitterName
}

// SetCommitterName allows setting the name of the Committer
func (ctx *Context) SetCommitterName(name string) {
	ctx.CommitterName = name
}

// GetOrganizationName returns name of the Organization
func (ctx *Context) GetOrganizationName() string {
	return ctx.repository.Organization
}

// SetOrganizationName allows setting the name of the Organization
func (ctx *Context) SetOrganizationName(name string) {
	ctx.repository.Organization = name
}

// GetRepositoryName returns name of the Repository
func (ctx *Context) GetRepositoryName() string {
	return ctx.repository.Name
}

// SetRepositoryName allows setting the name of the Repository
func (ctx *Context) SetRepositoryName(name string) {
	ctx.repository.Name = name
}

// GetAuthToken retrieves the Github Token from whichever source it can be found
// Options:
//  1. In the file .git/info/GITHUB_TOKEN
//  2. tbd...
//
func GetAuthToken() (gt string) {
	for range only.Once {

		fp, err := getTokenFilepathInGitInfo()
		if err != nil {
			break
		}
		var content []byte
		content, err = ioutil.ReadFile(fp)
		content = bytes.TrimSpace(content)
		invalidToken := len(content) != 40
		isError := err != nil

		if !isError && !invalidToken {
			gt = strings.TrimSpace(string(content))
			break
		}
		if isError {
			err = fmt.Errorf("GITHUB_TOKEN file not found or could not be loaded. Expected to be found at %s; %w",
				fp,
				err)
			break
		}
		if len(content) == 0 {
			err = fmt.Errorf("GITHUB_TOKEN file %s is empty; %w",
				fp,
				err)
			break
		}
		if invalidToken {
			err = fmt.Errorf("token found in %s appears invalid. 40 characters expected; %d characters found, with a value of '%s'; %w",
				fp,
				len(content),
				content,
				err)
			break
		}
	}
	return gt
}

func getTokenFilepathInGitInfo() (fp string, err error) {
	for range only.Once {
		dir, err := getCurrentDir()
		if err != nil {
			break
		}
		fp = strings.Replace(TokenFilepathInGitInfo, "/", string(os.PathSeparator), -1)
		fp = fmt.Sprintf("%s%c%s",
			dir,
			os.PathSeparator,
			fp)
	}
	return fp, err
}

func getCurrentDir() (dir string, err error) {
	dir, err = os.Getwd()
	if err != nil {
		err = fmt.Errorf("failed to access current directory; %w", err)
	}
	return dir, err
}
