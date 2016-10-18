package miniuser

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

//
//
//
func (obj *UserManager) newCursorFromSrc(cursorSrc string) *datastore.Cursor {
	c1, e := datastore.DecodeCursor(cursorSrc)
	if e != nil {
		return nil
	} else {
		return &c1
	}
}

func (obj *UserManager) makeCursorSrc(founds *datastore.Iterator) string {
	c, e := founds.Cursor()
	if e == nil {
		return c.String()
	} else {
		return ""
	}
}

//
//
func (obj *UserManager) FindUserWithNewOrder(ctx context.Context, cursorSrc string) ([]*User, string, string) {
	q := datastore.NewQuery(obj.userKind)
	q = q.Filter("ProjectId =", obj.projectId)
	q = q.Filter("State =", UserStatePublic)
	q = q.Limit(obj.limitOfFinding)
	return obj.FindUserFromQuery(ctx, q, cursorSrc)
}

func (obj *UserManager) FindUserWithPoint(ctx context.Context, cursorSrc string) ([]*User, string, string) {
	q := datastore.NewQuery(obj.userKind)
	q = q.Filter("ProjectId =", obj.projectId)
	q = q.Filter("State =", UserStatePublic)
	q = q.Order("-Point")
	q = q.Limit(obj.limitOfFinding)
	return obj.FindUserFromQuery(ctx, q, cursorSrc)
}

//
//
func (obj *UserManager) FindUserFromQuery(ctx context.Context, queryObj *datastore.Query, cursorSrc string) ([]*User, string, string) {
	cursor := obj.newCursorFromSrc(cursorSrc)
	if cursor != nil {
		queryObj = queryObj.Start(*cursor)
	}
	queryObj = queryObj.KeysOnly()

	var userObjList []*User

	founds := queryObj.Run(ctx)

	var cursorNext string = ""
	var cursorOne string = ""

	for i := 0; ; i++ {
		key, err := founds.Next(nil)
		if err != nil || err == datastore.Done {
			break
		} else {
			userObj := obj.newUserFromStringID(ctx, key.StringID())
			errLoadUserObj := userObj.loadFromDB(ctx)
			if errLoadUserObj != nil {
				log.Infof(ctx, "Failed LoadFromDB on FindUserFromQuery "+key.StringID())
			} else {
				userObjList = append(userObjList, userObj)
			}
		}
		if i == 0 {
			cursorOne = obj.makeCursorSrc(founds)
		}
	}
	cursorNext = obj.makeCursorSrc(founds)
	return userObjList, cursorOne, cursorNext
}
