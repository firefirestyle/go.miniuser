// gaeuser project gaeuser.go
package miniuser

import (
	"crypto/sha1"
	"encoding/json"
	"io"
	"time"

	"github.com/firefirestyle/go.miniprop"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"

	"strconv"
)

const (
	UserStatePublic  = "public"
	UserStatePrivate = "private"
	UserStateAll     = ""
)

const (
	TypeProjectId   = "ProjectId"
	TypeDisplayName = "DisplayName"
	TypeUserName    = "UserName"
	TypeCreated     = "Created"
	TypeLogined     = "Logined"
	TypeState       = "State"
	TypeInfo        = "Info"
	TypePoint       = "Point"
	TypeIconUrl     = "IconUrl"
	//
)

type GaeUserItem struct {
	ProjectId   string
	DisplayName string
	UserName    string    `datastore:",noindex"`
	Created     time.Time `datastore:",noindex"`
	Logined     time.Time
	State       string
	Info        string `datastore:",noindex"`
	Point       int
	IconUrl     string `datastore:",noindex"`
}

type User struct {
	gaeObject    *GaeUserItem
	gaeObjectKey *datastore.Key
	kind         string
	prop         map[string]map[string]interface{}
}

//
// -------------------------
//

func (obj *UserManager) newUserGaeObjectKey(ctx context.Context, userName string) *datastore.Key {
	return datastore.NewKey(ctx, obj.userKind, obj.MakeUserGaeObjectKeyStringId(userName), 0, nil)
}

func (obj *UserManager) newUserWithUserName(ctx context.Context) *User {
	var userObj *User = nil
	var err error = nil
	for {
		hashObj := sha1.New()
		now := time.Now().UnixNano()
		io.WriteString(hashObj, miniprop.MakeRandomId())
		io.WriteString(hashObj, strconv.FormatInt(now, 36))
		userName := string(hashObj.Sum(nil))
		userObj, err = obj.GetUserFromUserName(ctx, userName)
		if err == nil {
			break
		}
	}
	return userObj
}

func (obj *UserManager) newUser(ctx context.Context, userName string) *User {
	ret := new(User)
	ret.prop = make(map[string]map[string]interface{})
	ret.kind = obj.userKind
	ret.gaeObject = new(GaeUserItem)
	ret.gaeObject.ProjectId = obj.projectId

	ret.gaeObject.UserName = userName
	ret.gaeObjectKey = obj.newUserGaeObjectKey(ctx, userName)
	return ret
}

func (obj *UserManager) newUserFromStringID(ctx context.Context, stringId string) *User {
	ret := new(User)
	ret.prop = make(map[string]map[string]interface{})
	ret.kind = obj.userKind
	ret.gaeObject = new(GaeUserItem)
	ret.gaeObject.ProjectId = obj.projectId
	ret.gaeObjectKey = datastore.NewKey(ctx, obj.userKind, stringId, 0, nil)
	return ret
}

//
// -------------------------
//
func getStringFromProp(requestPropery map[string]interface{}, key string, defaultValue string) string {
	v := requestPropery[key]
	if v == nil {
		return defaultValue
	} else {
		return v.(string)
	}
}

func (userObj *User) SetUserFromsJson(ctx context.Context, source string) error {
	v := make(map[string]interface{})
	e := json.Unmarshal([]byte(source), &v)
	if e != nil {
		return e
	}
	//
	userObj.gaeObject.ProjectId = getStringFromProp(v, TypeProjectId, "")
	userObj.gaeObject.DisplayName = v[TypeDisplayName].(string)
	userObj.gaeObject.UserName = v[TypeUserName].(string)
	userObj.gaeObject.Created = time.Now()
	userObj.gaeObject.Logined = time.Now()
	userObj.gaeObject.Created = time.Unix(0, int64(v[TypeCreated].(float64))) //srcCreated
	userObj.gaeObject.Logined = time.Unix(0, int64(v[TypeLogined].(float64))) //srcLogin
	userObj.gaeObject.State = v[TypeState].(string)
	userObj.gaeObject.Info = v[TypeInfo].(string)
	userObj.gaeObject.Point = int(v[TypePoint].(float64))
	userObj.gaeObject.IconUrl = v[TypeIconUrl].(string)

	return nil
}

