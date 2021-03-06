package userModels

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/CodeCollaborate/CodeCollaborate/server/managers"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/user/requests"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/models"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/base/requests"
	"github.com/CodeCollaborate/CodeCollaborate/server/modules/project/models"
	"regexp"
)

type User struct {
	Id            string   `json:"-",bson:"_id"` // ID of object
	FirstName     string                         // User's First name
	LastName      string                         // User's Last name
	Email         string                         // Email of user; Unique
	Username      string                         // Username; Unique
	Password_Hash string   `json:"-"`            // Hashed Password
	Tokens        []string `json:"-"`            // Token after logged in.
}

func RegisterUser(wsConn *websocket.Conn, registrationRequest userRequests.UserRegisterRequest) {

	matched, err := regexp.MatchString("^[\\w]+$", registrationRequest.Username)
	if (!matched) {
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-106, registrationRequest.BaseRequest.Tag, nil))
		return
	}

	// Hash password using bcrypt
	pwHashBytes, err := bcrypt.GenerateFromPassword([]byte(registrationRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		managers.LogError("Failed to hash password", err)
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-101, registrationRequest.BaseRequest.Tag, nil))
		return
	}

	// Create new UserAuthData object
	userAuthData := new(User)
	userAuthData.Id = managers.NewObjectIdString()
	userAuthData.FirstName = registrationRequest.FirstName
	userAuthData.LastName = registrationRequest.LastName
	userAuthData.Username = registrationRequest.Username
	userAuthData.Email = registrationRequest.Email
	userAuthData.Password_Hash = string(pwHashBytes[:])

	// Get new DB connection
	session, collection := managers.GetMGoCollection("Users")
	defer session.Close()

	// Make sure email is unique
	index := mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = collection.EnsureIndex(index)
	if err != nil {
		managers.LogError("Failed to ensure email index", err)
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-101, registrationRequest.BaseRequest.Tag, nil))
		return
	}

	// Make sure username is unique
	index = mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = collection.EnsureIndex(index)
	if err != nil {
		managers.LogError("Failed to ensure username index", err)
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-101, registrationRequest.BaseRequest.Tag, nil))
		return
	}

	// Register new user
	err = collection.Insert(userAuthData)
	if err != nil {
		// Duplicate entry
		if mgo.IsDup(err) {
			managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-102, registrationRequest.BaseRequest.Tag, nil))
			return
		}
		managers.LogError("Error registering user", err)
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-101, registrationRequest.BaseRequest.Tag, nil))
		return
	}

	managers.SendWebSocketMessage(wsConn, baseModels.NewSuccessResponse(registrationRequest.BaseRequest.Tag, nil))
}

func LoginUser(wsConn *websocket.Conn, loginRequest userRequests.UserLoginRequest) {

	// Get new DB connection
	session, collection := managers.GetMGoCollection("Users")
	defer session.Close()

	user, err := GetUserByUsername(loginRequest.Username)
	if err != nil {
		// Could not find user
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-104, loginRequest.BaseRequest.Tag, nil))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password_Hash), []byte(loginRequest.Password)); err != nil {
		// Password did not match.
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-104, loginRequest.BaseRequest.Tag, nil))
		return
	}

	tokenBytes, err := bcrypt.GenerateFromPassword([]byte(loginRequest.Username + time.Now().String()), bcrypt.DefaultCost)
	if err != nil {
		managers.LogError("Failed to generate token", err)
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-103, loginRequest.BaseRequest.Tag, nil))
		return
	}

	token := string(tokenBytes[:])

	err = addToken(collection, user, token)
	if err != nil {
		managers.LogError("Failed to save token", err)
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-103, loginRequest.BaseRequest.Tag, nil))
		return
	}

	managers.SendWebSocketMessage(wsConn, baseModels.NewSuccessResponse(loginRequest.BaseRequest.Tag, map[string]interface{}{"Token": token}))
}

func LookupUser(wsConn *websocket.Conn, userLookupRequest userRequests.UserLookupRequest) {

	// Get new DB connection
	session, collection := managers.GetMGoCollection("Users")
	defer session.Close()

	user := User{}
	if err := collection.Find(bson.M{"username": userLookupRequest.LookupUsername}).One(&user); err != nil {
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-100, userLookupRequest.BaseRequest.Tag, nil))
		return;
	}

	data := map[string]interface{}{"User": user}

	managers.SendWebSocketMessage(wsConn, baseModels.NewSuccessResponse(userLookupRequest.BaseRequest.Tag, data))
}

func UserProjects(wsConn *websocket.Conn, userProjectsRequest userRequests.UserProjectsRequest) {

	// Get new DB connection
	session, collection := managers.GetMGoCollection("Projects")
	defer session.Close()


	var projects []projectModels.Project
	if err := collection.Find(bson.M{"permissions." + userProjectsRequest.BaseRequest.Username: bson.M{"$gt": 0}}).All(&projects); err != nil {
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-100, userProjectsRequest.BaseRequest.Tag, nil))
		return;
	}

	data := map[string]interface{}{"Projects": projects}

	managers.SendWebSocketMessage(wsConn, baseModels.NewSuccessResponse(userProjectsRequest.BaseRequest.Tag, data))
}

func CheckUserAuth(baseRequest baseRequests.BaseRequest) bool {

	// Get new DB connection
	session, collection := managers.GetMGoCollection("Users")
	defer session.Close()

	userAuthData := User{}
	if err := collection.Find(bson.M{"username": baseRequest.Username}).One(&userAuthData); err != nil {
		// Could not find user
		return false
	}

	for _, token := range userAuthData.Tokens {
		if token == baseRequest.Token {
			// Found matching token
			return true
		}
	}

	// No matching token found
	return false
}

func addToken(collection *mgo.Collection, userAuthData *User, token string) error {
	userAuthData.Tokens = append(userAuthData.Tokens, token)

	return collection.Update(bson.M{"username": userAuthData.Username}, bson.M{"$set": bson.M{"tokens": userAuthData.Tokens}})
}


func GetUserByUsername(username string) (*User, error) {
	// Get new DB connection
	session, collection := managers.GetMGoCollection("Users")
	defer session.Close()

	result := new(User)
	err := collection.Find(bson.M{"username": username}).One(&result)
	if err != nil {
		managers.LogError("Failed to retrieve User", err)
		return nil, err
	}

	return result, nil
}