package main

import (
	"github.com/sirupsen/logrus"
)

type Hub struct {
	clients []*Client
	register chan *Client
	unregister chan *Client
	broadcast chan Msg
}

func NewHub() *Hub {
	return &Hub{
		clients:      make([]*Client, 0),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		broadcast:    make(chan Msg),
	}
}

func (hub *Hub) Run() {
	for{
		select {
		case client := <-hub.register:
			hub.Regist(client)
		case message := <-hub.broadcast:
			hub.Broadcast(message)
		case client := <-hub.unregister:
			hub.Unregist(client)
		}
	}
}

func (hub *Hub) Regist(client *Client) {
	if len(hub.clients) == 0{
		hub.clients = make([]*Client, 0)
	}
	hub.clients = append(hub.clients, client)
	logrus.Infof("Regist succeed: client=%v", client)
}

func (hub *Hub) Unregist(client *Client) {
	//if _, ok := hub.clients[client]; ok {
	//
	//}
}

func (hub *Hub) Broadcast(msg Msg) {
	for _, client := range hub.clients {
		select {
		case client.receive <- msg:
		//default:
		//	panic("ppppppppp^-^")
		}

	}
	logrus.Infof("Broadcast succeed: clients=%v, msg=%v", hub.clients, msg)
}
