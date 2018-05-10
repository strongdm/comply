package model

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/nanobox-io/golang-scribble"
	"github.com/strongdm/comply/internal/config"
)

var dbSingletonOnce sync.Once
var dbSingleton *scribble.Driver

// DB provides a singleton reference to a local json cache; will panic if storage location is not writeable.
func DB() *scribble.Driver {
	dbSingletonOnce.Do(func() {
		if _, err := os.Stat(filepath.Join(config.ProjectRoot(), ".comply", "cache")); os.IsNotExist(err) {
			err = os.Mkdir(filepath.Join(config.ProjectRoot(), ".comply"), os.FileMode(0755))
			if err != nil {
				panic("could not create directory .comply: " + err.Error())
			}
			err = os.Mkdir(filepath.Join(config.ProjectRoot(), ".comply", "cache"), os.FileMode(0755))
			if err != nil {
				panic("could not create directory .comply/cache: " + err.Error())
			}
		}

		db, err := scribble.New(filepath.Join(config.ProjectRoot(), ".comply", "cache"), nil)
		if err != nil {
			panic("unable to load comply data: " + err.Error())
		}
		dbSingleton = db
	})
	return dbSingleton
}
