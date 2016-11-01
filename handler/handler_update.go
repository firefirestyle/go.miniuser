package handler

import (
	"net/http"

	//	"strings"

	//	miniuser "github.com/firefirestyle/go.miniuser/user"
	"io/ioutil"

	"github.com/firefirestyle/go.miniprop"
	"google.golang.org/appengine"
)

func (obj *UserHandler) HandleUpdateInfo(w http.ResponseWriter, r *http.Request) {
	outputProp := miniprop.NewMiniProp()
	v, _ := ioutil.ReadAll(r.Body)
	inputProp := miniprop.NewMiniPropFromJson(v)
	ctx := appengine.NewContext(r)
	userName := inputProp.GetString("userName", "")
	displayName := inputProp.GetString("displayName", "")
	content := inputProp.GetString("content", "")

	usrObj, userErr := obj.GetManager().GetUserFromRelayId(ctx, userName)
	if userErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		outputProp.SetInt("errorCode", 2001)
		outputProp.SetString("errorMessage", userErr.Error())
		w.Write(outputProp.ToJson())
		return
	}
	usrObj.SetDisplayName(displayName)
	usrObj.SetCont(content)
	nextUserObj, nextUserErr := obj.GetManager().SaveUserWithImmutable(ctx, usrObj)
	if nextUserErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		outputProp.SetInt("errorCode", 2002)
		outputProp.SetString("errorMessage", nextUserErr.Error())
		w.Write(outputProp.ToJson())
	} else {
		w.WriteHeader(http.StatusOK)
		//		outputProp.SetString("sign", nextUserObj.GetSign())
		//		outputProp.SetString("userName", nextUserObj.GetUserName())
		//		outputProp.SetString("key", nextUserObj.GetStringId())
		w.Write(nextUserObj.ToJsonPublic())
	}
}
