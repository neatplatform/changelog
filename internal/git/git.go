// Package git provides functionality to interact with Git repositories.
package git

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/go-git/go-git/v6"
	"github.com/neatplatform/craft/ui"
)

var (
	idPattern       = `[A-Za-z][0-9A-Za-z-]+[0-9A-Za-z]`
	domainPattern   = fmt.Sprintf(`%s\.[A-Za-z]{2,63}`, idPattern)
	repoPathPattern = fmt.Sprintf(`(%s/){1,20}(%s)`, idPattern, idPattern)
	httpsPattern    = fmt.Sprintf(`^https://(%s)/(%s)(\.git)?$`, domainPattern, repoPathPattern)
	sshPattern      = fmt.Sprintf(`^git@(%s):(%s)(\.git)?$`, domainPattern, repoPathPattern)
	httpsRE         = regexp.MustCompile(httpsPattern)
	sshRE           = regexp.MustCompile(sshPattern)
)

// Repo is a Git repository.
type Repo interface {
	GetRemote() (string, string, error)
}

type repo struct {
	ui  ui.UI
	git *git.Repository
}

// NewRepo creates a new instance of Repo.
func NewRepo(u ui.UI, path string) (Repo, error) {
	git, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{
		DetectDotGit: true,
	})

	if err != nil {
		return nil, err
	}

	u.Debugf(ui.Cyan, "Git repository found")

	return &repo{
		ui:  u,
		git: git,
	}, nil
}

// GetRemote returns the domain and path of the repository's origin remote URL.
// It returns an error if there is not exactly one remote with exactly one URL,
// or if the URL does not match a supported protocol (HTTPS or SSH).
func (r *repo) GetRemote() (string, string, error) {
	r.ui.Debugf(ui.Cyan, "Reading git remote URL ...")

	remotes, err := r.git.Remotes()
	if err != nil {
		return "", "", err
	}

	if len(remotes) == 0 {
		return "", "", errors.New("no git remotes found")
	}

	if len(remotes) > 1 {
		return "", "", errors.New("multiple git remotes are not supported")
	}

	urls := remotes[0].Config().URLs

	if len(urls) == 0 {
		return "", "", errors.New("no git remote urls found")
	}

	if len(urls) > 1 {
		return "", "", errors.New("multiple git remote urls are not supported")
	}

	// Parse the remote URL into a domain part a path part
	if matches := httpsRE.FindStringSubmatch(urls[0]); len(matches) == 6 {
		// Git remote url is using HTTPS protocol
		r.ui.Infof(ui.Green, "Git remote URL: %s", urls[0])
		return matches[1], matches[2], nil
	} else if matches := sshRE.FindStringSubmatch(urls[0]); len(matches) == 6 {
		// Git remote url is using SSH protocol
		r.ui.Infof(ui.Green, "Git remote URL: %s", urls[0])
		return matches[1], matches[2], nil
	}

	return "", "", fmt.Errorf("invalid git remote url: %s", urls[0])
}
