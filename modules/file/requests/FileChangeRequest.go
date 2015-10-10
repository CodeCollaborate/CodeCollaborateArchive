package fileRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/modules/base"
)

type FileChangeRequest struct {
	BaseRequest base.BaseRequest // Add, Update, Remove
	FileVersion int              // Version of file to be updated
	Changes     string           // Client-Computed changes (patch).
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
// func (message *FileChangeRequest) GetNotification() *base.WSNotification {
// 	notification := new(base.WSNotification)
// 	notification.Action = message.BaseMessage.Action
// 	notification.Resource = message.BaseMessage.Resource
// 	notification.ResId = message.BaseMessage.ResId
// 	notification.Data = map[string]interface{}{
// 		"changes": message.Changes,
// 	}
// 	return notification
// }
