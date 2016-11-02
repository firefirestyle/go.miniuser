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
	values := r.URL.Query()
	token := values.Get(inputProp.GetString("token", ""))
	loginResult := obj.GetSessionMgr().CheckLoginId(ctx, token, minisession.MakeAccessTokenConfigFromRequest(r))
	userName := loginResult.AccessTokenObj.GetUserName()
	if loginResult.IsLogin == false {
		userName = ""
	}
	obj.HandleGetBase(w, r, userName, "", "", false)
}

/*

 */
func (obj *UserHandler) HandleGetBase(w http.ResponseWriter, r *http.Request, userName string, sign string, key string, includePrivate bool) {
	ctx := appengine.NewContext(r)
	var usrObj *miniuser.User = nil
	var userErr error = nil

	if userName != "" {
		if sign == "" {
			usrObj, userErr = obj.GetManager().GetUserFromRelayId(ctx, userName)
		} else {
			usrObj, userErr = obj.GetManager().GetUserFromSign(ctx, userName, sign)
		}
	} else if key != "" {
		usrObj, userErr = obj.GetManager().GetUserFromKey(ctx, key)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrong Request"))
		return
	}

	if userErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not found User 1"))
		return
	} else {
		if key != "" || sign != "" {
			w.Header().Set("Cache-Control", "public, max-age=2592000")
		}
		if includePrivate == true {
			w.Write(usrObj.ToJson())
		} else {
			w.Write(usrObj.ToJsonPublic())
		}
		return
	}
}
