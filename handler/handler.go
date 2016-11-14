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
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	//
	"crypto/sha1"
)

type UserHandler struct {
	manager         *miniuser.UserManager
	relayIdMgr      *minipointer.PointerManager
	sessionMgr      *minisession.SessionManager
	blobHandler     *blobhandler.BlobHandler
	twitterHandler  *twitter.TwitterHandler
	facebookHandler *facebook.FacebookHandler
	onEvents        UserHandlerOnEvent
	completeFunc    func(w http.ResponseWriter, r *http.Request, outputProp *miniprop.MiniProp, hh *blobhandler.BlobHandler, blobObj *miniblob.BlobItem) error
}

type UserHandlerManagerConfig struct {
	RootGroup                  string
	UserKind                   string
	RelayIdKind                string
	SessionKind                string
	BlobKind                   string
	BlobPointerKind            string
	BlobSign                   string
	MemcachedOnlyInBlobPointer bool
	LengthHash                 int
	//
	//
	//MasterKey     []string
	//MasterUser    []string
	//MasterAccount []string
}

type UserHandlerOnEvent struct {
	OnGetUserRequestList       []func(w http.ResponseWriter, r *http.Request, h *UserHandler, o *miniprop.MiniProp) error
	OnGetUserFailedList        []func(w http.ResponseWriter, r *http.Request, h *UserHandler, o *miniprop.MiniProp)
	OnGetUserSuccessList       []func(w http.ResponseWriter, r *http.Request, h *UserHandler, i *miniuser.User, o *miniprop.MiniProp) error
	OnUpdateUserRequestList    []func(w http.ResponseWriter, r *http.Request, h *UserHandler, i *miniprop.MiniProp, o *miniprop.MiniProp) error
	OnUpdateUserFailedList     []func(w http.ResponseWriter, r *http.Request, h *UserHandler, i *miniprop.MiniProp, o *miniprop.MiniProp)
	OnUpdateUserBeforeSaveList []func(w http.ResponseWriter, r *http.Request, h *UserHandler, u *miniuser.User, i *miniprop.MiniProp, o *miniprop.MiniProp) error
	OnUpdateUserSuccessList    []func(w http.ResponseWriter, r *http.Request, h *UserHandler, u *miniuser.User, i *miniprop.MiniProp, o *miniprop.MiniProp) error
}

func NewUserHandler(callbackUrl string, //
	config UserHandlerManagerConfig) *UserHandler {
	if config.RootGroup == "" {
		config.RootGroup = "ffstyle"
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
		config.BlobSign = string(sha1.New().Sum([]byte("" + config.RootGroup + "&" + config.UserKind)))
	}
	//

	ret := &UserHandler{
		manager: miniuser.NewUserManager(miniuser.UserManagerConfig{
			RootGroup:       config.RootGroup,
			UserKind:        config.UserKind,
			UserPointerKind: config.RelayIdKind,
			LengthHash:      config.LengthHash,
		}),
		relayIdMgr: minipointer.NewPointerManager( //
			minipointer.PointerManagerConfig{
				Kind:      config.RelayIdKind,
				RootGroup: config.RootGroup,
			}),
		sessionMgr: minisession.NewSessionManager(minisession.SessionManagerConfig{
			Kind:      config.SessionKind,
			RootGroup: config.RootGroup,
		}),
		blobHandler: blobhandler.NewBlobHandler(callbackUrl, config.BlobSign, miniblob.BlobManagerConfig{
			RootGroup:              config.RootGroup,
			Kind:                   config.BlobKind,
			PointerKind:            config.BlobPointerKind,
			CallbackUrl:            callbackUrl,
			MemcachedOnlyInPointer: config.MemcachedOnlyInBlobPointer,
			HashLength:             10,
		}),
		onEvents: UserHandlerOnEvent{},
	}

	ret.blobHandler.AddOnBlobComplete(ret.OnBlobComplete)
	return ret
}

func (obj *UserHandler) GetPointerManager() *minipointer.PointerManager {
	return obj.relayIdMgr
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

func (obj *UserHandler) GetSessionMgr() *minisession.SessionManager {
	return obj.sessionMgr
}

func (obj *UserHandler) GetManager() *miniuser.UserManager {
	return obj.manager
}

func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}

func (obj *UserHandler) CheckLoginFromToken(r *http.Request, token string, useIp bool) minisession.CheckLoginIdResult {
	ctx := appengine.NewContext(r)
	return obj.GetSessionMgr().CheckLoginId(ctx, token, minisession.MakeAccessTokenConfigFromRequest(r), useIp)
}

func (obj *UserHandler) HandleError(w http.ResponseWriter, r *http.Request, outputProp *miniprop.MiniProp, errorCode int, errorMessage string) {
	//
	//
	if outputProp == nil {
		outputProp = miniprop.NewMiniProp()
	}
	if errorCode != 0 {
		outputProp.SetInt("errorCode", errorCode)
	}
	if errorMessage != "" {
		outputProp.SetString("errorMessage", errorMessage)
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write(outputProp.ToJson())
}
