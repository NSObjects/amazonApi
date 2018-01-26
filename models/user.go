package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type User struct {
	Id           int64      `json:"id"`
	Email        string     `json:"email"`
	Facebook     string     `json:"facebook"`
	Twitter      string     `json:"twitter"`
	Instagram    string     `json:"instagram"`
	Pinterest    string     `json:"pinterest"`
	Youtube      string     `json:"youtube"`
	ProfileUrl   string     `json:"profile_url"`
	ProfileId    string     `json:"profile_id"`
	Name         string     `json:"name"`
	Country      int        `json:"country"`
	HelpfulVotes int        `orm:"column(helpful_votes)" json:"helpful_votes"`
	Reviews      int        `orm:"column(reviews)" json:"reviews"`
	Created      time.Time  `orm:"auto_now_add;type(datetime)"`
	Updated      time.Time  `orm:"auto_now;type(datetime)"`
	Products     []*Product `orm:"reverse(many)"`
}

func init() {
	orm.RegisterModel(new(User))
}
