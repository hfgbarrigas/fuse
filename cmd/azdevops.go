// Package cmd is the entry point for cobra cli
package cmd

import (
	"fuse/internal/domain"
	"fuse/internal/workflows"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	organizationURL string
	projectName     string
)

var azdevopsCmd = &cobra.Command{
	Use:   "azdevops",
	Short: "Fuse for Azure DevOps",
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

		input := domain.AzureDevOpsInput{
			Common: domain.CommonInputs{
				RepositoryName: repoName,
				Pat:            pat,
				ContentDir:     contentDir,
				Concurrency:    concurrency,
			},
			PullRequest: domain.PullRequestMetadata{
				Title:        prTitle,
				AutoComplete: prAutoComplete,
				Enabled:      prEnabled,
			},
			OrganizationURL: organizationURL,
			ProjectName:     projectName,
		}

		return workflows.AzDevOpsFuse(&input)
	},
}

func init() {
	azdevopsCmd.Flags().StringVarP(&organizationURL, "orgUrl", "u", "", "Azure DevOps organization url.")
	azdevopsCmd.Flags().StringVarP(&projectName, "project", "p", "", "Azure DevOps project name.")

	_ = azdevopsCmd.MarkFlagRequired("organizationUrl")
	_ = azdevopsCmd.MarkFlagRequired("project")

	rootCmd.AddCommand(azdevopsCmd)
}
