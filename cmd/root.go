// Package cmd is the entry point for cobra cli
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	repoName         string
	pat              string
	tag              string
	contentDir       string
	commentDelimiter string
	concurrency      int8

	prettyLogging  bool
	logStackTraces bool
	logVerbose     bool

	prEnabled      bool
	prTitle        string
	prAutoComplete bool

	rootCmd = &cobra.Command{
		Use:   "fuse",
		Short: "Fuse.",
	}
)

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&repoName, "repoName", "r", "",
		"Remote repository name to be used as fuse target.")
	rootCmd.PersistentFlags().StringVarP(&pat, "pat", "a", "",
		"Personal access token to authenticate when performing actions.")
	rootCmd.PersistentFlags().StringVarP(&contentDir, "contentDir", "d", "",
		`Path to the directory that contains the content to be used by fuse in the target repository. 
				The path may be relative to the current execution process directory (where you execute fuse) or an absolute path.`)
	rootCmd.PersistentFlags().StringVarP(&tag, "tag", "v", "latest",
		"Git tag to be used when committing to master.")
	rootCmd.PersistentFlags().StringVarP(&commentDelimiter, "commentDelimiter", "h", "#",
		"Comment delimiter used for Fuse file mark.")

	rootCmd.PersistentFlags().BoolVarP(&prEnabled, "prEnabled", "i", false,
		"If enabled, fuse will work in a new branch and create the associated pull request with the changes.")
	rootCmd.PersistentFlags().StringVarP(&prTitle, "prTitle", "t", "Fuse Automated",
		"Pull request title.")
	rootCmd.PersistentFlags().BoolVarP(&prAutoComplete, "prAutocomplete", "y", false,
		"If enabled pull request will be auto completed. This depends on provider support.")

	rootCmd.PersistentFlags().Int8VarP(&concurrency, "concurrency", "c", 10,
		"Max concurrency allowed to process work items. Each work item represents a file to be patched.")
	rootCmd.PersistentFlags().BoolVarP(&prettyLogging, "pretty", "b", true,
		"Pretty print logging. If set to false json format will be used.")
	rootCmd.PersistentFlags().BoolVarP(&logStackTraces, "showStacks", "s", true,
		"If enabled and errors occur, stack traces will be shown.")
	rootCmd.PersistentFlags().BoolVarP(&logVerbose, "verbose", "l", false,
		"If enabled debug log level will be used.")

	_ = rootCmd.MarkFlagRequired("repoName")
	_ = rootCmd.MarkFlagRequired("pat")
	_ = rootCmd.MarkFlagRequired("contentDir")
}
