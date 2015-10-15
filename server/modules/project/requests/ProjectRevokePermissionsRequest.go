package projectRequests

import "github.com/CodeCollaborate/CodeCollaborate/server/modules/base/requests"

type ProjectRevokePermissionsRequest struct {
	BaseRequest baseRequests.BaseRequest // BaseMessage for Tag, Resource and Method
	RevokeUserId string           // User id
}

