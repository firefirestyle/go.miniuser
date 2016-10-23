package user

import (
	"encoding/json"
	"time"

	"github.com/firefirestyle/go.miniprop"
	"golang.org/x/net/context"
)

// ----
// json and copy
// ----
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
