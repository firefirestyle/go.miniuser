package handler

import (
	"net/http"

	//	"strings"

	//	miniuser "github.com/firefirestyle/go.miniuser/user"
	"io/ioutil"

	"github.com/firefirestyle/go.miniprop"
	"google.golang.org/appengine"
)

func (obj *UserHandler) HandleUpdateInfo(w http.ResponseWriter, r *http.Request) {
	outputProp := miniprop.NewMiniProp()
	v, _ := ioutil.ReadAll(r.Body)
	inputProp := miniprop.NewMiniPropFromJson(v)
	ctx := appengine.NewContext(r)
	userName := inputProp.GetString("userName", "")
	displayName := inputProp.GetString("displayName", "")
	content := inputProp.GetString("content", "")

	reqErr := obj.OnUpdateUserRequest(w, r, obj, inputProp, outputProp)
	if reqErr != nil {
		obj.OnUpdateUserFailed(w, r, obj, inputProp, outputProp)
		obj.HandleError(w, r, outputProp, 2001, reqErr.Error())
		return
	}
	usrObj, userErr := obj.GetManager().GetUserFromRelayId(ctx, userName)
	if userErr != nil {
		obj.OnUpdateUserFailed(w, r, obj, inputProp, outputProp)
		obj.HandleError(w, r, outputProp, 2002, userErr.Error())
		return
	}
	usrObj.SetDisplayName(displayName)
	usrObj.SetCont(content)
	defChec := obj.OnUpdateUserBeforeSave(w, r, obj, usrObj, inputProp, outputProp)
	if defChec != nil {
		obj.OnUpdateUserFailed(w, r, obj, inputProp, outputProp)
		obj.HandleError(w, r, outputProp, 2003, defChec.Error())
		return
	}
	nextUserObj, nextUserErr := obj.GetManager().SaveUserWithImmutable(ctx, usrObj)
	if nextUserErr != nil {
		obj.OnUpdateUserFailed(w, r, obj, inputProp, outputProp)
		obj.HandleError(w, r, outputProp, 2004, userErr.Error())
	} else {

		obj.OnUpdateUserSuccess(w, r, obj, usrObj, inputProp, outputProp)
		w.WriteHeader(http.StatusOK)
		w.Write(nextUserObj.ToJsonPublic())
	}
}
