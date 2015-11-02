package fileModels

import (
	"log"
	"time"

	"github.com/CodeCollaborate/CodeCollaborate/server/managers"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/file/requests"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/models"
	"github.com/gorilla/websocket"
)

type FileChange struct {
	Id       string    `bson:"_id"` // ID of object
	Changes  string                 // Client-Computed changes (patch).
	Version  int64                  // Version number
	FileId   string                 // id of file that was changed
	Username string                 // id of user that made the change
	Date     time.Time              // Date/Time change was made
}

func InsertChange(wsConn *websocket.Conn, fileChangeRequest fileRequests.FileChangeRequest) {

	// Check that file exists
	file, err := GetFileById(fileChangeRequest.BaseRequest.ResId);
	if err != nil {
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-300, fileChangeRequest.BaseRequest.Tag, nil))
		return
	}

	// Check that user is on latest version, then increment. Otherwise, throw error
	if (fileChangeRequest.FileVersion < file.Version) {
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-401, fileChangeRequest.BaseRequest.Tag, nil))
		return
	}
	fileChangeRequest.FileVersion++

	fileChange := new(FileChange)
	fileChange.Id = managers.NewObjectIdString()
	fileChange.Changes = fileChangeRequest.Changes
	fileChange.FileId = fileChangeRequest.BaseRequest.ResId
	fileChange.Version = fileChangeRequest.FileVersion
	fileChange.Username = fileChangeRequest.BaseRequest.Username
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
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-400, fileChangeRequest.BaseRequest.Tag, nil))
	}

	err = changesCollection.Insert(fileChange)
	if err != nil {
		if mgo.IsDup(err) {
			managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-401, fileChangeRequest.BaseRequest.Tag, nil))
			return
		}
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-400, fileChangeRequest.BaseRequest.Tag, nil))
		return
	}

	filesSession, filesCollection := managers.GetMGoCollection("Files")
	defer filesSession.Close()
	err = filesCollection.Update(bson.M{"_id": fileChangeRequest.BaseRequest.ResId}, bson.M{"$set": bson.M{"version": fileChangeRequest.FileVersion}})
	if err != nil {
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-400, fileChangeRequest.BaseRequest.Tag, nil))
		return
	}

	managers.SendWebSocketMessage(wsConn, baseModels.NewSuccessResponse(fileChangeRequest.BaseRequest.Tag, map[string]interface{}{"FileVersion": fileChangeRequest.FileVersion}))
	managers.NotifyProjectClients(file.Project, fileChangeRequest.GetNotification(fileChangeRequest.FileVersion), wsConn)
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
