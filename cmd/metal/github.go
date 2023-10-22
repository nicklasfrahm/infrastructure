package main

import (
	"github.com/google/go-github/v56/github"
)

var (
	client *github.Client
)

// NewGitHubClient returns a pointer to a configured GitHub client.
func NewGitHubClient() *github.Client {
	if client != nil {
		return client
	}

	client = github.NewClient(nil)

	return client
}
