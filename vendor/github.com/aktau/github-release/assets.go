package main

import (
	"time"
)

const (
	ASSET_DOWNLOAD_URI = "/repos/%s/%s/releases/assets/%d"
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

// findAssetID returns the asset ID if name can be found in assets,
// otherwise returns -1.
func findAssetID(assets []Asset, name string) int {
	for _, asset := range assets {
		if asset.Name == name {
			return asset.Id
		}
	}
	return -1
}
