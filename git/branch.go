package git

import (
	"fmt"
	"github.com/google/go-github/github"
	"path"
)

type Branch struct {
	Name  string `json:"name"`
	Label string `json:"label"`
	Ref   string `json:"ref"`
	//    SHA       string `json:"sha"`
	Repo      *Repository
	reference *github.Reference
}

type BranchArgs Branch

func NewBranch(args BranchArgs) *Branch {
	return &Branch{
		Name:  args.Name,
		Label: args.Label,
		Ref:   args.Ref,
		//        SHA:       args.SHA,
		Repo:      args.Repo,
		reference: args.reference,
	}
}

func NewBranchFromGoGithubRepositoryAndReference(r github.Repository, ref github.Reference) (b *Branch) {
	refName := ref.GetRef()
	return &Branch{
		Name:      path.Base(refName),
		Label:     fmt.Sprintf("%s:%s", r.Organization.Name, refName),
		Ref:       refName,
		Repo:      NewRepositoryFromGoGithubRepository(&r),
		reference: &ref,
	}
}

func NewBranchFromGoGithubPulLRequestBranch(ghb github.PullRequestBranch) *Branch {
	return &Branch{
		Label: ghb.GetLabel(),
		Ref:   ghb.GetRef(),
		//        SHA:   ghb.GetSHA(),
		Repo: NewRepositoryFromGoGithubRepository(ghb.GetRepo()),
	}
}

// GetGoGithubReference returns pointer to the github.Reference this Branch represents
func (b *Branch) GetGoGithubReference() *github.Reference {
	return b.reference
}

// GetSHA returns the SHA of the head commit for the branch
//func (b *Branch) GetSHA() (sha string) {
//    for range only.Once {
//        if b.SHA != "" {
//            sha = b.SHA
//            break
//        }
//        if b.reference == nil {
//            break
//        }
//        obj := b.reference.GetObject()
//        if obj == nil {
//            break
//        }
//        sha = obj.GetSHA()
//    }
//    return sha
//}
