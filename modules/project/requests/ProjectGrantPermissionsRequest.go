package projectRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/modules/base"
)

type ProjectGrantPermissionsRequest struct {
	BaseRequest base.BaseRequest // BaseMessage for Tag, Resource and Method
	GrantUserId     string           // User id
	PermissionLevel int              // Permissions level
}

