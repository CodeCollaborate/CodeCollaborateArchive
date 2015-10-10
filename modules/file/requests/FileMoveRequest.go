package fileRequests

import "github.com/CodeCollaborate/CodeCollaborate/modules/base"

type FileMoveRequest struct {
	BaseMessage base.BaseRequest // Add, Update, Remove
	FileId      string           // Id of file
	NewPath     string           // New File Name
}
