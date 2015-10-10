package fileRequests

import "github.com/CodeCollaborate/CodeCollaborate/modules/base"

type FileRenameRequest struct {
	BaseRequest base.BaseRequest // Add, Update, Remove
	NewFileName string           // New File Name
}
