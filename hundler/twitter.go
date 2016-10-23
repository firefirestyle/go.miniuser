package hundler

import (
	"net/http"

	"github.com/firefirestyle/go.minisession"

	"github.com/firefirestyle/go.minioauth/twitter"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

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
