package projectRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/modules/base"
)

type ProjectCreateRequest struct {
	BaseMessage base.BaseRequest  // BaseMessage for Tag, Resource and Method
	Name        string            // Project name
}

