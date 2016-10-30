package handler

import (

	//	"strings"

	"github.com/firefirestyle/go.minipointer"
	"github.com/firefirestyle/go.minisession"
	miniuser "github.com/firefirestyle/go.miniuser/user"
	"golang.org/x/net/context"
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
	}
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
