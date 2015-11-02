package fileModels

import (
	"github.com/CodeCollaborate/CodeCollaborate/server/managers"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/models"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/file/requests"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"os"
	"path/filepath"
	"github.com/gorilla/websocket"
)


/*
 *
 * TODO: Check permissions before running any of these!
 *
 */
type File struct {
	Id           string `bson:"_id"`           // ID of object
	Name         string                        // Name of file
	RelativePath string `bson:"relative_path"` // Path of file
	Version      int64                         // File version
	Project      string                        // Reference to Project object
	filePath     string `bson:"-",json:"-"`    // Temp filepath cached variable
}

func (file File) getPath() string {
	if (file.filePath == "") {
		// Change to use byte buffer for efficiency
		file.filePath = "files/" + file.Project + "/" + file.RelativePath + file.Name
	}
	return file.filePath
}

func CreateFile(wsConn *websocket.Conn, fileCreateRequest fileRequests.FileCreateRequest) {

	file := new(File)
	file.Id = managers.NewObjectIdString()
	file.Name = fileCreateRequest.Name
	file.RelativePath = fileCreateRequest.RelativePath
	file.Version = 0
	file.Project = fileCreateRequest.ProjectId

	session, collection := managers.GetMGoCollection("Files")
	defer session.Close()

	// Create indexes
	index := mgo.Index{
		Key:        []string{"name", "relative_path"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := collection.EnsureIndex(index)
	if err != nil {
		managers.LogError("Failed to ensure file name/path index", err)
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-301, fileCreateRequest.BaseRequest.Tag, nil))
		return
	}

	// Insert file record
	err = collection.Insert(file)
	if err != nil {
		if mgo.IsDup(err) {
			managers.LogError("Error creating file record", err)
			managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-305, fileCreateRequest.BaseRequest.Tag, nil))
			return
		}
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-301, fileCreateRequest.BaseRequest.Tag, nil))
		return
	}

	// Write file to disk

	fileCreateRequest.RelativePath = filepath.Clean(fileCreateRequest.RelativePath)
	if (fileCreateRequest.RelativePath[0:2] == ".." || filepath.IsAbs(fileCreateRequest.RelativePath)) {
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-308, fileCreateRequest.BaseRequest.Tag, nil))
		return
	}

	err = os.MkdirAll("files/" + fileCreateRequest.ProjectId + "/" + fileCreateRequest.RelativePath, os.ModeExclusive)
	if err != nil {
		managers.LogError("Failed to create file directory", err)
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-301, fileCreateRequest.BaseRequest.Tag, nil))
		return
	}
	err = ioutil.WriteFile(file.getPath(), fileCreateRequest.FileBytes, os.ModeExclusive)
	if err != nil {
		managers.LogError("Failed to write file", err)
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-301, fileCreateRequest.BaseRequest.Tag, nil))
		return
	}

	managers.SendWebSocketMessage(wsConn, baseModels.NewSuccessResponse(fileCreateRequest.BaseRequest.Tag, map[string]interface{}{"FileId": file.Id}))
	managers.NotifyProjectClients(file.Project, fileCreateRequest.GetNotification(file.Id), wsConn)
}

func RenameFile(wsConn *websocket.Conn, fileRenameRequest fileRequests.FileRenameRequest) {
	session, collection := managers.GetMGoCollection("Files")
	defer session.Close()

	// Check that file exists
	file, err := GetFileById(fileRenameRequest.BaseRequest.ResId);
	if err != nil {
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-300, fileRenameRequest.BaseRequest.Tag, nil))
		return
	}

	file.Version++;

	err = collection.Update(bson.M{"_id": fileRenameRequest.BaseRequest.ResId}, bson.M{"$set": bson.M{"name": fileRenameRequest.NewName, "version": file.Version}})
	if err != nil {
		if mgo.IsDup(err) {
			managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-306, fileRenameRequest.BaseRequest.Tag, nil))
			return
		}
		managers.LogError("Error renaming file", err)
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-302, fileRenameRequest.BaseRequest.Tag, nil))
		return
	}

	managers.SendWebSocketMessage(wsConn, baseModels.NewSuccessResponse(fileRenameRequest.BaseRequest.Tag, nil))
	managers.NotifyProjectClients(file.Project, fileRenameRequest.GetNotification(), wsConn)
}

