package fileRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/requests"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/models"
)

type FileRenameRequest struct {
	BaseRequest baseRequests.BaseRequest // Add, Update, Remove
	NewName     string                   // New File Name
}

func (message *FileRenameRequest) GetNotification() *baseModels.WSNotification {
	notification := new(baseModels.WSNotification)
	notification.Action = message.BaseRequest.Action
	notification.Resource = message.BaseRequest.Resource
	notification.ResId = message.BaseRequest.ResId
	notification.Data = map[string]interface{}{
		"NewName": message.NewName,
	}
	return notification
}
