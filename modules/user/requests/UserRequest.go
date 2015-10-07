package userRequests

import (
	"github.com/CodeCollaborate/CodeCollaborate/modules/base"
)

type UserRequest struct {
	BaseMessage base.BaseRequest // Add, Update, Remove
	Changes     string           // Client-Computed changes (patch).
	CommitHash  string           // Hash of last commit (if any)
}
