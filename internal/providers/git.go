// Package providers exposes third party communication channels
package providers

import (
	"fmt"
	"io/ioutil"
	"strings"

	"fuse/internal/process"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// GitClone will perform an os.Exec git clone to a temporary working directory.
// The caller is responsible to clean the dir when no longer needed.
// Returned string destination is only nil in case it wasn't possible to create the temporary directory.
// The repositoryURL should be in the form of https://remote_repository_web_url and when cloning
// basic authentication will be used, e.g: https://automated:pat_token@remote_url
func GitClone(repositoryURL, repositoryName, token string) (destination, gitCloneRoot string, err error) {
	// dev note: I did not use go-git because it has issues with azure devops. Check the following issues:
	// https://github.com/src-d/go-git/issues/335
	// https://github.com/src-d/go-git/issues/1058

	// this will yield https://automated:pat_token@remote_repository_web_url_without_https://
	templatedURL := fmt.Sprintf("https://automated:%s@%s", token, strings.Replace(repositoryURL, "https://", "", -1))

	log.Info().
		Str("repository", repositoryURL).
		Msg("Cloning repository.")

	// create temporary folder to clone the repo
	destination, err = ioutil.TempDir("/tmp", "tmp")
	if err != nil {
		return "", "", errors.Wrap(err, "Git error")
	}

	log.Debug().
		Str("command", strings.Join([]string{"git", "clone", repositoryURL, destination}, " ")).
		Send()

	_, stderr, err := process.ExecuteProcess(strings.Join([]string{"git", "clone", templatedURL}, " "), &destination)

	if err != nil {
		log.Error().
			Msg(stderr)
		return destination, "", errors.Wrap(err, "Git error")
	}

	gitCloneRoot = destination + "/" + repositoryName

	log.Info().
		Msg("Successfully cloned to " + destination)

	return destination, gitCloneRoot, nil
}

// CreateGitBranch will create a new branch with the given name
func CreateGitBranch(repositoryDir, branchName string) error {
	log.Info().
		Str("branchName", branchName).
		Msg("Creating git branch.")

	log.Debug().
		Str("command", strings.Join([]string{"git", "checkout", "-b", branchName}, " ")).
		Send()

	_, stderr, err := process.ExecuteProcess(strings.Join([]string{"git", "checkout", "-b", branchName}, " "), &repositoryDir)

	if err != nil {
		log.Error().
			Msg(stderr)
		return errors.Wrap(err, "Git error")
	}

	log.Info().Msg("Successfully created branch " + branchName)

	return nil
}

// GitCommit will add and commit any changes in the provided repository Dir
func GitCommit(repositoryDir string) error {
	log.Info().
		Msg("Committing git changes.")

	// add
	log.Debug().
		Str("command", strings.Join([]string{"git", "add", "-A"}, " ")).
		Send()

	_, stderr, err := process.ExecuteProcess(strings.Join([]string{"git", "add", "-A"}, " "), &repositoryDir)

	if err != nil {
		log.Error().
			Msg(stderr)
		return errors.Wrap(err, "Git error")
	}

	// commit
	log.Debug().
		Str("command", strings.Join([]string{"git", "commit", "-m", "Fuse automation"}, " ")).
		Send()

	_, stderr, err = process.ExecuteProcess(strings.Join([]string{"git", "commit", "-m",
		"'Fuse automation'"}, " "), &repositoryDir)

	if err != nil {
		log.Error().
			Msg(stderr)
		return errors.Wrap(err, "Git error")
	}

	log.Info().Msg("Successfully committed changes.")

	return nil
}

// GitPush will push any changes in the provided repository directory to the remote branch
func GitPush(repositoryDir, branch string) error {
	log.Info().
		Msg("Pushing git changes.")

	// add
	log.Debug().
		Str("command", strings.Join([]string{"git", "push", "-u", "origin", branch}, " ")).
		Send()

	_, stderr, err := process.ExecuteProcess(strings.Join([]string{"git", "push", "--follow-tags", "-u", "origin", branch}, " "),
		&repositoryDir)

	if err != nil {
		log.Error().
			Msg(stderr)
		return errors.Wrap(err, "Git error")
	}

	log.Info().Msg("Successfully pushed changes.")

	return nil
}

// GitTag creates a annotated git tag
func GitTag(repositoryDir, tag string) error {
	log.Info().
		Msg("Tagging git commit")

	_, stderr, err := process.ExecuteProcess(strings.Join([]string{"git", "tag", "-m", "Fuse release " + tag, tag}, " "), &repositoryDir)

	if err != nil {
		log.Error().
			Msg(stderr)
		return errors.Wrap(err, "Git error")
	}

	log.Info().Msg("Successfully tagged.")

	return nil
}

// ConfigureGit configures git user and email
func ConfigureGit(workingDir *string) error {
	log.Info().
		Str("user", "Fuse").
		Str("email", "fuse@dev.io").
		Msg("Configuring git user and email.")

	_, stderr, err := process.ExecuteProcess(strings.Join([]string{"git", "config", "--global", "user.name", "Fuse"}, " "), workingDir)

	if err != nil {
		log.Error().
			Msg(stderr)
		return errors.Wrap(err, "Git error")
	}

	_, stderr, err = process.ExecuteProcess(strings.Join([]string{"git", "config", "--global", "user.email", "fuse@dev.io"}, " "), workingDir)

	if err != nil {
		log.Error().
			Msg(stderr)
		return errors.Wrap(err, "Git error")
	}

	return nil
}
