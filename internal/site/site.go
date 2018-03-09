package site

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gohugoio/hugo/watcher"
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
		b, err := watcher.New(300 * time.Millisecond)
		if err != nil {
			panic(err)
		}
		b.Add("./templates/")
		b.Add("./policies/")
		b.Add("./procedures/")

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

	var wg sync.WaitGroup

	// PDF
	wg.Add(1)
	go func() {
		for {
			files, err := ioutil.ReadDir("./policies")
			if err != nil {
				panic(err)
			}

			for _, fileInfo := range files {
				// only non-README markdown files
				if !strings.HasSuffix(fileInfo.Name(), ".md") || strings.HasPrefix("README", strings.ToUpper(fileInfo.Name())) {
					continue
				}

				// only files that have been touched
				if !isNewer(fileInfo) {
					continue
				}
				recordModified(fileInfo)

				basename := strings.Replace(fileInfo.Name(), ".md", "", -1)

				ctx := context.Background()
				cli, err := client.NewEnvClient()
				if err != nil {
					panic(err)
				}

				_, err = cli.ImagePull(ctx, "jagregory/pandoc", types.ImagePullOptions{})
				if err != nil {
					panic(err)
				}

				pwd, err := os.Getwd()
				if err != nil {
					panic(err)
				}

				hc := &container.HostConfig{
					Binds: []string{pwd + ":/source"},
				}

				resp, err := cli.ContainerCreate(ctx, &container.Config{
					Image: "jagregory/pandoc",
					Cmd: []string{"--smart", "--toc", "-N", "--template=/source/templates/default.latex", "-o",
						fmt.Sprintf("/source/output/%s.pdf", basename),
						fmt.Sprintf("/source/policies/%s.md", basename),
					},
				}, hc, nil, "")
				if err != nil {
					panic(err)
				}

				if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
					panic(err)
				}

				cli.ContainerWait(ctx, resp.ID)

				_, err = cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
				if err != nil {
					panic(err)
				}

				// io.Copy(os.Stdout, out)
			}

			if !live {
				wg.Done()
				return
			}
			<-subscribe()
		}
	}()

	// HTML
	wg.Add(1)
	go func() {
		for {
			files, err := ioutil.ReadDir("./templates")
			if err != nil {
				panic(err)
			}
			for _, fileInfo := range files {
				if !strings.HasSuffix(fileInfo.Name(), ".ace") {
					continue
				}

				// only files that have been touched
				if !isNewer(fileInfo) {
					continue
				}
				recordModified(fileInfo)

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

				tpl, err := ace.Load("", filepath.Join("templates", basename), aceOpts)
				if err != nil {
					w.Write([]byte("<htmL><body>template error</body></html>"))
					fmt.Println(err)
				}

				err = tpl.Execute(w, values)
				if err != nil {
					w.Write([]byte("<htmL><body>template error</body></html>"))
					fmt.Println(err)
				}

				if live {
					w.Write([]byte(`<script>
			(function(){
				var ws = new WebSocket("ws://localhost:5122/ws")
				ws.onclose = function(e) {
					// reload!
					window.location=window.location
				}
			})()
			</script>`))
				}
				w.Close()
			}
			if !live {
				wg.Done()
				return
			}
			<-subscribe()
		}
	}()

	if live {
		open.Run("output/index.html")
	}

	wg.Wait()
	return nil
}

// var rootPathMu sync.Mutex
// var rootPath os.FileInfo

// func findRoot() {

// }
