package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/CodeCollaborate/CodeCollaborate/modules/base"
	"github.com/CodeCollaborate/CodeCollaborate/modules/user/models"
	"github.com/CodeCollaborate/CodeCollaborate/modules/user/requests"
	"strings"
	"github.com/CodeCollaborate/CodeCollaborate/managers"
	"github.com/CodeCollaborate/CodeCollaborate/modules/file/requests"
	"github.com/CodeCollaborate/CodeCollaborate/modules/project/requests"
	"github.com/CodeCollaborate/CodeCollaborate/modules/project/models"
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
		log.Println("Failed to upgrade connection:", err)
		return
	}

	webSockets = append(webSockets, wsConn)

	defer wsConn.Close()
	defer socketDisconnected(wsConn)

	for {
		messageType, message, err := wsConn.ReadMessage()
		var response base.WSResponse = base.NewFailResponse(-0, 0, nil)
		if err != nil {
			log.Println("Error reading from WebSocket:", err)
			break
		}

		// Deserialize data from json.
		var baseMessageObj base.BaseRequest
		if err := json.Unmarshal(message, &baseMessageObj); err != nil {

			response = base.NewFailResponse(-1, baseMessageObj.Tag, map[string]interface{}{"Error:":err})

		} else {
			if !(strings.Compare("User", baseMessageObj.Resource) == 0 && (strings.Compare("Register", baseMessageObj.Action) == 0 || strings.Compare("Login", baseMessageObj.Action) == 0)) && !userModels.CheckUserAuth(baseMessageObj) {
				response = base.NewFailResponse(-105, baseMessageObj.Tag, nil)
			} else {

				switch baseMessageObj.Resource {
				case "Project":
					switch baseMessageObj.Action {
					case "Create":

						// {"Resource":"Project", "Action":"Create", "UserId":"5615d78f4357413454000001", "Token": "$2a$10$FriLlb6m9GyxqxURN9YJj.8CmkefQF/uM454fSZY4LwazY.0X/nr2", "Name":"foo"}
						// Deserialize from JSON
						var projectCreateRequest projectRequests.ProjectCreateRequest
						if err := json.Unmarshal(message, &projectCreateRequest); err != nil {

							response = base.NewFailResponse(-1, baseMessageObj.Tag, nil)
							break
						}
						// Add BaseMessage reference
						projectCreateRequest.BaseMessage = baseMessageObj

						response = projectModels.CreateProject(projectCreateRequest)
					case "Rename":

						// {"Resource":"Project", "Action":"Rename", "UserId":"5615d78f4357413454000001", "Token": "$2a$10$FriLlb6m9GyxqxURN9YJj.8CmkefQF/uM454fSZY4LwazY.0X/nr2", "ProjectId": "5615d977435741340c000001", "NewName":"bar"}
						// Deserialize from JSON
						var projectRenameRequest projectRequests.ProjectRenameRequest
						if err := json.Unmarshal(message, &projectRenameRequest); err != nil {

							response = base.NewFailResponse(-1, baseMessageObj.Tag, nil)
							break
						}
						// Add BaseMessage reference
						projectRenameRequest.BaseMessage = baseMessageObj

						response = projectModels.RenameProject(projectRenameRequest)

					case "GrantPermissions":

						// {"Resource":"Project", "Action":"GrantPermissions", "UserId":"5615d78f4357413454000001", "Token": "$2a$10$FriLlb6m9GyxqxURN9YJj.8CmkefQF/uM454fSZY4LwazY.0X/nr2", "ProjectId": "5615d977435741340c000001", "GrantUserId":"5615ee9f4357410d10000001", "PermissionLevel":5}
						// Deserialize from JSON
						var projectGrantPermissionsRequest projectRequests.ProjectGrantPermissionsRequest
						if err := json.Unmarshal(message, &projectGrantPermissionsRequest); err != nil {

							response = base.NewFailResponse(-1, baseMessageObj.Tag, nil)
							break
						}
						// Add BaseMessage reference
						projectGrantPermissionsRequest.BaseMessage = baseMessageObj

						response = projectModels.GrantProjectPermissions(projectGrantPermissionsRequest)

					case "RevokePermissions":

						// {"Resource":"Project", "Action":"RevokePermissions", "UserId":"5615d78f4357413454000001", "Token": "$2a$10$FriLlb6m9GyxqxURN9YJj.8CmkefQF/uM454fSZY4LwazY.0X/nr2", "ProjectId": "5615d977435741340c000001", "RevokeUserId":"5615ee9f4357410d10000001"}
						// Deserialize from JSON
						var projectRevokePermissionsRequest projectRequests.ProjectRevokePermissionsRequest
						if err := json.Unmarshal(message, &projectRevokePermissionsRequest); err != nil {

							response = base.NewFailResponse(-1, baseMessageObj.Tag, nil)
							break
						}
						// Add BaseMessage reference
						projectRevokePermissionsRequest.BaseMessage = baseMessageObj

						response = projectModels.RevokeProjectPermissions(projectRevokePermissionsRequest)

					default:
						response = base.NewFailResponse(-3, baseMessageObj.Tag, nil)
						break
					}
				case "File":

					// eg: {"Tag": 112, "Action": "Update", "Resource": "File", "ResId": "511", "CommitHash": "4as5d4w5as", "Changes": "@@ -40,16 +40,17 @@\n almost i\n+t\n n shape", "UserId": "5615d78f4357413454000001", "Token": "$2a$10$FriLlb6m9GyxqxURN9YJj.8CmkefQF/uM454fSZY4LwazY.0X/nr2"}
					// Deserialize from JSON
					var fileMessageObj fileRequests.FileRequest
					if err := json.Unmarshal(message, &fileMessageObj); err != nil {

						response = base.NewFailResponse(-1, baseMessageObj.Tag, nil)
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
					for _, WSConnection := range webSockets {
						sendWebSocketMessage(WSConnection, websocket.TextMessage, notification)
					}
				case "User":
					switch baseMessageObj.Action {
					case "Register":

						// {"Resource":"User", "Action":"Register", "Username":"abcd", "Email":"abcd@efgh.edu", "Password":"abcd1234"}
						// Deserialize from JSON
						var userRegisterRequest userRequests.UserRegisterRequest
						if err := json.Unmarshal(message, &userRegisterRequest); err != nil {

							response = base.NewFailResponse(-1, baseMessageObj.Tag, nil)
							break
						}
						// Add BaseMessage reference
						userRegisterRequest.BaseMessage = baseMessageObj

						response = userModels.RegisterUser(userRegisterRequest)
					case "Login":

						// {"Resource":"User", "Action":"Login", "UsernameOREmail":"abcd", "Password":"abcd1234"}
						// Deserialize from JSON
						var userLoginRequest userRequests.UserLoginRequest
						if err := json.Unmarshal(message, &userLoginRequest); err != nil {

							response = base.NewFailResponse(-1, baseMessageObj.Tag, nil)
							break
						}
						// Add BaseMessage reference
						userLoginRequest.BaseMessage = baseMessageObj

						//Check username/pw, login if needed.
						response = userModels.LoginUser(userLoginRequest);

					default:
						response = base.NewFailResponse(-3, baseMessageObj.Tag, nil)
						break
					}
				default:
					// Invalid resource type
					response = base.NewFailResponse(-2, baseMessageObj.Tag, nil)
					break
				}
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

	managers.ConnectMGo()
	defer managers.GetPrimaryMGoSession().Close()

	http.HandleFunc("/ws/", handleWSConn)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
type Person struct {
	Name  string
	Phone string
}
