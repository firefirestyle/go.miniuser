package miniuser

import (
	"time"

	"errors"

	"github.com/firefirestyle/go.miniprop"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	//	"google.golang.org/appengine/log"
)

type GaeRelayIdItem struct {
	Name     string
	Id       string
	Type     string
	UserName string
	Info     string
	Update   time.Time
}

type RelayId struct {
	gaeObj *GaeRelayIdItem
	gaeKey *datastore.Key
	kind   string
}

const (
	TypeTwitter = "twitter"
)

func (obj *UserManager) LoginRegistFromTwitter(ctx context.Context, screenName string, userId string, oauthToken string) (bool, *RelayId, *User, error) {
	sessionObj := obj.GetRelayIdForTwitter(ctx, screenName, userId, oauthToken)
	needMake := false

	//
	// new userObj
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
		userObj.SetDisplayName(screenName)
	}

	//
	// set username
	sessionObj.gaeObj.UserName = userObj.GetUserName()

	//
	// save relayId
	_, err = datastore.Put(ctx, sessionObj.gaeKey, sessionObj.gaeObj)
	if err != nil {
		return needMake, nil, nil, errors.New("failed to save sessionobj : " + err.Error())
	}

	//
	// save user
	err = obj.SaveUser(ctx, userObj)
	if err != nil {
		return needMake, nil, nil, errors.New("failed to save userobj : " + err.Error())
	}
	return needMake, sessionObj, userObj, nil
}

func (obj *UserManager) GetRelayIdForTwitter(ctx context.Context, screenName string, userId string, oauthToken string) *RelayId {
	return obj.GetRelayIdWithNew(ctx, screenName, userId, TypeTwitter, map[string]string{"token": oauthToken})
}

func (obj *UserManager) NewRelayId(ctx context.Context, screenName string, //
	userId string, userIdType string, infos map[string]string) *RelayId {
	gaeKey := datastore.NewKey(ctx, obj.relayIdKind, obj.MakeRelayIdStringId(userId, userIdType), 0, nil)
	gaeObj := GaeRelayIdItem{
		Name: screenName,
		Id:   userId,
		Type: TypeTwitter,
	}
	propObj := miniprop.NewMiniPropFromJson([]byte(gaeObj.Info))
	for k, v := range infos {
		propObj.SetString(k, v)
	}
	gaeObj.Info = string(propObj.ToJson())
	gaeObj.Update = time.Now()
	return &RelayId{
		gaeObj: &gaeObj,
		gaeKey: gaeKey,
		kind:   obj.userKind,
	}
}

func (obj *UserManager) GetRelayId(ctx context.Context, identify string, identifyType string) (*RelayId, error) {
	gaeKey := datastore.NewKey(ctx, obj.relayIdKind, obj.MakeRelayIdStringId(identify, identifyType), 0, nil)
	gaeObj := GaeRelayIdItem{}
	err := datastore.Get(ctx, gaeKey, &gaeObj)
	if err != nil {
		return nil, err
	}
	return &RelayId{
		gaeObj: &gaeObj,
		gaeKey: gaeKey,
		kind:   obj.relayIdKind,
	}, nil
}

func (obj *UserManager) GetRelayIdWithNew(ctx context.Context, screenName string, userId string, userIdType string, infos map[string]string) *RelayId {
	relayObj, err := obj.GetRelayId(ctx, screenName, userIdType)
	if err != nil {
		relayObj = obj.NewRelayId(ctx, screenName, userId, userIdType, infos)
	}
	//
	propObj := miniprop.NewMiniPropFromJson([]byte(relayObj.gaeObj.Info))
	for k, v := range infos {
		propObj.SetString(k, v)
	}
	relayObj.gaeObj.Info = string(propObj.ToJson())
	relayObj.gaeObj.Update = time.Now()
	return relayObj
}

func (obj *UserManager) MakeRelayIdStringId(identify string, identifyType string) string {
	return obj.relayIdKind + ":" + obj.projectId + ":" + identifyType + ":" + identify
}
