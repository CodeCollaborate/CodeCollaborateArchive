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

	data := map[string]interface{}{
		"NewName": message.NewName,
	}
	return baseModels.NewNotification(message.BaseRequest, data)
}
