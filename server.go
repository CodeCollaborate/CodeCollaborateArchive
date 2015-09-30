package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	models "github.com/CodeCollaborate/CodeCollaborate/modules/base"
	"github.com/CodeCollaborate/CodeCollaborate/modules/file"
	"github.com/gorilla/websocket"
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
		var response models.WSResponse
		if err != nil {
			log.Println("Error reading from WebSocket:", err)
			break
		}

		// Deserialize data from json.
		// eg: {"Tag": 112, "Action": "Update", "Resource": "File", "ResId": 511, "CommitHash": "4as5d4w5as"}
		var baseMessageObj models.BaseMessage
		if err := json.Unmarshal(message, &baseMessageObj); err != nil {

			response = models.NewFailResponse(-101, baseMessageObj.Tag, "Error deserializing JSON to BaseMessage")

		} else {

			switch baseMessageObj.Resource {
			case "Project":
			case "File":
				// eg: {"Tag": 112, "Action": "Update", "Resource": "File", "ResId": 511, "CommitHash": "4as5d4w5as", "Changes": "@@ -40,16 +40,17 @@\n almost i\n+t\n n shape"}

				// Deserialize FileMessage from JSON
				var fileMessageObj file.FileMessage
				if err := json.Unmarshal(message, &fileMessageObj); err != nil {

					response = models.NewFailResponse(-101, baseMessageObj.Tag, "Error deserializing JSON to FileMessage")
					break
				}

				// Add BaseMessage reference
				fileMessageObj.BaseMessage = baseMessageObj

				// TODO: Do something.

				// Notify success.
				response = models.NewSuccessResponse(baseMessageObj.Tag, nil)

				// For debugging
				log.Println(fileMessageObj.ToString())

			default:
				// Invalid resource type
				response = models.NewFailResponse(-100, baseMessageObj.Tag, "Invalid resource type")
				break
			}
		}

		err = sendWebSocketResponse(c, mt, response)
		if err != nil {
			break
		}
	}
}

func sendWebSocketResponse(conn *websocket.Conn, messageType int, response interface{}) error {

	respBytes, err := json.Marshal(response)
	log.Println(string(respBytes[:]))

	if err != nil {
		log.Println("Error serializing response to JSON:", err)
		return err
	}

	err = conn.WriteMessage(messageType, respBytes)
	if err != nil {
		log.Println("Error writing to WebSocket:", err)
		return err
	}
	return nil
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
