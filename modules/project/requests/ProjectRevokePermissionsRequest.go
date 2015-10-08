package projectRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/modules/base"
)

type ProjectRevokePermissionsRequest struct {
	BaseMessage  base.BaseRequest // BaseMessage for Tag, Resource and Method
	ProjectId    string           // Project Id
	RevokeUserId string           // User id
}

