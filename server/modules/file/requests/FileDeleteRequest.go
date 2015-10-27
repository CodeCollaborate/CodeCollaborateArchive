package fileRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/models"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/requests"
)

type FileDeleteRequest struct {
	BaseRequest baseRequests.BaseRequest // Add, Update, Remove
}

func (message *FileDeleteRequest) GetNotification() *baseModels.WSNotification {

	return baseModels.NewNotification(message.BaseRequest, nil)
}
