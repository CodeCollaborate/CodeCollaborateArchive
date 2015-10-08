package projectRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/modules/base"
)

type ProjectRenameRequest struct {
	BaseMessage base.BaseRequest // BaseMessage for Tag, Resource and Method
	ProjectId   string           // Project Id
	NewName     string           // New Project Name
}

