package projectRequests

import "github.com/CodeCollaborate/CodeCollaborate/server/modules/base/requests"

type ProjectRenameRequest struct {
	BaseRequest baseRequests.BaseRequest // BaseMessage for Tag, Resource and Method
	NewName     string           // New Project Name
}

