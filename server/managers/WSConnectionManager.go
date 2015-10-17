package managers
import (
	"github.com/gorilla/websocket"
	"log"
	"encoding/json"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/models"
)

//var webSockets map[string][]*websocket.Conn
var webSockets []*websocket.Conn

func WebSocketConnected(conn *websocket.Conn) {
	webSockets = append(webSockets, conn)
}

func WebSocketDisconnected(conn *websocket.Conn) {
	for p, v := range webSockets {
		if v == conn {
			copy(webSockets[p:], webSockets[p + 1:])
			webSockets[len(webSockets) - 1] = nil // or the zero value of T
			webSockets = webSockets[:len(webSockets) - 1]
		}
	}
}

func NotifyAll(projectId string, notification *baseModels.WSNotification) {
	// Notify all connected clients
	// TODO: Change to use RabbitMQ or Redis
	for _, v := range webSockets {
		SendWebSocketMessage(v, notification)
	}
}

func SendWebSocketMessage(conn *websocket.Conn, message interface{}) error {

	respBytes, err := json.Marshal(message)
	log.Println(string(respBytes[:]))

	if err != nil {
		log.Println("Error serializing response to JSON:", err)
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, respBytes)
	if err != nil {
		log.Println("Error writing to WebSocket:", err)
		return err
	}
	return nil
}