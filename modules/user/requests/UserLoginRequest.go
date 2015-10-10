package userRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/modules/base"
)

type UserLoginRequest struct {
	BaseRequest base.BaseRequest  // BaseMessage for Tag, Resource and Method
	UsernameOREmail string            // Username or Email, doesn't matter
	Password string `bson:"-"` // Unhashed Password - WARNING: DO NOT SAVE OR PRINT.
}
