package hundler

import (
	"errors"

	//	"strings"

	//	"strconv"
	//	"time"

	"github.com/firefirestyle/go.minipointer"
	miniuser "github.com/firefirestyle/go.miniuser/user"
	"golang.org/x/net/context"
)

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
