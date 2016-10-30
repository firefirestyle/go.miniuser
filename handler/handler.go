package handler

import (

	//	"strings"

	"errors"
	"net/http"
	"strings"

	miniblob "github.com/firefirestyle/go.miniblob/blob"
	blobhandler "github.com/firefirestyle/go.miniblob/handler"
	"github.com/firefirestyle/go.minipointer"
	"github.com/firefirestyle/go.miniprop"
	"github.com/firefirestyle/go.minisession"
	miniuser "github.com/firefirestyle/go.miniuser/user"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type UserHandler struct {
	manager     *miniuser.UserManager
	relayIdMgr  *minipointer.PointerManager
	sessionMgr  *minisession.SessionManager
	blobHandler *blobhandler.BlobHandler
}

type UserHandlerManagerConfig struct {
	ProjectId       string
	UserKind        string
	RelayIdKind     string
	SessionKind     string
	BlobKind        string
	BlobPointerKind string
	BlobSign        string
}

type UserHandlerOnEvent struct {
	blobOnEvent blobhandler.BlobHandlerOnEvent
}

func NewUserHandler(callbackUrl string, config UserHandlerManagerConfig, onEvents UserHandlerOnEvent) *UserHandler {
	if config.ProjectId == "" {
		config.ProjectId = "ffstyle"
	}
	if config.UserKind == "" {
		config.UserKind = "ffuser"
	}
	if config.RelayIdKind == "" {
		config.UserKind = "ffuser-pointer"
	}
	if config.SessionKind == "" {
		config.SessionKind = "ffuser-session"
	}
	if config.BlobKind == "" {
		config.BlobKind = "ffuser-blob"
	}
	if config.BlobPointerKind == "" {
		config.BlobPointerKind = config.BlobKind + "-pointer"
	}
	if config.BlobSign == "" {
		config.BlobSign = miniprop.MakeRandomId()
	}
	//

	ret := &UserHandler{
		manager: miniuser.NewUserManager(miniuser.UserManagerConfig{
			ProjectId:       config.ProjectId,
			UserKind:        config.UserKind,
			UserPointerKind: config.RelayIdKind,
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
		//		blobHandler: blobHandlerObj,
	}
	completeFunc := onEvents.blobOnEvent.OnBlobComplete
	onEvents.blobOnEvent.OnBlobComplete = func(w http.ResponseWriter, r *http.Request, outputProp *miniprop.MiniProp, hh *blobhandler.BlobHandler, blobObj *miniblob.BlobItem) error {
		dir := r.URL.Query().Get("dir")
		if true == strings.HasPrefix(dir, "/user") {
			ctx := appengine.NewContext(r)
			userName := strings.Replace(dir, "/user/", "", -1)
			userMgrObj := ret
			userObj, userErr := userMgrObj.GetManager().GetUserFromRelayId(ctx, userName)
			if userErr != nil {
				outputProp.SetString("error", "not found user")
				return userErr
			}
			userObj.SetIconUrl("key://" + blobObj.GetBlobKey())
			userMgrObj.GetManager().SaveUserWithImmutable(ctx, userObj)
			return completeFunc(w, r, outputProp, hh, blobObj)
		} else {
			return errors.New("unsupport")
		}
	}

	ret.blobHandler = blobhandler.NewBlobHandler(callbackUrl, config.BlobSign, miniblob.BlobManagerConfig{
		ProjectId:   config.ProjectId,
		Kind:        config.BlobKind,
		PointerKind: config.BlobPointerKind,
		CallbackUrl: callbackUrl,
	}, blobhandler.BlobHandlerOnEvent{})
	return ret
}

func (obj *UserHandler) GetSessionMgr() *minisession.SessionManager {
	return obj.sessionMgr
}

func (obj *UserHandler) GetManager() *miniuser.UserManager {
	return obj.manager
}

func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}
