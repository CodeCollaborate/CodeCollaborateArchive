package fileRequests

import "github.com/CodeCollaborate/CodeCollaborate/modules/base"

type FileRenameRequest struct {
	BaseMessage base.BaseRequest // Add, Update, Remove
	FileId      string           // Id of file
	NewFileName string           // New File Name
}
