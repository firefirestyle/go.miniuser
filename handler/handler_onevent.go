package handler

import (
	"net/http"

	///	miniblob "github.com/firefirestyle/go.miniblob/blob"
	//	blobhandler "github.com/firefirestyle/go.miniblob/handler"
	//	"github.com/firefirestyle/go.minioauth/facebook"
	//	"github.com/firefirestyle/go.minioauth/twitter"
	//	"github.com/firefirestyle/go.minipointer"
	"github.com/firefirestyle/go.miniprop"
	"github.com/firefirestyle/go.miniuser/user"
	//	"github.com/firefirestyle/go.miniuser/handler"
	//	"github.com/firefirestyle/go.minisession"
	//	miniuser "github.com/firefirestyle/go.miniuser/user"
	//	"golang.org/x/net/context"
	//	"google.golang.org/appengine/log"
	//
	//	"crypto/sha1"
)

//
func (obj *UserHandler) AddOnGetUserRequest(f func(w http.ResponseWriter, r *http.Request, h *UserHandler, o *miniprop.MiniProp) error) {
	obj.onEvents.OnGetUserRequestList = append(obj.onEvents.OnGetUserRequestList, f)
}

func (obj *UserHandler) OnGetUserRequest(w http.ResponseWriter, r *http.Request, h *UserHandler, o *miniprop.MiniProp) error {
	for _, f := range obj.onEvents.OnGetUserRequestList {
		e := f(w, r, h, o)
		if e != nil {
			return e
		}
	}
	return nil
}

//
func (obj *UserHandler) AddOnGetUserFailed(f func(w http.ResponseWriter, r *http.Request, h *UserHandler, o *miniprop.MiniProp)) {
	obj.onEvents.OnGetUserFailedList = append(obj.onEvents.OnGetUserFailedList, f)
}

func (obj *UserHandler) OnGetUserFailed(w http.ResponseWriter, r *http.Request, h *UserHandler, o *miniprop.MiniProp) {
	for _, f := range obj.onEvents.OnGetUserFailedList {
		f(w, r, h, o)
	}
}

//
func (obj *UserHandler) AddOnGetUserSuccess(f func(w http.ResponseWriter, r *http.Request, h *UserHandler, i *user.User, o *miniprop.MiniProp) error) {
	obj.onEvents.OnGetUserSuccessList = append(obj.onEvents.OnGetUserSuccessList, f)
}

func (obj *UserHandler) OnGetUserSuccess(w http.ResponseWriter, r *http.Request, h *UserHandler, i *user.User, o *miniprop.MiniProp) error {
	for _, f := range obj.onEvents.OnGetUserSuccessList {
		e := f(w, r, h, i, o)
		if e != nil {
			return e
		}
	}
	return nil
}
