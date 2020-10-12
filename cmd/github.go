// Package cmd is the entry point for cobra cli
package cmd

import (
	"fuse/internal/domain"
	"fuse/internal/providers"
	"fuse/internal/workflow"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	owner string
)

var githubCmd = &cobra.Command{
	Use:   "github",
	Short: "Fuse for Github",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		ConfigLog()

		cmd.Flags().Visit(func(flag *pflag.Flag) {
			if flag.Value.String() == "" {
				panic("Found empty flag: " + flag.Name)
			}
		})
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		input := providers.GitHub{
			Common: domain.CommonInput{
				RepositoryName: repoName,
				Pat:            pat,
				ContentDir:     contentDir,
				Concurrency:    concurrency,
			},
			PullRequest: domain.PullRequestInput{
				Title:        prTitle,
				AutoComplete: prAutoComplete,
				Enabled:      prEnabled,
			},
			Owner: owner,
		}

		return workflow.AzDevOpsFuse(&input)
	},
}

func init() {
	githubCmd.Flags().StringVarP(&owner, "owner", "o", "", "Github owner")

	_ = githubCmd.MarkFlagRequired("owner")

	rootCmd.AddCommand(githubCmd)
}
