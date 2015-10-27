package baseModels

import "github.com/CodeCollaborate/CodeCollaborate/server/modules/base"

/************************************************
 Response Types
 ************************************************/
type WSResponse struct {
	Type    string                 // Response type
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
	baseResponse.Type = "Response"
	baseResponse.Status = Status
	baseResponse.Message = base.StatusCodes[Status]
	baseResponse.Tag = Tag
	baseResponse.Data = Data

	return baseResponse
}