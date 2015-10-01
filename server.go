package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/CodeCollaborate/CodeCollaborate/modules/base"
	"github.com/CodeCollaborate/CodeCollaborate/modules/file"
)

var addr = flag.String("addr", ":80", "http service address")

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

var webSockets []*websocket.Conn

func handleWSConn(responseWriter http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/ws/" {
		http.Error(responseWriter, "Not found", 404)
		return
	}
	if request.Method != "GET" {
		http.Error(responseWriter, "Method not allowed", 405)
		return
	}
	wsConn, err := upgrader.Upgrade(responseWriter, request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	webSockets = append(webSockets, wsConn)

	defer wsConn.Close()
	defer socketDisconnected(wsConn)

	for {
		messageType, message, err := wsConn.ReadMessage()
		var response base.WSResponse
		if err != nil {
			log.Println("Error reading from WebSocket:", err)
			break
		}

		// Deserialize data from json.
		// eg: {"Tag": 112, "Action": "Update", "Resource": "File", "ResId": 511, "CommitHash": "4as5d4w5as"}
		var baseMessageObj base.BaseRequest
		if err := json.Unmarshal(message, &baseMessageObj); err != nil {

			response = base.NewFailResponse(-101, baseMessageObj.Tag, nil)

		} else {

			switch baseMessageObj.Resource {
			case "Project":
			case "File":
				// eg: {"Tag": 112, "Action": "Update", "Resource": "File", "ResId": 511, "CommitHash": "4as5d4w5as", "Changes": "@@ -40,16 +40,17 @@\n almost i\n+t\n n shape"}

				// Deserialize FileMessage from JSON
				var fileMessageObj file.FileRequest
				if err := json.Unmarshal(message, &fileMessageObj); err != nil {

					response = base.NewFailResponse(-102, baseMessageObj.Tag, nil)
					break
				}

				// Add BaseMessage reference
				fileMessageObj.BaseMessage = baseMessageObj

				// TODO: Do something.

				// Notify success; return new version number.
				response = base.NewSuccessResponse(baseMessageObj.Tag, nil)

				// Notify all connected clients
				// TODO: Change to use RabbitMQ or Redis
				notification := fileMessageObj.GetNotification()
				for i := range webSockets {
					sendWebSocketMessage(webSockets[i], websocket.TextMessage, notification)
				}

			default:
				// Invalid resource type
				response = base.NewFailResponse(-100, baseMessageObj.Tag, nil)
				break
			}
		}

		err = sendWebSocketMessage(wsConn, messageType, response)
		if err != nil {
			break
		}
	}
}

func socketDisconnected(conn *websocket.Conn) {
	for p, v := range webSockets {
		if (v == conn) {
			copy(webSockets[p:], webSockets[p + 1:])
			webSockets[len(webSockets) - 1] = nil // or the zero value of T
			webSockets = webSockets[:len(webSockets) - 1]
		}
	}
}

func sendWebSocketMessage(conn *websocket.Conn, messageType int, message interface{}) error {

	respBytes, err := json.Marshal(message)
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

	http.HandleFunc("/ws/", handleWSConn)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
