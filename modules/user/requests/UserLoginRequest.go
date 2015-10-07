package userRequests

import (
)

type UserLoginRequest struct {
	UsernameOREmail string            // Username or Email, doesn't matter
	Password string `bson:"-"` // Unhashed Password - WARNING: DO NOT SAVE OR PRINT.
}
