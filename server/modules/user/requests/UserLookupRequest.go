package userRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/requests"
)

type UserLookupRequest struct {
	BaseRequest    baseRequests.BaseRequest // BaseMessage for Tag, Resource and Method
	LookupUsername string                   // username to lookup
}