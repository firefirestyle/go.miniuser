package handler

import (
	"net/http"

	"errors"

	"github.com/firefirestyle/go.minioauth/twitter"
	"github.com/firefirestyle/go.minipointer"
	"github.com/firefirestyle/go.minisession"
	miniuser "github.com/firefirestyle/go.miniuser/user"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

//
//
//
func (obj *UserHandler) GetTwitterHandlerObj(ctx context.Context, callbackAddr string, config twitter.TwitterOAuthConfig) *twitter.TwitterHandler {
	twitterHandlerObj := twitter.NewTwitterHandler( //
		callbackAddr, config, twitter.TwitterHundlerOnEvent{
			OnFoundUser: func(w http.ResponseWriter, r *http.Request, handler *twitter.TwitterHandler, accesssToken *twitter.SendAccessTokenResult) map[string]string {
				ctx := appengine.NewContext(r)

				//
				//
				_, _, userObj, err1 := obj.LoginRegistFromTwitter(ctx, //
					accesssToken.GetScreenName(), //
					accesssToken.GetUserID(),     //
					accesssToken.GetOAuthToken()) //
				if err1 != nil {
					return map[string]string{"errcode": "2", "errindo": err1.Error()}
				}
				//
				//
				tokenObj, err := obj.sessionMgr.Login(ctx, //
					userObj.GetUserName(), //
					minisession.MakeAccessTokenConfigFromRequest(r))
				if err != nil {
					return map[string]string{"errcode": "1"}
				} else {
					return map[string]string{"token": "" + tokenObj.GetLoginId(), "userName": userObj.GetUserName()}
				}
			},
		})

	return twitterHandlerObj
}

//
//
//
func (obj *UserHandler) LoginRegistFromTwitter(ctx context.Context, screenName string, userId string, oauthToken string) (bool, *minipointer.Pointer, *miniuser.User, error) {
	relayIdObj := obj.relayIdMgr.GetPointerForTwitter(ctx, screenName, userId, oauthToken)
	needMake := false

	//
	// new userObj
	var err error = nil
	var userObj *miniuser.User = nil
	var pointerObj *minipointer.Pointer = nil
	if relayIdObj.GetValue() != "" {
		needMake = true
		//		Debug(ctx, "LoginRegistFromTwitter (1) :"+relayIdObj.GetUserName())
		pointerObj = obj.relayIdMgr.GetPointerForRelayId(ctx, relayIdObj.GetValue())
		if pointerObj.GetValue() != "" {
			userObj, err = obj.GetManager().GetUserFromUserName(ctx, pointerObj.GetValue(), pointerObj.GetSign())
			if err != nil {
				userObj = nil
			}
		}
	}
	if userObj == nil {
		userObj = obj.GetManager().NewNewUser(ctx, "")
		userObj.SetDisplayName(screenName)
		//		Debug(ctx, "LoginRegistFromTwitter (2) :"+userObj.GetUserName())
		pointerObj = obj.relayIdMgr.GetPointerForRelayId(ctx, userObj.GetUserName())
		pointerObj.SetValue(userObj.GetUserName())
		pointerObj.SetSign("")
		Debug(ctx, "LoginRegistFromTwitter :")
		err := pointerObj.Save(ctx)
		if err != nil {
			return needMake, nil, nil, errors.New("failed to save pointreobj : " + err.Error())
		}
		//
		// set username
		relayIdObj.SetValue(pointerObj.GetValue())
	}

	//
	// save relayId
	err = relayIdObj.Save(ctx)
	if err != nil {
		return needMake, nil, nil, errors.New("failed to save sessionobj : " + err.Error())
	}

	//
	// save user
	err = obj.GetManager().SaveUser(ctx, userObj)
	if err != nil {
		return needMake, nil, nil, errors.New("failed to save userobj : " + err.Error())
	}
	return needMake, relayIdObj, userObj, nil
}
