package fileRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/requests"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/models"
)

type FileMoveRequest struct {
	BaseRequest baseRequests.BaseRequest // Add, Update, Remove
	NewPath     string                   // New File Name
}

func (message *FileMoveRequest) GetNotification() *baseModels.WSNotification {

	data := map[string]interface{}{
		"NewPath": message.NewPath,
	}
	return baseModels.NewNotification(message.BaseRequest, data)
}
