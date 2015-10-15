package fileRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/models"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/requests"
)

type FileDeleteRequest struct {
	BaseRequest baseRequests.BaseRequest // Add, Update, Remove
}

func (message *FileDeleteRequest) GetNotification() *baseModels.WSNotification {
	notification := new(baseModels.WSNotification)
	notification.Action = message.BaseRequest.Action
	notification.Resource = message.BaseRequest.Resource
	notification.ResId = message.BaseRequest.ResId
	notification.Data = nil
	return notification
}
