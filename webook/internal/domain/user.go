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
}
