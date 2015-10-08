package base

var StatusCodes = map[int]string{
	1:"Success",

	// 0s - General Errors
	-0:"Undefined Error",
	-1:"Invalid JSON Object",
	-2:"Invalid resource type",
	-3:"Invalid action",

	// 100s - User Errors
	-102:"Error creating user: Internal Error",
	-103:"Error creating user: Duplicate username",
	-104:"Error logging in: Internal Error",
	-105:"Error logging in: Invalid Username or Password",
	-106:"Invalid Token",

	// 200s - Project Errors
	-200:"Error deserializing JSON to ProjectRequest",

	// 300s - File Errors
	-300:"Error deserializing JSON to FileRequest",

	// 400s - Change Errors
	-400:"Error deserializing JSON to FileChangeRequest",
}