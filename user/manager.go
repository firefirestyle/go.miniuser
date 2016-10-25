package user

import (
	//	"github.com/firefirestyle/go.miniuser/relayid"
	//	"github.com/firefirestyle/go.miniprop"
	"strings"

	"errors"

	"strconv"

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

type UserKeyInfo struct {
	Kind      string
	ProjectId string
	UserName  string
	Sign      string
}

func (obj *UserManager) MakeUserGaeObjectKeyStringId(userName string, sign string) string {
	return "k:" + obj.userKind + ";p:" + obj.projectId + ";n:" + userName + ";s:" + sign + ";"
}

func (obj *UserManager) GetUserKeyInfo(stringId string) (*UserKeyInfo, error) {
	items := strings.Split(stringId, ";")
	if len(items) < 4 {
		return nil, errors.New("wrong id 1 : " + strconv.Itoa(len(items)))
	}
	ks := strings.Split(items[0], ":")
	if len(items) < 2 {
		return nil, errors.New("wrong id 2")
	}
	ps := strings.Split(items[1], ":")
	if len(items) < 2 {
		return nil, errors.New("wrong id 3")
	}
	ns := strings.Split(items[2], ":")
	if len(items) < 2 {
		return nil, errors.New("wrong id 4")
	}
	ss := strings.Split(items[3], ":")
	if len(items) < 2 {
		return nil, errors.New("wrong id")
	}

	return &UserKeyInfo{
		Kind:      ks[1],
		ProjectId: ps[1],
		UserName:  ns[1],
		Sign:      ss[1],
	}, nil
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
