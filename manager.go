package miniuser

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type UserManager struct {
	projectId      string
	userKind       string
	sessionKind    string
	limitOfFinding int
}

func NewUserManager(projectId string, userKind string, sessionKind string) *UserManager {
	obj := new(UserManager)
	obj.projectId = projectId
	obj.userKind = userKind
	obj.sessionKind = sessionKind
	obj.limitOfFinding = 10
	return obj
}

func (obj *UserManager) MakeUserGaeObjectKeyStringId(userName string) string {
	return obj.userKind + ":" + obj.projectId + ":" + userName
}

func (obj *UserManager) GetUserKind() string {
	return obj.userKind
}

func (obj *UserManager) GetLoginIdKind() string {
	return obj.sessionKind
}

func (obj *UserManager) NewNewUser(ctx context.Context) *User {
	return obj.newUserWithUserName(ctx)
}

func (obj *UserManager) GetUserFromUserName(ctx context.Context, userName string) (*User, error) {
	userObj := obj.newUser(ctx, userName)
	e := userObj.loadFromDB(ctx)
	return userObj, e
}

func (obj *UserManager) SaveUser(ctx context.Context, userObj *User) error {
	return userObj.pushToDB(ctx)
}

func (obj *UserManager) DeleteUser(ctx context.Context, userName string, passIdFromClient string) error {
	gaeKey := obj.newUserGaeObjectKey(ctx, userName)
	return datastore.Delete(ctx, gaeKey)
}
