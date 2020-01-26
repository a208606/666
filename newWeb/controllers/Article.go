package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	"newWeb/models"
	"path"
	"strconv"
	"time"
)

type ArticleControllers struct {
	beego.Controller
}

func (this *ArticleControllers) ShowIndex() {

	userName:=this.GetSession("userName")
	if userName == nil{
		this.Redirect("/login",302)
		return
	}else {
		this.Data["userName"]= userName.(string)
	}

	o := orm.NewOrm()
	qs := o.QueryTable("Article")

	var articles [] models.Article
	//
	_, err := qs.All(&articles)
	if err != nil {
		beego.Error("获取数据失败", err)
		return
	}

	TypeName := this.GetString("select")

	//获取总记录数
  var count int64
	if TypeName==""{
		count, _ = qs.RelatedSel("ArticleType").Count()
	}else {
		count, _ = qs.RelatedSel("ArticleType").Filter("ArticleType__TypeArticle",TypeName).Count()
	}

	pageIndex := 5
	//向上取值
	pageCount := math.Ceil(float64(count) / float64(pageIndex))
	//获取页码
	pageNum, err := this.GetInt("pageNum")
	if err != nil {
		pageNum = 1
	}

	//一页显示几条记录
	if TypeName == "" {
		_, err = qs.Limit(pageIndex, pageIndex*(pageNum-1)).RelatedSel("ArticleType").All(&articles)
	} else {
		_, err = qs.Limit(pageIndex, pageIndex*(pageNum-1)).
			RelatedSel("ArticleType").
			Filter("ArticleType__TypeArticle", TypeName).
			All(&articles)
	}

	//查询所有文章类型
	var articleType []models.ArticleType
	_, err = o.QueryTable("ArticleType").All(&articleType)
	this.Data["ArticleTypes"] = articleType
	//获取总页数
	this.Data["count"] = count
	//最后一页
	this.Data["pageCount"] = pageCount
	this.Data["pageNum"] = pageNum
	this.Data["TypeName"] = TypeName

	this.Data["articles"] = articles

	this.LayoutSections=make(map[string]string)
	this.LayoutSections["indexJs"]="indexJs.html"

	this.Layout="layout.html"
	this.TplName = "index.html"
}

func (this *ArticleControllers) ShowArticle() {
	o := orm.NewOrm()
	var articleType []models.ArticleType

	_, err := o.QueryTable("ArticleType").All(&articleType)

	if err != nil {
		beego.Error("", err)
		this.TplName = "add.html"
		return
	}

	this.Data["ArticleType"] = articleType
	this.Layout = "layout.html"
	this.TplName = "add.html"
}

