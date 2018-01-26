package models

import "github.com/astaxie/beego/orm"

type Product struct {
	UserProfile string   `orm:"-" json:"user_profile"`
	Id          int64    `json:"id"`
	Url         string   `json:"url"`
	CategoryId  int64    `json:"category_id"`
	Categorys   []string `orm:"-" json:"categorys"`
	Name        string   `json:"name"`
	UserId      int64    `json:"user_id"`
	ReviewUrl   string   `json:"review_url"`
}

func init() {
	orm.RegisterModel(new(Product))
}
