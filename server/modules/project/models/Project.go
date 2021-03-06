package projectModels

import (
	"github.com/CodeCollaborate/CodeCollaborate/server/managers"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/project/requests"
	"gopkg.in/mgo.v2/bson"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/models"
	"github.com/gorilla/websocket"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/file/models"
)

/**
Permissions for users:
	1 - read
	2 - read/comment
	3 - read/write
	5 - admin
	10 - owners

Admins and above can grant/revoke permissions for anyone at the same, or a lower permission level.

TODO: Accepts wildcard flag for Username: "*"
*/
type Project struct {
	Id          string         `bson:"_id"` // ID of object
	Name        string                      // Name of project
	ServerPath  string         `json:"-"`   // Path on server
	Permissions map[string]int              // Array of references to User objects

											//TODO: Add project versions, incremented on file creation, deletion, checked on ws connect
											//TODO: wildcard permissions, add once we make adding to projects a thing
}

func GetPermissionLevels(wsConn *websocket.Conn, projectGetPermissionLevelsRequest projectRequests.ProjectGetPermissionLevelsRequest) {

	managers.SendWebSocketMessage(wsConn, baseModels.NewSuccessResponse(projectGetPermissionLevelsRequest.BaseRequest.Tag, map[string]interface{}{"PermissionLevels": PermissionLevels}))
}

// Create new project
func CreateProject(wsConn *websocket.Conn, projectCreateRequest projectRequests.ProjectCreateRequest) {

	// Create new Project object
	project := new(Project)
	project.Id = managers.NewObjectIdString()
	project.Name = projectCreateRequest.Name
	project.ServerPath = project.Id
	project.Permissions = map[string]int{projectCreateRequest.BaseRequest.Username: 10} // Set creator to owner permissions

	// Get new DB connection
	session, collection := managers.GetMGoCollection("Projects")
	defer session.Close()

	// Create the project
	err := collection.Insert(project)
	if err != nil {
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-201, projectCreateRequest.BaseRequest.Tag, nil))
		return
	}

	managers.SendWebSocketMessage(wsConn, baseModels.NewSuccessResponse(projectCreateRequest.BaseRequest.Tag, map[string]interface{}{"ProjectId": project.Id}))
}

// Rename project (?)
func RenameProject(wsConn *websocket.Conn, projectRenameRequest projectRequests.ProjectRenameRequest) {

	// Get new DB connection
	session, collection := managers.GetMGoCollection("Projects")
	defer session.Close()

	// Rename the project
	err := collection.Update(bson.M{"_id": projectRenameRequest.BaseRequest.ResId}, bson.M{"$set": bson.M{"name": projectRenameRequest.NewName}})
	if err != nil {
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-202, projectRenameRequest.BaseRequest.Tag, nil))
		return
	}

	managers.SendWebSocketMessage(wsConn, baseModels.NewSuccessResponse(projectRenameRequest.BaseRequest.Tag, nil))
	managers.NotifyProjectClients(projectRenameRequest.BaseRequest.ResId, projectRenameRequest.GetNotification(), wsConn)
}

// Grant permission <Level> to <User>
//  - Check if user exists
//  - Grants permission level to user, overwriting if necessary.
func GrantProjectPermissions(wsConn *websocket.Conn, projectGrantPermissionsRequest projectRequests.ProjectGrantPermissionsRequest) {

	project, err := GetProjectById(projectGrantPermissionsRequest.BaseRequest.ResId)
	if err != nil {
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-200, projectGrantPermissionsRequest.BaseRequest.Tag, nil))
		return
	}

	if (!CheckUserHasPermissions(project, projectGrantPermissionsRequest.BaseRequest.Username, 5)) {
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-207, projectGrantPermissionsRequest.BaseRequest.Tag, nil))
		return
	}

	// Make sure that there is still an owner of the project.
	owner := ""
	for key, value := range project.Permissions {
		if value == 10 && key != projectGrantPermissionsRequest.GrantUsername {
			owner = key
		}
	}
	if owner == "" {
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-205, projectGrantPermissionsRequest.BaseRequest.Tag, nil))
		return
	}

	project.Permissions[projectGrantPermissionsRequest.GrantUsername] = projectGrantPermissionsRequest.PermissionLevel

	// Get new DB connection
	session, collection := managers.GetMGoCollection("Projects")
	defer session.Close()

	// Update permissions
	err = collection.Update(bson.M{"_id": projectGrantPermissionsRequest.BaseRequest.ResId}, bson.M{"$set": bson.M{"permissions": project.Permissions}})
	if err != nil {
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-202, projectGrantPermissionsRequest.BaseRequest.Tag, nil))
		return
	}

	managers.SendWebSocketMessage(wsConn, baseModels.NewSuccessResponse(projectGrantPermissionsRequest.BaseRequest.Tag, nil))
	managers.NotifyProjectClients(projectGrantPermissionsRequest.BaseRequest.ResId, projectGrantPermissionsRequest.GetNotification(), wsConn)
}

