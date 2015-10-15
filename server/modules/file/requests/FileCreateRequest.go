package fileRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/models"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/requests"
)

type FileCreateRequest struct {
	BaseRequest  baseRequests.BaseRequest // Add, Update, Remove
	Name         string           // Name of file
	RelativePath string           // Relative path w/in project
	ProjectId    string           // Owned by project with this id
}

func (message *FileCreateRequest) GetNotification(resourceId string) *baseModels.WSNotification {
	notification := new(baseModels.WSNotification)
	notification.Action = message.BaseRequest.Action
	notification.Resource = message.BaseRequest.Resource
	notification.ResId = resourceId
	notification.Data = nil
	return notification
}
