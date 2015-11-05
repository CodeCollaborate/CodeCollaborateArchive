package projectRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/requests"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/models"
)

type ProjectGrantPermissionsRequest struct {
	BaseRequest     baseRequests.BaseRequest // BaseMessage for Tag, Resource and Method
	GrantUsername   string                   // User id
	PermissionLevel int                      // Permissions level
}

func (message *ProjectGrantPermissionsRequest) GetNotification() *baseModels.WSNotification {

	data := map[string]interface{}{
		"GrantUsername": message.GrantUsername,
		"PermissionLevel": message.PermissionLevel,
	}
	return baseModels.NewNotification(message.BaseRequest, data)
}

