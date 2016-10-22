package miniuser

import (
	"net/http"

	"github.com/firefirestyle/go.miniprop"
	"google.golang.org/appengine"
)

type UserHandler struct {
	manager *UserManager
}

type UserHandlerOnEvent struct {
}

func NewUserHandler(config UserManagerConfig, onEvents UserHandlerOnEvent) *UserHandler {
	return &UserHandler{
		manager: NewUserManager(config),
	}
}

func (obj *UserHandler) GetManager() *UserManager {
	return obj.manager
}

func (obj *UserHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	values := r.URL.Query()
	userName := values.Get("userName")
	usrObj, userErr := obj.manager.GetUserFromUserName(ctx, userName)

	if userErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not found User"))
		return
	} else {
		cont, contErr := usrObj.ToJsonPublic()
		if contErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Not found User"))
			return
		} else {
			w.Write(cont)
			return
		}
	}
}

func (obj *UserHandler) HandleFind(w http.ResponseWriter, r *http.Request) {
	propObj := miniprop.NewMiniProp()
	ctx := appengine.NewContext(r)
	values := r.URL.Query()
	cursor := values.Get("cursor")
	mode := values.Get("keyOnly")
	keyOnly := false
	if mode != "0" {
		keyOnly = true
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
