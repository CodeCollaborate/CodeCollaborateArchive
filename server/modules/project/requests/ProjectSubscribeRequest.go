package projectRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/requests"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/models"
)

type ProjectSubscribeRequest struct {
	BaseRequest baseRequests.BaseRequest // BaseMessage for Tag, Resource and Method
}

func (message *ProjectSubscribeRequest) GetNotification() *baseModels.WSNotification {

	data := map[string]interface{}{
		"Username": message.BaseRequest.Username,
	}
	return baseModels.NewNotification(message.BaseRequest, data)
}
