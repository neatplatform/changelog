// Package generate implements the changelog generation logic.
package generate

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/neatplatform/craft/ui"

	"github.com/neatplatform/changelog/internal/changelog"
	"github.com/neatplatform/changelog/internal/changelog/markdown"
	"github.com/neatplatform/changelog/internal/service"
	"github.com/neatplatform/changelog/internal/service/github"
	"github.com/neatplatform/changelog/internal/service/gitlab"
	"github.com/neatplatform/changelog/spec"
)

// Generator is the changelog generator.
type Generator struct {
	ui        ui.UI
	repo      service.Repo
	processor changelog.Processor
}

// New creates a new changelog generator.
func New(s spec.Spec, u ui.UI) (*Generator, error) {
	if u == nil {
		u = ui.NewNop()
	}

	var repo service.Repo

	switch s.Repo.Platform {
	case spec.PlatformGitHub:
		parts := strings.Split(s.Repo.Path, "/")
		if len(parts) != 2 {
			return nil, errors.New("unexpected GitHub repository: cannot parse owner and repo")
		}
		repo = github.NewRepo(u, parts[0], parts[1], s.Repo.AccessToken)

	case spec.PlatformGitLab:
		repo = gitlab.NewRepo(u, s.Repo.Path, s.Repo.AccessToken)

	default:
		return nil, fmt.Errorf("unsupported repository platform: %q", s.Repo.Platform)
	}

	return &Generator{
		ui:        u,
		repo:      repo,
		processor: markdown.NewProcessor(u, s.General.Base, s.General.File),
	}, nil
}

// Generate generates changelogs for a Git repository.
func (g *Generator) Generate(ctx context.Context, s spec.Spec) (string, error) {
	// Parse the existing changelog if any.
	chlog, err := g.processor.Parse(changelog.ParseOptions{})
	if err != nil {
		return "", err
	}

	if err := g.repo.EnsurePermissions(ctx); err != nil {
		return "", err
	}

	/* -------------------- FETCH RELEASE BRANCH -------------------- */

	var branch service.Branch

	if s.Merges.Branch == "" {
		branch, err = g.repo.FetchDefaultBranch(ctx)
	} else {
		branch, err = g.repo.FetchBranch(ctx, s.Merges.Branch)
	}

	if err != nil {
		return "", err
	}

	/* -------------------- FETCH AND FILTER TAGS -------------------- */

	tags, err := g.repo.FetchTags(ctx)
	if err != nil {
		return "", err
	}

	g.ui.Infof(ui.Green, "Sorting and filtering git tags ...")

	sortedTags := tags.Sort()
	sortedTags = sortedTags.Exclude(s.Tags.Exclude...)

	if s.Tags.ExcludeRegex != "" {
		re, err := regexp.CompilePOSIX(s.Tags.ExcludeRegex)
		if err != nil {
			return "", err
		}
		sortedTags = sortedTags.ExcludeRegex(re)
	}

	newTags, err := g.resolveTags(s.Tags, sortedTags, chlog)
	if err != nil {
		return "", err
	}

	if len(newTags) == 0 {
		g.ui.Infof(ui.Green, "Changelog is up-to-date (no new tag or a future tag)")
		return "", nil
	}

	/* -------------------- RESOLVE GIT REVISION FOR COMPARISON -------------------- */

	var baseRev string
	if len(chlog.Existing) > 0 {
		baseRev = chlog.Existing[0].TagName
	} else {
		firstCommit, err := g.repo.FetchFirstCommit(ctx)
		if err != nil {
			return "", err
		}
		baseRev = firstCommit.Hash
	}

	/* -------------------- FETCH COMMITS FOR BRANCH AND TAGS -------------------- */

	// Construct a map of commit hashes to branch and tags names
	// We need to resolve the commit map with all sorted tags, so commits will not be misassigned to new tags
	commitMap, err := g.resolveCommitMap(ctx, branch, sortedTags)
	if err != nil {
		return "", err
	}

	/* -------------------- FETCH & ORGANIZE ISSUES AND MERGES -------------------- */

	// Fetch issues and merges since the last tag on changelog
	var since time.Time
	if len(chlog.Existing) > 0 {
		since = chlog.Existing[0].TagTime
	}

	issues, merges, err := g.repo.FetchIssuesAndMerges(ctx, since)
	if err != nil {
		return "", err
	}

	sortedIssues, sortedMerges := filterByLabels(s, issues, merges)
	g.ui.Infof(ui.Green, "Filtered issues (%d) and pull/merge requests (%d)", len(sortedIssues), len(sortedMerges))

	// We need to resolve the issue map with all sorted tags, so issues will not be misassigned to new tags
	possibleFutureTag := newTags[0]
	issueMap := resolveIssueMap(sortedIssues, sortedTags, possibleFutureTag)
	mergeMap := resolveMergeMap(sortedMerges, commitMap, possibleFutureTag)
	g.ui.Infof(ui.Green, "Partitioned issues and pull/merge requests by tag")

	chlog.New = g.resolveReleases(s, newTags, baseRev, issueMap, mergeMap)
	g.ui.Infof(ui.Green, "Grouped issues and pull/merge requests")

	/* -------------------- UPDATE THE CHANGELOG -------------------- */

	content, err := g.processor.Render(chlog)
	if err != nil {
		return "", err
	}

	if s.General.Print {
		fmt.Print(content)
	}

	return content, nil
}

