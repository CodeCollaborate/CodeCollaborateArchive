package fileRequests

import "github.com/CodeCollaborate/CodeCollaborate/modules/base"

type FileRenameRequest struct {
	BaseRequest base.BaseRequest // Add, Update, Remove
	NewName string           // New File Name
}

func (message *FileRenameRequest) GetNotification() *base.WSNotification {
	notification := new(base.WSNotification)
	notification.Action = message.BaseRequest.Action
	notification.Resource = message.BaseRequest.Resource
	notification.ResId = message.BaseRequest.ResId
	notification.Data = map[string]interface{}{
		"NewName": message.NewName,
	}
	return notification
}
