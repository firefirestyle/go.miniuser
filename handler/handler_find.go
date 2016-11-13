package handler

import (
	"net/http"

	//	"strings"

	"github.com/firefirestyle/go.miniprop"
	"github.com/firefirestyle/go.miniuser/user"
	"google.golang.org/appengine"
)

func (obj *UserHandler) HandleFind(w http.ResponseWriter, r *http.Request) {
	propObj := miniprop.NewMiniProp()
	ctx := appengine.NewContext(r)
	values := r.URL.Query()
	cursor := values.Get("cursor")
	mode := values.Get("mode")
	projectId := values.Get("group")
	keyOnly := true
	//if mode != "0" {
	//	keyOnly = false
	//}
	var foundObj *user.FoundUser = nil
	if mode == "-point" {
		foundObj = obj.manager.FindUserWithPoint(ctx, cursor, projectId, "Main", keyOnly)
	} else {
		foundObj = obj.manager.FindUserWithNewOrder(ctx, cursor, projectId, keyOnly)
	}
	propObj.SetPropStringList("", "keys", foundObj.UserIds)
	propObj.SetPropString("", "cursorOne", foundObj.CursorOne)
	propObj.SetPropString("", "cursorNext", foundObj.CursorNext)
	//if keyOnly == false {
	//	// todo
	//}
	w.Write(propObj.ToJson())
}
