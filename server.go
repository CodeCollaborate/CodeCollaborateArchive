package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"os"

	"github.com/CodeCollaborate/CodeCollaborate/managers"
	"github.com/CodeCollaborate/CodeCollaborate/modules/base"
	"github.com/CodeCollaborate/CodeCollaborate/modules/file/models"
	"github.com/CodeCollaborate/CodeCollaborate/modules/file/requests"
	"github.com/CodeCollaborate/CodeCollaborate/modules/project/models"
	"github.com/CodeCollaborate/CodeCollaborate/modules/project/requests"
	"github.com/CodeCollaborate/CodeCollaborate/modules/user/models"
	"github.com/CodeCollaborate/CodeCollaborate/modules/user/requests"
	"github.com/gorilla/websocket"
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
		var response = base.NewFailResponse(-0, 0, nil)
		if err != nil {
			log.Println("Error reading from WebSocket:", err)
			break
		}

		// Deserialize data from json.
		var baseMessageObj base.BaseRequest
		if err := json.Unmarshal(message, &baseMessageObj); err != nil {

			response = base.NewFailResponse(-1, baseMessageObj.Tag, map[string]interface{}{"Error:": err})

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

					case "Delete":
					// TODO

					default:
						response = base.NewFailResponse(-3, baseMessageObj.Tag, nil)
						break
					}
				case "File":
					// TODO: Do something.
					switch baseMessageObj.Action {

					case "Create":
						// {"Resource":"File", "Action":"Create", "UserId":"5615d78f4357413454000001", "Token": "$2a$10$FriLlb6m9GyxqxURN9YJj.8CmkefQF/uM454fSZY4LwazY.0X/nr2", "Name":"foo", "RelativePath":"test/path1/", "ProjectId":"5615d977435741340c000001"}
						// Deserialize from JSON
						var fileCreateRequest fileRequests.FileCreateRequest
						if err := json.Unmarshal(message, &fileCreateRequest); err != nil {
							response = base.NewFailResponse(-1, baseMessageObj.Tag, nil)
							break
						}
						// Add BaseMessage reference
						fileCreateRequest.BaseMessage = baseMessageObj
						response = fileModels.CreateFile(fileCreateRequest)

					case "Rename":
						// {"Resource":"File", "Action":"Rename", "UserId":"5615d78f4357413454000001", "Token": "$2a$10$FriLlb6m9GyxqxURN9YJj.8CmkefQF/uM454fSZY4LwazY.0X/nr2", "NewFileName":"foo2", "FileId":"561978c14357412bf8000001"}
						// Deserialize from JSON
						var fileRenameRequest fileRequests.FileRenameRequest
						if err := json.Unmarshal(message, &fileRenameRequest); err != nil {
							response = base.NewFailResponse(-1, baseMessageObj.Tag, nil)
							break
						}
						// Add BaseMessage reference
						fileRenameRequest.BaseMessage = baseMessageObj
						response = fileModels.RenameFile(fileRenameRequest)

					case "Move":
						// {"Resource":"File", "Action":"Move", "UserId":"5615d78f4357413454000001", "Token": "$2a$10$FriLlb6m9GyxqxURN9YJj.8CmkefQF/uM454fSZY4LwazY.0X/nr2", "NewPath":"test/path2/", "FileId":"561978c14357412bf8000001"}
						// Deserialize from JSON
						var fileMoveRequest fileRequests.FileMoveRequest
						if err := json.Unmarshal(message, &fileMoveRequest); err != nil {
							response = base.NewFailResponse(-1, baseMessageObj.Tag, nil)
							break
						}
						// Add BaseMessage reference
						fileMoveRequest.BaseMessage = baseMessageObj
						response = fileModels.MoveFile(fileMoveRequest)

					case "Delete":
						// {"Resource":"File", "Action":"Delete", "UserId":"5615d78f4357413454000001", "Token": "$2a$10$FriLlb6m9GyxqxURN9YJj.8CmkefQF/uM454fSZY4LwazY.0X/nr2", "FileId":"561978c14357412bf8000001"}
						// Deserialize from JSON
						var fileDeleteRequest fileRequests.FileDeleteRequest
						if err := json.Unmarshal(message, &fileDeleteRequest); err != nil {
							response = base.NewFailResponse(-1, baseMessageObj.Tag, nil)
							break
						}
						// Add BaseMessage reference
						fileDeleteRequest.BaseMessage = baseMessageObj
						response = fileModels.DeleteFile(fileDeleteRequest)

					case "Change":
						// eg: {"Tag": 112, "Action": "Update", "Resource": "File", "ResId": "511", "CommitHash": "4as5d4w5as", "Changes": "@@ -40,16 +40,17 @@\n almost i\n+t\n n shape", "UserId": "5615d78f4357413454000001", "Token": "$2a$10$FriLlb6m9GyxqxURN9YJj.8CmkefQF/uM454fSZY4LwazY.0X/nr2"}
						// Deserialize from JSON
						var fileChangeRequest fileRequests.FileChangeRequest
						if err := json.Unmarshal(message, &fileChangeRequest); err != nil {
							response = base.NewFailResponse(-1, baseMessageObj.Tag, nil)
							break
						}
						// Add BaseMessage reference
						fileChangeRequest.BaseRequest = baseMessageObj

						response = fileModels.InsertChange(fileChangeRequest)

						// Notify all connected clients
						// TODO: Change to use RabbitMQ or Redis
						// notification := fileChangeRequest.GetNotification()
						// for _, WSConnection := range webSockets {
						// 	sendWebSocketMessage(WSConnection, websocket.TextMessage, notification)
						// }

					default:
						response = base.NewFailResponse(-3, baseMessageObj.Tag, map[string]interface{}{"Action": baseMessageObj.Action})
						break
					}

					// // Notify success; return new version number.
					// response = base.NewSuccessResponse(baseMessageObj.Tag, nil)

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
						response = userModels.LoginUser(userLoginRequest)

						//TODO: maybe delete?

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
		if v == conn {
			copy(webSockets[p:], webSockets[p+1:])
			webSockets[len(webSockets)-1] = nil // or the zero value of T
			webSockets = webSockets[:len(webSockets)-1]
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

func handleHTTPConn(responseWriter http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.Error(responseWriter, "Not found", 404)
		return
	}
	if request.Method != "GET" {
		http.Error(responseWriter, "Method not allowed", 405)
		return
	}
	responseWriter.Header()
	fmt.Fprintf(responseWriter, "hello there")
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Could not get Working Directory: ", err)
	}
	log.Println("Running in directory:", pwd)

	managers.ConnectMGo()
	defer managers.GetPrimaryMGoSession().Close()

	http.HandleFunc("/ws/", handleWSConn)
	http.HandleFunc("/", handleHTTPConn)
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
