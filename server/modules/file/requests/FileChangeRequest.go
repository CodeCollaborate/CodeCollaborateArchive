package fileRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/requests"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/models"
)

type FileChangeRequest struct {
	BaseRequest baseRequests.BaseRequest // Add, Update, Remove
	FileVersion int                      // Version of file to be updated
	Change     string                   // Client-Computed changes (patch).
}

func (message *FileChangeRequest) GetNotification(fileVersion int) *baseModels.WSNotification {
	notification := new(baseModels.WSNotification)
	notification.Action = message.BaseRequest.Action
	notification.Resource = message.BaseRequest.Resource
	notification.ResId = message.BaseRequest.ResId
	notification.Data = map[string]interface{}{
		"Changes": message.Change,
		"FileVersion": fileVersion,
	}
	return notification
}


//
// func (message *FileChangeRequest) ToString() string {
//
// 	var buffer bytes.Buffer
//
// 	buffer.WriteString(message.BaseMessage.ToString())
// 	buffer.WriteString("\nChanges:\n")
// 	buffer.WriteString(message.Changes)
//
// 	return buffer.String()
// }
//