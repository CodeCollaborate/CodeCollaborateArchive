package userModels

import (
	"log"
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
)

type User struct {
	Id            string   `bson:"_id"` // ID of object
	Email         string   // Email of user
	Password      string   `json:"-",bson:"-"` // Unhashed Password
	Password_Hash string   `json:"-"` // Hashed Password
	Tokens        []string `json:"-"` // Token after logged in.
}

func RegisterUser(wsConn *websocket.Conn, registrationRequest userRequests.UserRegisterRequest){

	// Hash password using bcrypt
	pwHashBytes, err := bcrypt.GenerateFromPassword([]byte(registrationRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Failed to hash password:", err)
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-101, registrationRequest.BaseRequest.Tag, nil))
	}

	// Create new UserAuthData object
	userAuthData := new(User)
	userAuthData.Id = managers.NewObjectIdString()
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
		log.Println("Failed to ensure email index:", err)
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-101, registrationRequest.BaseRequest.Tag, nil))
	}

	// Register new user
	err = collection.Insert(userAuthData)
	if err != nil {
		// Duplicate entry
		if mgo.IsDup(err) {
			log.Println("Error registering user:", err)
			managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-102, registrationRequest.BaseRequest.Tag, nil))
		}
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-101, registrationRequest.BaseRequest.Tag, nil))
	}

	managers.SendWebSocketMessage(wsConn, baseModels.NewSuccessResponse(registrationRequest.BaseRequest.Tag, nil))
}

func LoginUser(wsConn *websocket.Conn, loginRequest userRequests.UserLoginRequest){

	// Get new DB connection
	session, collection := managers.GetMGoCollection("Users")
	defer session.Close()

	user := User{}
	if err := collection.Find(bson.M{"email": loginRequest.Email}).One(&user); err != nil {
		// Could not find user
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-104, loginRequest.BaseRequest.Tag, nil))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password_Hash), []byte(loginRequest.Password)); err != nil {
		// Password did not match.
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-104, loginRequest.BaseRequest.Tag, nil))
	}

	tokenBytes, err := bcrypt.GenerateFromPassword([]byte(loginRequest.Email +time.Now().String()), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Failed to generate token:", err)
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-103, loginRequest.BaseRequest.Tag, nil))
	}

	token := string(tokenBytes[:])

	err = addToken(collection, user, token)
	if err != nil {
		log.Println("Failed to save token:", err)
		managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-103, loginRequest.BaseRequest.Tag, nil))
	}

	managers.SendWebSocketMessage(wsConn, baseModels.NewSuccessResponse(loginRequest.BaseRequest.Tag, map[string]interface{}{"UserId": user.Id, "Token": token}))
}

func Subscribe(wsConn *websocket.Conn, subscriptionRequest userRequests.UserSubscribeRequest){

	toSubscribe := subscriptionRequest.Projects
	for _, project := range toSubscribe {
		proj, err := projectModels.GetProjectById(project)

		if err != nil {
			managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-200, subscriptionRequest.BaseRequest.Tag, nil))
		}

		for key, _ := range proj.Permissions {
			if key == subscriptionRequest.BaseRequest.UserId {
				if(!managers.WebSocketSubscribeProject(wsConn, project)){
					managers.SendWebSocketMessage(wsConn, baseModels.NewFailResponse(-206 , subscriptionRequest.BaseRequest.Tag, nil))
				}
			}

		}
	}

	managers.SendWebSocketMessage(wsConn, baseModels.NewSuccessResponse(subscriptionRequest.BaseRequest.Tag, nil))
}

func CheckUserAuth(baseRequest baseRequests.BaseRequest) bool {

	// Get new DB connection
	session, collection := managers.GetMGoCollection("Users")
	defer session.Close()

	userAuthData := User{}
	if err := collection.Find(bson.M{"_id": baseRequest.UserId}).One(&userAuthData); err != nil {
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

func addToken(collection *mgo.Collection, userAuthData User, token string) error {
	userAuthData.Tokens = append(userAuthData.Tokens, token)

	return collection.Update(bson.M{"email": userAuthData.Email}, bson.M{"$set": bson.M{"tokens": userAuthData.Tokens}})
}

func GetUserById(id string) (*User, error) {
	// Get new DB connection
	session, collection := managers.GetMGoCollection("Users")
	defer session.Close()

	result := new(User)
	err := collection.Find(bson.M{"_id": id}).One(&result)
	if err != nil {
		log.Println("Failed to retrieve User")
		log.Println(err)
		return nil, err
	}

	return result, nil
}