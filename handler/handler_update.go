package handler

import (
	"net/http"

	//	"strings"

	//	miniuser "github.com/firefirestyle/go.miniuser/user"
	"google.golang.org/appengine"
)

func (obj *UserHandler) HandleUpdateInfo(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	values := r.URL.Query()
	userName := values.Get("userName")
	displayName := values.Get("displayName")
	//	content := values.Get("content")

	usrObj, userErr := obj.GetManager().GetUserFromRelayId(ctx, userName)
	if userErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not found User 1"))
		return
	}
	usrObj.SetDisplayName(displayName)
	obj.GetManager().SaveUserWithImmutable(ctx, usrObj)
}
