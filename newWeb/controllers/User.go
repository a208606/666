package controllers

import (
	"encoding/base64"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"newWeb/models"
)

type UserControllers struct {
	beego.Controller
}

func (this *UserControllers) ShowRegister() {
	this.TplName = "register.html"
}
func (this *UserControllers) HandRegister() {
	name := this.GetString("userName")
	password := this.GetString("password")

	if name == "" || password == "" {
		beego.Error("账号密码不能为空")
		this.TplName = "register.html"
		return
	}

	var user models.User

	o := orm.NewOrm()
	user.Name = name
	user.Password = password
	id, err := o.Insert(&user)

	if err != nil || id == 0 {
		beego.Error("数据库插入数据失败", err)
		this.TplName = "register.html"
		return
	}
	fmt.Println("id", id)
	this.Redirect("/login", 302)
}

//登录
func (this *UserControllers) ShowLogin() {

	userName := this.Ctx.GetCookie("Name")
	dec, _ := base64.StdEncoding.DecodeString(userName)

	if userName != "" {
		this.Data["userName"] = string(dec)
		this.Data["checked"] = "checked"

	} else {
		this.Data["userName"] = ""
		this.Data["checked"] = ""
	}

	this.TplName = "login.html"
}

//密码比对
func (this *UserControllers) HandleLogin() {

	name := this.GetString("userName")
	password := this.GetString("password")
	var user models.User

	o := orm.NewOrm()
	user.Name = name
	err := o.Read(&user, "Name")
	if err != nil {
		beego.Error("用户名不存在")
		this.TplName = "login.html"
		return
	}
	if user.Password != password {
		beego.Error("密码错误")
		this.TplName = "login.html"
		return
	}
	//实现记住密码功能
	remember := this.GetString("remember")

	enc := base64.StdEncoding.EncodeToString([]byte(name))

	if remember != "" {
		this.Ctx.SetCookie("Name", enc, 60)
	} else {
		this.Ctx.SetCookie("Name", name, -1)
	}

	this.SetSession("userName",name)
	this.Redirect("/article/index", 302)
}
//退出
func (this*UserControllers)Logout()  {
	this.DelSession("userName")
	this.Redirect("/login",302)
}