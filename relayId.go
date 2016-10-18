package miniuser

import (
	"time"

	"errors"

	"github.com/firefirestyle/go.miniprop"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	//	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

const (
	TypeTwitter = "twitter"
)

const (
	TypeRelayIdProjectId = "ProjectId"
	TypeRelayIdName      = "Name"
	TypeRelayIdId        = "Id"
	TypeRelayIdType      = "Type"
	TypeRelayIdUserName  = "UserName"
	TypeRelayIdInfo      = "Info"
	TypeRelayIdUpdate    = "Update"
)

type GaeRelayIdItem struct {
	ProjectId string
	Name      string
	Id        string
	Type      string
	UserName  string
	Info      string
	Update    time.Time
}

type RelayId struct {
	gaeObj *GaeRelayIdItem
	gaeKey *datastore.Key
	kind   string
}

func (obj *RelayId) ToJson() []byte {
	propObj := miniprop.NewMiniProp()
	propObj.SetString(TypeRelayIdProjectId, obj.gaeObj.ProjectId)
	propObj.SetString(TypeRelayIdName, obj.gaeObj.Name)
	propObj.SetString(TypeRelayIdId, obj.gaeObj.Id)
	propObj.SetString(TypeRelayIdUserName, obj.gaeObj.UserName)
	propObj.SetString(TypeRelayIdInfo, obj.gaeObj.Info)
	propObj.SetTime(TypeInfo, obj.gaeObj.Update)
	return propObj.ToJson()
}

func (obj *RelayId) SetValueFromJson(data []byte) {
	propObj := miniprop.NewMiniPropFromJson(data)
	obj.gaeObj.ProjectId = propObj.GetString(TypeRelayIdProjectId, "")
	obj.gaeObj.Name = propObj.GetString(TypeRelayIdName, "")
	obj.gaeObj.Id = propObj.GetString(TypeRelayIdId, "")
	obj.gaeObj.UserName = propObj.GetString(TypeRelayIdUserName, "")
	obj.gaeObj.Info = propObj.GetString(TypeRelayIdInfo, "")
	obj.gaeObj.Update = propObj.GetTime(TypeInfo, time.Now())
}

func (obj *RelayId) UpdateMemcache(ctx context.Context) {
	userObjMemSource := obj.ToJson()
	userObjMem := &memcache.Item{
		Key:   obj.gaeKey.StringID(),
		Value: []byte(userObjMemSource), //
	}
	memcache.Set(ctx, userObjMem)
}

func (obj *RelayId) GetName() string {
	return obj.gaeObj.Name
}

func (obj *RelayId) GetId() string {
	return obj.gaeObj.Id
}

func (obj *RelayId) GetType() string {
	return obj.gaeObj.Type
}

func (obj *RelayId) GetUserName() string {
	return obj.gaeObj.UserName
}

func (obj *RelayId) SetUserName(v string) {
	obj.gaeObj.UserName = v
}

func (obj *RelayId) GetInfo() string {
	return obj.gaeObj.Info
}

func (obj *RelayId) GetUpdate() time.Time {
	return obj.gaeObj.Update
}

func (obj *RelayId) Save(ctx context.Context) error {
	_, err := datastore.Put(ctx, obj.gaeKey, obj.gaeObj)
	if err == nil {
		obj.UpdateMemcache(ctx)
	}
	return err
}

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
	sessionObj.SetUserName(userObj.GetUserName())

	//
	// save relayId
	err = sessionObj.Save(ctx)
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

func (obj *UserManager) NewRelayId(ctx context.Context, identify string, //
	userId string, identifyType string, infos map[string]string) *RelayId {
	gaeKey := obj.NewRelayIdGaeKey(ctx, userId, identifyType)
	gaeObj := GaeRelayIdItem{
		Name:      identify,
		Id:        userId,
		Type:      TypeTwitter,
		ProjectId: obj.projectId,
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
	gaeKey := obj.NewRelayIdGaeKey(ctx, identify, identifyType)
	gaeObj := GaeRelayIdItem{}

	//
	// mem
	memItemObj, errMemObj := memcache.Get(ctx, obj.MakeRelayIdStringId(identify, identifyType))
	if errMemObj == nil {
		ret := &RelayId{
			gaeObj: &gaeObj,
			gaeKey: gaeKey,
			kind:   obj.relayIdKind,
		}
		ret.SetValueFromJson(memItemObj.Value)
		return ret, nil
	}
	//
	// db
	err := datastore.Get(ctx, gaeKey, &gaeObj)
	if err != nil {
		return nil, err
	}
	ret := &RelayId{
		gaeObj: &gaeObj,
		gaeKey: gaeKey,
		kind:   obj.relayIdKind,
	}
	//
	//
	ret.UpdateMemcache(ctx)
	return ret, nil
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

func (obj *UserManager) NewRelayIdGaeKey(ctx context.Context, identify string, identifyType string) *datastore.Key {
	return datastore.NewKey(ctx, obj.relayIdKind, obj.MakeRelayIdStringId(identify, identifyType), 0, nil)
}

func (obj *UserManager) MakeRelayIdStringId(identify string, identifyType string) string {
	return obj.relayIdKind + ":" + obj.projectId + ":" + identifyType + ":" + identify
}
