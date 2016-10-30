package hundler

import (
	"net/http"

	//	"strings"

	"github.com/firefirestyle/go.minipointer"
	"github.com/firefirestyle/go.miniprop"
	"github.com/firefirestyle/go.minisession"
	miniuser "github.com/firefirestyle/go.miniuser/user"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type UserHandler struct {
	manager    *miniuser.UserManager
	relayIdMgr *minipointer.PointerManager
	sessionMgr *minisession.SessionManager
}

type UserHandlerManagerConfig struct {
	ProjectId   string
	UserKind    string
	RelayIdKind string
	SessionKind string
}

type UserHandlerOnEvent struct {
}

func NewUserHandler(config UserHandlerManagerConfig, onEvents UserHandlerOnEvent) *UserHandler {
	return &UserHandler{
		manager: miniuser.NewUserManager(miniuser.UserManagerConfig{
			ProjectId: config.ProjectId,
			UserKind:  config.UserKind,
		}),
		relayIdMgr: minipointer.NewPointerManager( //
			minipointer.PointerManagerConfig{
				Kind:      config.RelayIdKind,
				ProjectId: config.ProjectId,
			}),
		sessionMgr: minisession.NewSessionManager(minisession.SessionManagerConfig{
			Kind:      config.SessionKind,
			ProjectId: config.ProjectId,
		}),
	}
}

func (obj *UserHandler) GetSessionMgr() *minisession.SessionManager {
	return obj.sessionMgr
}

func (obj *UserHandler) GetManager() *miniuser.UserManager {
	return obj.manager
}

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
			usrObj, userErr = obj.GetUserFromRelayId(ctx, userName)
		} else {
			usrObj, userErr = obj.GetUserFromSign(ctx, userName, sign)
		}
	} else if key != "" {
		usrObj, userErr = obj.GetUserFromKey(ctx, key)
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
		cont, contErr := usrObj.ToJsonPublic()
		if contErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Not found User 2"))
			return
		} else {
			if key != "" || sign != "" {
				w.Header().Set("Cache-Control", "public, max-age=2592000")
			}
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
func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}
