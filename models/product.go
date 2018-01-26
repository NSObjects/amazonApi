package models

import "github.com/astaxie/beego/orm"

type Product struct {
	UserProfile string   `orm:"-" json:"user_profile"`
	Id          int64    `json:"id"`
	Url         string   `json:"url"`
	CategoryId  int64    `json:"category_id"`
	Categorys   []string `orm:"-" json:"categorys"`
	Name        string   `json:"name"`

	ReviewUrl string `json:"review_url"`
	User      *User  `orm:"rel(fk)"`
}

func init() {
	orm.RegisterModel(new(Product))
}
