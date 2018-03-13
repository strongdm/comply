package site

import (
	"net/http"
	"time"

	"github.com/gohugoio/hugo/watcher"
)

func watch() {
	b, err := watcher.New(300 * time.Millisecond)
	if err != nil {
		panic(err)
	}
	b.Add("./templates/")
	b.Add("./policies/")
	b.Add("./procedures/")

	b.Add("./.comply/")
	b.Add("./.comply/cache")
	b.Add("./.comply/cache/tickets")

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
