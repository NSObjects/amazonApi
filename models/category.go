package models

import "github.com/astaxie/beego/orm"

type Category struct {
	Id       int64
	ParentId int64
	Name     string
}

func init() {
	orm.RegisterModel(new(Category))
}
