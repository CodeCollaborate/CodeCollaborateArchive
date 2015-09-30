package file

import (
	"bytes"

	"github.com/CodeCollaborate/CodeCollaborate/modules/base"
)

type FileMessage struct {
	BaseMessage base.BaseMessage // Add, Update, Remove
	Changes     string           // Client-Computed changes (patch).
	CommitHash  string           // Hash of last commit (if any)
}

func (message *FileMessage) ToString() string {

	var buffer bytes.Buffer

	buffer.WriteString(message.BaseMessage.ToString())
	buffer.WriteString(" (")
	buffer.WriteString(message.CommitHash)
	buffer.WriteString(")")
	buffer.WriteString("\nChanges:\n")
	buffer.WriteString(message.Changes)

	return buffer.String()
}
