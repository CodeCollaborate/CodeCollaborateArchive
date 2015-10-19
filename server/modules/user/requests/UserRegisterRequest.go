package userRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/requests"
)

type UserRegisterRequest struct {
	BaseRequest baseRequests.BaseRequest // BaseMessage for Tag, Resource and Method
	Username    string                   // Username
	FirstName   string                   // User's First name
	LastName    string                   // User's Last name
	Email       string                   // Email of user
	Password    string `bson:"-"`        // Unhashed Password - WARNING: DO NOT SAVE OR PRINT.
}
