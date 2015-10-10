package projectRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/modules/base"
)

type ProjectRenameRequest struct {
	BaseRequest base.BaseRequest // BaseMessage for Tag, Resource and Method
	NewName     string           // New Project Name
}

