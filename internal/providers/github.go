// Package providers exposes third party communication channels
package providers

import (
	"context"
	"strconv"

	"golang.org/x/oauth2"

	"github.com/google/go-github/v32/github"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"fuse/internal/domain"
)

// GitHub encapsulates github metadata and bridges communication to github provider
type GitHub struct {
	Owner       string
	Common      domain.CommonInput
	PullRequest domain.PullRequestInput
}

// GetRepository will fetch the repository details on github
func (gh *GitHub) GetRepository() (*ProviderRepository, error) {
	log.Info().
		Str("repo", gh.Common.RepositoryName).
		Str("pat", gh.Common.Pat).
		Msg("Getting github repository")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: gh.Common.Pat})
	client := github.NewClient(oauth2.NewClient(ctx, ts))

	repository, _, err := client.Repositories.Get(ctx, gh.Owner, gh.Common.RepositoryName)

	if err != nil {
		return nil, errors.WithStack(errors.Wrap(err, "Github error"))
	}

	return &ProviderRepository{
		WebURL: *repository.HTMLURL,
		Name:   *repository.Name,
	}, nil
}

// CreatePullRequest creates a pull request on github
func (gh *GitHub) CreatePullRequest(sourceBranch *string) (*ProviderPullRequest, error) {
	log.Info().
		Str("prTitle", gh.PullRequest.Title).
		Str("prTarget", TargetBranch).
		Str("prSource", *sourceBranch).
		Str("repoName", gh.Common.RepositoryName).
		Msg("Creating GitHub pull request")

	// pr unique identifier
	prID := uuid.Must(uuid.NewRandom()).String()
	prTitle := "Fuse Automated - " + prID

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: gh.Common.Pat})
	client := github.NewClient(oauth2.NewClient(ctx, ts))
	targetBranch := TargetBranch

	ghpr, _, err := client.PullRequests.Create(ctx, gh.Owner, gh.Common.RepositoryName, &github.NewPullRequest{
		Title: &prTitle,
		Head:  sourceBranch,
		Base:  &targetBranch,
	})

	if err != nil {
		return nil, errors.Wrap(err, "Github Error")
	}

	return &ProviderPullRequest{
		PullRequestID:  strconv.FormatInt(*ghpr.ID, 10),
		PullRequestURL: *ghpr.HTMLURL,
	}, nil
}

// GetCommonInput returns common inputs provided by the user via cli
func (gh *GitHub) GetCommonInput() *domain.CommonInput {
	return &gh.Common
}

// GetPullRequestInput returns pull request inputs provided by the user via cli
func (gh *GitHub) GetPullRequestInput() *domain.PullRequestInput {
	return &gh.PullRequest
}
