package user

import (
	"strconv"

	"time"

	"errors"

	"golang.org/x/net/context"
)

//
//
func (obj *UserManager) SaveUserWithImmutable(ctx context.Context, userObj *User) error {
	// init
	userName := userObj.GetUserName()
	sign := strconv.Itoa(time.Now().Nanosecond())
	nextUserObj, _ := obj.GetUserFromUserName(ctx, userName, sign)

	replayObj := obj.pointerManager.GetPointerForRelayId(ctx, userName)
	currentSign := replayObj.GetSign()

	// copy
	userObj.CopyWithoutUserNameAndSign(ctx, nextUserObj)
	if nil != obj.SaveUser(ctx, nextUserObj) {
		return nil
	}
	replayObj.SetValue(nextUserObj.GetUserName())
	replayObj.SetSign(sign)
	replayObj.Save(ctx)
	//
	err1 := obj.SaveUser(ctx, nextUserObj)
	if nil != err1 {
		return err1
	}
	err2 := obj.DeleteUser(ctx, userObj.GetUserName(), currentSign)
	if nil != err2 {
		return err2
	}
	return nil
}

func (obj *UserManager) GetUserFromRelayId(ctx context.Context, userName string) (*User, error) {
	Debug(ctx, "SaveUserFromNamePointer :"+userName)

	pointerObj := obj.pointerManager.GetPointerForRelayId(ctx, userName)
	if pointerObj.GetValue() == "" {
		Debug(ctx, "SaveUserFromNamePointer err1 :"+userName)
		return nil, errors.New("not found")
	}
	return obj.GetUserFromUserName(ctx, pointerObj.GetValue(), pointerObj.GetSign())
}

func (obj *UserManager) GetUserFromSign(ctx context.Context, userName string, sign string) (*User, error) {
	Debug(ctx, " GetUserFromUserNameAndSign :"+userName+" : "+sign)

	return obj.GetUserFromUserName(ctx, userName, sign)
}

func (obj *UserManager) GetUserFromKey(ctx context.Context, stringId string) (*User, error) {
	Debug(ctx, "GetUserFromKey :"+stringId)
	keyInfo := obj.GetUserKeyInfo(stringId)
	Debug(ctx, "GetUserFromKey :"+keyInfo.UserName+" : "+keyInfo.Sign)
	return obj.GetUserFromUserName(ctx, keyInfo.UserName, keyInfo.Sign)
}
