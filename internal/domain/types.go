// Package domain contains domain types
package domain

// CommonInput are cli inputs that are common to all providers
type CommonInput struct {
	RepositoryName string
	Tag            string
	Pat            string
	ContentDir     string
	Concurrency    int8
}

// PullRequestInput are cli inputs related to pull requests
type PullRequestInput struct {
	Title        string
	Enabled      bool
	AutoComplete bool
}
