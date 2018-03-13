package site

import (
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/skratchdot/open-golang/open"
	"github.com/yosssi/ace"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var aceOpts = &ace.Options{
	DynamicReload: true,
	Indent:        "  ",
}

var watchChMu sync.Mutex
var watchCh chan struct{}

func subscribe() chan struct{} {
	watchChMu.Lock()
	defer watchChMu.Unlock()
	if watchCh == nil {
		watchCh = make(chan struct{})
	}
	return watchCh
}

func broadcast() {
	watchChMu.Lock()
	defer watchChMu.Unlock()
	close(watchCh)
	watchCh = nil
}

var lastModifiedMu sync.Mutex
var lastModified []os.FileInfo

func recordModified(fi os.FileInfo) {
	lastModifiedMu.Lock()
	defer lastModifiedMu.Unlock()

	replace := false
	replaceIdx := -1
	for i, f := range lastModified {
		if os.SameFile(f, fi) && f.ModTime().Before(fi.ModTime()) {
			replaceIdx = i
			replace = true
		}
	}
	if replace {
		lastModified[replaceIdx] = fi
	} else {
		lastModified = append(lastModified, fi)
	}
}

func isNewer(fi os.FileInfo) bool {
	lastModifiedMu.Lock()
	defer lastModifiedMu.Unlock()

	for _, f := range lastModified {
		if os.SameFile(f, fi) {
			if f.ModTime().Before(fi.ModTime()) {
				return true
			} else {
				return false
			}
		}
	}
	return true
}

func Build(output string, live bool) error {
	err := os.MkdirAll(output, os.FileMode(0755))
	if err != nil {
		panic(err)
	}

	if live {
		watch()
	}

	var wg sync.WaitGroup

	// // PDF
	// wg.Add(1)
	// go pdf(output, live, &wg)

	// HTML
	wg.Add(1)
	go html(output, live, &wg)

	if live {
		open.Run("output/index.html")
	}

	wg.Wait()
	return nil
}
