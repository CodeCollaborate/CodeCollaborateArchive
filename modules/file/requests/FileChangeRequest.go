package fileRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/modules/base"
)

type FileChangeRequest struct {
	BaseRequest base.BaseRequest // Add, Update, Remove
	FileVersion int              // Version of file to be updated
	Changes     string           // Client-Computed changes (patch).
}

func (message *FileChangeRequest) GetNotification(fileVersion int) *base.WSNotification {
	notification := new(base.WSNotification)
	notification.Action = message.BaseRequest.Action
	notification.Resource = message.BaseRequest.Resource
	notification.ResId = message.BaseRequest.ResId
	notification.Data = map[string]interface{}{
		"Changes": message.Changes,
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