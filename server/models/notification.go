package models

import (
	"time"

	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
)

type NotificationType int

const (
	NotificationDiscussion NotificationType = 1
	NotificationStar       NotificationType = 2
)

type UserNotification struct {
	UserID        string `bson:"_id"` // 通知先のUserID
	Notifications []Notification
}

func (u *UserNotification) TrimNotification(length int) {
	trimIdx := len(u.Notifications) - length
	if trimIdx > 0 {
		u.Notifications = u.Notifications[trimIdx:]
	}
}

type Notification struct {
	NotificationID   string           // 連番 time.Now.UnixNanoかな
	NotificationType NotificationType // 通知の区分
	SourceID         string           // 通知の元になったID
	Message          string           // 通知メッセージ
	DetailURL        string           // この通知の詳細を見たいときに飛ばすURL
	IsRead           bool             // 既読
	CreatedAt        time.Time        // 送信時間
}

// AppsContext appsのコレクション
type _NotificationTable struct {
}

func (_ _NotificationTable) Name() string {
	return "notification"
}

var _ modelsTable = (*_NotificationTable)(nil)

// NotificationTable appInfo
var NotificationTable = _NotificationTable{}

func (t _NotificationTable) withCollection(ctx context.Context, fn func(c *mgo.Collection)) {
	withDefaultCollection(ctx, t.Name(), fn)
}

func (t _NotificationTable) AddNotification(ctx context.Context, userID string, notification Notification) error {
	var result UserNotification
	var err error
	t.withCollection(ctx, func(c *mgo.Collection) {
		err = c.FindId(userID).One(&result)
		if err == mgo.ErrNotFound {
			result.UserID = userID
			result.Notifications = []Notification{}
			err = nil // NotFoundは新規なのでOK
		} else if err != nil {
			return
		}

		result.Notifications = append(result.Notifications, notification)
		result.TrimNotification(10)
	})

	if err != nil {
		return err
	}

	return t.Upsert(ctx, result)
}

func (t _NotificationTable) MustFindID(ctx context.Context, userID string) UserNotification {
	notification, err := t.FindID(ctx, userID)
	if err != nil {
		notification.UserID = userID
		notification.Notifications = []Notification{}
	}
	return notification
}

func (t _NotificationTable) FindID(ctx context.Context, userID string) (result UserNotification, err error) {
	t.withCollection(ctx, func(c *mgo.Collection) {
		err = c.FindId(userID).One(&result)
	})
	return
}

// Upsert 登録
func (t _NotificationTable) Upsert(ctx context.Context, un UserNotification) error {
	var err error
	t.withCollection(ctx, func(c *mgo.Collection) {
		var result interface{} // bson.M
		_, err = c.FindId(un.UserID).Apply(mgo.Change{
			Update: un,
			Upsert: true,
		}, &result)
	})
	return err
}
