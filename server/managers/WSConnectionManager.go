package managers

import (
	"github.com/gorilla/websocket"
	"log"
	"encoding/json"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/models"
)

var proj_webSockets = map[string][]*websocket.Conn{} // needed initialized
var webSockets_proj = map[*websocket.Conn][]string{}

func WebSocketSubscribeProject(conn *websocket.Conn, projectId string) bool {
	proj_webSockets[projectId] = append(proj_webSockets[projectId], conn)
	webSockets_proj[conn] = append(webSockets_proj[conn], projectId)
	// return false on already exists

	return true
}

func WebSocketDisconnected(conn *websocket.Conn) {

	for _, project := range webSockets_proj[conn] {
		for i, v := range proj_webSockets[project] {
			if v == conn {
				copy(proj_webSockets[project][i:], proj_webSockets[project][i + 1:])
				proj_webSockets[project][len(proj_webSockets[project]) - 1] = nil // or the zero value of T
				proj_webSockets[project] = proj_webSockets[project][:len(proj_webSockets[project]) - 1]
				if len(proj_webSockets[project]) == 0 {
					delete(proj_webSockets, project)
				}
			}
		}
	}

	delete(webSockets_proj, conn)

}

func NotifyProjectClients(projectId string, notification *baseModels.WSNotification) {
	// Notify all connected clients
	//	TODO: Change to use RabbitMQ or Redis

	for _, v := range proj_webSockets[projectId] {
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