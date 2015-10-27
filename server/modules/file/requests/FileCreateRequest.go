package fileRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/models"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/requests"
)

type FileCreateRequest struct {
	BaseRequest  baseRequests.BaseRequest // Add, Update, Remove
	Name         string                   // Name of file
	RelativePath string                   // Relative path w/in project
	ProjectId    string                   // Owned by project with this id
	FileBytes    []byte                   // Bytes of the file - binary
}

func (message *FileCreateRequest) GetNotification(fileId string) *baseModels.WSNotification {

	notification := baseModels.NewNotification(message.BaseRequest, nil)
	notification.ResId = fileId;

	return notification
}
