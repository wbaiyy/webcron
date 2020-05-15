package http

import (
	"fmt"
	"github.com/astaxie/beego"
	"html/template"
	"net/http"
	"github.com/wbaiyy/webcron-source/app/controllers"
)

const VERSION = "1.0.0"

func InitHttp(done chan bool) {
	setRoutes()
	// 生产环境不输出debug日志
	if beego.AppConfig.String("runmode") == "prod" {
		beego.SetLevel(beego.LevelInformational)
	}
	beego.AppConfig.Set("version", VERSION)
	beego.BConfig.WebConfig.Session.SessionOn = true

	go func() {
		defer func() {
			fmt.Println("http done")
			done<- true
		}()
		beego.Run(fmt.Sprintf("%s:%s", beego.AppConfig.String("httpaddr"), beego.AppConfig.String("httpport")))
	}()
}

func setRoutes() {
	// 设置默认404页面
	beego.ErrorHandler("404", func(rw http.ResponseWriter, r *http.Request) {
		t, _ := template.New("404.html").ParseFiles(beego.BConfig.WebConfig.ViewsPath + "/error/404.html")
		data := make(map[string]interface{})
		data["content"] = "page not found"
		t.Execute(rw, data)
	})

	// 路由设置
	beego.Router("/", &controllers.MainController{}, "*:Index")
	beego.Router("/login", &controllers.MainController{}, "*:Login")
	beego.Router("/logout", &controllers.MainController{}, "*:Logout")
	beego.Router("/profile", &controllers.MainController{}, "*:Profile")
	beego.Router("/gettime", &controllers.MainController{}, "*:GetTime")
	beego.AutoRouter(&controllers.HelpController{})
	beego.AutoRouter(&controllers.TaskController{})
	beego.AutoRouter(&controllers.PTaskController{})
	beego.AutoRouter(&controllers.GroupController{})
	beego.AutoRouter(&controllers.ApiController{})
	beego.Router("/api/batch-opt", &controllers.ApiController{}, "*:Batch0pt")
}
