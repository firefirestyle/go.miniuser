package handler

/*
import (
	"net/http"

	//	"strings"

	//	"errors"

	//"io/ioutil"

	//	miniblob "github.com/firefirestyle/go.miniblob/blob"
	//	blobhandler "github.com/firefirestyle/go.miniblob/handler"
	"github.com/firefirestyle/go.miniprop"
	"google.golang.org/appengine"
)


func (obj *UserHandler) HandleRegistAsMaster(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	inputProp := miniprop.NewMiniPropFromJsonReader(r.Body)
	token := inputProp.GetString("token", "")
	loginCheckResult := obj.CheckLogin(r, token, true)
	if loginCheckResult.IsLogin == false {
		obj.HandleError(w, r, miniprop.NewMiniProp(), 0, "failed to check login status")
		return
	}
	userName := loginCheckResult.AccessTokenObj.GetUserName()
	founded := obj.GetManager().FindAuthPointer(ctx, userName)
	for _, k := range founded.Keys {
		if k == "" {
			kInfo := obj.GetPointerManager().GetKeyInfoFromStringId(k)
			kInfo.Identify
		}
	}
}
*/
