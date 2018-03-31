package render

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/yosssi/ace"
)

func html(output string, live bool, wg *sync.WaitGroup) {
	for {
		files, err := ioutil.ReadDir("./templates")
		if err != nil {
			panic(err)
		}
		for _, fileInfo := range files {
			if !strings.HasSuffix(fileInfo.Name(), ".ace") {
				continue
			}

			basename := strings.Replace(fileInfo.Name(), ".ace", "", -1)
			w, err := os.Create(filepath.Join(output, fmt.Sprintf("%s.html", basename)))
			if err != nil {
				panic(err)
			}

			tpl, err := ace.Load("", filepath.Join("templates", basename), aceOpts)
			if err != nil {
				w.Write([]byte("<htmL><body>template error</body></html>"))
				fmt.Println(err)
			}

			err = tpl.Execute(w, loadValues())
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
}
