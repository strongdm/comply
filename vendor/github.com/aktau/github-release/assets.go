package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/github-release/github-release/github"
)

const (
	// GET /repos/:owner/:repo/releases/assets/:id
	// DELETE /repos/:owner/:repo/releases/assets/:id
	ASSET_URI = "/repos/%s/%s/releases/assets/%d"

	// API: https://developer.github.com/v3/repos/releases/#list-assets-for-a-release
	// GET /repos/:owner/:repo/releases/:id/assets
	ASSET_RELEASE_LIST_URI = "/repos/%s/%s/releases/%d/assets"
)

type Asset struct {
	Url         string    `json:"url"`
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	ContentType string    `json:"content_type"`
	State       string    `json:"state"`
	Size        uint64    `json:"size"`
	Downloads   uint64    `json:"download_count"`
	Created     time.Time `json:"created_at"`
	Published   time.Time `json:"published_at"`
}

// findAsset returns the asset if an asset with name can be found in assets,
// otherwise returns nil.
func findAsset(assets []Asset, name string) *Asset {
	for _, asset := range assets {
		if asset.Name == name {
			return &asset
		}
	}
	return nil
}

// Delete sends a HTTP DELETE request for the given asset to Github. Returns
// nil if the asset was deleted OR there was nothing to delete.
func (a *Asset) Delete(user, repo, token string) error {
	URL := nvls(EnvApiEndpoint, github.DefaultBaseURL) +
		fmt.Sprintf(ASSET_URI, user, repo, a.Id)
	resp, err := github.DoAuthRequest("DELETE", URL, "application/json", token, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to delete asset %s (ID: %d), HTTP error: %b", a.Name, a.Id, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete asset %s (ID: %d), status: %s", a.Name, a.Id, resp.Status)
	}
	return nil
}
