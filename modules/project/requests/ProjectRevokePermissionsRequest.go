package projectRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/modules/base"
)

type ProjectRevokePermissionsRequest struct {
	BaseRequest base.BaseRequest // BaseMessage for Tag, Resource and Method
	RevokeUserId string           // User id
}

