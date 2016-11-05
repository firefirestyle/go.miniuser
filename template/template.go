package template

import (
	"net/http"

	"errors"

	blobhandler "github.com/firefirestyle/go.miniblob/handler"
	"github.com/firefirestyle/go.minioauth/twitter"
	"github.com/firefirestyle/go.miniprop"
	"github.com/firefirestyle/go.minisession"

	userhundler "github.com/firefirestyle/go.miniuser/handler"
	//
	"github.com/firefirestyle/go.minioauth/facebook"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	//
)

const (
	UrlTwitterTokenUrlRedirect  = "/api/v1/twitter/tokenurl/redirect"
	UrlTwitterTokenCallback     = "/api/v1/twitter/tokenurl/callback"
	UrlFacebookTokenUrlRedirect = "/api/v1/facebook/tokenurl/redirect"
	UrlFacebookTokenCallback    = "/api/v1/facebook/tokenurl/callback"
	UrlUserGet                  = "/api/v1/user/get"
	UrlUserFind                 = "/api/v1/user/find"
	UrlUserBlobGet              = "/api/v1/user/getblob"
	UrlUserRequestBlobUrl       = "/api/v1/user/requestbloburl"
	UrlUserCallbackBlobUrl      = "/api/v1/user/callbackbloburl"
	UrlMeLogout                 = "/api/v1/me/logout"
	UrlMeUpdate                 = "/api/v1/me/update"
	UrlMeGet                    = "/api/v1/me/get"
)

type UserTemplateConfig struct {
	GroupName                string
	KindBaseName             string
	TwitterConsumerKey       string
	TwitterConsumerSecret    string
	TwitterAccessToken       string
	TwitterAccessTokenSecret string
	FacebookAppSecret        string
	FacebookAppId            string
}
type UserTemplate struct {
	config         UserTemplateConfig
	userHandlerObj *userhundler.UserHandler
}

func NewUserTemplate(config UserTemplateConfig) *UserTemplate {
	if config.GroupName != "" {
		config.GroupName = "FFS"
	}
	if config.KindBaseName != "" {
		config.KindBaseName = "FFSUser"
	}

	return &UserTemplate{
		config: config,
	}
}

func (tmpObj *UserTemplate) CheckLogin(r *http.Request, input *miniprop.MiniProp) minisession.CheckLoginIdResult {
	ctx := appengine.NewContext(r)
	token := input.GetString("token", "")
	Debug(ctx, "CheckLogin ++>"+token)
	return tmpObj.GetUserHundlerObj(ctx).GetSessionMgr().CheckLoginId(ctx, token, minisession.MakeAccessTokenConfigFromRequest(r))
}

func (tmpObj *UserTemplate) GetUserHundlerObj(ctx context.Context) *userhundler.UserHandler {
	if tmpObj.userHandlerObj == nil {
		v := appengine.DefaultVersionHostname(ctx)
		if v == "127.0.0.1:8080" {
			v = "localhost:8080"
		}
		tmpObj.userHandlerObj = userhundler.NewUserHandler(UrlUserCallbackBlobUrl,
			userhundler.UserHandlerManagerConfig{ //
				ProjectId: tmpObj.config.GroupName,
				UserKind:  tmpObj.config.KindBaseName,
			})
		tmpObj.userHandlerObj.GetBlobHandler().GetBlobHandleEvent().OnBlobRequest =
			append(tmpObj.userHandlerObj.GetBlobHandler().GetBlobHandleEvent().OnBlobRequest,
				func(w http.ResponseWriter, r *http.Request, input *miniprop.MiniProp, output *miniprop.MiniProp, h *blobhandler.BlobHandler) (string, map[string]string, error) {
					ret := tmpObj.CheckLogin(r, input)
					if ret.IsLogin == false {
						return "", map[string]string{}, errors.New("Failed in token check")
					}
					return ret.AccessTokenObj.GetLoginId(), map[string]string{}, nil
				})
		tmpObj.userHandlerObj.AddFacebookSession(facebook.FacebookOAuthConfig{
			ConfigFacebookAppSecret: tmpObj.config.FacebookAppSecret,
			ConfigFacebookAppId:     tmpObj.config.FacebookAppId,
			SecretSign:              appengine.VersionID(ctx),
			CallbackUrl:             "http://" + v + "" + UrlFacebookTokenCallback,
		})
		tmpObj.userHandlerObj.AddTwitterSession(twitter.TwitterOAuthConfig{
			ConsumerKey:       tmpObj.config.TwitterConsumerKey,
			ConsumerSecret:    tmpObj.config.TwitterConsumerSecret,
			AccessToken:       tmpObj.config.TwitterAccessToken,
			AccessTokenSecret: tmpObj.config.TwitterAccessTokenSecret,
			CallbackUrl:       "http://" + appengine.DefaultVersionHostname(ctx) + "" + UrlTwitterTokenCallback,
			SecretSign:        appengine.VersionID(ctx),
		})
	}
	return tmpObj.userHandlerObj
}

func (tmpObj *UserTemplate) InitUserApi() {
	// twitter
	http.HandleFunc(UrlTwitterTokenUrlRedirect, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		tmpObj.GetUserHundlerObj(appengine.NewContext(r)).HandleTwitterRequestToken(w, r)
	})
	http.HandleFunc(UrlTwitterTokenCallback, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		tmpObj.GetUserHundlerObj(appengine.NewContext(r)).HandleTwitterCallbackToken(w, r)
	})
	// facebook
	http.HandleFunc(UrlFacebookTokenUrlRedirect, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		tmpObj.GetUserHundlerObj(appengine.NewContext(r)).HandleFacebookRequestToken(w, r)
	})
	http.HandleFunc(UrlFacebookTokenCallback, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		tmpObj.GetUserHundlerObj(appengine.NewContext(r)).HandleFacebookCallbackToken(w, r)
	})
	// user
	http.HandleFunc(UrlUserGet, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		tmpObj.GetUserHundlerObj(appengine.NewContext(r)).HandleGet(w, r)
	})
	http.HandleFunc(UrlUserFind, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		tmpObj.GetUserHundlerObj(appengine.NewContext(r)).HandleFind(w, r)
	})
	http.HandleFunc(UrlUserRequestBlobUrl, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		tmpObj.GetUserHundlerObj(appengine.NewContext(r)).HandleBlobRequestToken(w, r)
	})
	http.HandleFunc(UrlUserCallbackBlobUrl, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		tmpObj.GetUserHundlerObj(appengine.NewContext(r)).HandleBlobUpdated(w, r)
	})
	http.HandleFunc(UrlUserBlobGet, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		tmpObj.GetUserHundlerObj(appengine.NewContext(r)).HandleBlobGet(w, r)
	})

	// me
	http.HandleFunc(UrlMeLogout, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		tmpObj.GetUserHundlerObj(appengine.NewContext(r)).HandleLogout(w, r)
	})

	http.HandleFunc(UrlMeUpdate, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		tmpObj.GetUserHundlerObj(appengine.NewContext(r)).HandleUpdateInfo(w, r)
	})

	http.HandleFunc(UrlMeGet, func(w http.ResponseWriter, r *http.Request) {
		Debug(appengine.NewContext(r), "UrlMeGet called -->")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		tmpObj.GetUserHundlerObj(appengine.NewContext(r)).HandleGetMe(w, r)
	})
}

func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}
