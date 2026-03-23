package generate

import (
	"context"
	"time"

	"github.com/neatplatform/changelog/internal/changelog"
	"github.com/neatplatform/changelog/internal/service"
)

type (
	GetRemoteMock struct {
		OutDomain string
		OutPath   string
		OutError  error
	}

	MockGitRepo struct {
		GetRemoteIndex int
		GetRemoteMocks []GetRemoteMock
	}
)

func (m *MockGitRepo) GetRemote() (string, string, error) {
	if m.GetRemoteIndex >= len(m.GetRemoteMocks) {
		panic("GetRemote called more times than expected")
	}

	i := m.GetRemoteIndex
	m.GetRemoteIndex++

	return m.GetRemoteMocks[i].OutDomain, m.GetRemoteMocks[i].OutPath, m.GetRemoteMocks[i].OutError
}

type (
	FutureTagMock struct {
		InName string
		OutTag service.Tag
	}

	CompareURLMock struct {
		InBase    string
		InHead    string
		OutString string
	}

	EnsurePermissionsMock struct {
		InContext context.Context
		OutError  error
	}

	FetchDefaultBranchMock struct {
		InContext context.Context
		OutBranch service.Branch
		OutError  error
	}

	FetchBranchMock struct {
		InContext context.Context
		InName    string
		OutBranch service.Branch
		OutError  error
	}

	FetchTagsMock struct {
		InContext context.Context
		OutTags   service.Tags
		OutError  error
	}

	FetchFirstCommitMock struct {
		InContext context.Context
		OutCommit service.Commit
		OutError  error
	}

	FetchParentCommitsMock struct {
		InContext  context.Context
		InHash     string
		OutCommits service.Commits
		OutError   error
	}

	FetchIssuesAndMergesMock struct {
		InContext context.Context
		InSince   time.Time
		OutIssues service.Issues
		OutMerges service.Merges
		OutError  error
	}

	MockRepo struct {
		EnsurePermissionsIndex int
		EnsurePermissionsMocks []EnsurePermissionsMock

		FutureTagIndex int
		FutureTagMocks []FutureTagMock

		CompareURLIndex int
		CompareURLMocks []CompareURLMock

		FetchFirstCommitIndex int
		FetchFirstCommitMocks []FetchFirstCommitMock

		FetchBranchIndex int
		FetchBranchMocks []FetchBranchMock

		FetchDefaultBranchIndex int
		FetchDefaultBranchMocks []FetchDefaultBranchMock

		FetchTagsIndex int
		FetchTagsMocks []FetchTagsMock

		FetchIssuesAndMergesIndex int
		FetchIssuesAndMergesMocks []FetchIssuesAndMergesMock

		FetchParentCommitsIndex int
		FetchParentCommitsMocks []FetchParentCommitsMock
	}
)

func (m *MockRepo) FutureTag(name string) service.Tag {
	if m.FutureTagIndex >= len(m.FutureTagMocks) {
		panic("FutureTag called more times than expected")
	}

	i := m.FutureTagIndex
	m.FutureTagIndex++

	m.FutureTagMocks[i].InName = name
	return m.FutureTagMocks[i].OutTag
}

func (m *MockRepo) CompareURL(base, head string) string {
	if m.CompareURLIndex >= len(m.CompareURLMocks) {
		panic("CompareURL called more times than expected")
	}

	i := m.CompareURLIndex
	m.CompareURLIndex++

	m.CompareURLMocks[i].InBase = base
	m.CompareURLMocks[i].InHead = head
	return m.CompareURLMocks[i].OutString
}

func (m *MockRepo) EnsurePermissions(ctx context.Context) error {
	if m.EnsurePermissionsIndex >= len(m.EnsurePermissionsMocks) {
		panic("EnsurePermissions called more times than expected")
	}

	i := m.EnsurePermissionsIndex
	m.EnsurePermissionsIndex++

	m.EnsurePermissionsMocks[i].InContext = ctx
	return m.EnsurePermissionsMocks[i].OutError
}

func (m *MockRepo) FetchDefaultBranch(ctx context.Context) (service.Branch, error) {
	if m.FetchDefaultBranchIndex >= len(m.FetchDefaultBranchMocks) {
		panic("FetchDefaultBranch called more times than expected")
	}

	i := m.FetchDefaultBranchIndex
	m.FetchDefaultBranchIndex++

	m.FetchDefaultBranchMocks[i].InContext = ctx
	return m.FetchDefaultBranchMocks[i].OutBranch, m.FetchDefaultBranchMocks[i].OutError
}

