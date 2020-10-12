// Package providers exposes third party communication channels
package providers

import (
	"context"

	"fuse/internal/domain"

	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"github.com/microsoft/azure-devops-go-api/azuredevops/identity"
	"github.com/microsoft/azure-devops-go-api/azuredevops/webapi"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// AzDevOpsGetRepository will fetch the repository details on azure devops.
func AzDevOpsGetRepository(input *domain.AzureDevOpsInput) (*git.GitRepository, error) {
	connection := azuredevops.NewPatConnection(input.OrganizationURL, input.Common.Pat)

	ctx := context.Background()

	gitClient, err := git.NewClient(ctx, connection)

	if err != nil {
		return nil, errors.Wrap(err, "AzureDevOps error")
	}

	gitRepo, err := gitClient.GetRepository(ctx, git.GetRepositoryArgs{
		RepositoryId: &input.Common.RepositoryName,
		Project:      &input.ProjectName,
	})

	if err != nil {
		return nil, errors.WithStack(errors.Wrap(err, "AzureDevOps error"))
	}

	return gitRepo, nil
}

// AzDevOpsCreatePullRequest creates a pull request
func AzDevOpsCreatePullRequest(input *domain.AzureDevOpsInput, branchName string) (*git.GitPullRequest, error) {
	targetBranch := "master"

	log.Info().
		Str("prTitle", input.PullRequest.Title).
		Str("prTarget", targetBranch).
		Str("prSource", branchName).
		Str("repoName", input.Common.RepositoryName).
		Msg("Creating pull request")

	// pr unique identifier
	prID := uuid.Must(uuid.NewRandom()).String()
	ctx := context.Background()

	connection := azuredevops.NewPatConnection(input.OrganizationURL, input.Common.Pat)
	gitClient, err := git.NewClient(ctx, connection)

	if err != nil {
		return nil, errors.Wrap(err, "AzureDevOps error")
	}

	refPrefix := "refs/heads/"
	prTarget := refPrefix + targetBranch
	prSource := refPrefix + branchName

	prMetadata := git.GitPullRequest{
		ArtifactId:    &prID,
		SourceRefName: &prSource,
		TargetRefName: &prTarget,
		Title:         &input.PullRequest.Title,
	}

	// if auto complete is on, set the pr auto completion
	if input.PullRequest.AutoComplete {
		identityClient, err := identity.NewClient(ctx, connection)

		if err != nil {
			return nil, errors.Wrap(err, "AzureDevOps error")
		}

		self, err := identityClient.GetSelf(ctx, identity.GetSelfArgs{})

		if err != nil {
			return nil, errors.Wrap(err, "AzureDevOps error")
		}

		selfID := self.Id.String()
		prMetadata.AutoCompleteSetBy = &webapi.IdentityRef{
			Id: &selfID,
		}
	}

	return gitClient.CreatePullRequest(ctx, git.CreatePullRequestArgs{
		GitPullRequestToCreate: &prMetadata,
		RepositoryId:           &input.Common.RepositoryName,
		Project:                &input.ProjectName,
	})
}
