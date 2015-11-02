package fileRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/requests"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/models"
)

type FileChangeRequest struct {
	BaseRequest baseRequests.BaseRequest // Add, Update, Remove
	FileVersion int64                    // Version of file to be updated
	Changes     string                   // Client-Computed changes (patch).
}

func (message *FileChangeRequest) GetNotification(fileVersion int64) *baseModels.WSNotification {

	data := map[string]interface{}{
		"Changes": message.Changes,
		"FileVersion": fileVersion,
	}
	return baseModels.NewNotification(message.BaseRequest, data)
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