package user

import (
	//	"github.com/firefirestyle/go.miniuser/relayid"
	//	"github.com/firefirestyle/go.miniprop"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type UserManagerConfig struct {
	ProjectId string
	UserKind  string
}

type UserManager struct {
	projectId      string
	userKind       string
	limitOfFinding int
}

func NewUserManager(config UserManagerConfig) *UserManager {
	obj := new(UserManager)
	obj.projectId = config.ProjectId
	obj.userKind = config.UserKind
	obj.limitOfFinding = 10

	return obj
}

func (obj *UserManager) MakeUserGaeObjectKeyStringId(userName string, sign string) string {
	return "k:" + obj.userKind + ";p:" + obj.projectId + ";n:" + userName + ";s:" + sign + ";"
}

func (obj *UserManager) GetUserKind() string {
	return obj.userKind
}

func (obj *UserManager) NewNewUser(ctx context.Context, sign string) *User {
	return obj.newUserWithUserName(ctx, sign)
}

func (obj *UserManager) GetUserFromUserName(ctx context.Context, userName string, sign string) (*User, error) {
	userObj := obj.newUser(ctx, userName, sign)
	Debug(ctx, "GetUserFromUserName A:"+userName+":"+sign)

	e := userObj.loadFromDB(ctx)
	if e != nil {
		Debug(ctx, "GetUserFromUserName A:E:"+userName+":"+sign)

	}
	return userObj, e
}

func (obj *UserManager) SaveUser(ctx context.Context, userObj *User) error {
	Debug(ctx, "GetUserFromUserName ZZ:"+userObj.gaeObjectKey.StringID())
	return userObj.pushToDB(ctx)
}

func (obj *UserManager) DeleteUser(ctx context.Context, userName string, sign string) error {
	Debug(ctx, "----------------dlete:"+userName+":"+sign+"----------------dlete:")
	gaeKey := obj.newUserGaeObjectKey(ctx, userName, sign)
	return datastore.Delete(ctx, gaeKey)
}

func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}
