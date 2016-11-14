package template

import (
	"net/http"

	"errors"

	miniblob "github.com/firefirestyle/go.miniblob/blob"
	blobhandler "github.com/firefirestyle/go.miniblob/handler"
	"github.com/firefirestyle/go.minioauth/twitter"
	"github.com/firefirestyle/go.miniprop"
	"github.com/firefirestyle/go.minisession"

	"sync"

	"github.com/firefirestyle/go.minioauth/facebook"
	userhundler "github.com/firefirestyle/go.miniuser/handler"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
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
	GroupName       string
	KindBaseName    string
	PrivateKey      string
	AllowInvalidSSL bool
	//
	MasterKey     []string
	MasterUser    []string
	MasterAccount []string

	//
	TwitterConsumerKey       string
	TwitterConsumerSecret    string
	TwitterAccessToken       string
	TwitterAccessTokenSecret string
	FacebookAppSecret        string
	FacebookAppId            string

	MemcachedOnlyInBlobPointer bool
}

type UserTemplate struct {
	config         UserTemplateConfig
	userHandlerObj *userhundler.UserHandler
	initOpt        func(context.Context)
	m              *sync.Mutex
}

func NewUserTemplate(config UserTemplateConfig) *UserTemplate {
	if config.GroupName == "" {
		config.GroupName = "FFS"
	}
	if config.KindBaseName == "" {
		config.KindBaseName = "FFSUser"
	}

	return &UserTemplate{
		config:  config,
		initOpt: func(context.Context) {},
		m:       new(sync.Mutex),
	}
}

func (tmpObj *UserTemplate) SetInitFunc(f func(ctx context.Context)) {
	tmpObj.m.Lock()
	defer tmpObj.m.Unlock()
	tmpObj.initOpt = f
}

func (tmpObj *UserTemplate) InitalizeTemplate(ctx context.Context) {

	if tmpObj.initOpt == nil {
		return
	}
	tmpObj.m.Lock()
	defer tmpObj.m.Unlock()
	tmpObj.GetUserHundlerObj(ctx)
	if tmpObj.initOpt != nil {
		tmpObj.initOpt(ctx)
	}
	tmpObj.initOpt = nil
}

func (tmpObj *UserTemplate) CheckLogin(r *http.Request, input *miniprop.MiniProp, useIp bool) minisession.CheckLoginIdResult {
	ctx := appengine.NewContext(r)
	token := input.GetString("token", "")
	Debug(ctx, "CheckLogin ++>"+token)
	return tmpObj.CheckLoginFromToken(r, token, useIp)
}

func (tmpObj *UserTemplate) CheckLoginFromToken(r *http.Request, token string, useIp bool) minisession.CheckLoginIdResult {
	ctx := appengine.NewContext(r)
	return tmpObj.GetUserHundlerObj(ctx).GetSessionMgr().CheckLoginId(ctx, token, minisession.MakeAccessTokenConfigFromRequest(r), useIp)
}

func (tmpObj *UserTemplate) GetUserHundlerObj(ctx context.Context) *userhundler.UserHandler {
	if tmpObj.userHandlerObj == nil {
		v := appengine.DefaultVersionHostname(ctx)
		scheme := "https"
		if v == "127.0.0.1:8080" {
			v = "localhost:8080"
			scheme = "http"
		}

		tmpObj.userHandlerObj = userhundler.NewUserHandler(UrlUserCallbackBlobUrl,
			userhundler.UserHandlerManagerConfig{ //
				RootGroup:                  tmpObj.config.GroupName,
				UserKind:                   tmpObj.config.KindBaseName,
				BlobSign:                   tmpObj.config.PrivateKey,
				MemcachedOnlyInBlobPointer: tmpObj.config.MemcachedOnlyInBlobPointer,
				LengthHash:                 9,
			})
		tmpObj.userHandlerObj.GetBlobHandler().AddOnBlobRequest(
			func(w http.ResponseWriter, r *http.Request, input *miniprop.MiniProp, output *miniprop.MiniProp, h *blobhandler.BlobHandler) (map[string]string, error) {
				ret := tmpObj.CheckLogin(r, input, true)
				if ret.IsLogin == false {
					return map[string]string{}, errors.New("Failed in token check")
				}
				return map[string]string{"tk": ret.AccessTokenObj.GetLoginId()}, nil
			})
		tmpObj.userHandlerObj.GetBlobHandler().AddOnBlobComplete(func(w http.ResponseWriter, r *http.Request, p *miniprop.MiniProp, h *blobhandler.BlobHandler, i *miniblob.BlobItem) error {
			pp := tmpObj.CheckLoginFromToken(r, r.FormValue("tk"), false)
			if pp.IsLogin == true {
				return nil
			} else {
				return errors.New("errors:" + r.FormValue("tk"))
			}
		})
		tmpObj.userHandlerObj.AddFacebookSession(facebook.FacebookOAuthConfig{
			ConfigFacebookAppSecret: tmpObj.config.FacebookAppSecret,
			ConfigFacebookAppId:     tmpObj.config.FacebookAppId,
			SecretSign:              appengine.VersionID(ctx),
			CallbackUrl:             "" + scheme + "://" + v + "" + UrlFacebookTokenCallback,
			AllowInvalidSSL:         tmpObj.config.AllowInvalidSSL,
		})
		tmpObj.userHandlerObj.AddTwitterSession(twitter.TwitterOAuthConfig{
			ConsumerKey:       tmpObj.config.TwitterConsumerKey,
			ConsumerSecret:    tmpObj.config.TwitterConsumerSecret,
			AccessToken:       tmpObj.config.TwitterAccessToken,
			AccessTokenSecret: tmpObj.config.TwitterAccessTokenSecret,
			CallbackUrl:       "" + scheme + "://" + appengine.DefaultVersionHostname(ctx) + "" + UrlTwitterTokenCallback,
			SecretSign:        appengine.VersionID(ctx),
			AllowInvalidSSL:   tmpObj.config.AllowInvalidSSL,
		})
	}
	return tmpObj.userHandlerObj
}

