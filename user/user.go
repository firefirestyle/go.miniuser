// gaeuser project gaeuser.go
package user

import (
	"crypto/sha1"
	"encoding/json"
	"io"
	"time"

	"encoding/base64"

	"strconv"
	"strings"

	"github.com/firefirestyle/go.miniprop"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
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
	TypePublicInfo  = "PublicInfo"
	TypePoint       = "Point"
	TypeIconUrl     = "IconUrl"
	TypePrivateInfo = "PrivateInfo"
)

type GaeUserItem struct {
	ProjectId   string
	DisplayName string
	UserName    string    `datastore:",noindex"`
	Created     time.Time `datastore:",noindex"`
	Logined     time.Time
	State       string
	PublicInfo  string `datastore:",noindex"`
	PrivateInfo string `datastore:",noindex"`
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
		userName := string(base64.StdEncoding.EncodeToString(hashObj.Sum(nil)))
		userObj, err = obj.GetUserFromUserName(ctx, userName)
		if err != nil {
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
	userObj.SetUserFromsMap(ctx, v)
	return nil
}

func (userObj *User) CopyWithoutuserName(ctx context.Context, copyObj *User) {
	itemInfo := userObj.ToMapAll()
	itemInfo[TypeUserName] = copyObj.GetUserName()
	copyObj.SetUserFromsMap(ctx, itemInfo)
}

func (userObj *User) SetUserFromsMap(ctx context.Context, v map[string]interface{}) {
	propObj := miniprop.NewMiniPropFromMap(v)
	userObj.gaeObject.ProjectId = propObj.GetString(TypeProjectId, "")
	userObj.gaeObject.DisplayName = propObj.GetString(TypeDisplayName, "")
	userObj.gaeObject.UserName = propObj.GetString(TypeUserName, "")
	userObj.gaeObject.Created = propObj.GetTime(TypeCreated, time.Now()) //srcCreated
	userObj.gaeObject.Logined = propObj.GetTime(TypeLogined, time.Now()) //time.Unix(0, int64(v[TypeLogined].(float64))) //srcLogin
	userObj.gaeObject.State = propObj.GetString(TypeState, "")
	userObj.gaeObject.PublicInfo = propObj.GetString(TypePublicInfo, "")
	userObj.gaeObject.PrivateInfo = propObj.GetString(TypePrivateInfo, "")
	userObj.gaeObject.Point = propObj.GetInt(TypePoint, 0)
	userObj.gaeObject.IconUrl = propObj.GetString(TypeIconUrl, "")
}

func (obj *User) ToMapPublic() map[string]interface{} {
	return map[string]interface{}{
		TypeProjectId:   obj.gaeObject.ProjectId,
		TypeDisplayName: obj.gaeObject.DisplayName,        //
		TypeUserName:    obj.gaeObject.UserName,           //
		TypeCreated:     obj.gaeObject.Created.UnixNano(), //
		TypeLogined:     obj.gaeObject.Logined.UnixNano(), //
		TypeState:       obj.gaeObject.State,              //
		TypePoint:       obj.gaeObject.Point,              //
		TypeIconUrl:     obj.gaeObject.IconUrl,            //
		TypePublicInfo:  obj.gaeObject.PublicInfo}
}

func (obj *User) ToMapAll() map[string]interface{} {
	v := obj.ToMapPublic()
	v[TypePrivateInfo] = obj.gaeObject.PrivateInfo
	return v
}

func (obj *User) ToJson() ([]byte, error) {
	vv, e := json.Marshal(obj.ToMapAll())
	return vv, e
}
func (obj *User) ToJsonPublic() ([]byte, error) {
	vv, e := json.Marshal(obj.ToMapPublic())
	return vv, e
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
func (obj *User) GetOriginalUserName() string {
	return strings.Split(obj.gaeObject.UserName, "::sign::")[0]
}

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
	return obj.gaeObject.PublicInfo
}

func (obj *User) SetInfo(v string) {
	obj.gaeObject.PublicInfo = v
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
