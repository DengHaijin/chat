package main

import (
	"bytes"
	"git.byted.org/ee/byteview/svc/ultron/util/json"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	BUFF_LEN = 256
)

type Client struct {
	ID      IDType
	conn    *websocket.Conn
	receive chan Msg
	hub     *Hub
}

func (client *Client) WaitForMessage() {
	defer func() {
		client.hub.unregister <- client
		close(client.receive)
	}()
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Errorf("error: err=%v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))
		client.hub.broadcast <- Msg{
			UserID:     client.ID,
			Content:    string(message),
		}
	}
}

func (client *Client) WriteMessage() {
	defer close(client.receive)
	for {
		select {
		case message := <-client.receive:
			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				logrus.Errorf("get writer failed, err=%v", err)
			}
			byteMsg, _ := json.Marshal(message)
			w.Write(byteMsg)
			w.Close()
		}
	}
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Errorf("bad request: err=%v", err)
	}
	client := &Client{
		ID:      GenerateID(),
		conn:    conn,
		receive: make(chan Msg),
		hub:     hub,
	}
	hub.register <- client
	go client.WaitForMessage()
	go client.WriteMessage()
}
