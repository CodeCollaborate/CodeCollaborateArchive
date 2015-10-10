package fileModels

import (
	"log"
	"time"

	"github.com/CodeCollaborate/CodeCollaborate/managers"
	"github.com/CodeCollaborate/CodeCollaborate/modules/base"
	"github.com/CodeCollaborate/CodeCollaborate/modules/file/requests"
	"gopkg.in/mgo.v2"
)

type FileChange struct {
	Id      string    `bson:"_id"` // ID of object
	Changes string    // Client-Computed changes (patch).
	Version int       // Version number
	File    string    // id of file that was changed
	User    string    // id of user that made the change
	Date    time.Time // Date/Time change was made
}

func InsertChange(fileChangeRequest fileRequests.FileChangeRequest) base.WSResponse {

	fileChange := new(FileChange)
	fileChange.Id = managers.NewObjectIdString()
	fileChange.Changes = fileChangeRequest.Changes
	fileChange.File = fileChangeRequest.BaseRequest.ResId
	fileChange.Version = fileChangeRequest.FileVersion
	fileChange.User = fileChangeRequest.BaseRequest.UserId
	fileChange.Date = time.Now().UTC()

	session, collection := managers.GetMGoCollection("Changes")
	defer session.Close()

	index := mgo.Index{
		Key:        []string{"file", "version"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := collection.EnsureIndex(index)
	if err != nil {
		log.Println("Failed to ensure changes index:", err)
		return base.NewFailResponse(-400, fileChangeRequest.BaseRequest.Tag, nil)
	}

	err = collection.Insert(fileChange)
	if err != nil {
		if mgo.IsDup(err) {
			return base.NewFailResponse(-401, fileChangeRequest.BaseRequest.Tag, nil)
		}
		return base.NewFailResponse(-400, fileChangeRequest.BaseRequest.Tag, nil)
	}

	return base.NewSuccessResponse(fileChangeRequest.BaseRequest.Tag, nil)

}