// resolveTags determines the new tags that should be added to the changelog.
// sortedTags are expected to be sorted from the most recent to the least recent.
// Similarly, chlog.Existing are expected to be sorted from the most recent to the least recent.
// The return value is the list of new tags for generating changelog for them.
func (g *Generator) resolveTags(s spec.Tags, sortedTags service.Tags, chlog *changelog.Changelog) (service.Tags, error) {
	g.ui.Debugf(ui.Cyan, "Resolving new tags for changelog ...")

	mapFunc := func(t service.Tag) string {
		return t.Name
	}

	// Select those tags that are not in changelog
	newTags, _ := sortedTags.Select(func(t service.Tag) bool {
		for _, release := range chlog.Existing {
			if t.Name == release.TagName {
				return false
			}
		}
		return true
	})

	// Resolve the from tag
	if from := s.From; from != "" {
		i := newTags.Index(from)
		if i == -1 {
			return nil, fmt.Errorf("from-tag can be one of %v", newTags.Map(mapFunc))
		}
		// new tags are also sorted from the most recent to the least recent
		newTags = newTags[:i+1]
	}

	// Resolve the to tag
	if to := s.To; to != "" {
		i := newTags.Index(to)
		if i == -1 {
			return nil, fmt.Errorf("to-tag can be one of %v", newTags.Map(mapFunc))
		}
		// new tags are also sorted from the most recent to the least recent
		newTags = newTags[i:]
	}

	// Resolve the future tag
	// The future tag should be the most recent tag (at index zero) if any
	if future := s.Future; future != "" {
		if _, ok := sortedTags.Find(future); ok {
			return nil, fmt.Errorf("future tag cannot be same as an existing tag: %s", future)
		}

		futureTag := g.repo.FutureTag(future)
		newTags = append(service.Tags{futureTag}, newTags...)
	}

	g.ui.Infof(ui.Green, "Resolved new tags for changelog: %v", newTags.Map(mapFunc))

	return newTags, nil
}

// resolveCommitMap returns a map of commit hashes to revisions.
// A revision includes a branch name and a list of tags.
// The resulting map lets us to know what is the branch and all the tags than any given commit falls into.
func (g *Generator) resolveCommitMap(ctx context.Context, branch service.Branch, sortedTags service.Tags) (commitMap, error) {
	commitMap := commitMap{}

	// Resolve which commits are in the branch
	branchCommits, err := g.repo.FetchParentCommits(ctx, branch.Commit.Hash)
	if err != nil {
		return nil, err
	}

	for _, c := range branchCommits {
		if rev, ok := commitMap[c.Hash]; ok {
			rev.Branch = branch.Name
		} else {
			commitMap[c.Hash] = &revisions{
				Branch: branch.Name,
			}
		}
	}

	// Resolve which commits are in the each tag
	// sortedTags are sorted from the most recent to the least recent
	for _, tag := range sortedTags {
		// The first tag can be a future tag without a commit
		if !tag.Commit.IsZero() {
			tagCommits, err := g.repo.FetchParentCommits(ctx, tag.Commit.Hash)
			if err != nil {
				return nil, err
			}

			for _, c := range tagCommits {
				if rev, ok := commitMap[c.Hash]; ok {
					rev.Tags = append(rev.Tags, tag.Name)
				} else {
					commitMap[c.Hash] = &revisions{
						Tags: []string{tag.Name},
					}
				}
			}
		}
	}

	return commitMap, nil
}

