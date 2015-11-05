package userRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/requests"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/models"
)

type UserSubscribeRequest struct {
	BaseRequest baseRequests.BaseRequest // BaseMessage for Tag, Resource and Method
	Projects    []string                 // array of projects stored locally on the client
}

func (message *UserSubscribeRequest) GetNotification() *baseModels.WSNotification {

	data := map[string]interface{}{
		"Username": message.BaseRequest.Username,
	}
	return baseModels.NewNotification(message.BaseRequest, data)
}
