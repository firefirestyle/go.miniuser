package user

import (
	"github.com/firefirestyle/go.minipointer"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type UserManagerConfig struct {
	RootGroup       string
	UserKind        string
	UserPointerKind string
	LengthHash      int
	LimitOfFinding  int
}

type UserManager struct {
	config         UserManagerConfig
	pointerManager *minipointer.PointerManager
}

func NewUserManager(config UserManagerConfig) *UserManager {
	obj := new(UserManager)
	if config.RootGroup == "" {
		config.RootGroup = "FFUser"
	}
	if config.UserKind == "" {
		config.UserKind = "FFUser"
	}
	if config.UserPointerKind == "" {
		config.UserPointerKind = config.UserKind + "-pointer"
	}
	if config.LimitOfFinding <= 0 {
		config.LimitOfFinding = 20
	}
	obj.config = config

	obj.pointerManager = minipointer.NewPointerManager(minipointer.PointerManagerConfig{
		RootGroup: config.RootGroup,
		Kind:      config.UserPointerKind,
	})

	return obj
}

func (obj *UserManager) GetUserKind() string {
	return obj.config.UserKind
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
