package render

import (
	"net/http"
	"os"
	"sync"
	"time"

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
var lastModified = make(map[string]time.Time)

func recordModified(path string, t time.Time) {
	lastModifiedMu.Lock()
	defer lastModifiedMu.Unlock()

	previous, ok := lastModified[path]
	if !ok || t.After(previous) {
		lastModified[path] = t
	}
}

func isNewer(path string, t time.Time) bool {
	lastModifiedMu.Lock()
	defer lastModifiedMu.Unlock()

	previous, ok := lastModified[path]
	if !ok {
		return true
	}

	// is tested after previous? Then isNewer is true.
	return t.After(previous)
}

func Build(output string, live bool) error {
	err := os.RemoveAll(output)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(output, os.FileMode(0755))
	if err != nil {
		panic(err)
	}

	if live {
		watch()
	}

	var wg sync.WaitGroup

	// PDF
	wg.Add(1)
	go pdf(output, live, &wg)

	// HTML
	wg.Add(1)
	go html(output, live, &wg)

	if live {
		open.Run("output/index.html")
	}

	wg.Wait()
	return nil
}
