package handler

import (
	"net/http"

	//	"strings"

	//	miniuser "github.com/firefirestyle/go.miniuser/user"
	//"google.golang.org/appengine"
	"io/ioutil"

	"github.com/firefirestyle/go.miniprop"
	"github.com/firefirestyle/go.minisession"
	"google.golang.org/appengine"
)

func (obj *UserHandler) HandleTwitterRequestToken(w http.ResponseWriter, r *http.Request) {
	obj.twitterHandler.HandleLoginEntry(w, r)
}

func (obj *UserHandler) HandleTwitterCallbackToken(w http.ResponseWriter, r *http.Request) {
	obj.twitterHandler.HandleLoginExit(w, r)
}

func (obj *UserHandler) HandleFacebookRequestToken(w http.ResponseWriter, r *http.Request) {
	obj.facebookHandler.HandleLoginEntry(w, r)
}

func (obj *UserHandler) HandleFacebookCallbackToken(w http.ResponseWriter, r *http.Request) {
	obj.facebookHandler.HandleLoginExit(w, r)
}

func (obj *UserHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	propObj := miniprop.NewMiniPropFromJson(bodyBytes)
	token := propObj.GetString("token", "")

	obj.sessionMgr.Logout(appengine.NewContext(r), token, minisession.MakeAccessTokenConfigFromRequest(r))
}
