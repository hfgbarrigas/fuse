// Package workflow contains the entry point to start the fuse process for the available provider implementation
package workflow

import (
	"fuse/internal/core"
	"fuse/internal/providers"
	"os"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Fuse kicks off the patching workflow
func Fuse(provider providers.Provider) error {
	log.Info().Msg("Fusing")
	branchName := providers.TargetBranch

	if provider.GetPullRequestInput().Enabled {
		branchName = uuid.Must(uuid.NewRandom()).String()
	}

	gitCloneRoot, err := layoutStage(provider, branchName)

	if err != nil {
		return logErrAndReturn(err)
	}

	// I'm ok if this errors and the folder is not removed. If this ends up not ok, return this error
	defer os.RemoveAll(*gitCloneRoot)

	// start the crawling and diffing process
	diffsChannel, err := core.Crawl(provider.GetCommonInput().ContentDir, *gitCloneRoot,
		provider.GetCommonInput().CommentDelimiter, provider.GetCommonInput().Concurrency)

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

		if providers.TargetBranch == branchName {
			err = providers.GitTag(*gitCloneRoot, provider.GetCommonInput().Tag)

			if err != nil {
				return logErrAndReturn(err)
			}
		}

		err = providers.GitPush(*gitCloneRoot, branchName)

		if err != nil {
			return logErrAndReturn(err)
		}

		// create the associated pull request if fuse was configured to do so
		if provider.GetPullRequestInput().Enabled {
			pr, err := provider.CreatePullRequest(&branchName)

			if err != nil {
				return logErrAndReturn(err)
			}
			log.Info().
				Uint32("totalDiffs", diffs.WithDiffs).
				Str("pullRequestID", pr.PullRequestID).
				Str("pullRequestURL", pr.PullRequestURL).
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

func layoutStage(provider providers.Provider, branchName string) (*string, error) {
	defaultWorkingDir := "/tmp"
	gitRepo, err := provider.GetRepository()

	if err != nil {
		return nil, err
	}

	err = providers.ConfigureGit(&defaultWorkingDir)

	if err != nil {
		return nil, err
	}

	_, gitCloneRoot, err := providers.GitClone(gitRepo.WebURL, gitRepo.Name, provider.GetCommonInput().Pat)

	if err != nil {
		return nil, err
	}

	// Only create the branch if we're told to use a branch
	if provider.GetPullRequestInput().Enabled && branchName != "" {
		err = providers.CreateGitBranch(gitCloneRoot, branchName)

		if err != nil {
			return nil, err
		}
	}

	return &gitCloneRoot, nil
}
