package base

import (
	"bytes"
	"strconv"
)

/************************************************
 Response Types
 ************************************************/
type WSResponse struct {
	Status int         // Status code; 1 = success, negative numbers indicate error code
	Tag    int64       // Request tag
	Data   interface{} // Any other data.
}

func NewSuccessResponse(Tag int64, Data interface{}) WSResponse {
	return newBaseResponse(1, Tag, Data)
}

func NewFailResponse(Status int, Tag int64, Data interface{}) WSResponse {
	return newBaseResponse(Status, Tag, Data)
}

func newBaseResponse(Status int, Tag int64, Data interface{}) WSResponse {
	baseResponse := WSResponse{}
	baseResponse.Status = Status
	baseResponse.Tag = Tag
	baseResponse.Data = Data

	return baseResponse
}

type BaseMessage struct {
	Tag      int64  // Request tag
	Action   string // Add, Update, Remove
	Resource string // Project vs file
	ResId    int64  // Id of resource
}

func (message *BaseMessage) ToString() string {

	var buffer bytes.Buffer

	buffer.WriteString(message.Action)
	buffer.WriteString(" ")
	buffer.WriteString(message.Resource)
	buffer.WriteString(": ")
	buffer.WriteString(strconv.FormatInt(message.ResId, 10))

	return buffer.String()
}
