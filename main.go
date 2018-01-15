package main

import (
	"net/http"

	"fmt"
	"time"

	_ "amazonApi/models"
	"strconv"

	"amazonApi/models"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.GET("/user", func(c echo.Context) error {

		type Data struct {
			Id         string `json:"id"`
			Email      string `json:"email"`
			Facebook   string `json:"facebook"`
			Twitter    string `json:"twitter"`
			Instagram  string `json:"instagram"`
			Pinterest  string `json:"pinterest"`
			Youtube    string `json:"youtube"`
			ProfileUrl string `json:"profile_url"`
			ProfileId  string `json:"profile_id"`
			Name       string `json:"name"`
			Country    string `json:"country"`
			Total      string `json:"review_count"`
		}
		var userJson = struct {
			Code  int    `json:"code"`
			Total int    `json:"total"`
			Datas []Data `json:"datas"`
		}{}

		o := orm.NewOrm()

		size, err := strconv.Atoi(c.QueryParam("size"))
		if err != nil {
			size = 15
		}
		page, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil {
			page = 0
		}
		if page > 0 {
			page -= 1
		}

		name := c.QueryParam("name")
		var maps []orm.Params
		if name == "" {

			_, err = o.Raw("select user.id,user.email,user.facebook,user.twitter,user.instagram,user.profile_url,user.pinterest,user.youtube,user.country,user.name,count(*) "+
				"from user,product "+
				"where user.id = product.user_id "+
				"Group by user.id "+
				"order by count(*) "+
				"Desc limit ? offset ?", size, page*size).Values(&maps)
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}
		} else {

			v := "select user.id,user.email,user.facebook,user.twitter,user.instagram,user.profile_url,user.pinterest,user.youtube,user.country,user.name,count(*) " +
				"from user,product where user.id = product.user_id and user.name like"
			v += " '%"
			v += fmt.Sprintf("%s", name)
			v += "%' Group by user.id order by count(*) Desc "
			v += fmt.Sprintf("limit %d offset %d", size, page*size)
			_, err = o.Raw(v).Values(&maps)
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}
		}

		for _, v := range maps {
			data := Data{
				Id:         v["id"].(string),
				Email:      v["email"].(string),
				Facebook:   v["facebook"].(string),
				Twitter:    v["twitter"].(string),
				Instagram:  v["instagram"].(string),
				Pinterest:  v["pinterest"].(string),
				Youtube:    v["youtube"].(string),
				ProfileUrl: v["profile_url"].(string),
				Name:       v["name"].(string),
				Country:    v["country"].(string),
				Total:      v["count(*)"].(string),
			}
			userJson.Datas = append(userJson.Datas, data)
		}

		if name == "" {
			_, err = o.Raw("select count(*) from user").Values(&maps)
		} else {
			v := "select count(*) from user where user.name like"
			v += " '%"
			v += fmt.Sprintf("%s", name)
			v += "%'"
			_, err = o.Raw(v).Values(&maps)
		}

		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		if len(maps) > 0 {
			if s, ok := maps[0]["count(*)"].(string); ok {
				count, err := strconv.Atoi(s)
				if err == nil {
					userJson.Total = count
				} else {
					fmt.Println(err)
				}
			}

		}

		userJson.Code = 200

		return c.JSON(http.StatusOK, userJson)
	})

	e.GET("/user/:id", func(c echo.Context) error {
		id := c.Param("id")

		type Data struct {
			Id           string `json:"id"`
			CategoryName string `json:"name"`
			ReviewCount  string `json:"review_count"`
		}
		var category = struct {
			Code  int    `json:"code"`
			Total int    `json:"total"`
			Datas []Data `json:"datas"`
		}{}
		o := orm.NewOrm()
		size, err := strconv.Atoi(c.QueryParam("size"))
		if err != nil {
			size = 15
		}
		page, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil {
			page = 0
		}
		if page > 0 {
			page -= 1
		}

		var maps []orm.Params

		_, err = o.Raw("select category.id,category.name,count(*) "+
			"from user,product,category "+
			"where user.id= ? and product.`user_id` = ? and product.`category_id` = category.id "+
			"GROUP BY category_id "+
			"order by count(*) desc limit ? offset ?", id, id, size, page*size).Values(&maps)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		for _, v := range maps {
			data := Data{
				Id:           v["id"].(string),
				CategoryName: v["name"].(string),
				ReviewCount:  v["count(*)"].(string),
			}
			category.Datas = append(category.Datas, data)
		}

		_, err = o.Raw("select count(distinct category.id) from user,product,category where user.id= ? and product.`user_id` = ? and product.`category_id` = category.id", id, id).Values(&maps)

		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		if len(maps) > 0 {
			if s, ok := maps[0]["count(distinct category.id)"].(string); ok {
				count, err := strconv.Atoi(s)
				if err == nil {
					category.Total = count
				} else {
					fmt.Println(err)
				}
			}

		}

		return c.JSON(http.StatusOK, category)
	})

	e.GET("/category", func(context echo.Context) error {
		type Data struct {
			Id              string `json:"id"`
			Name            string `json:"name"`
			ReviewUserCount string `json:"review_user_count"`
			CategoryLink    string `json:"category_link"`
		}
		var category = struct {
			Code  int    `json:"code"`
			Total int    `json:"total"`
			Datas []Data `json:"datas"`
		}{}
		o := orm.NewOrm()
		size, err := strconv.Atoi(context.QueryParam("size"))
		if err != nil {
			size = 15
		}
		page, err := strconv.Atoi(context.QueryParam("page"))
		if err != nil {
			page = 0
		}
		if page > 0 {
			page -= 1
		}

		name := context.QueryParam("name")
		var maps []orm.Params
		if name == "" {
			_, err = o.Raw("select  category.parent_id, category.name,category.id,count(*) "+
				"from category,product where category.id=product.category_id "+
				"GROUP BY category.id "+
				"ORDER BY  COUNT(*) desc limit ? offset ?", size, page*size).Values(&maps)
			if err != nil {
				return context.String(http.StatusBadRequest, err.Error())
			}
		} else {
			v := "select category.parent_id,category.name,category.id,count(*) from category,product where category.id=product.category_id and " +
				"category.name like"
			v += " '%"
			v += fmt.Sprintf("%s", name)
			v += "%' GROUP BY category.id ORDER BY  COUNT(*) desc "
			v += fmt.Sprintf("limit %d offset %d", size, page*size)
			_, err = o.Raw(v).Values(&maps)
			if err != nil {
				return context.String(http.StatusBadRequest, err.Error())
			}
		}

		for _, v := range maps {
			data := Data{
				Id:              v["id"].(string),
				Name:            v["name"].(string),
				ReviewUserCount: v["count(*)"].(string),
			}
			var categoryLink []string
			categoryLink = append(categoryLink, v["name"].(string))
			var parentId int64
			for {
				if parentId == 0 {
					if pid, ok := v["parent_id"].(string); ok {
						if id, err := strconv.Atoi(pid); err == nil {
							var category models.Category
							err := o.QueryTable("category").Filter("id", id).One(&category)
							if err == nil {
								categoryLink = append(categoryLink, category.Name)
								if category.Id != 0 && category.ParentId != 0 {
									parentId = category.ParentId
								} else {
									break
								}
							} else {
								break
							}
						} else {
							break
						}

					} else {
						break
					}
				} else {
					var category models.Category
					err := o.QueryTable("category").Filter("id", parentId).One(&category)
					if err == nil {
						categoryLink = append(categoryLink, category.Name)
						if category.Id != 0 && category.ParentId != 0 {
							parentId = category.ParentId
						} else {
							break
						}
					} else {
						break
					}
				}

			}

			for i := len(categoryLink) - 1; i >= 0; i-- {
				if i == 0 {
					data.CategoryLink += categoryLink[i]
				} else {
					data.CategoryLink += fmt.Sprintf("%s->", categoryLink[i])
				}

			}

			category.Datas = append(category.Datas, data)
		}

		if name == "" {
			_, err = o.Raw("select count(distinct category.id) from category,product where category.id=product.category_id").Values(&maps)
		} else {
			v := "select count(distinct category.id) from category where category.name like"
			v += " '%"
			v += fmt.Sprintf("%s", name)
			v += "%'"
			_, err = o.Raw(v).Values(&maps)
		}

		if err != nil {
			return context.String(http.StatusBadRequest, err.Error())
		}

		if len(maps) > 0 {
			if s, ok := maps[0]["count(distinct category.id)"].(string); ok {
				count, err := strconv.Atoi(s)
				if err == nil {
					category.Total = count
				} else {
					fmt.Println(err)
				}
			}

		}

		return context.JSON(http.StatusOK, category)
	})

	e.GET("/category/:id", func(c echo.Context) error {
		id := c.Param("id")

		type Data struct {
			Id         string `json:"id"`
			Email      string `json:"email"`
			Facebook   string `json:"facebook"`
			Twitter    string `json:"twitter"`
			Instagram  string `json:"instagram"`
			Pinterest  string `json:"pinterest"`
			Youtube    string `json:"youtube"`
			ProfileUrl string `json:"profile_url"`
			ProfileId  string `json:"profile_id"`
			Name       string `json:"name"`
			Country    string `json:"country"`
			Total      string `json:"review_count"`
		}
		var category = struct {
			Code  int    `json:"code"`
			Total int    `json:"total"`
			Datas []Data `json:"datas"`
		}{}
		o := orm.NewOrm()
		size, err := strconv.Atoi(c.QueryParam("size"))
		if err != nil {
			size = 15
		}
		page, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil {
			page = 0
		}
		if page > 0 {
			page -= 1
		}

		var maps []orm.Params

		_, err = o.Raw("select user.id,user.email,user.facebook,user.twitter,user.instagram,user.profile_url,user.pinterest,user.youtube,user.country,user.name,count(*) "+
			"from user,product "+
			"where user.id = product.user_id and product.`category_id` = ? "+
			"GROUP BY user.id ORDER BY  COUNT(*) "+
			"desc limit ? offset ?", id, size, page*size).Values(&maps)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		for _, v := range maps {
			data := Data{
				Id:         v["id"].(string),
				Email:      v["email"].(string),
				Facebook:   v["facebook"].(string),
				Twitter:    v["twitter"].(string),
				Instagram:  v["instagram"].(string),
				Pinterest:  v["pinterest"].(string),
				Youtube:    v["youtube"].(string),
				ProfileUrl: v["profile_url"].(string),
				Name:       v["name"].(string),
				Country:    v["country"].(string),
				Total:      v["count(*)"].(string),
			}
			category.Datas = append(category.Datas, data)
		}

		_, err = o.Raw("select count(DISTINCT user.id) from user,product where user.id = product.user_id and product.`category_id` = ? ", id).Values(&maps)

		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		if len(maps) > 0 {
			if s, ok := maps[0]["count(DISTINCT user.id)"].(string); ok {
				count, err := strconv.Atoi(s)
				if err == nil {
					category.Total = count
				} else {
					fmt.Println(err)
				}
			}

		}

		return c.JSON(http.StatusOK, category)
	})

	e.Logger.Fatal(e.Start(":9527"))

}

func init() {
	local, err := time.LoadLocation("Asia/Shanghai")

	if err != nil {
		fmt.Println(err)
	}
	time.Local = local
	err = orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/amazon?parseTime=true&loc=Asia%2FShanghai", 30, 30)
	if err != nil {
		fmt.Println(err)
	}
}
