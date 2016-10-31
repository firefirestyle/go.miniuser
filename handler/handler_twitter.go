package handler

import (
	"net/http"

	//	"strings"

	//	miniuser "github.com/firefirestyle/go.miniuser/user"
	//"google.golang.org/appengine"
)

func (obj *UserHandler) HandleTwitterRequestToken(w http.ResponseWriter, r *http.Request) {
	obj.twitterHandler.HandleLoginEntry(w, r)
}

func (obj *UserHandler) HandleTwitterCallbackToken(w http.ResponseWriter, r *http.Request) {
	obj.twitterHandler.HandleLoginExit(w, r)
}
