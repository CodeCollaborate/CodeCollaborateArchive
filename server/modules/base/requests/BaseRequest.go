package baseRequests

import (
	"bytes"
)

type BaseRequest struct {
	Tag      int64  // Request tag
	Action   string // Add, Update, Remove
	Resource string // Project vs file
	ResId    string // Id of resource
	Username string // Username
	Token    string // Token
}

func (message *BaseRequest) ToString() string {

	var buffer bytes.Buffer

	buffer.WriteString(message.Action)
	buffer.WriteString(" ")
	buffer.WriteString(message.Resource)
	buffer.WriteString(": ")
	buffer.WriteString(message.ResId)

	return buffer.String()
}