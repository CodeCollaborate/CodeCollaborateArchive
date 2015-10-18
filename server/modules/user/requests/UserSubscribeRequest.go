package userRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/requests"
)

type UserSubscribeRequest struct {
	BaseRequest baseRequests.BaseRequest  // BaseMessage for Tag, Resource and Method
	UsernameOREmail string            	  // Username or Email, doesn't matter
	Projects []string 						  // array of projects stored locally on the client
}