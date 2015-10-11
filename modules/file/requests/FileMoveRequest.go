package fileRequests

import "github.com/CodeCollaborate/CodeCollaborate/modules/base"

type FileMoveRequest struct {
	BaseRequest base.BaseRequest // Add, Update, Remove
	NewPath     string           // New File Name
}

func (message *FileMoveRequest) GetNotification() *base.WSNotification {
	notification := new(base.WSNotification)
	notification.Action = message.BaseRequest.Action
	notification.Resource = message.BaseRequest.Resource
	notification.ResId = message.BaseRequest.ResId
	notification.Data = map[string]interface{}{
		"NewPath": message.NewPath,
	}
	return notification
}
