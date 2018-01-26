package models

import "github.com/astaxie/beego/orm"

type Product struct {
	UserProfile string `orm:"-" json:"user_profile"`
	Id          int64
	Url         string
	CategoryId  int64
	Categorys   []string `orm:"-" json:"categorys"`
	Name        string
	UserId      int64
	ReviewUrl   string
}

func init() {
	orm.RegisterModel(new(Product))
}
