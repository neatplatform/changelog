// Package gitlab provides functionality to interact with GitLab repositories.
package gitlab

import (
	"context"
	"net/http"
	"time"

	"github.com/neatplatform/craft/ui"

	"github.com/neatplatform/changelog/internal/service"
)

// repo implements the service.Repo interface for GitLab.
type repo struct {
	ui          ui.UI
	client      *http.Client
	path        string
	accessToken string
}

// NewRepo creates a new GitLab repository.
func NewRepo(ui ui.UI, path, accessToken string) service.Repo {
	transport := &http.Transport{}
	client := &http.Client{
		Transport: transport,
	}

	return &repo{
		ui:          ui,
		client:      client,
		path:        path,
		accessToken: accessToken,
	}
}

// FutureTag returns a tag that does not exist yet.
func (r *repo) FutureTag(name string) service.Tag {
	return service.Tag{}
}

// CompareURL returns a URL for comparing two revisions.
func (r *repo) CompareURL(base, head string) string {
	return ""
}

// EnsurePermissions ensures the client has all the required permissions.
func (r *repo) EnsurePermissions(context.Context) error {
	return nil
}

// FetchDefaultBranch retrieves the default branch.
func (r *repo) FetchDefaultBranch(ctx context.Context) (service.Branch, error) {
	return service.Branch{}, nil
}

// FetchBranch retrieves a branch by name.
func (r *repo) FetchBranch(ctx context.Context, name string) (service.Branch, error) {
	return service.Branch{}, nil
}

// FetchTags retrieves all tags.
func (r *repo) FetchTags(ctx context.Context) (service.Tags, error) {
	return service.Tags{}, nil
}

// FetchFirstCommit retrieves the firist/initial commit.
func (r *repo) FetchFirstCommit(ctx context.Context) (service.Commit, error) {
	return service.Commit{}, nil
}

// FetchParentCommits retrieves all parent commits of a given commit hash.
func (r *repo) FetchParentCommits(ctx context.Context, hash string) (service.Commits, error) {
	return service.Commits{}, nil
}

// FetchIssuesAndMerges retrieves all closed issues and merged merge requests.
func (r *repo) FetchIssuesAndMerges(ctx context.Context, since time.Time) (service.Issues, service.Merges, error) {
	return service.Issues{}, service.Merges{}, nil
}
