package userRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/requests"
)

type UserLoginRequest struct {
	BaseRequest baseRequests.BaseRequest // BaseMessage for Tag, Resource and Method
	Username    string                   // Username
	Password    string `bson:"-"`        // Unhashed Password - WARNING: DO NOT SAVE OR PRINT.
}
