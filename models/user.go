package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type User struct {
	Id         int64
	Email      string
	Facebook   string
	Twitter    string
	Instagram  string
	Pinterest  string
	Youtube    string
	ProfileUrl string
	ProfileId  string
	Name       string
	Country    int
	Created    time.Time `orm:"auto_now_add;type(datetime)"`
	Updated    time.Time `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(User))
}
