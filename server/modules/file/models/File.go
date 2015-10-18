package fileModels

import (
	"log"

	"github.com/CodeCollaborate/CodeCollaborate/server/managers"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/models"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/file/requests"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"os"
)

type File struct {
	Id           string `bson:"_id"`           // ID of object
	Name         string                        // Name of file
	RelativePath string `bson:"relative_path"` // Path of file
	Version      int                           // File version
	Project      string                        // Reference to Project object
	filePath     string `bson:"-",json:"-"`
}

func (file File) getPath() string {
	if (file.filePath == "") {
		// Change to use byte buffer for efficiency
		file.filePath = "files/" + file.Project + "/" + file.RelativePath + file.Name
	}
	return file.filePath
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
		log.Println("Failed to ensure file name/path index:", err)
		return baseModels.NewFailResponse(-301, fileCreateRequest.BaseRequest.Tag, nil)
	}

	// Insert file record
	err = collection.Insert(file)
	if err != nil {
		if mgo.IsDup(err) {
			log.Println("Error creating file record:", err)
			return baseModels.NewFailResponse(-305, fileCreateRequest.BaseRequest.Tag, nil)
		}
		return baseModels.NewFailResponse(-301, fileCreateRequest.BaseRequest.Tag, nil)
	}

	// Write file to disk
	err = os.MkdirAll("files/" + fileCreateRequest.ProjectId + "/" + fileCreateRequest.RelativePath, os.ModeExclusive)
	if err != nil {
		log.Println("Failed to create file directory:", err)
		return baseModels.NewFailResponse(-301, fileCreateRequest.BaseRequest.Tag, nil)
	}
	err = ioutil.WriteFile(file.getPath(), fileCreateRequest.FileBytes, os.ModeExclusive)
	if err != nil {
		log.Println("Failed to write file:", err)
		return baseModels.NewFailResponse(-301, fileCreateRequest.BaseRequest.Tag, nil)
	}

	managers.NotifyProjectClients(file.Project, fileCreateRequest.GetNotification(file.Id))

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

	managers.NotifyProjectClients(file.Project, fileRenameRequest.GetNotification())

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

	managers.NotifyProjectClients(file.Project, fileMoveRequest.GetNotification())

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

	err = os.Remove(file.getPath())

	err = collection.Remove(bson.M{"_id": fileDeleteRequest.BaseRequest.ResId})
	if err != nil {
		return baseModels.NewFailResponse(-304, fileDeleteRequest.BaseRequest.Tag, nil)
	}

	managers.NotifyProjectClients(file.Project, fileDeleteRequest.GetNotification())

	return baseModels.NewSuccessResponse(fileDeleteRequest.BaseRequest.Tag, nil)
}

func PullFile(filePullRequest fileRequests.FilePullRequest) baseModels.WSResponse {

	// Check that file exists
	file, err := GetFileById(filePullRequest.BaseRequest.ResId);
	if err != nil {
		return baseModels.NewFailResponse(-300, filePullRequest.BaseRequest.Tag, nil)
	}

	// Read file from disk
	if _, err := os.Stat(file.getPath()); os.IsNotExist(err) {
		return baseModels.NewSuccessResponse(filePullRequest.BaseRequest.Tag, map[string]interface{}{"FileBytes": "", "Changes": ""})
	}
	fileBytes, err := ioutil.ReadFile(file.getPath())
	if err != nil {
		log.Println("Failed to read from file:", err)
		return baseModels.NewFailResponse(-301, filePullRequest.BaseRequest.Tag, nil)
	}

	changes, err := GetChangesByFile(file.Id)
	if err != nil {
		log.Println("Failed to retrieve changes:", err)
		return baseModels.NewFailResponse(-402, filePullRequest.BaseRequest.Tag, nil)
	}

	return baseModels.NewSuccessResponse(filePullRequest.BaseRequest.Tag, map[string]interface{}{"FileBytes": fileBytes, "Changes": changes})
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

