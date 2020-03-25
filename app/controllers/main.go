package controllers

import (
	"encoding/base64"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/utils"
	"log"
	"runtime"
	"strconv"
	"strings"
	"time"
	"webcron-source/app/jobs"
	"webcron-source/app/libs"
	"webcron-source/app/models"
)

type SsoPerson struct {
	DepartmentId string `json:"department_id"`
	UserNo	string `json:"user_no"`
	Username string
	Name string
	Email string
	Phone string
	IsGroupLeader string `json:"is_group_leader"`
	Created string
	Modified string
	Status string   //0正常， 1禁用
}

type MainController struct {
	BaseController
}

// 首页
func (this *MainController) Index() {
	this.Data["pageTitle"] = "系统概况"

	// 即将执行的任务
	entries := jobs.GetEntries(30)
	jobList := make([]map[string]interface{}, len(entries))
	for k, v := range entries {
		row := make(map[string]interface{})
		job := v.Job.(*jobs.Job)
		row["task_id"] = job.GetId()
		row["task_name"] = job.GetName()
		row["next_time"] = beego.Date(v.Next, "Y-m-d H:i:s")
		jobList[k] = row
	}

	// 最近执行的日志
	logs, _ := models.TaskLogGetList(1, 20)
	recentLogs := make([]map[string]interface{}, len(logs))
	for k, v := range logs {
		task, err := models.TaskGetById(v.TaskId)
		taskName := ""
		if err == nil {
			taskName = task.TaskName
		}
		row := make(map[string]interface{})
		row["task_name"] = taskName
		row["id"] = v.Id
		row["start_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
		row["process_time"] = float64(v.ProcessTime) / 1000
		row["ouput_size"] = libs.SizeFormat(float64(len(v.Output)))
		row["output"] = beego.Substr(v.Output, 0, 100)
		row["status"] = v.Status
		recentLogs[k] = row
	}

	// 最近执行失败的日志
	logs, _ = models.TaskLogGetList(1, 20, "status__lt", 0)
	errLogs := make([]map[string]interface{}, len(logs))
	for k, v := range logs {
		task, err := models.TaskGetById(v.TaskId)
		taskName := ""
		if err == nil {
			taskName = task.TaskName
		}
		row := make(map[string]interface{})
		row["task_name"] = taskName
		row["id"] = v.Id
		row["start_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
		row["process_time"] = float64(v.ProcessTime) / 1000
		row["ouput_size"] = libs.SizeFormat(float64(len(v.Output)))
		row["error"] = beego.Substr(v.Error, 0, 100)
		row["status"] = v.Status
		errLogs[k] = row
	}

	this.Data["recentLogs"] = recentLogs
	this.Data["errLogs"] = errLogs
	this.Data["jobs"] = jobList
	this.Data["cpuNum"] = runtime.NumCPU()
	this.display()
}

// 个人信息
func (this *MainController) Profile() {
	beego.ReadFromRequest(&this.Controller)
	user, _ := models.UserGetById(this.userId)

	if this.isPost() {
		flash := beego.NewFlash()
		user.Email = this.GetString("email")
		user.Update()
		password1 := this.GetString("password1")
		password2 := this.GetString("password2")
		if password1 != "" {
			if len(password1) < 6 {
				flash.Error("密码长度必须大于6位")
				flash.Store(&this.Controller)
				this.redirect(beego.URLFor(".Profile"))
			} else if password2 != password1 {
				flash.Error("两次输入的密码不一致")
				flash.Store(&this.Controller)
				this.redirect(beego.URLFor(".Profile"))
			} else {
				user.Salt = string(utils.RandomCreateBytes(10))
				user.Password = libs.Md5([]byte(password1 + user.Salt))
				user.Update()
			}
		}
		flash.Success("修改成功！")
		flash.Store(&this.Controller)
		this.redirect(beego.URLFor(".Profile"))
	}

	this.Data["pageTitle"] = "个人信息"
	this.Data["user"] = user
	this.display()
}

// 登录
func (this *MainController) LoginSource() {
	if this.userId > 0 {
		this.redirect("/")
	}
	beego.ReadFromRequest(&this.Controller)
	if this.isPost() {
		flash := beego.NewFlash()

		username := strings.TrimSpace(this.GetString("username"))
		password := strings.TrimSpace(this.GetString("password"))
		remember := this.GetString("remember")
		if username != "" && password != "" {
			user, err := models.UserGetByName(username)
			errorMsg := ""
			if err != nil || user.Password != libs.Md5([]byte(password+user.Salt)) {
				errorMsg = "帐号或密码错误"
			} else if user.Status == -1 {
				errorMsg = "该帐号已禁用"
			} else {
				user.LastIp = this.getClientIp()
				user.LastLogin = time.Now().Unix()
				models.UserUpdate(user)

				authkey := libs.Md5([]byte(this.getClientIp() + "|" + user.Password + user.Salt))
				if remember == "yes" {
					this.Ctx.SetCookie("auth", strconv.Itoa(user.Id)+"|"+authkey, 7*86400)
				} else {
					this.Ctx.SetCookie("auth", strconv.Itoa(user.Id)+"|"+authkey)
				}

				this.redirect(beego.URLFor("TaskController.List"))
			}
			flash.Error(errorMsg)
			flash.Store(&this.Controller)
			this.redirect(beego.URLFor("MainController.Login"))
		}
	}

	this.TplName = "main/login.html"
}

// 登录
func (this *MainController) Login() {
	if this.userId > 0 {
		this.redirect("/")
	}
	sid := this.GetString("sid", "")
	if sid == "" {
		params := make(map[string]string)
		referer := this.Ctx.Request.Referer()
		if referer == "" {
			referer = beego.URLFor("PTask.List")
		}

		params["struli"] = base64.StdEncoding.EncodeToString(
			[]byte(this.getHost() + beego.URLFor("MainController.Login") + "|" + referer))
 		this.redirect(this.getSsoUrl("login/index/sso", params))
	}

	response := httplib.Get(this.getSsoUrl("login/index/checksso", map[string]string{"sid": sid}))

	ssoPerson := SsoPerson{}
	body, _ := response.String()
	bodyBytes, _:= base64.StdEncoding.DecodeString(body)
	body = string(bodyBytes)
	if "fbd" == body || "" == body {
		this.loginException("登录失败,请稍后重试")
		return
	}

	err := json.Unmarshal(bodyBytes, &ssoPerson)
	if err != nil {
		log.Println("[json.unmarshal] sso user info error:", err)
	}
	if ssoPerson.Status == "1"{
		this.loginException("登录异常：账号已锁定")
		return
	}
	user, err := models.UserGetByName(ssoPerson.Username)
	now := time.Now().Unix()
	ip := this.getClientIp()
	var salt = ""
	var password = ""
	var userId int
	if user == nil {
		salt = libs.GetSlat(ssoPerson.Username, now, ip)
		password = libs.Md5([]byte("123456" + salt))
		id, err := models.UserAdd(&models.User{
			UserName : ssoPerson.Username,
			Email: ssoPerson.Email,
			LastLogin : now,
			LastIp : ip,
			Status: 0,
			Salt : salt,
			Password: password,
		})
		if err != nil {
			log.Println("[models.UserAdd] add user error:", err)
			this.loginException("登录异常：插入新账号失败")

		}
		userId = int(id)
	} else {
		user.LastLogin = now
		user.LastIp = ip
		salt = user.Salt
		userId = user.Id
		password = user.Password
		err := models.UserUpdate(user)
		if err != nil {
			log.Println("[models.UserUpdate] update user error:", err)
		}
	}

	authkey := libs.GetCookieAuthKey(password, salt)
	cookieTime, err := beego.AppConfig.Int("login.cooikeTime")
	if err != nil {
		cookieTime = 4 * 3600
	}

	this.Ctx.SetCookie("auth", strconv.Itoa(userId)+"|"+authkey, cookieTime)
	this.redirect(beego.URLFor("PTaskController.List"))
}

// 退出登录
func (this *MainController) Logout() {
	this.Ctx.SetCookie("auth", "")
	ssoLogoutUrl := this.getSsoUrl("login/index/loginout",
		map[string]string{"returnurl" : base64.StdEncoding.EncodeToString([]byte(this.getHost() + beego.URLFor("MainController.Login")))})
	this.redirect(ssoLogoutUrl)
}

// 获取系统时间
func (this *MainController) GetTime() {
	out := make(map[string]interface{})
	out["time"] = time.Now().UnixNano() / int64(time.Millisecond)
	this.jsonResult(out)
}

// 登录异常
func (this *MainController) loginException(message string) {
	this.Data["message"] = message

	this.Data["redirect"] = this.getSsoUrl("login/index/loginout",
		map[string]string{"returnurl": base64.StdEncoding.EncodeToString(
			[]byte(this.getHost() + beego.URLFor("MainController.Index")))})
	this.TplName = "error/message.html"
}
