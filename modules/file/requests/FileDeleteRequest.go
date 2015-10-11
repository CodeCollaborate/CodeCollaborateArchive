package fileRequests

import "github.com/CodeCollaborate/CodeCollaborate/modules/base"

type FileDeleteRequest struct {
	BaseRequest base.BaseRequest // Add, Update, Remove
}

func (message *FileDeleteRequest) GetNotification() *base.WSNotification {
	notification := new(base.WSNotification)
	notification.Action = message.BaseRequest.Action
	notification.Resource = message.BaseRequest.Resource
	notification.ResId = message.BaseRequest.ResId
	notification.Data = nil
	return notification
}
