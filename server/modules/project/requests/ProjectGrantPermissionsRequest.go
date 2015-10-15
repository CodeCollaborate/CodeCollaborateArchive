package projectRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/requests"
)

type ProjectGrantPermissionsRequest struct {
	BaseRequest baseRequests.BaseRequest // BaseMessage for Tag, Resource and Method
	GrantUserId     string           // User id
	PermissionLevel int              // Permissions level
}

