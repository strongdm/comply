package site

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gohugoio/hugo/watcher"
	"github.com/gorilla/websocket"
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

// Asset: func(name string) ([]byte, error) {
// 	switch name {
// 	case "index.ace":
// 		return []byte(indexACE), nil
// 	}
// 	return nil, errors.New("template not found: " + name)
// }}

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

func Build(input, output string, live bool) error {
	err := os.MkdirAll(output, os.FileMode(0755))
	if err != nil {
		panic(err)
	}

	if live {
		b, err := watcher.New(300 * time.Millisecond)
		if err != nil {
			panic(err)
		}
		b.Add(input)
		b.Add("./site/")

		go func() {
			for {
				select {
				case e := <-b.Errors:
					panic(e)
				case <-b.Events:
					broadcast()
				}
			}
		}()

		serveWs := func(w http.ResponseWriter, r *http.Request) {
			ws, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				return
			}
			<-subscribe()
			time.Sleep(500 * time.Millisecond)
			ws.Close()
		}

		http.HandleFunc("/ws", serveWs)
		go http.ListenAndServe("127.0.0.1:5122", nil)
	}

	files, err := ioutil.ReadDir("./site")
	if err != nil {
		panic(err)
	}

	for {
		for _, fileInfo := range files {
			if !strings.HasSuffix(fileInfo.Name(), ".ace") {
				continue
			}

			basename := strings.Replace(fileInfo.Name(), ".ace", "", -1)

			w, err := os.Create(filepath.Join(output, fmt.Sprintf("%s.html", basename)))
			if err != nil {
				panic(err)
			}

			values := make(map[string]interface{})

			values["Title"] = "Acme Compliance Program"
			values["Procedures"] = []string{
				"Jump",
				"Sit",
				"Squat",
			}

			tpl, err := ace.Load("", filepath.Join("site", basename), aceOpts)
			if err != nil {
				w.Write([]byte("<htmL><body>template error</body></html>"))
				fmt.Println(err)
			}

			err = tpl.Execute(w, values)
			if err != nil {
				w.Write([]byte("<htmL><body>template error</body></html>"))
				fmt.Println(err)
			}

			w.Write([]byte(`<script>
			(function(){
				var ws = new WebSocket("ws://localhost:5122/ws")
				ws.onclose = function(e) {
					// reload!
					window.location=window.location
				}
			})()
			</script>`))
			w.Close()
		}
		<-subscribe()
	}

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
	return nil
}
