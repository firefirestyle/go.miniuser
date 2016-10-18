package miniuser

import (
	"time"

	"errors"

	"github.com/firefirestyle/go.miniprop"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type GaeSessionItem struct {
	Name     string
	Id       string
	Type     string
	UserName string
	Info     string
	Update   time.Time
}

type Session struct {
	gaeObj *GaeSessionItem
	gaeKey *datastore.Key
}

const (
	TypeTwitter = "twitter"
)

func (obj *UserManager) LoginRegistFromTwitter(ctx context.Context, screenName string, userId string, oauthToken string) (bool, *Session, *User, error) {
	sessionObj := obj.GetSessionTwitter(ctx, screenName, userId, oauthToken)
	needMake := false
	var err error = nil
	var userObj *User = nil

	if sessionObj.gaeObj.Name != "" {
		needMake = true
		userObj, err = obj.GetUserFromUserName(ctx, sessionObj.gaeObj.Name)
		if err != nil {
			userObj = nil
		}
	}

	if userObj == nil {
		userObj = obj.NewNewUser(ctx)
	}
	//
	sessionObj.gaeObj.UserName = userObj.GetUserName()
	_, err = datastore.Put(ctx, sessionObj.gaeKey, sessionObj.gaeObj)
	if err != nil {
		return needMake, nil, nil, errors.New("failed to save sessionobj")
	}
	err = obj.SaveUser(ctx, userObj)
	if err != nil {
		return needMake, nil, nil, errors.New("failed to save userobj")
	}
	return needMake, sessionObj, userObj, nil
}

func (obj *UserManager) GetSessionTwitter(ctx context.Context, screenName string, userId string, oauthToken string) *Session {
	return obj.GetSession(ctx, screenName, userId, TypeTwitter, map[string]string{"token": oauthToken})
}

func (obj *UserManager) GetSession(ctx context.Context, screenName string, userId string, userIdType string, infos map[string]string) *Session {
	gaeKey := datastore.NewKey(ctx, obj.sessionKind, obj.MakeSessionId(screenName, userIdType), 0, nil)
	gaeObj := GaeSessionItem{}
	err := datastore.Get(ctx, gaeKey, &gaeObj)
	if err != nil {
		gaeObj = GaeSessionItem{
			Name: screenName,
			Id:   userId,
			Type: TypeTwitter,
		}

	}
	propObj := miniprop.NewMiniPropFromJson([]byte(gaeObj.Info))
	for k, v := range infos {
		propObj.SetString(k, v)
	}
	gaeObj.Info = string(propObj.ToJson())
	gaeObj.Update = time.Now()
	return &Session{
		gaeObj: &gaeObj,
		gaeKey: gaeKey,
	}
}

func (obj *UserManager) MakeSessionId(identify string, identifyType string) string {
	return obj.sessionKind + ":" + obj.projectId + ":" + identifyType + ":" + identify
}
