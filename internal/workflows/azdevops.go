// Package workflows contains the entry point to start the fuse process for the available provider implementation
package workflows

import (
	"fuse/internal/core"
	"fuse/internal/domain"
	"fuse/internal/providers"
	"os"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// AzDevOpsFuse kicks off the patching workflow for azdevops
func AzDevOpsFuse(input *domain.AzureDevOpsInput) error {
	log.Info().Msg("Azure DevOps workflow.")
	branchName := "master"

	if input.PullRequest.Enabled {
		branchName = uuid.Must(uuid.NewRandom()).String()
	}

	gitCloneRoot, err := layoutStage(input, branchName)

	if err != nil {
		return logErrAndReturn(err)
	}

	// I'm ok if this errors and the folder is not removed. If this ends up not ok, return this error
	defer os.RemoveAll(*gitCloneRoot)

	// start the crawling and diffing process
	diffsChannel, err := core.Crawl(input.Common.ContentDir, *gitCloneRoot, input.Common.Concurrency)

	if err != nil {
		return logErrAndReturn(err)
	}

	diffs := <-diffsChannel

	// only proceed with pushing changes we have any and we didn't find any error
	if diffs.Error == 0 && diffs.WithDiffs > 0 {
		err = providers.GitCommit(*gitCloneRoot)

		if err != nil {
			return logErrAndReturn(err)
		}

		err = providers.GitPush(*gitCloneRoot, branchName)

		if err != nil {
			return logErrAndReturn(err)
		}

		// create the associated pull request if fuse was configured to do so
		if input.PullRequest.Enabled {
			pr, err := providers.AzDevOpsCreatePullRequest(input, branchName)

			if err != nil {
				return logErrAndReturn(err)
			}
			log.Info().
				Uint32("total_diffs", diffs.WithDiffs).
				Int("pullRequestId", *pr.PullRequestId).
				Str("artifactId", *pr.ArtifactId).
				Msg("Pull request created")
		}

		log.Info().
			Uint32("total_diffs", diffs.WithDiffs).
			Msg("Pushed changes to master")

		return nil
	} else if diffs.Error > 0 {
		log.Error().
			Uint32("total_errors", diffs.Error).
			Msg("Cannot proceed due to errors.")

		return errors.New("crawling process encountered errors")
	}

	log.Info().
		Msg("No changes to commit")

	return nil
}

func logErrAndReturn(err error) error {
	log.Error().
		Stack().
		Err(err).
		Send()
	return err
}

func layoutStage(input *domain.AzureDevOpsInput, branchName string) (*string, error) {
	defaultWorkingDir := "/tmp"
	gitRepo, err := providers.AzDevOpsGetRepository(input)

	if err != nil {
		return nil, err
	}

	err = providers.ConfigureGit(&defaultWorkingDir)

	if err != nil {
		return nil, err
	}

	_, gitCloneRoot, err := providers.GitClone(*gitRepo.WebUrl, *gitRepo.Name, input.Common.Pat)

	if err != nil {
		return nil, err
	}

	// Only create the branch if we're told to use a branch
	if input.PullRequest.Enabled && branchName != "" {
		err = providers.CreateGitBranch(gitCloneRoot, branchName)

		if err != nil {
			return nil, err
		}
	}

	return &gitCloneRoot, nil
}
