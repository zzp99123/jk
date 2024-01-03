// 用户
package domain

import (
	"time"
)

type User struct {
	Id              int64
	Email           string
	Password        string
	Phone           string
	Ctime           time.Time
	Birthday        time.Time //生日
	PersonalProfile string    //个人简介
	Nickname        string    //昵称
	// 不要使用组合，因为你将来可能还有 DingDingInfo 之类的
	WechatInfo WechatInfo
}
