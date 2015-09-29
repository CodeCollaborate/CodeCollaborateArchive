package main

import (
	"flag"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"encoding/json"
	models "github.com/obsessiveorange/CodeCollaborate/modules/models"
)

var addr = flag.String("addr", ":80", "http service address")

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},

} // use default options

func echo(responseWriter http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/ws/" {
		http.Error(responseWriter, "Not found", 404)
		return
	}
	if request.Method != "GET" {
		http.Error(responseWriter, "Method not allowed", 405)
		return
	}
	c, err := upgrader.Upgrade(responseWriter, request, nil)
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

		// Deserialize data from json.
		// eg: {   "Action": "testAction",   "Resource": "testResouce",   "Id": 123,   "CommitHash": "4as5d4w5as" }
		var messageObj models.Message
		if err := json.Unmarshal(message, &messageObj); err != nil {
			panic(err)
		}
		log.Println(messageObj.ToString())

		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	http.HandleFunc("/ws/", echo)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
