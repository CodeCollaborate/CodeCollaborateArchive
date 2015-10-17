package fileModels

import (
	"log"
	"time"

	"github.com/CodeCollaborate/CodeCollaborate/server/managers"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/file/requests"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/models"
)

type FileChange struct {
	Id      string    `bson:"_id"` // ID of object
	Changes string                 // Client-Computed changes (patch).
	Version int                    // Version number
	File    string                 // id of file that was changed
	User    string                 // id of user that made the change
	Date    time.Time              // Date/Time change was made
}

func InsertChange(fileChangeRequest fileRequests.FileChangeRequest) baseModels.WSResponse {

	// Check that file exists
	file, err := GetFileById(fileChangeRequest.BaseRequest.ResId);
	if err != nil {
		return baseModels.NewFailResponse(-300, fileChangeRequest.BaseRequest.Tag, nil)
	}

	// Check that user is on latest version, then increment. Otherwise, throw error
	if (fileChangeRequest.FileVersion < file.Version) {
		return baseModels.NewFailResponse(-401, fileChangeRequest.BaseRequest.Tag, nil)
	}
	fileChangeRequest.FileVersion++

	fileChange := new(FileChange)
	fileChange.Id = managers.NewObjectIdString()
	fileChange.Changes = fileChangeRequest.Changes
	fileChange.File = fileChangeRequest.BaseRequest.ResId
	fileChange.Version = fileChangeRequest.FileVersion
	fileChange.User = fileChangeRequest.BaseRequest.UserId
	fileChange.Date = time.Now().UTC()

	changesSession, changesCollection := managers.GetMGoCollection("Changes")
	defer changesSession.Close()

	index := mgo.Index{
		Key:        []string{"file", "version"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = changesCollection.EnsureIndex(index)
	if err != nil {
		log.Println("Failed to ensure changes index:", err)
		return baseModels.NewFailResponse(-400, fileChangeRequest.BaseRequest.Tag, nil)
	}

	err = changesCollection.Insert(fileChange)
	if err != nil {
		if mgo.IsDup(err) {
			return baseModels.NewFailResponse(-401, fileChangeRequest.BaseRequest.Tag, nil)
		}
		return baseModels.NewFailResponse(-400, fileChangeRequest.BaseRequest.Tag, nil)
	}

	filesSession, filesCollection := managers.GetMGoCollection("Files")
	defer filesSession.Close()
	err = filesCollection.Update(bson.M{"_id": fileChangeRequest.BaseRequest.ResId}, bson.M{"$set": bson.M{"version": fileChangeRequest.FileVersion}})
	if err != nil {
		return baseModels.NewFailResponse(-400, fileChangeRequest.BaseRequest.Tag, nil)
	}

	managers.NotifyAll(file.Project, fileChangeRequest.GetNotification(fileChangeRequest.FileVersion))

	return baseModels.NewSuccessResponse(fileChangeRequest.BaseRequest.Tag, nil)

}

func GetChangeById(id string) (*FileChange, error) {
	// Get new DB connection
	session, collection := managers.GetMGoCollection("Changes")
	defer session.Close()

	result := new(FileChange)
	err := collection.Find(bson.M{"_id": id}).One(&result)
	if err != nil {
		log.Println("Failed to retrieve FileChange")
		log.Println(err)
		return nil, err
	}

	return result, nil
}

func GetChangesByFile(fileId string) ([]FileChange, error) {
	// Get new DB connection
	session, collection := managers.GetMGoCollection("Changes")
	defer session.Close()

	var result []FileChange
	err := collection.Find(bson.M{"file": fileId}).Sort("version").All(&result)
	if err != nil {
		log.Println("Failed to retrieve FileChanges")
		log.Println(err)
		return nil, err
	}

	return result, nil
}