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
	TypeRootGroup   = "RootGroup"
	TypeDisplayName = "DisplayName"
	TypeUserName    = "UserName"
	TypeCreated     = "Created"
	TypeUpdated     = "Updated"
	TypeState       = "State"
	TypeTag         = "Tag"
	TypePublicInfo  = "PublicInfo"
	TypePointNames  = "PointNames"
	TypePointValues = "PointValues"
	TypePropNames   = "PropNames"
	TypePropValues  = "PropValues"
	TypeIconUrl     = "IconUrl"
	TypePrivateInfo = "PrivateInfo"
	TypeSign        = "Sign"
	TypeCont        = "Cont"
)

type GaeUserItem struct {
	RootGroup   string
	DisplayName string
	UserName    string
	Created     time.Time
	Updated     time.Time
	State       string
	PublicInfo  string    `datastore:",noindex"`
	PrivateInfo string    `datastore:",noindex"`
	Tags        []string  `datastore:"Tags.Tag"`
	PointNames  []string  `datastore:"Points.Name"`
	PointValues []float64 `datastore:"Points.Value"`
	PropNames   []string  `datastore:"Props.Name"`
	PropValues  []string  `datastore:"Props.Value"`
	//Point   int
	IconUrl string `datastore:",noindex"`
	Sign    string `datastore:",noindex"`
	Cont    string `datastore:",noindex"`
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
	ret.gaeObject.RootGroup = obj.rootGroup
	ret.gaeObject.Sign = sign
	ret.gaeObject.UserName = userName
	ret.gaeObjectKey = obj.newUserGaeObjectKey(ctx, userName, sign)
	Debug(ctx, "GetUserFromUserName B:"+sign+":==:"+ret.gaeObjectKey.StringID())
	return ret
}

func (obj *UserManager) newUserFromStringID(ctx context.Context, stringId string) *User {
	ret := new(User)
	ret.prop = make(map[string]map[string]interface{})
	ret.kind = obj.userKind
	ret.gaeObject = new(GaeUserItem)
	ret.gaeObject.RootGroup = obj.rootGroup
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
	return obj.gaeObject.Updated
}

func (obj *User) GetPublicInfo() string {
	return obj.gaeObject.PublicInfo
}

func (obj *User) SetPublicInfo(v string) {
	obj.gaeObject.PublicInfo = v
}

func (obj *User) GetPrivateInfo() string {
	return obj.gaeObject.PrivateInfo
}

func (obj *User) SetPrivateInfo(v string) {
	obj.gaeObject.PrivateInfo = v
}

func (obj *User) GetPoint(name string) float64 {
	index := -1
	for i, v := range obj.gaeObject.PointNames {
		if v == name {
			index = i
			break
		}
	}
	if index < 0 {
		return 0
	}
	return obj.gaeObject.PointValues[index]
}

func (obj *User) SetPoint(name string, v float64) {
	index := -1
	for i, iv := range obj.gaeObject.PointNames {
		if iv == name {
			index = i
			break
		}
	}
	if index == -1 {
		obj.gaeObject.PointValues = append(obj.gaeObject.PointValues, v)
		obj.gaeObject.PointNames = append(obj.gaeObject.PointNames, name)
	} else {
		obj.gaeObject.PointValues[index] = v
	}
}

func (obj *User) GetProp(name string) string {
	index := -1
	for i, v := range obj.gaeObject.PropNames {
		if v == name {
			index = i
			break
		}
	}
	if index < 0 {
		return ""
	}
	return obj.gaeObject.PropValues[index]
}

func (obj *User) SetProp(name, v string) {
	index := -1
	for i, iv := range obj.gaeObject.PropNames {
		if iv == name {
			index = i
			break
		}
	}
	if index == -1 {
		obj.gaeObject.PropValues = append(obj.gaeObject.PropValues, v)
		obj.gaeObject.PropNames = append(obj.gaeObject.PropNames, name)
	} else {
		obj.gaeObject.PropValues[index] = v
	}
}

func (obj *User) SetStatus(v string) {
	obj.gaeObject.State = v
}

func (obj *User) GetStatus() string {
	return obj.gaeObject.State
}

func (obj *User) SetCont(v string) {
	obj.gaeObject.Cont = v
}

func (obj *User) GetCont() string {
	return obj.gaeObject.Cont
}

func (obj *User) GetSign() string {
	return obj.gaeObject.Sign
}

func (obj *User) GetStringId() string {
	return obj.gaeObjectKey.StringID()
}

func (obj *User) GetTags() []string {
	ret := make([]string, 0)
	for _, v := range obj.gaeObject.Tags {
		//		ret = append(ret, v.Tag)
		ret = append(ret, v)
	}
	return ret
}

func (obj *User) SetTags(vs []string) {
	//	obj.gaeObject.Tags = make([]Tag, 0)
	obj.gaeObject.Tags = make([]string, 0)
	for _, v := range vs {
		//		obj.gaeObject.Tags = append(obj.gaeObject.Tags, Tag{Tag: v})
		obj.gaeObject.Tags = append(obj.gaeObject.Tags, v)
	}
}
