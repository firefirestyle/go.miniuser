package hundler

import (
	"errors"

	//	"strings"

	"strconv"
	"time"

	"github.com/firefirestyle/go.minipointer"
	miniuser "github.com/firefirestyle/go.miniuser/user"
	"golang.org/x/net/context"
)

func (obj *UserHandler) SaveUserWithImmutable(ctx context.Context, userObj *miniuser.User) error {
	// init
	userName := userObj.GetUserName()
	sign := strconv.Itoa(time.Now().Nanosecond())
	nextUserObj, _ := obj.manager.GetUserFromUserName(ctx, userName, sign)

	replayObj := obj.relayIdMgr.GetPointerForRelayId(ctx, userName)
	currentSign := replayObj.GetSign()

	// copy
	userObj.CopyWithoutUserNameAndSign(ctx, nextUserObj)
	if nil != obj.manager.SaveUser(ctx, nextUserObj) {
		return nil
	}
	replayObj.SetValue(nextUserObj.GetUserName())
	replayObj.SetSign(sign)
	replayObj.Save(ctx)
	//
	err1 := obj.manager.SaveUser(ctx, nextUserObj)
	if nil != err1 {
		return err1
	}
	err2 := obj.manager.DeleteUser(ctx, userObj.GetUserName(), currentSign)
	if nil != err2 {
		return err2
	}
	return nil
}

func (obj *UserHandler) GetUserFromRelayId(ctx context.Context, userName string) (*miniuser.User, error) {
	Debug(ctx, "SaveUserFromNamePointer :"+userName)

	pointerObj := obj.relayIdMgr.GetPointerForRelayId(ctx, userName)
	if pointerObj.GetValue() == "" {
		Debug(ctx, "SaveUserFromNamePointer err1 :"+userName)
		return nil, errors.New("not found")
	}
	return obj.GetManager().GetUserFromUserName(ctx, pointerObj.GetValue(), pointerObj.GetSign())
}

func (obj *UserHandler) GetUserFromSign(ctx context.Context, userName string, sign string) (*miniuser.User, error) {
	Debug(ctx, " GetUserFromUserNameAndSign :"+userName+" : "+sign)

	return obj.GetManager().GetUserFromUserName(ctx, userName, sign)
}
func (obj *UserHandler) GetUserFromKey(ctx context.Context, stringId string) (*miniuser.User, error) {
	Debug(ctx, "GetUserFromKey :"+stringId)
	keyInfo, err := obj.GetManager().GetUserKeyInfo(stringId)
	if err != nil {
		Debug(ctx, "GetUserFromKey : err "+err.Error())
		return nil, err
	}
	Debug(ctx, "GetUserFromKey :"+keyInfo.UserName+" : "+keyInfo.Sign)
	return obj.GetManager().GetUserFromUserName(ctx, keyInfo.UserName, keyInfo.Sign)
}

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
