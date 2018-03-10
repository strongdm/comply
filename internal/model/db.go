package model

import (
	"os"
	"path/filepath"

	"github.com/nanobox-io/golang-scribble"
	"github.com/strongdm/comply/internal/config"
)

// DB provides a singleton reference to a local json cache
func DB() *scribble.Driver {
	if _, err := os.Stat(filepath.Join(config.ProjectRoot(), ".comply", "cache")); os.IsNotExist(err) {
		err = os.Mkdir(filepath.Join(config.ProjectRoot(), ".comply"), os.FileMode(0755))
		if err != nil {
			panic(err)
		}
		err = os.Mkdir(filepath.Join(config.ProjectRoot(), ".comply", "cache"), os.FileMode(0755))
		if err != nil {
			panic(err)
		}
	}

	db, err := scribble.New(filepath.Join(config.ProjectRoot(), ".comply", "cache"), nil)
	if err != nil {
		panic("unable to load comply data: " + err.Error())
	}
	return db
}