func MoveFile(wsConn *websocket.Conn, fileMoveRequest fileRequests.FileMoveRequest) {
	session, collection := managers.GetMGoCollection("Files")
	defer session.Close()

	// Check that file exists
	file, err := GetFileById(fileMoveRequest.BaseRequest.ResId);
	if err != nil {
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-300, fileMoveRequest.BaseRequest.Tag, nil))
		return
	}

	file.Version++;

	err = collection.Update(bson.M{"_id": fileMoveRequest.BaseRequest.ResId}, bson.M{"$set": bson.M{"relative_path": fileMoveRequest.NewPath, "version": file.Version}})
	if err != nil {
		if mgo.IsDup(err) {
			managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-307, fileMoveRequest.BaseRequest.Tag, nil))
			return
		}
		managers.LogError("Error renaming file", err)
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-303, fileMoveRequest.BaseRequest.Tag, nil))
		return
	}

	managers.SendWebSocketMessage(wsConn, baseModels.NewSuccessResponse(fileMoveRequest.BaseRequest.Tag, nil))
	managers.NotifyProjectClients(file.Project, fileMoveRequest.GetNotification(), wsConn)
}

func DeleteFile(wsConn *websocket.Conn, fileDeleteRequest fileRequests.FileDeleteRequest) {
	session, collection := managers.GetMGoCollection("Files")
	defer session.Close()

	// Check that file exists
	file, err := GetFileById(fileDeleteRequest.BaseRequest.ResId);
	if err != nil {
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-300, fileDeleteRequest.BaseRequest.Tag, nil))
		return
	}

	err = os.Remove(file.getPath())

	err = collection.Remove(bson.M{"_id": fileDeleteRequest.BaseRequest.ResId})
	if err != nil {
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-304, fileDeleteRequest.BaseRequest.Tag, nil))
		return
	}

	managers.SendWebSocketMessage(wsConn, baseModels.NewSuccessResponse(fileDeleteRequest.BaseRequest.Tag, nil))
	managers.NotifyProjectClients(file.Project, fileDeleteRequest.GetNotification(), wsConn)
}

func PullFile(wsConn *websocket.Conn, filePullRequest fileRequests.FilePullRequest) {

	// Check that file exists
	file, err := GetFileById(filePullRequest.BaseRequest.ResId);
	if err != nil {
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-300, filePullRequest.BaseRequest.Tag, nil))
		return
	}

	// Read file from disk
	if _, err := os.Stat(file.getPath()); os.IsNotExist(err) {
		managers.SendWebSocketMessage(wsConn, baseModels.NewSuccessResponse(filePullRequest.BaseRequest.Tag, map[string]interface{}{"FileBytes": "", "Changes": ""}))
		return
	}
	fileBytes, err := ioutil.ReadFile(file.getPath())
	if err != nil {
		managers.LogError("Failed to read from file", err)
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-301, filePullRequest.BaseRequest.Tag, nil))
		return
	}

	changes, err := GetChangesByFile(file.Id)
	if err != nil {
		managers.LogError("Failed to retrieve changes", err)
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-402, filePullRequest.BaseRequest.Tag, nil))
		return
	}

	managers.SendWebSocketMessage(wsConn, baseModels.NewSuccessResponse(filePullRequest.BaseRequest.Tag, map[string]interface{}{"FileBytes": fileBytes, "Changes": changes}))
}

func GetFileById(id string) (*File, error) {
	// Get new DB connection
	session, collection := managers.GetMGoCollection("Files")
	defer session.Close()

	result := new(File)
	err := collection.Find(bson.M{"_id": id}).One(&result)
	if err != nil {
		managers.LogError("Failed to retrieve File", err)
		return nil, err
	}

	return result, nil
}

//func GetFileByPathName(path string, name string) (*File, error) {
//	// Get new DB connection
//	session, collection := managers.GetMGoCollection("Files")
//	defer session.Close()
//
//	result := new(File)
//	err := collection.Find(bson.M{"$or": []interface{}{bson.M{"relative_path": path}, bson.M{"name": name}}}).One(&result)
//	if err != nil {
//		log.Println("Failed to retrieve File")
//		log.Println(err)
//		return nil, err
//	}
//
//	return result, nil
//}

