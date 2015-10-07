package userRequests

import (
)

type UserRegisterRequest struct {
	Username string            // Username
	Email    string            // Email of user
	Password string `bson:"-"` // Unhashed Password - WARNING: DO NOT SAVE OR PRINT.
}
