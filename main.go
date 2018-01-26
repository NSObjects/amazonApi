package main

import (
	"amazonApi/models"
	"fmt"
	"net/http"
	"strconv"
	"time"

	//"os"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

func main() {

	e := echo.New()
	//orm.Debug = true
	//orm.DebugLog = orm.NewLog(os.Stdout)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.GET("/user", func(c echo.Context) error {

		var userJson = struct {
			Code  int           `json:"code"`
			Total int64         `json:"total"`
			Datas []models.User `json:"datas"`
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
		sort := c.QueryParam("sort")
		if sort == "" {
			sort = "helpful_votes"
		}

		sort = "-" + sort

		country := c.QueryParam("country")

		name := c.QueryParam("name")
		var users []models.User
		if name == "" {
			qs := o.QueryTable("user")
			if country != "" {
				qs = qs.Filter("country", country)
			}
			_, err = qs.OrderBy(sort).Limit(size, page*size).All(&users)
			if err != nil {
				fmt.Println(err)
			}
			userJson.Total, _ = o.QueryTable("user").Count()
		} else {
			//sql := "select DISTINCT user.profile_url,user.id,user.email,user.facebook,user.twitter,user.instagram,user.profile_url,user.pinterest,user.youtube,user.country,user.name,user.helpful_votes, user.reviews from user,product where product.name like "
			//sql += "'%"
			//sql += name + "%' and user.id = product.user_id "
			//if country != "" {
			//	sql += fmt.Sprintf(" and user.country = %s ", country)
			//}
			//
			//sql += fmt.Sprintf("order by %s", sort)
			//sql += fmt.Sprintf(" limit %d offset %d", size, page*size)
			//_, err := o.Raw(sql).QueryRows(&users)
			//if err != nil {
			//	fmt.Println(err)name__icontains
			//}

			_, err := o.QueryTable("user").Filter("Products__name__icontains", name).RelatedSel().All(&users)
			if err != nil {
				fmt.Println(err)
			}

			userJson.Total = int64(len(users))
		}

		userJson.Datas = users
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

	e.GET("/product", func(context echo.Context) error {
		userId := context.QueryParam("userId")

		if userId == "" {
			return context.JSON(http.StatusOK, "user id is nil")
		}

		name := context.QueryParam("name")
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
		var products []models.Product

		o := orm.NewOrm()

		_, err = o.QueryTable("product").
			Filter("user_id", userId).
			Filter("name__icontains", name).
			Limit(size, page*size).
			All(&products)
		if err != nil {
			return context.JSON(http.StatusOK, err)
		}
		count, err := o.QueryTable("product").
			Filter("user_id", userId).
			Filter("name__icontains", name).Count()

		var j struct {
			Data  []models.Product `json:"data"`
			Total int64            `json:"total"`
		}
		j.Total = count
		j.Data = products
		return context.JSON(http.StatusOK, &j)
	})

	e.Logger.Fatal(e.Start(":9527"))

}

func init() {

	viper.SetConfigType("yaml")
	viper.SetConfigFile("./config/app.yaml")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	local, err := time.LoadLocation("Asia/Shanghai")

	if err != nil {
		fmt.Println(err)
	}
	time.Local = local

	if viper.Get("runmodel") == "dev" {
		err = orm.RegisterDataBase("default", "mysql", "root:123456@tcp(192.168.12.137:3306)/amazon?parseTime=true&loc=Asia%2FShanghai")
	} else {
		err = orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1.137:3306)/amazon?parseTime=true&loc=Asia%2FShanghai")
	}

	if err != nil {
		fmt.Println(err)
	}
}
