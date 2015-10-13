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

	managers.WebSocketConnected(wsConn)

	defer wsConn.Close()
	defer managers.WebSocketDisconnected(wsConn)

	for {
		// messageType, message, err := wsConn.ReadMessage()
		messageType, message, err := wsConn.ReadMessage()

		if(messageType == websocket.BinaryMessage){
			// file upload!
			
		}
		var response = base.NewFailResponse(-0, 0, nil)
		if err != nil {
			log.Println("Error reading from WebSocket:", err)
			break
		}

		// Deserialize data from json.
		var baseRequestObj base.BaseRequest
		if err := json.Unmarshal(message, &baseRequestObj); err != nil {

			response = base.NewFailResponse(-1, baseRequestObj.Tag, map[string]interface{}{"Error:": err})

		} else {
			if !(strings.Compare("User", baseRequestObj.Resource) == 0 && (strings.Compare("Register", baseRequestObj.Action) == 0 || strings.Compare("Login", baseRequestObj.Action) == 0)) && !userModels.CheckUserAuth(baseRequestObj) {
				response = base.NewFailResponse(-105, baseRequestObj.Tag, nil)
			} else {

				switch baseRequestObj.Resource {
				case "Project":
					switch baseRequestObj.Action {
					case "Create":

						// {"Resource":"Project", "Action":"Create", "UserId":"561986674357413b14000001", "Token": "$2a$10$gifm6Vrfn2vBBCX7qvaQzu.Pvttotyu1pRW5V6X7RnhYYiQCUHh4e", "Name":"foo"}
						// Deserialize from JSON
						var projectCreateRequest projectRequests.ProjectCreateRequest
						if err := json.Unmarshal(message, &projectCreateRequest); err != nil {
							response = base.NewFailResponse(-1, baseRequestObj.Tag, nil)
							break
						}
						// Add BaseRequest reference
						projectCreateRequest.BaseRequest = baseRequestObj
						response = projectModels.CreateProject(projectCreateRequest)

					case "Rename":

						// {"Resource":"Project", "Action":"Rename", "ResId": "561987174357413b14000002", "UserId":"561986674357413b14000001", "Token": "$2a$10$gifm6Vrfn2vBBCX7qvaQzu.Pvttotyu1pRW5V6X7RnhYYiQCUHh4e", "NewName":"bar"}
						// Deserialize from JSON
						var projectRenameRequest projectRequests.ProjectRenameRequest
						if err := json.Unmarshal(message, &projectRenameRequest); err != nil {
							response = base.NewFailResponse(-1, baseRequestObj.Tag, nil)
							break
						}
						// Add BaseRequest reference
						projectRenameRequest.BaseRequest = baseRequestObj
						response = projectModels.RenameProject(projectRenameRequest)

					case "GrantPermissions":

						// {"Resource":"Project", "Action":"GrantPermissions", "ResId": "561987174357413b14000002", "UserId":"561986674357413b14000001", "Token": "$2a$10$gifm6Vrfn2vBBCX7qvaQzu.Pvttotyu1pRW5V6X7RnhYYiQCUHh4e", "GrantUserId":"561987434357413b14000003", "PermissionLevel":5}
						// Deserialize from JSON
						var projectGrantPermissionsRequest projectRequests.ProjectGrantPermissionsRequest
						if err := json.Unmarshal(message, &projectGrantPermissionsRequest); err != nil {

							response = base.NewFailResponse(-1, baseRequestObj.Tag, nil)
							break
						}
						// Add BaseRequest reference
						projectGrantPermissionsRequest.BaseRequest = baseRequestObj

						response = projectModels.GrantProjectPermissions(projectGrantPermissionsRequest)

					case "RevokePermissions":

						// {"Resource":"Project", "Action":"RevokePermissions", "ResId": "561987174357413b14000002", "UserId":"561986674357413b14000001", "Token": "$2a$10$gifm6Vrfn2vBBCX7qvaQzu.Pvttotyu1pRW5V6X7RnhYYiQCUHh4e", "RevokeUserId":"561987434357413b14000003"}
						// Deserialize from JSON
						var projectRevokePermissionsRequest projectRequests.ProjectRevokePermissionsRequest
						if err := json.Unmarshal(message, &projectRevokePermissionsRequest); err != nil {

							response = base.NewFailResponse(-1, baseRequestObj.Tag, nil)
							break
						}
						// Add BaseRequest reference
						projectRevokePermissionsRequest.BaseRequest = baseRequestObj

						response = projectModels.RevokeProjectPermissions(projectRevokePermissionsRequest)

					case "Delete":
					// TODO

					default:
						response = base.NewFailResponse(-3, baseRequestObj.Tag, nil)
						break
					}
				case "File":
					// TODO: Do something.
					switch baseRequestObj.Action {

					case "Create":
						// {"Resource":"File", "Action":"Create", "UserId":"561986674357413b14000001", "Token": "$2a$10$gifm6Vrfn2vBBCX7qvaQzu.Pvttotyu1pRW5V6X7RnhYYiQCUHh4e", "Name":"foo", "RelativePath":"test/path1/", "ProjectId":"561987174357413b14000002"}
						// Deserialize from JSON
						var fileCreateRequest fileRequests.FileCreateRequest
						if err := json.Unmarshal(message, &fileCreateRequest); err != nil {
							response = base.NewFailResponse(-1, baseRequestObj.Tag, nil)
							break
						}
						// Add BaseRequest reference
						fileCreateRequest.BaseRequest = baseRequestObj
						response = fileModels.CreateFile(fileCreateRequest)

					case "Rename":
						// {"Resource":"File", "Action":"Rename", "ResId":"561987a84357413b14000006", "UserId":"561986674357413b14000001", "Token": "$2a$10$gifm6Vrfn2vBBCX7qvaQzu.Pvttotyu1pRW5V6X7RnhYYiQCUHh4e", "NewName":"foo2"}
						// Deserialize from JSON
						var fileRenameRequest fileRequests.FileRenameRequest
						if err := json.Unmarshal(message, &fileRenameRequest); err != nil {
							response = base.NewFailResponse(-1, baseRequestObj.Tag, nil)
							break
						}
						// Add BaseRequest reference
						fileRenameRequest.BaseRequest = baseRequestObj
						response = fileModels.RenameFile(fileRenameRequest)

					case "Move":
						// {"Resource":"File", "Action":"Move", "ResId":"561987a84357413b14000006", "UserId":"561986674357413b14000001", "Token": "$2a$10$gifm6Vrfn2vBBCX7qvaQzu.Pvttotyu1pRW5V6X7RnhYYiQCUHh4e", "NewPath":"test/path2/"}
						// Deserialize from JSON
						var fileMoveRequest fileRequests.FileMoveRequest
						if err := json.Unmarshal(message, &fileMoveRequest); err != nil {
							response = base.NewFailResponse(-1, baseRequestObj.Tag, nil)
							break
						}
						// Add BaseRequest reference
						fileMoveRequest.BaseRequest = baseRequestObj
						response = fileModels.MoveFile(fileMoveRequest)

					case "Delete":
						// {"Resource":"File", "Action":"Delete", "ResId":"561987a84357413b14000006", "UserId":"561986674357413b14000001", "Token": "$2a$10$gifm6Vrfn2vBBCX7qvaQzu.Pvttotyu1pRW5V6X7RnhYYiQCUHh4e"}
						// Deserialize from JSON
						var fileDeleteRequest fileRequests.FileDeleteRequest
						if err := json.Unmarshal(message, &fileDeleteRequest); err != nil {
							response = base.NewFailResponse(-1, baseRequestObj.Tag, nil)
							break
						}
						// Add BaseRequest reference
						fileDeleteRequest.BaseRequest = baseRequestObj
						response = fileModels.DeleteFile(fileDeleteRequest)

					case "Change":
						// {"Tag": 112, "Action": "Change", "Resource": "File", "ResId": "561987a44357413b14000004", "FileVersion":0, "Changes": "@@ -40,16 +40,17 @@\n almost i\n+t\n n shape", "UserId": "561986674357413b14000001", "Token": "$2a$10$gifm6Vrfn2vBBCX7qvaQzu.Pvttotyu1pRW5V6X7RnhYYiQCUHh4e"}
						// Deserialize from JSON
						var fileChangeRequest fileRequests.FileChangeRequest
						if err := json.Unmarshal(message, &fileChangeRequest); err != nil {
							response = base.NewFailResponse(-1, baseRequestObj.Tag, nil)
							break
						}
						// Add BaseRequest reference
						fileChangeRequest.BaseRequest = baseRequestObj

						response = fileModels.InsertChange(fileChangeRequest)

					// Notify all connected clients
					// TODO: Change to use RabbitMQ or Redis
					// notification := fileChangeRequest.GetNotification()
					// for _, WSConnection := range webSockets {
					// 	sendWebSocketMessage(WSConnection, websocket.TextMessage, notification)
					// }

					default:
						response = base.NewFailResponse(-3, baseRequestObj.Tag, map[string]interface{}{"Action": baseRequestObj.Action})
						break
					}

				// // Notify success; return new version number.
				// response = base.NewSuccessResponse(baseRequestObj.Tag, nil)

				case "User":
					switch baseRequestObj.Action {
					case "Register":

						// {"Resource":"User", "Action":"Register", "Username":"abcd", "Email":"abcd@efgh.edu", "Password":"abcd1234"}
						// Deserialize from JSON
						var userRegisterRequest userRequests.UserRegisterRequest
						if err := json.Unmarshal(message, &userRegisterRequest); err != nil {

							response = base.NewFailResponse(-1, baseRequestObj.Tag, nil)
							break
						}
						// Add BaseRequest reference
						userRegisterRequest.BaseRequest = baseRequestObj

						response = userModels.RegisterUser(userRegisterRequest)
					case "Login":

						// {"Resource":"User", "Action":"Login", "UsernameOREmail":"abcd", "Password":"abcd1234"}
						// Deserialize from JSON
						var userLoginRequest userRequests.UserLoginRequest
						if err := json.Unmarshal(message, &userLoginRequest); err != nil {

							response = base.NewFailResponse(-1, baseRequestObj.Tag, nil)
							break
						}
						// Add BaseRequest reference
						userLoginRequest.BaseRequest = baseRequestObj

						//Check username/pw, login if needed.
						response = userModels.LoginUser(userLoginRequest)

					//TODO: maybe delete?

					//TODO: Change PW

					default:
						response = base.NewFailResponse(-3, baseRequestObj.Tag, nil)
						break
					}
				default:
					// Invalid resource type
					response = base.NewFailResponse(-2, baseRequestObj.Tag, nil)
					break
				}
			}
		}

		err = managers.SendWebSocketMessage(wsConn, response)
		if err != nil {
			break
		}
	}
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
