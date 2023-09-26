package main

import (
	"log/slog"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	conn      *websocket.Conn
	writeChan chan []byte
}

func newClient(conn *websocket.Conn) *Client {
	c := &Client{
		conn:      conn,
		writeChan: make(chan []byte, 5),
	}

	go c.read()
	go c.write()

	return c
}

func (c *Client) read() {
	defer c.conn.Close()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			slog.Error("conn.ReadMessage() error", "err", err)
			return
		}

		slog.Info("read", "msg", string(msg))
	}
}

func (c *Client) write() {
	defer c.conn.Close()

	defer func() {
		n := len(c.writeChan)
		for i := 0; i < n; i++ {
			<-c.writeChan
		}
	}()

	for {
		select {
		case msg := <-c.writeChan:
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				slog.Error("conn.WriteMessage() error", "err", err)
				return
			}
		}
	}
}