func (obj *User) ToJson() (string, error) {
	v := map[string]interface{}{
		TypeProjectId:   obj.gaeObject.ProjectId,
		TypeDisplayName: obj.GetDisplayName(),        //
		TypeUserName:    obj.GetUserName(),           //
		TypeCreated:     obj.GetCreated().UnixNano(), //string(srcCreated),   //obj.GetCreated().UnixNano().(int64), //srcCreated,           //
		TypeLogined:     obj.GetLogined().UnixNano(), //
		TypeState:       obj.GetStatus(),             //
		TypeInfo:        obj.GetInfo(),               //
		TypePoint:       obj.GetPoint(),              //
		TypeIconUrl:     obj.GetIconUrl(),            //
	}
	vv, e := json.Marshal(v)
	return string(vv), e
}

func (obj *User) updateMemcache(ctx context.Context) error {
	userObjMemSource, err_toJson := obj.ToJson()
	if err_toJson == nil {
		userObjMem := &memcache.Item{
			Key:   obj.gaeObjectKey.StringID(),
			Value: []byte(userObjMemSource), //
		}
		memcache.Set(ctx, userObjMem)
	}
	return err_toJson
}

func (obj *User) deleteMemcache(ctx context.Context) error {
	return memcache.Delete(ctx, obj.gaeObjectKey.StringID())
}

//
//
func (obj *User) GetUserName() string {
	return obj.gaeObject.UserName
}

func (obj *User) GetDisplayName() string {
	return obj.gaeObject.DisplayName
}

func (obj *User) SetDisplayName(v string) {
	obj.gaeObject.DisplayName = v
}

func (obj *User) GetHaveIcon() bool {
	if obj.gaeObject.IconUrl == "" {
		return false
	} else {
		return true
	}
}

func (obj *User) SetIconUrl(v string) {
	obj.gaeObject.IconUrl = v
}

func (obj *User) GetIconUrl() string {
	return obj.gaeObject.IconUrl
}

func (obj *User) GetCreated() time.Time {
	return obj.gaeObject.Created
}

func (obj *User) GetLogined() time.Time {
	return obj.gaeObject.Logined
}

func (obj *User) GetInfo() string {
	return obj.gaeObject.Info
}

func (obj *User) SetInfo(v string) {
	obj.gaeObject.Info = v
}

func (obj *User) GetPoint() int {
	return obj.gaeObject.Point
}

func (obj *User) SetPoint(v int) {
	obj.gaeObject.Point = v
}

func (obj *User) SetStatus(v string) {
	obj.gaeObject.State = v
}

func (obj *User) GetStatus() string {
	return obj.gaeObject.State
}

//
//
//

func (obj *User) loadFromDB(ctx context.Context) error {
	item, err := memcache.Get(ctx, obj.gaeObjectKey.StringID())
	if err == nil {
		err1 := obj.SetUserFromsJson(ctx, string(item.Value))
		if err1 == nil {
			return nil
		} else {
			log.Infof(ctx, ">>> Failed Load UseObj Json On LoadFronDB :: %s", err1.Error())
		}
	}
	//
	//
	err_loaded := datastore.Get(ctx, obj.gaeObjectKey, obj.gaeObject)
	if err_loaded != nil {
		return err_loaded
	}

	obj.updateMemcache(ctx)

	return nil
}

func (obj *User) pushToDB(ctx context.Context) error {
	_, e := datastore.Put(ctx, obj.gaeObjectKey, obj.gaeObject)
	obj.updateMemcache(ctx)
	return e
}
