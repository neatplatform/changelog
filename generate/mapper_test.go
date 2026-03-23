package generate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/neatplatform/changelog/internal/changelog"
	"github.com/neatplatform/changelog/internal/service"
	"github.com/neatplatform/changelog/spec"
)

func TestFilterByLabels(t *testing.T) {
	tests := []struct {
		name           string
		s              spec.Spec
		issues         service.Issues
		merges         service.Merges
		expectedIssues service.Issues
		expectedMerges service.Merges
	}{
		{
			name: "None",
			s: spec.Spec{
				Issues: spec.Issues{
					Selection: spec.SelectionNone,
				},
				Merges: spec.Merges{
					Selection: spec.SelectionNone,
				},
			},
			issues:         service.Issues{issue1, issue2},
			merges:         service.Merges{merge1, merge2},
			expectedIssues: service.Issues{},
			expectedMerges: service.Merges{},
		},
		{
			name: "AllWithIncludeLabels",
			s: spec.Spec{
				Issues: spec.Issues{
					Selection:     spec.SelectionAll,
					IncludeLabels: []string{"bug"},
				},
				Merges: spec.Merges{
					Selection:     spec.SelectionAll,
					IncludeLabels: []string{"enhancement"},
				},
			},
			issues:         service.Issues{issue1, issue2},
			merges:         service.Merges{merge1, merge2},
			expectedIssues: service.Issues{issue1},
			expectedMerges: service.Merges{merge1, merge2},
		},
		{
			name: "AllWithExcludeLabels",
			s: spec.Spec{
				Issues: spec.Issues{
					Selection:     spec.SelectionAll,
					ExcludeLabels: []string{"invalid"},
				},
				Merges: spec.Merges{
					Selection:     spec.SelectionAll,
					ExcludeLabels: []string{"enhancement"},
				},
			},
			issues:         service.Issues{issue1, issue2},
			merges:         service.Merges{merge1, merge2},
			expectedIssues: service.Issues{issue1},
			expectedMerges: service.Merges{merge2},
		},
		{
			name: "LabeledWithIncludeLabels",
			s: spec.Spec{
				Issues: spec.Issues{
					Selection:     spec.SelectionLabeled,
					IncludeLabels: []string{"bug"},
				},
				Merges: spec.Merges{
					Selection:     spec.SelectionLabeled,
					IncludeLabels: []string{"enhancement"},
				},
			},
			issues:         service.Issues{issue1, issue2},
			merges:         service.Merges{merge1, merge2},
			expectedIssues: service.Issues{issue1},
			expectedMerges: service.Merges{merge1},
		},
		{
			name: "LabeledWithExcludeLabels",
			s: spec.Spec{
				Issues: spec.Issues{
					Selection:     spec.SelectionLabeled,
					ExcludeLabels: []string{"invalid"},
				},
				Merges: spec.Merges{
					Selection:     spec.SelectionLabeled,
					ExcludeLabels: []string{"enhancement"},
				},
			},
			issues:         service.Issues{issue1, issue2},
			merges:         service.Merges{merge1, merge2},
			expectedIssues: service.Issues{issue1},
			expectedMerges: service.Merges{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			issues, merges := filterByLabels(tc.s, tc.issues, tc.merges)

			assert.Equal(t, tc.expectedIssues, issues)
			assert.Equal(t, tc.expectedMerges, merges)
		})
	}
}

func TestResolveIssueMap(t *testing.T) {
	futureTag := service.Tag{
		Name: "v0.1.4",
		Time: time.Now(),
	}

	tests := []struct {
		name             string
		issues           service.Issues
		sortedTags       service.Tags
		futureTag        service.Tag
		expectedIssueMap issueMap
	}{
		{
			name:       "OK",
			issues:     service.Issues{issue1, issue2},
			sortedTags: service.Tags{tag3, tag2, tag1},
			futureTag:  futureTag,
			expectedIssueMap: issueMap{
				"v0.1.4": service.Issues{issue2},
				"v0.1.3": service.Issues{issue1},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			issueMap := resolveIssueMap(tc.issues, tc.sortedTags, tc.futureTag)

			assert.Equal(t, tc.expectedIssueMap, issueMap)
		})
	}
}

func TestResolveMergeMap(t *testing.T) {
	futureTag := service.Tag{
		Name: "v0.1.4",
	}

	cm := commitMap{
		"20c5414eccaa147f2d6644de4ca36f35293fa43e": &revisions{
			Branch: "main",
		},
		"c414d1004154c6c324bd78c69d10ee101e676059": &revisions{
			Branch: "main",
			Tags:   []string{"v0.1.3"},
		},
		"0251a422d2038967eeaaaa5c8aa76c7067fdef05": &revisions{
			Branch: "main",
			Tags:   []string{"v0.1.3", "v0.1.2"},
		},
		"25aa2bdbaf10fa30b6db40c2c0a15d280ad9f378": &revisions{
			Branch: "main",
			Tags:   []string{"v0.1.3", "v0.1.2", "v0.1.1"},
		},
	}

	tests := []struct {
		name             string
		merges           service.Merges
		commitMap        commitMap
		futureTag        service.Tag
		expectedMergeMap mergeMap
	}{
		{
			name:      "OK",
			merges:    service.Merges{merge1, merge2},
			commitMap: cm,
			futureTag: futureTag,
			expectedMergeMap: mergeMap{
				"v0.1.4": service.Merges{merge2},
				"v0.1.3": service.Merges{merge1},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mergeMap := resolveMergeMap(tc.merges, tc.commitMap, tc.futureTag)

			assert.Equal(t, tc.expectedMergeMap, mergeMap)
		})
	}
}

func TestToIssueGroup(t *testing.T) {
	tests := []struct {
		name               string
		title              string
		issues             service.Issues
		expectedIssueGroup changelog.IssueGroup
	}{
		{
			name:   "OK",
			title:  "Enhancements",
			issues: service.Issues{issue1, issue2},
			expectedIssueGroup: changelog.IssueGroup{
				Title:  "Enhancements",
				Issues: []changelog.Issue{changelogIssue1, changelogIssue2},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			issueGroup := toIssueGroup(tc.title, tc.issues)

			assert.Equal(t, tc.expectedIssueGroup, issueGroup)
		})
	}
}

func TestToMergeGroup(t *testing.T) {
	tests := []struct {
		name               string
		title              string
		merges             service.Merges
		expectedMergeGroup changelog.MergeGroup
	}{
		{
			name:   "OK",
			title:  "Enhancements",
			merges: service.Merges{merge1, merge2},
			expectedMergeGroup: changelog.MergeGroup{
				Title:  "Enhancements",
				Merges: []changelog.Merge{changelogMerge1, changelogMerge2},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mergeGroup := toMergeGroup(tc.title, tc.merges)

			assert.Equal(t, tc.expectedMergeGroup, mergeGroup)
		})
	}
}
