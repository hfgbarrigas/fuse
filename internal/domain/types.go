// Package domain contains domain types
package domain

// CommonInputs are cli inputs that are common to all providers
type CommonInputs struct {
	RepositoryName string
	Pat            string
	ContentDir     string
	Concurrency    int8
}

// PullRequestMetadata are cli inputs related to pull requests
type PullRequestMetadata struct {
	Title        string
	Enabled      bool
	AutoComplete bool
}

// AzureDevOpsInput encapsulates user provided input
type AzureDevOpsInput struct {
	OrganizationURL string
	ProjectName     string
	Common          CommonInputs
	PullRequest     PullRequestMetadata
}