func (g *Generator) resolveReleases(s spec.Spec, sortedTags service.Tags, baseRev string, im issueMap, cm mergeMap) []changelog.Release {
	releases := []changelog.Release{}

	for i, tag := range sortedTags {
		releaseURL := s.Content.GetReleaseURL(tag.Name)

		var compareURL string
		if j := i + 1; j < len(sortedTags) {
			compareURL = g.repo.CompareURL(sortedTags[j].Name, tag.Name)
		} else {
			compareURL = g.repo.CompareURL(baseRev, tag.Name)
		}

		// Every tag represents a new release
		release := changelog.Release{
			TagName:    tag.Name,
			TagURL:     tag.WebURL,
			TagTime:    tag.Time,
			ReleaseURL: releaseURL,
			CompareURL: compareURL,
		}

		// Group issues for the current tag
		if issues, ok := im[tag.Name]; ok {
			unselected := issues

			switch s.Issues.Grouping {
			case spec.GroupingMilestone:
				milestones := issues.Milestones()
				g.ui.Debugf(ui.Cyan, "Grouping issues by milestones %s ...", milestones)

				for _, milestone := range milestones {
					f := func(i service.Issue) bool {
						return i.Milestone == milestone
					}

					selected, _ := issues.Select(f)
					_, unselected = unselected.Select(f)

					if len(selected) > 0 {
						title := fmt.Sprintf("Milestone %s", milestone)
						issueGroup := toIssueGroup(title, selected)
						release.IssueGroups = append(release.IssueGroups, issueGroup)
					}
				}

			case spec.GroupingLabel:
				g.ui.Debugf(ui.Cyan, "Grouping issues by labels ...")

				for _, group := range s.Issues.LabelGroups() {
					f := func(i service.Issue) bool {
						return i.Labels.Any(group.Labels...)
					}

					selected, _ := issues.Select(f)
					_, unselected = unselected.Select(f)

					if len(selected) > 0 {
						issueGroup := toIssueGroup(group.Title, selected)
						release.IssueGroups = append(release.IssueGroups, issueGroup)
					}
				}
			}

			if len(unselected) > 0 {
				issueGroup := toIssueGroup("Closed Issues", unselected)
				release.IssueGroups = append(release.IssueGroups, issueGroup)
			}
		}

		// Group merges for the current tag
		if merges, ok := cm[tag.Name]; ok {
			unselected := merges

			switch s.Merges.Grouping {
			case spec.GroupingMilestone:
				milestones := merges.Milestones()
				g.ui.Debugf(ui.Cyan, "Grouping merges by milestones %s ...", milestones)

				for _, milestone := range milestones {
					f := func(m service.Merge) bool {
						return m.Milestone == milestone
					}

					selected, _ := merges.Select(f)
					_, unselected = unselected.Select(f)

					if len(selected) > 0 {
						title := fmt.Sprintf("Milestone %s", milestone)
						mergeGroup := toMergeGroup(title, selected)
						release.MergeGroups = append(release.MergeGroups, mergeGroup)
					}
				}

			case spec.GroupingLabel:
				g.ui.Debugf(ui.Cyan, "Grouping merges by labels ...")

				for _, group := range s.Merges.LabelGroups() {
					f := func(m service.Merge) bool {
						return m.Labels.Any(group.Labels...)
					}

					selected, _ := merges.Select(f)
					_, unselected = unselected.Select(f)

					if len(selected) > 0 {
						mergeGroup := toMergeGroup(group.Title, selected)
						release.MergeGroups = append(release.MergeGroups, mergeGroup)
					}
				}
			}

			if len(unselected) > 0 {
				mergeGroup := toMergeGroup("Merged Changes", unselected)
				release.MergeGroups = append(release.MergeGroups, mergeGroup)
			}
		}

		releases = append(releases, release)
	}

	return releases
}
