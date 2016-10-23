package user

import (
	"crypto/sha1"
	"io"
	"time"

	"encoding/base32"

	"strconv"

	"github.com/firefirestyle/go.miniprop"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
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

// ----
// new object
// ----

func (obj *UserManager) newUserGaeObjectKey(ctx context.Context, userName string, sign string) *datastore.Key {
	Debug(ctx, "=================== make sign :"+userName+":"+sign+"=======================")
	return datastore.NewKey(ctx, obj.userKind, obj.MakeUserGaeObjectKeyStringId(userName, sign), 0, nil)
}

func (obj *UserManager) newUserWithUserName(ctx context.Context, sign string) *User {
	var userObj *User = nil
	var err error = nil
	for {
		hashObj := sha1.New()
		now := time.Now().UnixNano()
		io.WriteString(hashObj, miniprop.MakeRandomId())
		io.WriteString(hashObj, strconv.FormatInt(now, 36))
		userName := string(base32.StdEncoding.EncodeToString(hashObj.Sum(nil)))
		userObj, err = obj.GetUserFromUserName(ctx, userName, sign)
		if err != nil {
			break
		}
	}
	return userObj
}

func (obj *UserManager) newUser(ctx context.Context, userName string, sign string) *User {
	ret := new(User)
	ret.prop = make(map[string]map[string]interface{})
	ret.kind = obj.userKind
	ret.gaeObject = new(GaeUserItem)
	ret.gaeObject.ProjectId = obj.projectId

	ret.gaeObject.UserName = userName
	ret.gaeObjectKey = obj.newUserGaeObjectKey(ctx, userName, sign)
	Debug(ctx, "GetUserFromUserName B:"+ret.gaeObjectKey.StringID())
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

// ----
// getter setter
// ----

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
