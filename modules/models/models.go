package models

import (
	"bytes"
	"strconv"
)

type Message struct {
	Action     string // Add, Update, Remove
	Resource   string // Project vs file
	Id         int64  // Id of Resource
	CommitHash string // Hash of last commit (if any)
}

func (message *Message) ToString() string {

	var buffer bytes.Buffer

	buffer.WriteString(message.Action)
	buffer.WriteString(" ")
	buffer.WriteString(message.Resource)
	buffer.WriteString(": ")
	buffer.WriteString(strconv.FormatInt(message.Id, 10))
	buffer.WriteString(" from commit ")
	buffer.WriteString(message.CommitHash)

	return buffer.String()
}