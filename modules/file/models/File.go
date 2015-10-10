package fileModels

import (
	"log"

	"github.com/CodeCollaborate/CodeCollaborate/managers"
	"github.com/CodeCollaborate/CodeCollaborate/modules/base"
	"github.com/CodeCollaborate/CodeCollaborate/modules/file/requests"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type File struct {
	Id           string `bson:"_id"` // ID of object
	Name         string // Name of file
	RelativePath string `bson:"relative_path"` // Path of file
	Version      int    // File version
	Project      string // Reference to Project object
}

func CreateFile(fileCreateRequest fileRequests.FileCreateRequest) base.WSResponse {

	file := new(File)
	file.Id = managers.NewObjectIdString()
	file.Name = fileCreateRequest.Name
	file.RelativePath = fileCreateRequest.RelativePath
	file.Version = 0
	file.Project = fileCreateRequest.ProjectId

	session, collection := managers.GetMGoCollection("Files")
	defer session.Close()

	index := mgo.Index{
		Key:        []string{"name", "relative_path"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := collection.EnsureIndex(index)
	if err != nil {
		log.Println("Failed to ensure username index:", err)
		return base.NewFailResponse(-301, fileCreateRequest.BaseRequest.Tag, nil)
	}

	err = collection.Insert(file)
	if err != nil {
		if mgo.IsDup(err) {
			log.Println("Error registering user:", err)
			return base.NewFailResponse(-305, fileCreateRequest.BaseRequest.Tag, nil)
		}
		return base.NewFailResponse(-301, fileCreateRequest.BaseRequest.Tag, nil)
	}

	return base.NewSuccessResponse(fileCreateRequest.BaseRequest.Tag, map[string]interface{}{"FileId": file.Id})

}

func RenameFile(fileRenameRequest fileRequests.FileRenameRequest) base.WSResponse {
	session, collection := managers.GetMGoCollection("Files")
	defer session.Close()

	// Check that file exists
	// Check that new path/file does not exist - or check if it becomes a duplicate?

	err := collection.Update(bson.M{"_id": fileRenameRequest.BaseRequest.ResId}, bson.M{"$set": bson.M{"name": fileRenameRequest.NewFileName}})
	if err != nil {
		return base.NewFailResponse(-302, fileRenameRequest.BaseRequest.Tag, nil)
	}

	return base.NewSuccessResponse(fileRenameRequest.BaseRequest.Tag, nil)
}

func MoveFile(fileMoveRequest fileRequests.FileMoveRequest) base.WSResponse {
	session, collection := managers.GetMGoCollection("Files")
	defer session.Close()

	// Check that file exists
	// Check that new path/file does not exist - or check if it becomes a duplicate?

	err := collection.Update(bson.M{"_id": fileMoveRequest.BaseRequest.ResId}, bson.M{"$set": bson.M{"relative_path": fileMoveRequest.NewPath}})
	if err != nil {
		return base.NewFailResponse(-303, fileMoveRequest.BaseRequest.Tag, nil)
	}

	return base.NewSuccessResponse(fileMoveRequest.BaseRequest.Tag, nil)
}

func DeleteFile(fileDeleteRequest fileRequests.FileDeleteRequest) base.WSResponse {
	session, collection := managers.GetMGoCollection("Files")
	defer session.Close()

	// Check that file exists

	err := collection.Remove(bson.M{"_id": fileDeleteRequest.BaseRequest.ResId})
	if err != nil {
		return base.NewFailResponse(-304, fileDeleteRequest.BaseRequest.Tag, nil)
	}

	return base.NewSuccessResponse(fileDeleteRequest.BaseRequest.Tag, nil)
}
