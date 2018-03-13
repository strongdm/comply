package site

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/strongdm/comply/internal/model"
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

			values := make(map[string]interface{})

			values["Title"] = "Acme Compliance Program"
			values["Procedures"] = []string{
				"Jump",
				"Sit",
				"Squat",
			}

			rt, err := model.DB().ReadAll("tickets")
			if err == nil {
				ts := model.Tickets(rt)
				var total, open, oldestDays int
				for _, t := range ts {
					total++
					if t.State == model.Open {
						if t.CreatedAt != nil {
							oldestDays = int(time.Since(*t.CreatedAt).Hours() / float64(24))
						}
						open++
					}

				}

				values["OldestDays"] = strconv.Itoa(oldestDays)
				values["Total"] = strconv.Itoa(total)
				values["Open"] = strconv.Itoa(open)
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
}
