package projectModels

import (
	"log"

	"github.com/CodeCollaborate/CodeCollaborate/managers"
	"github.com/CodeCollaborate/CodeCollaborate/modules/base"
	"github.com/CodeCollaborate/CodeCollaborate/modules/project/requests"
	"gopkg.in/mgo.v2/bson"
)

/**
Permissions for users:
	1 - read
	2 - write
	5 - admin
	10 - owners

Admins and above can grant/revoke permissions for anyone at the same, or a lower permission level.

Accepts wildcard flag for Username: "*"
*/
type Project struct {
	Id          string         `bson:"_id"` // ID of object
	Name        string         // Name of project
	ServerPath  string         // Path on server
	Permissions map[string]int // Array of references to User objects

	//TODO: Add project versions, incremented on file creation, deletion, checked on ws connect
}

// Create new project
func CreateProject(projectCreateRequest projectRequests.ProjectCreateRequest) base.WSResponse {

	// Create new Project object
	project := new(Project)
	project.Id = managers.NewObjectIdString()
	project.Name = projectCreateRequest.Name
	project.ServerPath = project.Id
	project.Permissions = map[string]int{projectCreateRequest.BaseMessage.UserId: 10} // Set creator to owner permissions

	// Get new DB connection
	session, collection := managers.GetMGoCollection("Projects")
	defer session.Close()

	// Create the project
	err := collection.Insert(project)
	if err != nil {
		return base.NewFailResponse(-201, projectCreateRequest.BaseMessage.Tag, nil)
	}

	return base.NewSuccessResponse(projectCreateRequest.BaseMessage.Tag, map[string]interface{}{"ProjectId": project.Id})
}

// Rename project (?)
func RenameProject(projectRenameRequest projectRequests.ProjectRenameRequest) base.WSResponse {

	// Get new DB connection
	session, collection := managers.GetMGoCollection("Projects")
	defer session.Close()

	// Rename the project
	err := collection.Update(bson.M{"_id": projectRenameRequest.ProjectId}, bson.M{"$set": bson.M{"name": projectRenameRequest.NewName}})
	if err != nil {
		return base.NewFailResponse(-202, projectRenameRequest.BaseMessage.Tag, nil)
	}

	return base.NewSuccessResponse(projectRenameRequest.BaseMessage.Tag, nil)
}

// Delete project (?)

// Grant permission <Level> to <User>
//  - Check if user exists
//  - Grants permission level to user, overwriting if necessary.
func GrantProjectPermissions(projectGrantPermissionsRequest projectRequests.ProjectGrantPermissionsRequest) base.WSResponse {

	project, err := GetProjectById(projectGrantPermissionsRequest.ProjectId)
	if err != nil {
		return base.NewFailResponse(-200, projectGrantPermissionsRequest.BaseMessage.Tag, nil)
	}
	project.Permissions[projectGrantPermissionsRequest.GrantUserId] = projectGrantPermissionsRequest.PermissionLevel

	// Get new DB connection
	session, collection := managers.GetMGoCollection("Projects")
	defer session.Close()

	// Create the project
	err = collection.Update(bson.M{"_id": projectGrantPermissionsRequest.ProjectId}, bson.M{"$set": bson.M{"permissions": project.Permissions}})
	if err != nil {
		return base.NewFailResponse(-202, projectGrantPermissionsRequest.BaseMessage.Tag, nil)
	}

	return base.NewSuccessResponse(projectGrantPermissionsRequest.BaseMessage.Tag, nil)
}

// Revoke permission for <User>
//  - Check if user has permissions
//  - Revokes permissions entirely; removes entry.
func RevokeProjectPermissions(projectRevokePermissionsRequest projectRequests.ProjectRevokePermissionsRequest) base.WSResponse {

	project, err := GetProjectById(projectRevokePermissionsRequest.ProjectId)
	if err != nil {
		return base.NewFailResponse(-200, projectRevokePermissionsRequest.BaseMessage.Tag, nil)
	}

	// Make sure that there is still an owner of the project.
	owner := ""
	for key, value := range project.Permissions {
		if value == 10 && key != projectRevokePermissionsRequest.RevokeUserId {
			owner = key
		}
	}
	if owner == "" {
		return base.NewFailResponse(-205, projectRevokePermissionsRequest.BaseMessage.Tag, nil)
	}

	delete(project.Permissions, projectRevokePermissionsRequest.RevokeUserId)

	// Get new DB connection
	session, collection := managers.GetMGoCollection("Projects")
	defer session.Close()

	// Create the project
	err = collection.Update(bson.M{"_id": projectRevokePermissionsRequest.ProjectId}, bson.M{"$set": bson.M{"permissions": project.Permissions}})
	if err != nil {
		return base.NewFailResponse(-202, projectRevokePermissionsRequest.BaseMessage.Tag, nil)
	}

	return base.NewSuccessResponse(projectRevokePermissionsRequest.BaseMessage.Tag, nil)
}

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
