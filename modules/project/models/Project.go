package project

import (
	"github.com/CodeCollaborate/CodeCollaborate/modules/user/models"
)

type Project struct {
	_id               string            // ID of object
	Name              string            // Name of project
	Server_Path       string            // Path on server
	Read_Permissions  []userModels.User // Array of references to User objects
	Write_Permissions []userModels.User // Array of references to User objects
}

/**
	Permissions for users:
		1 - read
		2 - write
		5 - admin
		10 - owners

	Admins and above can grant/revoke permissions for anyone at the same, or a lower permission level.

	Accepts wildcard flag for Username: "*"
 */
type Permission struct {
	Username        string // Username of grantee
	PermissionLevel int    // Permission level
}

// Create new project


// Rename project (?)


// Delete project (?)


// Grant permission <Level> to <User>
//  - Check if user exists
//  - Grants permission level to user, overwriting if necessary.


// Revoke permission for <User>
//  - Check if user has permissions
//  - Revokes permissions entirely; removes entry.