package fileModels

import (
	"github.com/CodeCollaborate/CodeCollaborate/modules/base"
	"github.com/CodeCollaborate/CodeCollaborate/modules/user/models"
)

type FileContentChange struct {
	Id      string  `bson:"_id"` // ID of object
	Changes string               // Client-Computed changes (patch).
	Version int                  // Version number
	File    File                 // Reference to File object
	User    userModels.User      // Reference to User object
}

func (fileContentChange *FileContentChange) GetNotification() *base.WSNotification {
	notification := new(base.WSNotification)
	notification.Action = "Update"
	notification.Resource = "File"
	notification.ResId = fileContentChange.File.Id
	notification.Data = map[string]interface{}{
		"changes":fileContentChange.Changes,
	}
	return notification
}
