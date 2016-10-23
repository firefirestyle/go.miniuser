package miniuser

import (
	//	"github.com/firefirestyle/go.miniuser/relayid"
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

func (obj *UserManager) MakeUserGaeObjectKeyStringId(userName string) string {
	return obj.userKind + ":" + obj.projectId + ":" + userName
}

func (obj *UserManager) GetUserKind() string {
	return obj.userKind
}

func (obj *UserManager) NewNewUser(ctx context.Context) *User {
	return obj.newUserWithUserName(ctx)
}

func (obj *UserManager) GetUserFromUserName(ctx context.Context, userName string) (*User, error) {
	userObj := obj.newUser(ctx, userName)
	Debug(ctx, "GetUserFromUserName :"+userName)

	e := userObj.loadFromDB(ctx)
	return userObj, e
}

func (obj *UserManager) SaveUser(ctx context.Context, userObj *User) error {
	return userObj.pushToDB(ctx)
}

func (obj *UserManager) DeleteUser(ctx context.Context, userName string) error {
	gaeKey := obj.newUserGaeObjectKey(ctx, userName)
	return datastore.Delete(ctx, gaeKey)
}

func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}
