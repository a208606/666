package routers

import (

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"newWeb/controllers"
	_ "newWeb/models"
)

func init() {

	beego.InsertFilter("/article/*",beego.BeforeExec,FilterFunc)
	beego.Router("/", &controllers.MainController{})
	beego.Router("/register", &controllers.UserControllers{}, "get:ShowRegister;post:HandRegister")
	beego.Router("/login",&controllers.UserControllers{},"get:ShowLogin;post:HandleLogin")

	beego.Router("/article/index",&controllers.ArticleControllers{},"get,post:ShowIndex")
	//添加文章表
	beego.Router("/article/addArticle",&controllers.ArticleControllers{},"get:ShowArticle;post:HandleArticle")
	//获取内容详情
	beego.Router("/article/content",&controllers.ArticleControllers{},"get:ShowContent")
	//编辑内容
	beego.Router("/article/update",&controllers.ArticleControllers{},"get:ShowUpdate;post:HandleUpdate")
	//删除内容
	beego.Router("/article/delete",&controllers.ArticleControllers{},"get:ShowDelete")
	//添加文章类型
	beego.Router("/article/addType",&controllers.ArticleControllers{},"get:ShowAddType;post:HandleAddType")
	//删除文章类型
	beego.Router("/article/deleteType",&controllers.ArticleControllers{},"get:DeleteType")
	//退出
	beego.Router("/article/logout",&controllers.UserControllers{},"get:Logout")

}

func FilterFunc(ctx *context.Context)  {

	userName :=ctx.Input.Session("userName")
	if userName == nil{
		ctx.Redirect(302,"/login")
		return
	}


}