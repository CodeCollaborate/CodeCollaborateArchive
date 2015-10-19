package base

var StatusCodes = map[int]string{
	1: "Success",

	// 0s - General Errors
	-0: "Undefined Error",
	-1: "Invalid JSON Object",
	-2: "Invalid resource type",
	-3: "Invalid action",

	// 100s - User Errors
	-100: "No such user found",
	-101: "Error creating user: Internal Error",
	-102: "Error creating user: Duplicate username",
	-103: "Error logging in: Internal Error",
	-104: "Error logging in: Invalid Username or Password",
	-105: "Invalid Token",

	// 200s - Project Errors
	-200: "No such project found",
	-201: "Error creating Project: Internal Error",
	-202: "Error renaming Project: Internal Error",
	-203: "Error granting permissions: Internal Error",
	-204: "Error revoking permissions: Internal Error",
	-205: "Error revoking permissions: Must have an owner",
	-206: "Error subscribing to project",

	// 300s - File Errors
	-300: "No such file found",
	-301: "Error creating File: Internal Error",
	-302: "Error renaming File: Internal Error",
	-303: "Error moving File: Internal Error",
	-304: "Error deleting File: Internal Error",
	-305: "Error creating File: Duplicate File",
	-306: "Error renaming File: Duplicate File",
	-307: "Error moving File: Duplicate File",
	-308: "Error creating File: Invalid File Path",

	// 400s - Change Errors
	-400: "Error inserting Change: Internal Error",
	-401: "Error inserting Change: Duplicate Version Number",
	-402: "Error reading Change: Internal Error",
	// -420:"Error, too blazed"
}
