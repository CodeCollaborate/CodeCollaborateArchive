package fileModels

import (
	"log"

	"github.com/CodeCollaborate/CodeCollaborate/server/managers"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/models"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/file/requests"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type File struct {
	Id           string `bson:"_id"`           // ID of object
	Name         string                        // Name of file
	RelativePath string `bson:"relative_path"` // Path of file
	Version      int                           // File version
	Project      string                        // Reference to Project object
}

func CreateFile(fileCreateRequest fileRequests.FileCreateRequest) baseModels.WSResponse {

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
		return baseModels.NewFailResponse(-301, fileCreateRequest.BaseRequest.Tag, nil)
	}

	err = collection.Insert(file)
	if err != nil {
		if mgo.IsDup(err) {
			log.Println("Error registering user:", err)
			return baseModels.NewFailResponse(-305, fileCreateRequest.BaseRequest.Tag, nil)
		}
		return baseModels.NewFailResponse(-301, fileCreateRequest.BaseRequest.Tag, nil)
	}

	managers.NotifyAll(file.Project, fileCreateRequest.GetNotification(file.Id))

	return baseModels.NewSuccessResponse(fileCreateRequest.BaseRequest.Tag, map[string]interface{}{"FileId": file.Id})

}

func RenameFile(fileRenameRequest fileRequests.FileRenameRequest) baseModels.WSResponse {
	session, collection := managers.GetMGoCollection("Files")
	defer session.Close()

	// Check that file exists
	file, err := GetFileById(fileRenameRequest.BaseRequest.ResId);
	if err != nil {
		return baseModels.NewFailResponse(-300, fileRenameRequest.BaseRequest.Tag, nil)
	}

	file.Version++;

	err = collection.Update(bson.M{"_id": fileRenameRequest.BaseRequest.ResId}, bson.M{"$set": bson.M{"name": fileRenameRequest.NewName, "version": file.Version}})
	if err != nil {
		if mgo.IsDup(err) {
			log.Println("Error registering user:", err)
			return baseModels.NewFailResponse(-306, fileRenameRequest.BaseRequest.Tag, nil)
		}
		return baseModels.NewFailResponse(-302, fileRenameRequest.BaseRequest.Tag, nil)
	}

	managers.NotifyAll(file.Project, fileRenameRequest.GetNotification())

	return baseModels.NewSuccessResponse(fileRenameRequest.BaseRequest.Tag, nil)
}

func MoveFile(fileMoveRequest fileRequests.FileMoveRequest) baseModels.WSResponse {
	session, collection := managers.GetMGoCollection("Files")
	defer session.Close()

	// Check that file exists
	file, err := GetFileById(fileMoveRequest.BaseRequest.ResId);
	if err != nil {
		return baseModels.NewFailResponse(-300, fileMoveRequest.BaseRequest.Tag, nil)
	}

	file.Version++;

	err = collection.Update(bson.M{"_id": fileMoveRequest.BaseRequest.ResId}, bson.M{"$set": bson.M{"relative_path": fileMoveRequest.NewPath, "version": file.Version}})
	if err != nil {
		if mgo.IsDup(err) {
			log.Println("Error registering user:", err)
			return baseModels.NewFailResponse(-307, fileMoveRequest.BaseRequest.Tag, nil)
		}
		return baseModels.NewFailResponse(-303, fileMoveRequest.BaseRequest.Tag, nil)
	}

	managers.NotifyAll(file.Project, fileMoveRequest.GetNotification())

	return baseModels.NewSuccessResponse(fileMoveRequest.BaseRequest.Tag, nil)
}

func DeleteFile(fileDeleteRequest fileRequests.FileDeleteRequest) baseModels.WSResponse {
	session, collection := managers.GetMGoCollection("Files")
	defer session.Close()

	// Check that file exists
	file, err := GetFileById(fileDeleteRequest.BaseRequest.ResId);
	if err != nil {
		return baseModels.NewFailResponse(-300, fileDeleteRequest.BaseRequest.Tag, nil)
	}

	err = collection.Remove(bson.M{"_id": fileDeleteRequest.BaseRequest.ResId})
	if err != nil {
		return baseModels.NewFailResponse(-304, fileDeleteRequest.BaseRequest.Tag, nil)
	}

	managers.NotifyAll(file.Project, fileDeleteRequest.GetNotification())

	return baseModels.NewSuccessResponse(fileDeleteRequest.BaseRequest.Tag, nil)
}

func GetFileById(id string) (*File, error) {
	// Get new DB connection
	session, collection := managers.GetMGoCollection("Files")
	defer session.Close()

	result := new(File)
	err := collection.Find(bson.M{"_id": id}).One(&result)
	if err != nil {
		log.Println("Failed to retrieve File")
		log.Println(err)
		return nil, err
	}

	return result, nil
}

func GetFileByPathNamee(path string, name string) (*File, error) {
	// Get new DB connection
	session, collection := managers.GetMGoCollection("Files")
	defer session.Close()

	result := new(File)
	err := collection.Find(bson.M{"$or": []interface{}{bson.M{"relative_path": path}, bson.M{"name": name}}}).One(&result)
	if err != nil {
		log.Println("Failed to retrieve File")
		log.Println(err)
		return nil, err
	}

	return result, nil
}

