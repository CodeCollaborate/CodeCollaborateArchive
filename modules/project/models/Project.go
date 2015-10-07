package project

import (
	"github.com/CodeCollaborate/CodeCollaborate/modules/user/models"
)

type Project struct {
	_id               string // ID of object
	Name              string // Name of project
	Server_Path       string // Path on server
	Read_Permissions  []userModels.User    // Array of references to User objects
	Write_Permissions []userModels.User // Array of references to User objects
}
