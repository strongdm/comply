package main

import (
	"fmt"

	"github.com/aktau/github-release/github"
)

const (
	TAGS_URI = "/repos/%s/%s/tags"
)

type Tag struct {
	Name       string `json:"name"`
	Commit     Commit `json:"commit"`
	ZipBallUrl string `json:"zipball_url"`
	TarBallUrl string `json:"tarball_url"`
}

func (t *Tag) String() string {
	return t.Name + " (commit: " + t.Commit.Url + ")"
}

// Get the tags associated with a repo.
func Tags(user, repo, token string) ([]Tag, error) {
	var tags []Tag
	return tags, github.Client{Token: token, BaseURL: EnvApiEndpoint}.Get(fmt.Sprintf(TAGS_URI, user, repo), &tags)
}
