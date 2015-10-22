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
	notification := new(baseModels.WSNotification)
	notification.Action = message.BaseRequest.Action
	notification.Resource = message.BaseRequest.Resource
	notification.ResId = message.BaseRequest.ResId
	notification.Data = map[string]interface{}{
		"RevokeUserId": message.RevokeUserId,
	}
	return notification
}

