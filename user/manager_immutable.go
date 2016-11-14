package user

import (
	"strconv"

	"time"

	"errors"

	"github.com/firefirestyle/go.minipointer"
	"golang.org/x/net/context"
)

//
//
func (obj *UserManager) SaveUserWithImmutable(ctx context.Context, userObj *User) (*User, error) {
	// init
	userName := userObj.GetUserName()
	sign := strconv.Itoa(time.Now().Nanosecond())
	nextUserObj, _ := obj.GetUserFromUserName(ctx, userName, sign)

	replayObj := obj.pointerManager.GetPointerWithNewForRelayId(ctx, userName)
	currentSign := replayObj.GetSign()

	// copy
	userObj.CopyWithoutUserNameAndSign(ctx, nextUserObj)
	if nil != obj.SaveUser(ctx, nextUserObj) {
		return nextUserObj, nil
	}
	replayObj.SetValue(nextUserObj.GetUserName())
	replayObj.SetSign(sign)
	obj.pointerManager.Save(ctx, replayObj)
	//
	err1 := obj.SaveUser(ctx, nextUserObj)
	if nil != err1 {
		return nil, err1
	}
	err2 := obj.DeleteUser(ctx, userObj.GetUserName(), currentSign)
	if nil != err2 {
		return nil, err2
	}
	return nextUserObj, nil
}

func (obj *UserManager) GetPointerFromUserName(ctx context.Context, userName string) (*minipointer.Pointer, error) {
	return obj.pointerManager.GetPointer(ctx, userName, minipointer.TypePointer)
}

func (obj *UserManager) GetUserFromRelayId(ctx context.Context, userName string) (*User, error) {
	Debug(ctx, "SaveUserFromNamePointer :"+userName)

	pointerObj, pointerErr := obj.pointerManager.GetPointer(ctx, userName, minipointer.TypePointer)
	if pointerErr != nil {
		Debug(ctx, "SaveUserFromNamePointer err1 :"+userName)
		return nil, errors.New(pointerErr.Error())
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
