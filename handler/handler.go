package handler

import (
	"net/http"

	miniblob "github.com/firefirestyle/go.miniblob/blob"
	blobhandler "github.com/firefirestyle/go.miniblob/handler"
	"github.com/firefirestyle/go.minioauth/facebook"
	"github.com/firefirestyle/go.minioauth/twitter"
	"github.com/firefirestyle/go.minipointer"
	"github.com/firefirestyle/go.miniprop"
	"github.com/firefirestyle/go.minisession"
	miniuser "github.com/firefirestyle/go.miniuser/user"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

type UserHandler struct {
	manager         *miniuser.UserManager
	relayIdMgr      *minipointer.PointerManager
	sessionMgr      *minisession.SessionManager
	blobHandler     *blobhandler.BlobHandler
	twitterHandler  *twitter.TwitterHandler
	facebookHandler *facebook.FacebookHandler
	completeFunc    func(w http.ResponseWriter, r *http.Request, outputProp *miniprop.MiniProp, hh *blobhandler.BlobHandler, blobObj *miniblob.BlobItem) error
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
}

func NewUserHandler(callbackUrl string, //
	config UserHandlerManagerConfig) *UserHandler {
	if config.ProjectId == "" {
		config.ProjectId = "ffstyle"
	}
	if config.UserKind == "" {
		config.UserKind = "ffuser"
	}
	if config.RelayIdKind == "" {
		config.RelayIdKind = config.UserKind + "-pointer"
	}
	if config.SessionKind == "" {
		config.SessionKind = config.UserKind + "-session"
	}
	if config.BlobKind == "" {
		config.BlobKind = config.UserKind + "-blob"
	}
	if config.BlobPointerKind == "" {
		config.BlobPointerKind = config.UserKind + "-blob-pointer"
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
		blobHandler: blobhandler.NewBlobHandler(callbackUrl, config.BlobSign, miniblob.BlobManagerConfig{
			ProjectId:   config.ProjectId,
			Kind:        config.BlobKind,
			PointerKind: config.BlobPointerKind,
			CallbackUrl: callbackUrl,
		}),
	}

	ret.blobHandler.GetBlobHandleEvent().OnBlobComplete = ret.OnBlobComplete
	return ret
}

func (obj *UserHandler) GetBlobHandler() *blobhandler.BlobHandler {
	return obj.blobHandler
}

func (obj *UserHandler) AddTwitterSession(twitterConfig twitter.TwitterOAuthConfig) {
	obj.twitterHandler = obj.NewTwitterHandlerObj(twitterConfig)
}

func (obj *UserHandler) AddFacebookSession(facebookConfig facebook.FacebookOAuthConfig) {
	obj.facebookHandler = obj.NewFacebookHandlerObj(facebookConfig)
}

func (obj *UserHandler) GetUserHandleEvent() {
	//obj.
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
