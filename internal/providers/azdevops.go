// Package providers exposes third party communication channels
package providers

import (
	"context"
	"strconv"

	"fuse/internal/domain"

	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"github.com/microsoft/azure-devops-go-api/azuredevops/identity"
	"github.com/microsoft/azure-devops-go-api/azuredevops/webapi"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// AzureDevOps encapsulates azure devops metadata and bridges communication to az devops provider
type AzureDevOps struct {
	OrganizationURL string
	ProjectName     string
	Common          domain.CommonInput
	PullRequest     domain.PullRequestInput
}

// GetRepository will fetch the repository details on azure devops
func (az *AzureDevOps) GetRepository() (*ProviderRepository, error) {
	log.Info().
		Str("repoName", az.Common.RepositoryName).
		Msg("Getting repository from azure devops")

	connection := azuredevops.NewPatConnection(az.OrganizationURL, az.Common.Pat)

	ctx := context.Background()

	gitClient, err := git.NewClient(ctx, connection)

	if err != nil {
		return nil, errors.Wrap(err, "AzureDevOps error")
	}

	gitRepo, err := gitClient.GetRepository(ctx, git.GetRepositoryArgs{
		RepositoryId: &az.Common.RepositoryName,
		Project:      &az.ProjectName,
	})

	if err != nil {
		return nil, errors.WithStack(errors.Wrap(err, "AzureDevOps error"))
	}

	return &ProviderRepository{
		WebURL: *gitRepo.WebUrl,
		Name:   *gitRepo.Name,
	}, nil
}

// CreatePullRequest creates a pull request on azure devops
func (az *AzureDevOps) CreatePullRequest(sourceBranch *string) (*ProviderPullRequest, error) {
	log.Info().
		Str("prTitle", az.PullRequest.Title).
		Str("prTarget", TargetBranch).
		Str("prSource", *sourceBranch).
		Str("repoName", az.Common.RepositoryName).
		Msg("Creating GitHub pull request")

	// pr unique identifier
	prID := uuid.Must(uuid.NewRandom()).String()
	ctx := context.Background()

	connection := azuredevops.NewPatConnection(az.OrganizationURL, az.Common.Pat)
	gitClient, err := git.NewClient(ctx, connection)

	if err != nil {
		return nil, errors.Wrap(err, "AzureDevOps error")
	}

	refPrefix := "refs/heads/"
	prTarget := refPrefix + TargetBranch
	prSource := refPrefix + *sourceBranch

	prMetadata := git.GitPullRequest{
		ArtifactId:    &prID,
		SourceRefName: &prSource,
		TargetRefName: &prTarget,
		Title:         &az.PullRequest.Title,
	}

	// if auto complete is on, set the pr auto completion
	if az.PullRequest.AutoComplete {
		identityClient, err2 := identity.NewClient(ctx, connection)

		if err2 != nil {
			return nil, errors.Wrap(err2, "AzureDevOps error")
		}

		self, err2 := identityClient.GetSelf(ctx, identity.GetSelfArgs{})

		if err2 != nil {
			return nil, errors.Wrap(err, "AzureDevOps error")
		}

		selfID := self.Id.String()
		prMetadata.AutoCompleteSetBy = &webapi.IdentityRef{
			Id: &selfID,
		}
	}

	azPr, err := gitClient.CreatePullRequest(ctx, git.CreatePullRequestArgs{
		GitPullRequestToCreate: &prMetadata,
		RepositoryId:           &az.Common.RepositoryName,
		Project:                &az.ProjectName,
	})

	if err != nil {
		return nil, errors.Wrap(err, "AzureDevOps error")
	}

	return &ProviderPullRequest{
		PullRequestID:  strconv.Itoa(*azPr.PullRequestId),
		PullRequestURL: *azPr.RemoteUrl,
	}, nil
}

// GetCommonInput returns common inputs provided by the user via cli
func (az *AzureDevOps) GetCommonInput() *domain.CommonInput {
	return &az.Common
}

// GetPullRequestInput returns pull request inputs provided by the user via cli
func (az *AzureDevOps) GetPullRequestInput() *domain.PullRequestInput {
	return &az.PullRequest
}
