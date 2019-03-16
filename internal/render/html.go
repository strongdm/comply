package render

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/skratchdot/open-golang/open"
	"github.com/yosssi/ace"
)

const websocketReloader = `<script>
(function(){
	var ws = new WebSocket("ws://localhost:%d/ws")
	if (location.host != "") {
		ws = new WebSocket("ws://"+location.host+"/ws")
	}
	var connected = false
	ws.onopen = function(e) {
		connected = true
	}
	ws.onclose = function(e) {
		// reload!
		if (connected) {
			window.location.reload(true)
		}
	}
})()
</script>`

func html(output string, live bool, errCh chan error, wg *sync.WaitGroup) {
	opened := false

	for {
		files, err := ioutil.ReadDir(filepath.Join(".", "templates"))
		if err != nil {
			errCh <- errors.Wrap(err, "unable to open template directory")
			return
		}

		_, data, err := loadWithStats()
		if err != nil {
			errCh <- errors.Wrap(err, "unable to load data")
			return
		}

		for _, fileInfo := range files {
			if !strings.HasSuffix(fileInfo.Name(), ".ace") {
				continue
			}

			basename := strings.Replace(fileInfo.Name(), ".ace", "", -1)
			outputFilename := filepath.Join(output, fmt.Sprintf("%s.html", basename))

			w, err := os.Create(outputFilename)
			if err != nil {
				errCh <- errors.Wrap(err, "unable to create HTML file")
				return
			}

			fmt.Printf("%s -> %s\n", filepath.Join("templates", fileInfo.Name()), outputFilename)

			tpl, err := ace.Load("", filepath.Join("templates", basename), aceOpts)
			if err != nil {
				w.Write([]byte("<htmL><body>template error</body></html>"))
				fmt.Println(err)
			}

			err = tpl.Execute(w, data)
			if err != nil {
				w.Write([]byte("<htmL><body>template error</body></html>"))
				fmt.Println(err)
			}

			if live {
				w.Write([]byte(fmt.Sprintf(websocketReloader, ServePort)))
			}
			w.Close()
		}

		if live {
			if !opened {
				opened = true
				open.Run(fmt.Sprintf("http://127.0.0.1:%d/", ServePort))
			}
		} else {
			wg.Done()
			return
		}

		<-subscribe()
	}
}
