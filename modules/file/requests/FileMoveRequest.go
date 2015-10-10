package fileRequests

import "github.com/CodeCollaborate/CodeCollaborate/modules/base"

type FileMoveRequest struct {
	BaseRequest base.BaseRequest // Add, Update, Remove
	NewPath     string           // New File Name
}
