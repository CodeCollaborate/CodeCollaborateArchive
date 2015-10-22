package projectModels

import (
	"log"

	"github.com/CodeCollaborate/CodeCollaborate/server/managers"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/project/requests"
	"gopkg.in/mgo.v2/bson"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/models"
)

/**
Permissions for users:
	1 - read
	2 - write
	5 - admin
	10 - owners

Admins and above can grant/revoke permissions for anyone at the same, or a lower permission level.

TODO: Accepts wildcard flag for Username: "*"
*/
type Project struct {
	Id          string         `bson:"_id"` // ID of object
	Name        string         // Name of project
	ServerPath  string         // Path on server
	Permissions map[string]int // Array of references to User objects

	//TODO: Add project versions, incremented on file creation, deletion, checked on ws connect
	//TODO: wildcard permissions, add once we make adding to projects a thing
}

// Create new project
func CreateProject(projectCreateRequest projectRequests.ProjectCreateRequest) baseModels.WSResponse {

	// Create new Project object
	project := new(Project)
	project.Id = managers.NewObjectIdString()
	project.Name = projectCreateRequest.Name
	project.ServerPath = project.Id
	project.Permissions = map[string]int{projectCreateRequest.BaseRequest.UserId: 10} // Set creator to owner permissions

	// Get new DB connection
	session, collection := managers.GetMGoCollection("Projects")
	defer session.Close()

	// Create the project
	err := collection.Insert(project)
	if err != nil {
		return baseModels.NewFailResponse(-201, projectCreateRequest.BaseRequest.Tag, nil)
	}

	return baseModels.NewSuccessResponse(projectCreateRequest.BaseRequest.Tag, map[string]interface{}{"ProjectId": project.Id})
}

// Rename project (?)
func RenameProject(projectRenameRequest projectRequests.ProjectRenameRequest) baseModels.WSResponse {

	// Get new DB connection
	session, collection := managers.GetMGoCollection("Projects")
	defer session.Close()

	// Rename the project
	err := collection.Update(bson.M{"_id": projectRenameRequest.BaseRequest.ResId}, bson.M{"$set": bson.M{"name": projectRenameRequest.NewName}})
	if err != nil {
		return baseModels.NewFailResponse(-202, projectRenameRequest.BaseRequest.Tag, nil)
	}

	managers.NotifyProjectClients(projectRenameRequest.BaseRequest.ResId, projectRenameRequest.GetNotification())

	return baseModels.NewSuccessResponse(projectRenameRequest.BaseRequest.Tag, nil)
}

// Grant permission <Level> to <User>
//  - Check if user exists
//  - Grants permission level to user, overwriting if necessary.
func GrantProjectPermissions(projectGrantPermissionsRequest projectRequests.ProjectGrantPermissionsRequest) baseModels.WSResponse {

	project, err := GetProjectById(projectGrantPermissionsRequest.BaseRequest.ResId)
	if err != nil {
		return baseModels.NewFailResponse(-200, projectGrantPermissionsRequest.BaseRequest.Tag, nil)
	}

	if(!CheckUserHasPermissions(project, projectGrantPermissionsRequest.BaseRequest.UserId, 5)){
		return baseModels.NewFailResponse(-207, projectGrantPermissionsRequest.BaseRequest.Tag, nil )
	}

	// Make sure that there is still an owner of the project.
	owner := ""
	for key, value := range project.Permissions {
		if value == 10 && key != projectGrantPermissionsRequest.GrantUserId {
			owner = key
		}
	}
	if owner == "" {
		return baseModels.NewFailResponse(-205, projectGrantPermissionsRequest.BaseRequest.Tag, nil)
	}

	project.Permissions[projectGrantPermissionsRequest.GrantUserId] = projectGrantPermissionsRequest.PermissionLevel

	// Get new DB connection
	session, collection := managers.GetMGoCollection("Projects")
	defer session.Close()

	// Update permissions
	err = collection.Update(bson.M{"_id": projectGrantPermissionsRequest.BaseRequest.ResId}, bson.M{"$set": bson.M{"permissions": project.Permissions}})
	if err != nil {
		return baseModels.NewFailResponse(-202, projectGrantPermissionsRequest.BaseRequest.Tag, nil)
	}

	managers.NotifyProjectClients(projectGrantPermissionsRequest.BaseRequest.ResId, projectGrantPermissionsRequest.GetNotification())

	return baseModels.NewSuccessResponse(projectGrantPermissionsRequest.BaseRequest.Tag, nil)
}

// Revoke permission for <User>
//  - Check if user has permissions
//  - Revokes permissions entirely; removes entry.
func RevokeProjectPermissions(projectRevokePermissionsRequest projectRequests.ProjectRevokePermissionsRequest) baseModels.WSResponse {

	project, err := GetProjectById(projectRevokePermissionsRequest.BaseRequest.ResId)
	if err != nil {
		return baseModels.NewFailResponse(-200, projectRevokePermissionsRequest.BaseRequest.Tag, nil)
	}

	if(!CheckUserHasPermissions(project, projectRevokePermissionsRequest.BaseRequest.UserId, 5)){
		return baseModels.NewFailResponse(-207, projectRevokePermissionsRequest.BaseRequest.Tag, nil )
	}

	// Make sure that there is still an owner of the project.
	owner := ""
	for key, value := range project.Permissions {
		if value == 10 && key != projectRevokePermissionsRequest.RevokeUserId {
			owner = key
		}
	}
	if owner == "" {
		return baseModels.NewFailResponse(-205, projectRevokePermissionsRequest.BaseRequest.Tag, nil)
	}

	delete(project.Permissions, projectRevokePermissionsRequest.RevokeUserId)

	// Get new DB connection
	session, collection := managers.GetMGoCollection("Projects")
	defer session.Close()

	// Update permissions
	err = collection.Update(bson.M{"_id": projectRevokePermissionsRequest.BaseRequest.ResId}, bson.M{"$set": bson.M{"permissions": project.Permissions}})
	if err != nil {
		return baseModels.NewFailResponse(-202, projectRevokePermissionsRequest.BaseRequest.Tag, nil)
	}

	managers.NotifyProjectClients(projectRevokePermissionsRequest.BaseRequest.ResId, projectRevokePermissionsRequest.GetNotification())

	return baseModels.NewSuccessResponse(projectRevokePermissionsRequest.BaseRequest.Tag, nil)
}

// Delete project (?)

func GetProjectById(id string) (*Project, error) {
	// Get new DB connection
	session, collection := managers.GetMGoCollection("Projects")
	defer session.Close()

	result := new(Project)
	err := collection.Find(bson.M{"_id": id}).One(&result)
	if err != nil {
		log.Println("Failed to retrieve Project")
		log.Println(err)
		return nil, err
	}

	return result, nil
}

func CheckUserHasPermissions( project *Project, userId string, permissionsLevel int) bool{
	if(project.Permissions[userId] >= permissionsLevel){
		return true;
	}
	return false;
}
