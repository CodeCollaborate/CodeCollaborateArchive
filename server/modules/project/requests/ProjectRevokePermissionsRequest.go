package projectRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/requests"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/models"
)

type ProjectRevokePermissionsRequest struct {
	BaseRequest baseRequests.BaseRequest // BaseMessage for Tag, Resource and Method
	RevokeUserId string           // User id
}

func (message *ProjectRevokePermissionsRequest) GetNotification() *baseModels.WSNotification {

	data := map[string]interface{}{
		"RevokeUserId": message.RevokeUserId,
	}
	return baseModels.NewNotification(message.BaseRequest, data)
}