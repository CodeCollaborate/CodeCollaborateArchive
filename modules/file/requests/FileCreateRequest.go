package fileRequests

import "github.com/CodeCollaborate/CodeCollaborate/modules/base"

type FileCreateRequest struct {
	BaseRequest  base.BaseRequest // Add, Update, Remove
	Name         string           // Name of file
	RelativePath string           // Relative path w/in project
	ProjectId    string           // Owned by project with this id
}

func (message *FileCreateRequest) GetNotification(resourceId string) *base.WSNotification {
	notification := new(base.WSNotification)
	notification.Action = message.BaseRequest.Action
	notification.Resource = message.BaseRequest.Resource
	notification.ResId = resourceId
	notification.Data = nil
	return notification
}
