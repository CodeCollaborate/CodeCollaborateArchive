package projectRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/requests"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/models"
)

type ProjectRenameRequest struct {
	BaseRequest baseRequests.BaseRequest // BaseMessage for Tag, Resource and Method
	NewName     string           // New Project Name
}

func (message *ProjectRenameRequest) GetNotification() *baseModels.WSNotification {

	data := map[string]interface{}{
		"NewName": message.NewName,
	}
	return baseModels.NewNotification(message.BaseRequest, data)
}


