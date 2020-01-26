package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type User struct {
	Id       int
	Name     string
	Password string
	Articles []*Article `orm:"reverse(many)"`
}

type Article struct {
	Id          int          `orm:"pk;auto;unique"`
	Title       string       `orm:"unique;size(40)"`
	Content     string       `orm:"size(999)"`
	Img         string       `orm:"null"`
	Time        time.Time    `orm:"type(datetime);auto_now_add"`
	ReadCount   int          `orm:"default(0)"`

	ArticleType *ArticleType `orm:"rel(fk);null;on_delete(set_null)"`
	Users       []*User      `orm:"rel(m2m)"`
}

type ArticleType struct {
	Id          int
	TypeArticle string `orm:"unique"`
	Article     []*Article `orm:"reverse(many)"`
}

func init() {
	err := orm.RegisterDataBase("default", "mysql", "root:111@tcp(10.0.2.18:3306)/newsWeb")
	if err != nil {
		beego.Error("数据库丽娜姐失败", err)
		return
	}

	orm.RegisterModel(new(User), new(Article), new(ArticleType))

	err = orm.RunSyncdb("default", false, true)
	if err != nil {
		beego.Error("数据库运行同步发生失败", err)
		return
	}
}
