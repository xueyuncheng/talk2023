package main

import (
	"log/slog"
	"net/http"
	"os"
	"sync"
)

func main() {
	initRoute()
	if err := http.ListenAndServe(":8090", nil); err != nil {
		slog.Error("http.ListenAndServe() error", "err", err)
		os.Exit(1)
	}
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
	mlock := &sync.Mutex{}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			slog.Error("upgrader.Upgrade() error", "err", err)
			return
		}

		client := newClient(conn)
		mlock.Lock()
		m[client] = true
		mlock.Unlock()
	})

	go func() {
		for {
			id := <-ch
			mlock.Lock()
			for c := range m {
				if c.IsClosed {
					continue
				}
				c.writeChan <- []byte(id)
			}
			mlock.Unlock()
		}
	}()
}
