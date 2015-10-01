package file

import (
	"bytes"
	"github.com/CodeCollaborate/CodeCollaborate/modules/base"
)

type FileRequest struct {
	BaseMessage base.BaseRequest // Add, Update, Remove
	Changes     string           // Client-Computed changes (patch).
	CommitHash  string           // Hash of last commit (if any)
}

func (message *FileRequest) ToString() string {

	var buffer bytes.Buffer

	buffer.WriteString(message.BaseMessage.ToString())
	buffer.WriteString(" (")
	buffer.WriteString(message.CommitHash)
	buffer.WriteString(")")
	buffer.WriteString("\nChanges:\n")
	buffer.WriteString(message.Changes)

	return buffer.String()
}

func (message *FileRequest) GetNotification() *base.WSNotification {
	notification := new(base.WSNotification)
	notification.Action = message.BaseMessage.Action
	notification.Resource = message.BaseMessage.Resource
	notification.ResId = message.BaseMessage.ResId
	notification.Data = map[string]interface{}{
		"changes":message.Changes,
	}
	return notification
}
