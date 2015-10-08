package projectRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/modules/base"
)

type ProjectGrantPermissionsRequest struct {
	BaseMessage     base.BaseRequest // BaseMessage for Tag, Resource and Method
	ProjectId       string           // Project Id
	GrantUserId     string           // User id
	PermissionLevel int              // Permissions level
}

