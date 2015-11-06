package scrunching

import (
	"os/exec"
	"gopkg.in/mgo.v2/bson"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/file/models"
	"os"
	"io/ioutil"
	"github.com/CodeCollaborate/CodeCollaborate/server/managers"
)

func ScrunchProject(projectID string) {
	files, err := fileModels.GetFilesByProjectId(projectID)
	if err {
		managers.LogError("Error retreaving project files on WS disconnect", err)
	} else {
		for _, file := range files {
			ScrunchDB(file.Id)
		}

	}
}

func ScrunchDB(fileID string) {

	if fileID != "" {
		managers.LogWarn("scrunchDB: passed an empty string")
		return
	}
	go goScrunchDB(fileID)
	return
}

func goScrunchDB(fileID string) {

	cmd := exec.Command("java", "-jar", "Scrunching.jar", fileID)
	err := cmd.Start()
	err = cmd.Wait()
	if err != nil {
		var message string
		switch err.Error() {
		case "exit status 60":
			message = "scrunchDB: Improper number of arguments" // will never get here
			managers.LogWarn(message)
		case "exit status 61":
			message = "scrunchDB: File not found in db with fileID: " + fileID
			managers.LogWarn(message)
			return
		case "exit status 62":
			message = "File not found on disk " + fileID
			managers.LogWarn(message)
			return correctFileNotFound(fileID)
		case "exit status 63":
			message = "scrunchDB: Unknown error while reading file: " + fileID
			managers.LogError(message, err)
		case "exit status 64":
			message = "scrunchDB: Can't apply patch for fileID" + fileID + ": Run Scrunching.jar directly to get bad patch key"
			//TODO: figure out if there's a way to get patch key from Scrunching.java @ line 103
			managers.LogError(message, err)
		case "exit status 65":
			message = "scrunchDB: Unable to compile patch for filID: " + fileID
			//TODO: figure out if there's a way to get patch key from Scrunching.java @ line 158
			managers.LogError(message, err)
		} 

		return
	}
}

func correctFileNotFound(fileID string) {
	session, collection := managers.GetMGoCollection("Files")
	defer session.Close()

	file := new(fileModels.File)
	err := collection.Find(bson.M{"file_id": fileID}).One(&file)
	if err != nil {
		managers.LogError("Scrunching create file: Failed to retrieve file location", err)
		return
	}

	err = os.MkdirAll("files/" + file.Project + "/" + file.GetPath(), os.ModeExclusive)
	if err != nil {
		managers.LogError("Failed to create file directory", err)
		return
	}
	err = ioutil.WriteFile(file.GetPath(), []byte{}, os.ModeExclusive)
	if err != nil {
		managers.LogError("Failed to write file", err)
		return
	}

	goScrunchDB(fileID)
}