//上传文章
func (this *ArticleControllers) HandleArticle() {

	Title := this.GetString("articleName")
	content := this.GetString("content")
	typeName := this.GetString("select")

	if Title == "" || content == "" {
		beego.Error("文章不能为空")
		this.Data["errmsg"] = "文章不能为空"
		this.TplName = "add.html"
		return
	}

	file, header, err := this.GetFile("uploadName")

	if err != nil {
		beego.Error("获取图片失败")
		this.Data["errmsg"] = "获取图片失败"
		this.TplName = "add.html"
		return
	}

	defer file.Close()

	if header.Size < 50000 {
		beego.Error("插入图片失败")
		this.Data["errmsg"] = "上传图片过大限制5000KB"
		this.TplName = "add.html"
		return
	}

	ext := path.Ext(header.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".txt" {
		beego.Error("图片格式错误 jpg , png")
		this.Data["errmsg"] = "文件格式错误 jpg ,png"
		this.TplName = "add.html"
		return
	}

	fileName := time.Now().Format("2006-01-02-15-4-5-3333")

	err = this.SaveToFile("uploadName", "./static/img/"+fileName+ext)
	if err != nil {
		beego.Error("插入图片失败", err)
		this.Data["errmsg"] = "上传图片失败"
		this.TplName = "add.html"
		return
	}

	var article models.Article

	o := orm.NewOrm()
	article.Title = Title
	article.Content = content
	article.Img = "/static/img/" + fileName + ext
	var articleType models.ArticleType

	articleType.TypeArticle = typeName
	err = o.Read(&articleType, "TypeArticle")

	if err != nil {
		beego.Error("插入数据失败", err)
		this.TplName = "add.html"
		return
	}
	article.ArticleType = &articleType

	_, err = o.Insert(&article)

	if err != nil {
		beego.Error("插入数据失败", err)
		this.Redirect("index.html", 302)
		return
	}

	this.Layout="layout.html"
	this.Redirect("index.html", 302)

}

//获取详情
func (this *ArticleControllers) ShowContent() {

	id, err := this.GetInt("id")
	if err != nil {
		beego.Error("获取Id失败", err)
		this.Redirect("index", 302)
		return
	}

	o := orm.NewOrm()
	var article models.Article

	article.Id = id

	err = o.Read(&article)

	//_,err=o.LoadRelated(&article,"Users")

	var users []models.User
	_,err=o.QueryTable("User").Filter("Articles__Article__Id",id).Distinct().All(&users)
	this.Data["users"] = users



	if err != nil {
		beego.Error("查询Id失败", err)
		this.Redirect("index", 302)
		return
	}
	//更新阅读次数
	article.ReadCount += 1

	IdNum, err := o.Update(&article)
	if err != nil {
		beego.Error("更新数据d失败", err)
		this.Redirect("index", 302)
		return
	}

	fmt.Println("ID", IdNum)
	userName:=this.GetSession("userName")

	var user models.User
	user.Name = userName.(string)

	_=o.Read(&user,"Name")
	m2m:=o.QueryM2M(&article ,"Users")

	 i ,err:=m2m.Add(user)
	fmt.Println("m2m.Add", i)



	//返回数据
	this.Data["article"] = article
	this.Layout="layout.html"
	this.TplName = "content.html"

}

//编辑更新内容
func (this *ArticleControllers) ShowUpdate() {

	id, err := this.GetInt("id")
	if err != nil {
		beego.Error("Update获取Id失败", err)
		this.Data["errmsg"] = "Update获取Id失败"
		this.Redirect("/index", 302)
		return
	}
	o := orm.NewOrm()
	var article models.Article
	article.Id = id

	err = o.Read(&article)
	if err != nil {
		beego.Error("Update读取失败", err)
		this.Data["errmsg"] = "Update读取article失败"
		this.Redirect("/index", 302)
		return
	}

	this.Data["update"] = article
	this.Layout="layout.html"
	this.TplName = "update.html"

}

//编辑更新
func (this *ArticleControllers) HandleUpdate() {
	title := this.GetString("articleName")
	content := this.GetString("content")

	id, _ := this.GetInt("id")

	if title == "" || content == "" {
		beego.Error("文章不能为空")
		this.Data["errmsg"] = "文章不能为空"
		this.Redirect("/update?id="+strconv.Itoa(id), 302)
		return
	}

	file, header, err := this.GetFile("uploadname")

	if err != nil {
		beego.Error("获取图片失败")
		this.Redirect("/update?id="+strconv.Itoa(id), 302)
		return
	}
	defer file.Close()

	if header.Size <= 50000 {
		beego.Error("图片太大了")
		this.Data["errmsg"] = "图片太大了"
		this.Redirect("/update?id="+strconv.Itoa(id), 302)
		return
	}

	ext := path.Ext(header.Filename)

	date := time.Now().Format("2006-01-02-15:04:05-1111")

	err = this.SaveToFile("uploadname", "./static/"+date+ext)

	if err != nil {
		fmt.Println("this.SaveToFile", err)
		this.Redirect("/update?id="+strconv.Itoa(id), 302)
		return
	}

	o := orm.NewOrm()

	var article models.Article

	article.Id = id

	err = o.Read(&article)
	if err != nil {
		fmt.Println("o.Read(&article)", err)
		this.Redirect("/update?id="+strconv.Itoa(id), 302)
		return
	}

	article.Title = title
	article.Content = content
	article.Img = "/static/" + date + ext

	_, err = o.Update(&article)

	page := this.GetString("pageNum")

	this.Data["articles"] = article

	this.Redirect("/index?pageNum=page"+page, 302)

}

//删除内容
func (this *ArticleControllers) ShowDelete() {
	id, err := this.GetInt("id")
	if err != nil {
		beego.Error("ShowDelete获取Id失败", err)
		this.Data["errmsg"] = "ShowDelete获取Id失败"
		this.Redirect("/index", 302)
		return
	}

	o := orm.NewOrm()
	var article models.Article
	article.Id = id

	err = o.Read(&article, "Id")
	if err != nil {
		beego.Error("ShowDelete 读取数据库", err)
		this.Data["errmsg"] = "ShowDelete读取数据库"
		this.Redirect("/index", 302)
		return
	}

	num, err := o.Delete(&article, "Id")
	if err != nil {
		beego.Error("ShowDelete 读取数据库", err)
		this.Data["errmsg"] = "ShowDelete读取数据库"
		this.Redirect("/index", 302)
		return
	}

	this.Data["errmsg"] = num
	this.Redirect("/article/index", 302)

}

//展示类型
func (this *ArticleControllers) ShowAddType() {

	o := orm.NewOrm()
	var articleType []models.ArticleType

	qs := o.QueryTable("ArticleType")

	_, err := qs.All(&articleType)

	if err != nil {
		beego.Error("获取Id失败", err)
		this.Redirect("/article/addType", 302)
		return
	}

	this.Data["addType"] = articleType

	this.LayoutSections=make(map[string]string)
	this.LayoutSections["indexJs"]="indexJs.html"

	this.Layout="layout.html"
	this.TplName = "addType.html"

}

//添加类型
func (this *ArticleControllers) HandleAddType() {
	//获取信息
	typeName := this.GetString("typeName")

	if typeName == "" {
		beego.Error("文章类型能为空")
		this.Redirect("/article/addType", 302)
		return
	}

	o := orm.NewOrm()

	var articleType models.ArticleType

	articleType.TypeArticle = typeName

	_, err := o.Insert(&articleType)

	if err != nil || articleType.TypeArticle == "" {
		beego.Error("插入文章类型数据失败", err)
		this.Redirect("/article/addType", 302)
		return
	}

	//this.Data["addType"] = articleType
	this.Redirect("/article/addType", 302)

}

//删除文章类型
func (this *ArticleControllers) DeleteType() {
	id, err := this.GetInt("id")
	if err != nil {
		beego.Error("获取Id失败", err)
		this.Redirect("/article/addType", 302)
		return
	}

	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.Id = id
	err = o.Read(&articleType)
	if err != nil {
		beego.Error("获取Id失败", err)
		this.Redirect("/article/addType", 302)
		return
	}
	num, err := o.Delete(&articleType, "Id")

	if err != nil {
		beego.Error("删除失败", err)
		this.Redirect("/article/addType", 302)
		return
	}

	this.Redirect("/article/addType", 302)
	beego.Info("删除文章类型成功", num)
}
