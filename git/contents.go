package git

import (
	"github.com/google/go-github/github"
	"github.com/mikeschinkel/go-only"
)

type Contents []*Content

// FetchContents retrieves slice of *Content instances
func FetchContents(ctx ContextPropertyGetter, path string) (chs Contents, err error) {
	for range only.Once {
		rcs, err1 := FetchGoGithubRepositoryContentList(ctx, path)
		if err1 != nil {
			err = err1
			break
		}
		chs = GetContentsFromRepositoryContentList(rcs)
	}
	return chs, err
}

// GetRepositoryContentListFromContents returns slice of *github.RepositoryContent instances
func GetRepositoryContentListFromContents(chs Contents) (rcs []*github.RepositoryContent) {
	rcs = make([]*github.RepositoryContent, len(chs))
	for i, ch := range chs {
		rcs[i] = ch.RepositoryContent
	}
	return rcs
}

// GetContentsFromRepositoryContentList returns slice of *Content instances
func GetContentsFromRepositoryContentList(rcs []*github.RepositoryContent) (chs Contents) {
	chs = make(Contents, len(rcs))
	for i, rc := range rcs {
		chs[i] = NewContent(rc)
	}
	return chs
}
