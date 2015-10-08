package userModels

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"gopkg.in/mgo.v2"
	"github.com/CodeCollaborate/CodeCollaborate/modules/user/requests"
	"gopkg.in/mgo.v2/bson"
	"time"
	"github.com/CodeCollaborate/CodeCollaborate/modules/base"
	"github.com/CodeCollaborate/CodeCollaborate/managers"
)

type User struct {
	Id            string  `bson:"_id"` // ID of object
	Username      string               // Username
	Email         string               // Email of user
	Password      string  `json:"-"`   // Unhashed Password
	Password_Hash string  `json:"-"`   // Hashed Password
	Tokens        []string  `json:"-"` // Token after logged in.
}

func RegisterUser(registrationRequest userRequests.UserRegisterRequest) base.WSResponse {

	// Hash password using bcrypt
	pwHashBytes, err := bcrypt.GenerateFromPassword([]byte(registrationRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Failed to hash password:", err)
		return base.NewFailResponse(-101, registrationRequest.BaseMessage.Tag, nil)
	}

	// Create new UserAuthData object
	userAuthData := new(User)
	userAuthData.Id = managers.NewObjectIdString()
	userAuthData.Username = registrationRequest.Username
	userAuthData.Email = registrationRequest.Email
	userAuthData.Password_Hash = string(pwHashBytes[:])

	// Get new DB connection
	session, collection := managers.GetMGoCollection("Users")
	defer session.Close()

	// Make sure index is unique
	index := mgo.Index{
		Key: []string{"username"},
		Unique:true,
		DropDups:true,
		Background:true,
		Sparse:true,
	}
	err = collection.EnsureIndex(index);
	if err != nil {
		log.Println("Failed to ensure username index:", err)
		return base.NewFailResponse(-101, registrationRequest.BaseMessage.Tag, nil)
	}

	// Register new user
	err = collection.Insert(userAuthData)
	if err != nil {
		if !mgo.IsDup(err) {
			log.Println("Error registering user:", err)
			return base.NewFailResponse(-101, registrationRequest.BaseMessage.Tag, nil)
		}
		// Duplicate entry
		return base.NewFailResponse(-102, registrationRequest.BaseMessage.Tag, nil)
	}

	return base.NewSuccessResponse(registrationRequest.BaseMessage.Tag, nil)
}

func LoginUser(loginRequest userRequests.UserLoginRequest) base.WSResponse {

	// Get new DB connection
	session, collection := managers.GetMGoCollection("Users")
	defer session.Close()

	user := User{}
	if err := collection.Find(bson.M{"$or": []interface{}{bson.M{"username": loginRequest.UsernameOREmail}, bson.M{"email": loginRequest.UsernameOREmail}}}).One(&user); err != nil {
		// Could not find user
		return base.NewFailResponse(-104, loginRequest.BaseMessage.Tag, nil)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password_Hash), []byte(loginRequest.Password)); err != nil {
		// Password did not match.
		return base.NewFailResponse(-104, loginRequest.BaseMessage.Tag, nil)
	}

	tokenBytes, err := bcrypt.GenerateFromPassword([]byte(loginRequest.UsernameOREmail + time.Now().String()), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Failed to generate token:", err)
		return base.NewFailResponse(-103, loginRequest.BaseMessage.Tag, nil)
	}

	token := string(tokenBytes[:])

	err = addToken(collection, user, token)
	if err != nil {
		log.Println("Failed to save token:", err)
		return base.NewFailResponse(-103, loginRequest.BaseMessage.Tag, nil)
	}

	return base.NewSuccessResponse(loginRequest.BaseMessage.Tag, map[string]interface{}{"UserId": user.Id, "Token":token})
}

func CheckUserAuth(baseRequest base.BaseRequest) bool {

	// Get new DB connection
	session, collection := managers.GetMGoCollection("Users")
	defer session.Close()

	userAuthData := User{}
	if err := collection.Find(bson.M{"_id": baseRequest.UserId}).One(&userAuthData); err != nil {
		// Could not find user
		return false
	}

	for _, token := range userAuthData.Tokens {
		if (token == baseRequest.Token) {
			// Found matching token
			return true
		}
	}

	// No matching token found
	return false
}

func addToken(collection *mgo.Collection, userAuthData User, token string) error {
	userAuthData.Tokens = append(userAuthData.Tokens, token)

	return collection.Update(bson.M{"username": userAuthData.Username}, bson.M{"$set": bson.M{"tokens": userAuthData.Tokens}})
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