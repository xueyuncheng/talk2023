package main

import (
	"log/slog"
	"net/http"
	"sync"
)

func main() {
	initRoute()
	http.ListenAndServe(":8080", nil)
}

func initRoute() {
	ch := make(chan string, 100)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		ch <- id

		if _, err := w.Write([]byte(id)); err != nil {
			slog.Error("w.Write() error", "err", err)
			return
		}
	})

	m := make(map[*Client]bool)
	l := &sync.Mutex{}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			slog.Error("upgrader.Upgrade() error", "err", err)
			return
		}

		client := newClient(conn)
		l.Lock()
		m[client] = true
		l.Unlock()
	})

	go func() {
		for {
			id := <-ch
			l.Lock()
			for c := range m {
				c.writeChan <- []byte(id)
			}
			l.Unlock()
		}
	}()
}