func (m *MockRepo) FetchBranch(ctx context.Context, name string) (service.Branch, error) {
	if m.FetchBranchIndex >= len(m.FetchBranchMocks) {
		panic("FetchBranch called more times than expected")
	}

	i := m.FetchBranchIndex
	m.FetchBranchIndex++

	m.FetchBranchMocks[i].InContext = ctx
	m.FetchBranchMocks[i].InName = name
	return m.FetchBranchMocks[i].OutBranch, m.FetchBranchMocks[i].OutError
}

func (m *MockRepo) FetchTags(ctx context.Context) (service.Tags, error) {
	if m.FetchTagsIndex >= len(m.FetchTagsMocks) {
		panic("FetchTags called more times than expected")
	}

	i := m.FetchTagsIndex
	m.FetchTagsIndex++

	m.FetchTagsMocks[i].InContext = ctx
	return m.FetchTagsMocks[i].OutTags, m.FetchTagsMocks[i].OutError
}

func (m *MockRepo) FetchFirstCommit(ctx context.Context) (service.Commit, error) {
	if m.FetchFirstCommitIndex >= len(m.FetchFirstCommitMocks) {
		panic("FetchFirstCommit called more times than expected")
	}

	i := m.FetchFirstCommitIndex
	m.FetchFirstCommitIndex++

	m.FetchFirstCommitMocks[i].InContext = ctx
	return m.FetchFirstCommitMocks[i].OutCommit, m.FetchFirstCommitMocks[i].OutError
}

func (m *MockRepo) FetchParentCommits(ctx context.Context, hash string) (service.Commits, error) {
	if m.FetchParentCommitsIndex >= len(m.FetchParentCommitsMocks) {
		panic("FetchParentCommits called more times than expected")
	}

	i := m.FetchParentCommitsIndex
	m.FetchParentCommitsIndex++

	m.FetchParentCommitsMocks[i].InContext = ctx
	m.FetchParentCommitsMocks[i].InHash = hash
	return m.FetchParentCommitsMocks[i].OutCommits, m.FetchParentCommitsMocks[i].OutError
}

func (m *MockRepo) FetchIssuesAndMerges(ctx context.Context, since time.Time) (service.Issues, service.Merges, error) {
	if m.FetchIssuesAndMergesIndex >= len(m.FetchIssuesAndMergesMocks) {
		panic("FetchIssuesAndMerges called more times than expected")
	}

	i := m.FetchIssuesAndMergesIndex
	m.FetchIssuesAndMergesIndex++

	m.FetchIssuesAndMergesMocks[i].InContext = ctx
	m.FetchIssuesAndMergesMocks[i].InSince = since
	return m.FetchIssuesAndMergesMocks[i].OutIssues, m.FetchIssuesAndMergesMocks[i].OutMerges, m.FetchIssuesAndMergesMocks[i].OutError
}

type (
	ParseMock struct {
		InParseOptions changelog.ParseOptions
		OutChangelog   *changelog.Changelog
		OutError       error
	}

	RenderMock struct {
		InChangelog *changelog.Changelog
		OutContent  string
		OutError    error
	}

	MockChangelogProcessor struct {
		ParseIndex int
		ParseMocks []ParseMock

		RenderIndex int
		RenderMocks []RenderMock
	}
)

func (m *MockChangelogProcessor) Parse(opts changelog.ParseOptions) (*changelog.Changelog, error) {
	if m.ParseIndex >= len(m.ParseMocks) {
		panic("Parse called more times than expected")
	}

	i := m.ParseIndex
	m.ParseIndex++

	m.ParseMocks[i].InParseOptions = opts
	return m.ParseMocks[i].OutChangelog, m.ParseMocks[i].OutError
}

func (m *MockChangelogProcessor) Render(chlog *changelog.Changelog) (string, error) {
	if m.RenderIndex >= len(m.RenderMocks) {
		panic("Render called more times than expected")
	}

	i := m.RenderIndex
	m.RenderIndex++

	m.RenderMocks[i].InChangelog = chlog
	return m.RenderMocks[i].OutContent, m.RenderMocks[i].OutError
}
