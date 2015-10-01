package base

import (
	"bytes"
	"strconv"
)

/************************************************
 Response Types
 ************************************************/
type WSResponse struct {
	Status  int                    // Status code; 1 = success, negative numbers indicate error code
	Tag     int64                  // Request tag
	Message string                 // Message
	Data    map[string]interface{} // Any other data.
}

func NewSuccessResponse(Tag int64, Data map[string]interface{}) WSResponse {
	return newBaseResponse(1, Tag, Data)
}

func NewFailResponse(Status int, Tag int64, Data map[string]interface{}) WSResponse {
	return newBaseResponse(Status, Tag, Data)
}

func newBaseResponse(Status int, Tag int64, Data map[string]interface{}) WSResponse {
	baseResponse := WSResponse{}
	baseResponse.Status = Status
	baseResponse.Message = StatusCodes[Status]
	baseResponse.Tag = Tag
	baseResponse.Data = Data

	return baseResponse
}

type BaseRequest struct {
	Tag      int64  // Request tag
	Action   string // Add, Update, Remove
	Resource string // Project vs file
	ResId    int64  // Id of resource
}

func (message *BaseRequest) ToString() string {

	var buffer bytes.Buffer

	buffer.WriteString(message.Action)
	buffer.WriteString(" ")
	buffer.WriteString(message.Resource)
	buffer.WriteString(": ")
	buffer.WriteString(strconv.FormatInt(message.ResId, 10))

	return buffer.String()
}

type WSNotification struct {
	Action   string                 // Add, Update, Remove
	Resource string                 // Project vs file
	ResId    int64                  // Id of resource
	Data     map[string]interface{} // Any other data
}
