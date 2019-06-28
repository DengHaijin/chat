package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"path/filepath"
)

var addr = flag.String("addr", ":8080", "http service address")

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fp, _ := filepath.Abs("./frontend/home.html")
	http.ServeFile(w, r, fp)
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, []byte("you mean " + string(message) + "?"))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	hub := NewHub()
	go hub.Run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		serveWs(hub, writer, request)
	})
	log.Fatal(http.ListenAndServe(*addr, nil))
}
