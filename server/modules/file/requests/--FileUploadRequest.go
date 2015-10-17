package fileRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/requests"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/models"
)

type FileUploadRequest struct {
	BaseRequest baseRequests.BaseRequest // Add, Update, Remove
	FileBytes   []byte                   // New File Name

	// TEMPORARY TESTING VARIABLE - PULL FROM FILE LATER
	FileName    string
}

func (message *FileUploadRequest) GetNotification() *baseModels.WSNotification {
	notification := new(baseModels.WSNotification)
	notification.Action = message.BaseRequest.Action
	notification.Resource = message.BaseRequest.Resource
	notification.ResId = message.BaseRequest.ResId
	notification.Data = nil
	return notification
}
