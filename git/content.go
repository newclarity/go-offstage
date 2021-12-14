package git

import (
	"github.com/google/go-github/github"
	"github.com/mikeschinkel/go-only"
)

type Content struct {
	*github.RepositoryContent
}

func NewContent(rc *github.RepositoryContent) *Content {
	return &Content{
		RepositoryContent: rc,
	}
}

// GetContent retrieves content for namespace YAML file
func (ch *Content) GetContent() (content string, err error) {
	for range only.Once {
		if ch == nil {
			err = ErrorContentNotInitialized
			break
		}
		content, err = ch.RepositoryContent.GetContent()
	}
	if err != nil {
		err = ErrorGettingFileContents.Wrap(err,
			ch.GetURL())
	}
	return content, err
}
