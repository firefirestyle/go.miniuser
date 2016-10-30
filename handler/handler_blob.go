package handler

import (
	"net/http"

	"strings"

	"google.golang.org/appengine"
)

func (obj *UserHandler) GetUserNameFromDir(dir string) string {
	if false == strings.HasPrefix(dir, "/user/") {
		return ""
	}
	t1 := strings.Replace(dir, "/user/", "", 1)
	t2 := strings.Index(t1, "/")
	if t2 == -1 {
		t2 = len(t1)
	}

	return t1[0:t2]
}

func (obj *UserHandler) GetDirFromDir(dir string) string {
	if false == strings.HasPrefix(dir, "/user/") {
		return ""
	}
	t1 := strings.Replace(dir, "/user/", "", 1)
	t2 := strings.Index(t1, "/")
	if t2 == -1 {
		t2 = 0
	}

	return t1[t2:len(t1)]
}

func (obj *UserHandler) HandleBlobRequestToken(w http.ResponseWriter, r *http.Request) {
	//
	// load param from json
	articleId := r.URL.Query().Get("userName")
	dir := r.URL.Query().Get("dir")
	name := r.URL.Query().Get("file")
	//
	// todo check articleId

	//
	//
	obj.blobHandler.HandleBlobRequestTokenFromParams(w, r, "/user/"+articleId+"/"+dir, name)
}

func (obj *UserHandler) HandleBlobUpdated(w http.ResponseWriter, r *http.Request) {
	//
	ctx := appengine.NewContext(r)
	Debug(ctx, "callbeck AAAA")
	obj.blobHandler.HandleUploaded(w, r)
}

func (obj *UserHandler) HandleBlobGet(w http.ResponseWriter, r *http.Request) {
	//
	ctx := appengine.NewContext(r)
	Debug(ctx, "callbeck AAAA")
	obj.blobHandler.HandleGet(w, r)
}
