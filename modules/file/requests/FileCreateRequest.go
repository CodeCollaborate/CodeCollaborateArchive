package fileRequests

import "github.com/CodeCollaborate/CodeCollaborate/modules/base"

type FileCreateRequest struct {
	BaseRequest base.BaseRequest // Add, Update, Remove
	Name         string           // Name of file
	RelativePath string           // Relative path w/in project
	ProjectId    string           // Owned by project with this id
}
