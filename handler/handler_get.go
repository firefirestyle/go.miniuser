package handler

import (
	"net/http"

	//	"strings"

	miniuser "github.com/firefirestyle/go.miniuser/user"
	"google.golang.org/appengine"
)

func (obj *UserHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	values := r.URL.Query()
	userName := values.Get("userName")
	sign := values.Get("sign")
	key := values.Get("key")
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
		cont := usrObj.ToJsonPublic()
		if key != "" || sign != "" {
			w.Header().Set("Cache-Control", "public, max-age=2592000")
		}
		w.Write(cont)
		return
	}
}
