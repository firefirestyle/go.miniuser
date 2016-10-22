package miniuser

import (
	"errors"

	"strings"

	"strconv"
	"time"

	"github.com/firefirestyle/go.miniuser/relayid"
	"golang.org/x/net/context"
)

func (obj *UserManager) SaveUserFromSession(ctx context.Context, userObj *User) error {
	// init
	userName := strings.Split(userObj.GetUserName(), "::sign::")[0]
	nextUserObj := obj.newUser(ctx, userName+"::sign::"+strconv.Itoa(time.Now().Nanosecond()))
	Debug(ctx, "SaveUserFromSession :"+userName)
	replayObj := obj.relayIdMgr.GetRelayIdAsPointer(ctx, userName)

	// copy
	userObj.CopyWithoutuserName(ctx, nextUserObj)
	replayObj.SetUserName(nextUserObj.GetUserName())
	replayObj.Save(ctx)

	//
	obj.SaveUser(ctx, nextUserObj)
	obj.DeleteUser(ctx, userObj.GetUserName())
	return userObj.pushToDB(ctx)
}

func (obj *UserManager) GetUserFromUserNamePointer(ctx context.Context, userName string) (*User, error) {
	Debug(ctx, "SaveUserFromNamePointer :"+userName)
	pointerObj := obj.relayIdMgr.GetRelayIdAsPointer(ctx, userName)
	if pointerObj.GetUserName() == "" {
		return nil, errors.New("not found")
	}
	return obj.GetUserFromUserName(ctx, pointerObj.GetId())
}

func (obj *UserManager) LoginRegistFromTwitter(ctx context.Context, screenName string, userId string, oauthToken string) (bool, *relayid.RelayId, *User, error) {
	relayIdObj := obj.relayIdMgr.GetRelayIdForTwitter(ctx, screenName, userId, oauthToken)
	needMake := false

	//
	// new userObj
	var err error = nil
	var userObj *User = nil
	var pointerObj *relayid.RelayId = nil
	if relayIdObj.GetUserName() != "" {
		needMake = true
		Debug(ctx, "LoginRegistFromTwitter (1) :"+relayIdObj.GetUserName())
		pointerObj = obj.relayIdMgr.GetRelayIdAsPointer(ctx, relayIdObj.GetUserName())
		if pointerObj.GetUserName() != "" {
			userObj, err = obj.GetUserFromUserName(ctx, pointerObj.GetUserName())
			if err != nil {
				userObj = nil
			}
		}
	}
	if userObj == nil {
		userObj = obj.NewNewUser(ctx)
		userObj.SetDisplayName(screenName)
		Debug(ctx, "LoginRegistFromTwitter (2) :"+userObj.GetUserName())
		pointerObj = obj.relayIdMgr.GetRelayIdAsPointer(ctx, userObj.GetUserName())
		pointerObj.SetUserName(userObj.GetUserName())
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
	err = obj.SaveUser(ctx, userObj)
	if err != nil {
		return needMake, nil, nil, errors.New("failed to save userobj : " + err.Error())
	}
	return needMake, relayIdObj, userObj, nil
}
