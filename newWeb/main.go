package main

import (
	"github.com/astaxie/beego"
	_ "newWeb/models"
	_ "newWeb/routers"
)

func main() {
	err := beego.AddFuncMap("prePage", PrePage)
	err = beego.AddFuncMap("nextPage", NextPage)
	if err != nil {
		beego.Error("函数映射失败", err)
	}
	beego.Run()
}

func PrePage(pageNum int) int {
	if pageNum <= 1 {
		return 1
	}
	return pageNum - 1
}

func NextPage(pageNum int, pageCount float64) int {

	if pageNum >= int(pageCount) {
		return int(pageCount)
	}
	return pageNum + 1
}
