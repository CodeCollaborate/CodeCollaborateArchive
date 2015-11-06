package managers

import (
	"github.com/gorilla/websocket"
	"encoding/json"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/models"
	"github.com/CodeCollaborate/CodeCollaborate/server/managers/models"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/project/requests"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/file/models"
)

var proj_wsConn = map[string][]*models.WSConnection{} // maps projectId to WSConnection instances
var wsConn_proj = map[*models.WSConnection][]string{} // maps WSConnection instances to array of projectIds
var webSocket_wsConn = map[*websocket.Conn]*models.WSConnection{} // maps websocket connection to WSConnection instances

func WebSocketSubscribeProject(conn *websocket.Conn, username string, projectId string) bool {

	var wsConn *models.WSConnection

	if (webSocket_wsConn[conn] == nil) {
		wsConn = new(models.WSConnection)
		wsConn.Username = username
		wsConn.WSConn = conn

		webSocket_wsConn[conn] = wsConn
	} else if (webSocket_wsConn[conn].Username == username && webSocket_wsConn[conn].WSConn == conn) {
		return true
	} else {
		wsConn = webSocket_wsConn[conn];
	}

	if (proj_wsConn[projectId] == nil) {
		proj_wsConn[projectId] = []*models.WSConnection{}
	}
	if (wsConn_proj[wsConn] == nil) {
		wsConn_proj[wsConn] = []string{}
	}

	proj_wsConn[projectId] = append(proj_wsConn[projectId], wsConn)
	wsConn_proj[wsConn] = append(wsConn_proj[wsConn], projectId)

	return true
}

func WebSocketDisconnected(conn *websocket.Conn) {

	// if user has subscribed to a project during this connection, remove them from those projects
	if (webSocket_wsConn[conn] != nil) {

		wsConn := webSocket_wsConn[conn];

		for _, project := range wsConn_proj[wsConn] {
			for i, v := range proj_wsConn[project] {
				if v == wsConn {
					copy(proj_wsConn[project][i:], proj_wsConn[project][i + 1:])
					proj_wsConn[project][len(proj_wsConn[project]) - 1] = nil // or the zero value of T
					proj_wsConn[project] = proj_wsConn[project][:len(proj_wsConn[project]) - 1]
					if len(proj_wsConn[project]) == 0 {
						files, err := fileModels.GetFilesByProjectId(project)
						if err {
							LogError("Error retreaving project files on WS disconnect", err)
						} else {
							for _, file := range files {
								scrunchDB(file.Id)
							}
							delete(proj_wsConn, project)
						}
					}
				}
			}
		}

		delete(wsConn_proj, wsConn)
	}

	// remove session from the active sessions.
	delete(webSocket_wsConn, conn)
}

func NotifyProjectClients(projectId string, notification *baseModels.WSNotification, wsConn *websocket.Conn) {
	// Notify all connected clients
	//	TODO: Change to use RabbitMQ or Redis

	for _, value := range proj_wsConn[projectId] {
		if (value.WSConn != wsConn) {
			SendWebSocketMessage(value.WSConn, notification)
		}
	}
}

func GetSubscribedClients(conn *websocket.Conn, getConnectedClientsRequest projectRequests.ProjectGetSubscribedClientsRequest) {

	var users []string = make([]string, 0, len(proj_wsConn[getConnectedClientsRequest.BaseRequest.ResId]))

	for _, value := range proj_wsConn[getConnectedClientsRequest.BaseRequest.ResId] {

		users = append(users, value.Username)
	}

	SendWebSocketMessage(conn, baseModels.NewSuccessResponse(getConnectedClientsRequest.BaseRequest.Tag, map[string]interface{}{"SubscribedUsers": users}))
}

func SendWebSocketMessage(conn *websocket.Conn, message interface{}) error {
	respBytes, err := json.Marshal(message)
	LogDebug(string(respBytes[:]))

	if err != nil {
		LogError("Error serializing response to JSON:", err)
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, respBytes)
	if err != nil {
		LogError("Error writing to WebSocket:", err)
		return err
	}
	return nil
}
