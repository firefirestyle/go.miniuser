package hundler

import (
	"errors"

	//	"strings"

	"strconv"
	"time"

	"github.com/firefirestyle/go.miniuser/relayid"
	miniuser "github.com/firefirestyle/go.miniuser/user"
	"golang.org/x/net/context"
)

func (obj *UserHandler) SaveUserWithImmutable(ctx context.Context, userObj *miniuser.User) error {
	// init
	userName := userObj.GetUserName()
	sign := strconv.Itoa(time.Now().Nanosecond())
	Debug(ctx, "SaveUserWithImmutable SIGN 1:"+sign)
	nextUserObj, _ := obj.manager.GetUserFromUserName(ctx, userName, sign)

	Debug(ctx, "SaveUserFromSession :"+userName+":"+nextUserObj.GetUserName()+":"+userObj.GetUserName())
	replayObj := obj.relayIdMgr.GetRelayIdAsPointer(ctx, userName)
	currentSign := replayObj.GetSign()

	// copy
	userObj.CopyWithoutuserName(ctx, nextUserObj)
	if nil != obj.manager.SaveUser(ctx, nextUserObj) {
		Debug(ctx, "SaveUserFromSession Delete failed 1: "+userObj.GetUserName())
		return nil
	}
	replayObj.SetUserName(nextUserObj.GetUserName())
	replayObj.SetSign(sign)
	replayObj.Save(ctx)
	Debug(ctx, "SaveUserWithImmutable SIGN 2:"+sign)
	//
	err1 := obj.manager.SaveUser(ctx, nextUserObj)
	if nil != err1 {
		Debug(ctx, "SaveUserWithImmutable SIGN ERR:"+err1.Error())
	}
	if nil != obj.manager.DeleteUser(ctx, userObj.GetUserName(), currentSign) {
		Debug(ctx, "SaveUserFromSession Delete failed 2: "+userObj.GetUserName())
	}
	return nil
}

func (obj *UserHandler) GetUserFromUserNameAndRelayId(ctx context.Context, userName string) (*miniuser.User, error) {
	Debug(ctx, "SaveUserFromNamePointer :"+userName)

	pointerObj := obj.relayIdMgr.GetRelayIdAsPointer(ctx, userName)
	if pointerObj.GetUserName() == "" {
		Debug(ctx, "SaveUserFromNamePointer err1 :"+userName)
		return nil, errors.New("not found")
	}
	return obj.GetManager().GetUserFromUserName(ctx, pointerObj.GetUserName(), pointerObj.GetSign())
}

func (obj *UserHandler) LoginRegistFromTwitter(ctx context.Context, screenName string, userId string, oauthToken string) (bool, *relayid.RelayId, *miniuser.User, error) {
	relayIdObj := obj.relayIdMgr.GetRelayIdForTwitter(ctx, screenName, userId, oauthToken)
	needMake := false

	//
	// new userObj
	var err error = nil
	var userObj *miniuser.User = nil
	var pointerObj *relayid.RelayId = nil
	if relayIdObj.GetUserName() != "" {
		needMake = true
		//		Debug(ctx, "LoginRegistFromTwitter (1) :"+relayIdObj.GetUserName())
		pointerObj = obj.relayIdMgr.GetRelayIdAsPointer(ctx, relayIdObj.GetUserName())
		if pointerObj.GetUserName() != "" {
			userObj, err = obj.GetManager().GetUserFromUserName(ctx, pointerObj.GetUserName(), pointerObj.GetSign())
			if err != nil {
				userObj = nil
			}
		}
	}
	if userObj == nil {
		userObj = obj.GetManager().NewNewUser(ctx, "")
		userObj.SetDisplayName(screenName)
		//		Debug(ctx, "LoginRegistFromTwitter (2) :"+userObj.GetUserName())
		pointerObj = obj.relayIdMgr.GetRelayIdAsPointer(ctx, userObj.GetUserName())
		pointerObj.SetUserName(userObj.GetUserName())
		pointerObj.SetSign("")
		Debug(ctx, "LoginRegistFromTwitter :")
		err := pointerObj.Save(ctx)
		if err != nil {
			return needMake, nil, nil, errors.New("failed to save pointreobj : " + err.Error())
		}
		//
		// set username
		relayIdObj.SetUserName(pointerObj.GetUserName())
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
