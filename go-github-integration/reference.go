package gghi

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/mikeschinkel/go-only"
)

// CreateOrFetchReference creates or retrieves a reference at or to the remote GitHub repository
func CreateOrFetchReference(ctx ContextPropertyGetter, base *github.Reference, branch string) (ref *github.Reference, created bool, err error) {
	ref, err = FetchReference(ctx, branch)
	if err != nil {
		ref, err = CreateReference(ctx,
			*base,
			branch)
		created = err == nil
	}
	return ref, created, err
}

// CreateReference creates a reference at the remote GitHub repository
// The reference will be the `branch` name specified and based on `base`.
func CreateReference(ctx ContextPropertyGetter, base github.Reference, branch string) (ref *github.Reference, err error) {

	client, err1 := ctx.GetGithubClient()

	for range only.Once {

		if err1 != nil {
			err = err1
			break
		}

		path := GetRefsHeadsPath(branch)
		ref, _, err = client.Git.CreateRef(
			ctx.Unwrap(),
			ctx.GetOrganizationName(),
			ctx.GetRepositoryName(),
			&github.Reference{
				Ref:    &path,
				Object: base.Object,
			},
		)

		if err != nil {
			err = ErrorReferenceAlreadyExists.Wrap(err, branch)
			break
		}

	}
	return ref, err
}

type referenceMap = map[string]*github.Reference

var references = make(referenceMap, 0)

func GetRefsHeadsPath(name string) string {
	return fmt.Sprintf("refs/heads/%s", name)
}

// FetchReference retrieves a pointer to a github.Reference instance for this branch
func FetchReference(ctx ContextPropertyGetter, branch string) (ref *github.Reference, err error) {

	path := GetRefsHeadsPath(branch)

	for range only.Once {
		var ok bool
		ref, ok = references[branch]
		if ok {
			break
		}

		client, err1 := ctx.GetGithubClient()
		if err1 != nil {
			err = err1
			break
		}

		ref, _, err = client.Git.GetRef(
			ctx.Unwrap(),
			ctx.GetOrganizationName(),
			ctx.GetRepositoryName(),
			path,
		)

		if err != nil {
			// TODO Ensure this error is `branch does not exist` vs. a different error
			ref = &github.Reference{
				Ref: &path,
			}
			err = ErrorUnexpectedReference.Wrap(err,
				*ref.Ref,
				path,
				ctx.GetRepoURL(),
			)
			break
		}
		references[branch] = ref

	}
	if path != *ref.Ref {
		err = ErrorFetchingReference.Wrap(err,
			branch,
			ctx.GetRepoURL(),
		)
		ref = nil
	}
	return ref, err
}
