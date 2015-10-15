package userRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/requests"
)

type UserLoginRequest struct {
	BaseRequest baseRequests.BaseRequest  // BaseMessage for Tag, Resource and Method
	UsernameOREmail string            // Username or Email, doesn't matter
	Password string `bson:"-"` // Unhashed Password - WARNING: DO NOT SAVE OR PRINT.
}