// Revoke permission for <User>
//  - Check if user has permissions
//  - Revokes permissions entirely; removes entry.
func RevokeProjectPermissions(wsConn *websocket.Conn, projectRevokePermissionsRequest projectRequests.ProjectRevokePermissionsRequest) {

	project, err := GetProjectById(projectRevokePermissionsRequest.BaseRequest.ResId)
	if err != nil {
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-200, projectRevokePermissionsRequest.BaseRequest.Tag, nil))
		return
	}

	if (!CheckUserHasPermissions(project, projectRevokePermissionsRequest.BaseRequest.Username, 5)) {
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-207, projectRevokePermissionsRequest.BaseRequest.Tag, nil))
		return
	}

	// Make sure that there is still an owner of the project.
	owner := ""
	for key, value := range project.Permissions {
		if value == 10 && key != projectRevokePermissionsRequest.RevokeUsername {
			owner = key
		}
	}
	if owner == "" {
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-205, projectRevokePermissionsRequest.BaseRequest.Tag, nil))
		return
	}

	delete(project.Permissions, projectRevokePermissionsRequest.RevokeUsername)

	// Get new DB connection
	session, collection := managers.GetMGoCollection("Projects")
	defer session.Close()

	// Update permissions
	err = collection.Update(bson.M{"_id": projectRevokePermissionsRequest.BaseRequest.ResId}, bson.M{"$set": bson.M{"permissions": project.Permissions}})
	if err != nil {
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-202, projectRevokePermissionsRequest.BaseRequest.Tag, nil))
		return
	}

	managers.SendWebSocketMessage(wsConn, baseModels.NewSuccessResponse(projectRevokePermissionsRequest.BaseRequest.Tag, nil))
	managers.NotifyProjectClients(projectRevokePermissionsRequest.BaseRequest.ResId, projectRevokePermissionsRequest.GetNotification(), wsConn)
}

func GetCollaborators(wsConn *websocket.Conn, getCollaboratorsRequest projectRequests.ProjectGetCollaboratorsRequest) {

	project, err := GetProjectById(getCollaboratorsRequest.BaseRequest.ResId);

	if err != nil {
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-200, getCollaboratorsRequest.BaseRequest.Tag, nil))
		return
	}

	managers.SendWebSocketMessage(wsConn, baseModels.NewSuccessResponse(getCollaboratorsRequest.BaseRequest.Tag, map[string]interface{}{"Collaborators":project.Permissions}))
}

func GetFiles(wsConn *websocket.Conn, projectFilesRequest projectRequests.ProjectGetFilesRequest) {

	projects, err := fileModels.GetFilesByProjectId(projectFilesRequest.BaseRequest.ResId);

	if err != nil {
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-200, projectFilesRequest.BaseRequest.Tag, nil))
		return
	}

	managers.SendWebSocketMessage(wsConn, baseModels.NewSuccessResponse(projectFilesRequest.BaseRequest.Tag, map[string]interface{}{"Files":projects}))
}

func Subscribe(wsConn *websocket.Conn, subscriptionRequest projectRequests.ProjectSubscribeRequest) {

	projectId := subscriptionRequest.BaseRequest.ResId

	proj, err := GetProjectById(projectId)

	if err != nil {
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-200, subscriptionRequest.BaseRequest.Tag, nil))
		return
	}

	hasPermission := false
	for key, value := range proj.Permissions {
		if ((key == subscriptionRequest.BaseRequest.Username || key == "*") && value >= 1) {
			if (!managers.WebSocketSubscribeProject(wsConn, subscriptionRequest.BaseRequest.Username, projectId)) {
				managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-206, subscriptionRequest.BaseRequest.Tag, nil))
				return
			}
			hasPermission = true
		}
	}

	if (!hasPermission){
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-208, subscriptionRequest.BaseRequest.Tag, nil))
		return
	}

	managers.NotifyProjectClients(projectId, subscriptionRequest.GetNotification(), wsConn)
	managers.SendWebSocketMessage(wsConn, baseModels.NewSuccessResponse(subscriptionRequest.BaseRequest.Tag, nil))
}

func Unsubscribe(wsConn *websocket.Conn, unsubscriptionRequest projectRequests.ProjectUnsubscribeRequest) {

//	project := subscriptionRequest.BaseRequest.ResId
//
//	proj, err := GetProjectById(project)
//
//	if err != nil {
//		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-200, subscriptionRequest.BaseRequest.Tag, nil))
//		return
//	}
//
//	// TODO: Add fail message if permission denied
//	for key, _ := range proj.Permissions {
//		if key == subscriptionRequest.BaseRequest.Username {
//			if (!managers.WebSocketSubscribeProject(wsConn, subscriptionRequest.BaseRequest.Username, project)) {
//				managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-206, subscriptionRequest.BaseRequest.Tag, nil))
//				return
//			}
//		}
//	}
//
//	managers.NotifyProjectClients(project, subscriptionRequest.GetNotification(), wsConn)
//	managers.SendWebSocketMessage(wsConn, baseModels.NewSuccessResponse(subscriptionRequest.BaseRequest.Tag, nil))

	managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-4, unsubscriptionRequest.BaseRequest.Tag, nil))
}

// Delete project (?)

func GetProjectById(id string) (*Project, error) {
	// Get new DB connection
	session, collection := managers.GetMGoCollection("Projects")
	defer session.Close()

	result := new(Project)
	err := collection.Find(bson.M{"_id": id}).One(&result)
	if err != nil {
		managers.LogError("Failed to retrieve Project", err)
		return nil, err
	}

	return result, nil
}

func CheckUserHasPermissions(project *Project, username string, permissionsLevel int) bool {
	if (project.Permissions[username] >= permissionsLevel) {
		return true;
	}
	return false;
}
