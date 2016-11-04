package handler

import (
	"net/http"

	"strings"

	"errors"

	"io/ioutil"

	miniblob "github.com/firefirestyle/go.miniblob/blob"
	blobhandler "github.com/firefirestyle/go.miniblob/handler"
	"github.com/firefirestyle/go.miniprop"
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

func (obj *UserHandler) MakeDir(userName string, dir string) string {
	if dir == "" {
		return "/user/" + userName
	}
	if strings.HasPrefix(dir, "/") {
		dir = dir[1:]
	}
	if strings.HasSuffix(dir, "/") {
		dir = dir[0 : len(dir)-1]
	}

	return "/user/" + userName + "/" + dir
}

func (obj *UserHandler) HandleBlobRequestToken(w http.ResponseWriter, r *http.Request) {
	//
	// load param from json
	params, _ := ioutil.ReadAll(r.Body)
	inputPropObj := miniprop.NewMiniPropFromJson(params)
	token := inputPropObj.GetString("token", "")
	dir := inputPropObj.GetString("dir", "")
	name := inputPropObj.GetString("file", "")
	userName := inputPropObj.GetString("userName", "")
	if token != "" {
		keyInfo, _ := obj.GetSessionMgr().MakeLoginIdInfoFromLoginId(token)
		//		if keyInfo != nil {
		userName = keyInfo.UserName
		Debug(appengine.NewContext(r), ">keyinfo>:"+keyInfo.UserName)
		//	}
	}
	//
	//
	Debug(appengine.NewContext(r), ">>"+dir+">>"+name+">>"+userName+"::"+token)
	obj.blobHandler.HandleBlobRequestTokenFromParams(w, r, obj.MakeDir(userName, dir), name, inputPropObj)
}

func (obj *UserHandler) HandleBlobUpdated(w http.ResponseWriter, r *http.Request) {
	//
	ctx := appengine.NewContext(r)
	Debug(ctx, "callbeck AAAA")
	obj.blobHandler.HandleUploaded(w, r)
}

func (obj *UserHandler) HandleBlobGet(w http.ResponseWriter, r *http.Request) {
	//
	values := r.URL.Query()
	key := values.Get("key")
	dir := values.Get("dir")
	file := values.Get("file")
	userName := values.Get("userName")

	obj.blobHandler.HandleGetBase(w, r, key, obj.MakeDir(userName, dir), file)
}

func (userMgrObj *UserHandler) OnBlobComplete(w http.ResponseWriter, r *http.Request, outputProp *miniprop.MiniProp, hh *blobhandler.BlobHandler, blobObj *miniblob.BlobItem) error {
	dir := r.URL.Query().Get("dir")
	if true == strings.HasPrefix(dir, "/user") {
		ctx := appengine.NewContext(r)
		userName := userMgrObj.GetUserNameFromDir(dir)
		Debug(ctx, "dir::"+dir+";;username::"+userName)

		userObj, userErr := userMgrObj.GetManager().GetUserFromRelayId(ctx, userName)
		if userErr != nil {
			outputProp.SetString("error", "not found user")
			return userErr
		}
		userObj.SetIconUrl("key://" + blobObj.GetBlobKey())
		userMgrObj.GetManager().SaveUserWithImmutable(ctx, userObj)
		if userMgrObj.completeFunc != nil {
			return userMgrObj.completeFunc(w, r, outputProp, hh, blobObj)
		} else {
			return nil
		}
	} else {
		return errors.New("unsupport")
	}
}
