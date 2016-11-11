package handler

import (
	"net/http"

	//	"strings"

	"github.com/firefirestyle/go.miniprop"
	"github.com/firefirestyle/go.minisession"
	miniuser "github.com/firefirestyle/go.miniuser/user"
	"google.golang.org/appengine"
)

func (obj *UserHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	userName := values.Get("userName")
	sign := values.Get("sign")
	key := values.Get("key")
	obj.HandleGetBase(w, r, userName, sign, key, false)
}

func (obj *UserHandler) HandleGetMe(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	inputProp := miniprop.NewMiniPropFromJsonReader(r.Body)
	token := inputProp.GetString("token", "")
	loginResult := obj.GetSessionMgr().CheckLoginId(ctx, token, minisession.MakeAccessTokenConfigFromRequest(r))
	userName := loginResult.AccessTokenObj.GetUserName()
	if loginResult.IsLogin == false {
		userName = ""
	}
	obj.HandleGetBase(w, r, userName, "", "", false)
}

func (obj *UserHandler) HandleGetBase(w http.ResponseWriter, r *http.Request, userName string, sign string, key string, includePrivate bool) {
	ctx := appengine.NewContext(r)
	var usrObj *miniuser.User = nil
	var userErr error = nil

	outputProp := miniprop.NewMiniProp()
	reqErr := obj.OnGetUserRequest(w, r, obj, outputProp)
	if reqErr != nil {
		obj.OnGetUserFailed(w, r, obj, outputProp)
		obj.HandleError(w, r, outputProp, 2001, reqErr.Error())
		return
	}
	if userName != "" {
		if sign == "" {
			usrObj, userErr = obj.GetManager().GetUserFromRelayId(ctx, userName)
		} else {
			usrObj, userErr = obj.GetManager().GetUserFromSign(ctx, userName, sign)
		}
	} else if key != "" {
		usrObj, userErr = obj.GetManager().GetUserFromKey(ctx, key)
	} else {
		obj.OnGetUserFailed(w, r, obj, outputProp)
		obj.HandleError(w, r, outputProp, 2002, "wrong request")
		return
	}

	if userErr != nil {
		obj.OnGetUserFailed(w, r, obj, outputProp)
		obj.HandleError(w, r, outputProp, 2002, reqErr.Error())
		return
	}
	//
	//
	if key != "" || sign != "" {
		w.Header().Set("Cache-Control", "public, max-age=2592000")
	}
	if includePrivate == true {
		outputProp = miniprop.NewMiniPropFromMap(usrObj.ToMapAll())
	} else {
		outputProp = miniprop.NewMiniPropFromMap(usrObj.ToMapPublic())
	}
	errSuc := obj.OnGetUserSuccess(w, r, obj, usrObj, outputProp)
	if errSuc != nil {
		obj.OnGetUserFailed(w, r, obj, outputProp)
		obj.HandleError(w, r, outputProp, 2002, errSuc.Error())
		return
	}
	w.Write(outputProp.ToJson())
	return
}

func (obj *UserHandler) CheckLogin(r *http.Request, token string) minisession.CheckLoginIdResult {
	return obj.GetSessionMgr().CheckLoginId(appengine.NewContext(r), token, minisession.MakeAccessTokenConfigFromRequest(r))
}
