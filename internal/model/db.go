package model

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/nanobox-io/golang-scribble"
	"github.com/strongdm/comply/internal/config"
)

var dbSingletonOnce sync.Once
var dbSingleton *scribble.Driver

// DB provides a singleton reference to a local json cache
func DB() *scribble.Driver {
	dbSingletonOnce.Do(func() {
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
		dbSingleton = db
	})
	return dbSingleton
}

func Tickets(rawTickets []string) []*Ticket {
	var tickets []*Ticket
	for _, rt := range rawTickets {
		t := &Ticket{}
		err := json.Unmarshal([]byte(rt), t)
		if err != nil {
			panic("Malformed ticket JSON")
		}
		tickets = append(tickets, t)
	}
	return tickets
}
