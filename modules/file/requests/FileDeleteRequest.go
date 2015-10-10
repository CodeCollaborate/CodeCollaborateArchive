package fileRequests

import "github.com/CodeCollaborate/CodeCollaborate/modules/base"

type FileDeleteRequest struct {
	BaseMessage base.BaseRequest // Add, Update, Remove
	FileId      string           // Id of file
}
