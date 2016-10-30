package handler

import (
	"net/http"

	//	"strings"

	"github.com/firefirestyle/go.miniprop"
	"google.golang.org/appengine"
)

func (obj *UserHandler) HandleFind(w http.ResponseWriter, r *http.Request) {
	propObj := miniprop.NewMiniProp()
	ctx := appengine.NewContext(r)
	values := r.URL.Query()
	cursor := values.Get("cursor")
	mode := values.Get("keyOnly")
	keyOnly := true
	if mode != "0" {
		keyOnly = false
	}

	foundObj := obj.manager.FindUserWithNewOrder(ctx, cursor, keyOnly)
	propObj.SetPropStringList("", "keys", foundObj.UserIds)
	propObj.SetPropString("", "cursorOne", foundObj.CursorOne)
	propObj.SetPropString("", "cursorOne", foundObj.CursorNext)
	if keyOnly == false {
		// todo
	}
	w.Write(propObj.ToJson())
}
