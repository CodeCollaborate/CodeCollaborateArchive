package projectRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/requests"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/models"
)

type ProjectRevokePermissionsRequest struct {
	BaseRequest baseRequests.BaseRequest // BaseMessage for Tag, Resource and Method
	RevokeUsername string           // User id
}

func (message *ProjectRevokePermissionsRequest) GetNotification() *baseModels.WSNotification {

	data := map[string]interface{}{
		"RevokeUserEmail": message.RevokeUsername,
	}
	return baseModels.NewNotification(message.BaseRequest, data)
}