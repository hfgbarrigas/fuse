// Package providers exposes third party communication channels
package providers

import (
	"fuse/internal/domain"
)

// Provider defines the necessary
type Provider interface {
	GetRepository() (*ProviderRepository, error)
	CreatePullRequest(sourceBranch *string) (*ProviderPullRequest, error)
	GetCommonInput() *domain.CommonInput
	GetPullRequestInput() *domain.PullRequestInput
}

// ProviderRepository encapsulates data about a provider repository
type ProviderRepository struct {
	WebURL string
	Name   string
}

// ProviderPullRequest encapsulates data about a provider pull request
type ProviderPullRequest struct {
	PullRequestID  string
	PullRequestURL string
}

// TargetBranch defines the default target branch used on provider pull requests
const TargetBranch = "master"