func (tmpObj *UserTemplate) InitUserApi() {
	// twitter
	http.HandleFunc(UrlTwitterTokenUrlRedirect, func(w http.ResponseWriter, r *http.Request) {
		tmpObj.InitalizeTemplate(appengine.NewContext(r))
		w.Header().Add("Access-Control-Allow-Origin", "*")
		tmpObj.GetUserHundlerObj(appengine.NewContext(r)).HandleTwitterRequestToken(w, r)
	})

	http.HandleFunc(UrlTwitterTokenCallback, func(w http.ResponseWriter, r *http.Request) {
		tmpObj.InitalizeTemplate(appengine.NewContext(r))
		w.Header().Add("Access-Control-Allow-Origin", "*")
		tmpObj.GetUserHundlerObj(appengine.NewContext(r)).HandleTwitterCallbackToken(w, r)
	})

	// facebook
	http.HandleFunc(UrlFacebookTokenUrlRedirect, func(w http.ResponseWriter, r *http.Request) {
		tmpObj.InitalizeTemplate(appengine.NewContext(r))
		w.Header().Add("Access-Control-Allow-Origin", "*")
		tmpObj.GetUserHundlerObj(appengine.NewContext(r)).HandleFacebookRequestToken(w, r)
	})

	http.HandleFunc(UrlFacebookTokenCallback, func(w http.ResponseWriter, r *http.Request) {
		tmpObj.InitalizeTemplate(appengine.NewContext(r))
		w.Header().Add("Access-Control-Allow-Origin", "*")
		tmpObj.GetUserHundlerObj(appengine.NewContext(r)).HandleFacebookCallbackToken(w, r)
	})

	// user
	http.HandleFunc(UrlUserGet, func(w http.ResponseWriter, r *http.Request) {
		tmpObj.InitalizeTemplate(appengine.NewContext(r))
		w.Header().Add("Access-Control-Allow-Origin", "*")
		tmpObj.GetUserHundlerObj(appengine.NewContext(r)).HandleGet(w, r)
	})

	http.HandleFunc(UrlUserFind, func(w http.ResponseWriter, r *http.Request) {
		tmpObj.InitalizeTemplate(appengine.NewContext(r))
		w.Header().Add("Access-Control-Allow-Origin", "*")
		tmpObj.GetUserHundlerObj(appengine.NewContext(r)).HandleFind(w, r)
	})

	http.HandleFunc(UrlUserRequestBlobUrl, func(w http.ResponseWriter, r *http.Request) {
		tmpObj.InitalizeTemplate(appengine.NewContext(r))
		w.Header().Add("Access-Control-Allow-Origin", "*")
		tmpObj.GetUserHundlerObj(appengine.NewContext(r)).HandleBlobRequestToken(w, r)
	})

	http.HandleFunc(UrlUserCallbackBlobUrl, func(w http.ResponseWriter, r *http.Request) {
		tmpObj.InitalizeTemplate(appengine.NewContext(r))
		w.Header().Add("Access-Control-Allow-Origin", "*")
		tmpObj.GetUserHundlerObj(appengine.NewContext(r)).HandleBlobUpdated(w, r)
	})

	http.HandleFunc(UrlUserBlobGet, func(w http.ResponseWriter, r *http.Request) {
		tmpObj.InitalizeTemplate(appengine.NewContext(r))
		w.Header().Add("Access-Control-Allow-Origin", "*")
		tmpObj.GetUserHundlerObj(appengine.NewContext(r)).HandleBlobGet(w, r)
	})

	// me
	http.HandleFunc(UrlMeLogout, func(w http.ResponseWriter, r *http.Request) {
		tmpObj.InitalizeTemplate(appengine.NewContext(r))
		w.Header().Add("Access-Control-Allow-Origin", "*")
		tmpObj.GetUserHundlerObj(appengine.NewContext(r)).HandleLogout(w, r)
	})

	http.HandleFunc(UrlMeUpdate, func(w http.ResponseWriter, r *http.Request) {
		tmpObj.InitalizeTemplate(appengine.NewContext(r))
		w.Header().Add("Access-Control-Allow-Origin", "*")
		tmpObj.GetUserHundlerObj(appengine.NewContext(r)).HandleUpdateInfo(w, r)
	})

	http.HandleFunc(UrlMeGet, func(w http.ResponseWriter, r *http.Request) {
		tmpObj.InitalizeTemplate(appengine.NewContext(r))
		w.Header().Add("Access-Control-Allow-Origin", "*")
		tmpObj.GetUserHundlerObj(appengine.NewContext(r)).HandleGetMe(w, r)
	})
}

func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}
