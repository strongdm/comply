package model

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/strongdm/comply/internal/config"
)

func TestSaveGet(t *testing.T) {
	dir := os.TempDir()
	config.SetProjectRoot(dir)
	f, err := os.Create(filepath.Join(dir, "config.yml"))
	if err != nil {
		panic(err)
	}
	f.Close()

	name := "Do something excellent"
	err = DB().Write("tickets", "100", &Ticket{ID: "100", Name: name})
	if err != nil {
		panic(err)
	}

	ticket := &Ticket{}
	err = DB().Read("tickets", "100", ticket)
	if err != nil {
		panic(err)
	}

	if ticket.Name != name {
		t.Error("failed to read ticket")
	}
}
