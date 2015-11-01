package userRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/requests"
)

type UserLookupRequest struct {
	BaseRequest    baseRequests.BaseRequest // BaseMessage for Tag, Resource and Method
	LookupEmail string                   // Name of user to lookup Id for
}