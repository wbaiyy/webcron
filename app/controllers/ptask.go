package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"time"
	"github.com/wbaiyy/webcron-source/app/goreman"
	"github.com/wbaiyy/webcron-source/app/libs"
	"github.com/wbaiyy/webcron-source/app/models"
)

type PTaskController struct {
	BaseController
}

func (this *PTaskController) Add()  {
	if this.isPost() {
		PTask := new(models.PTask)
		PTask.Name = this.GetString("goreman_name")
		PTask.Command = this.GetString("goreman_command")
		PTask.Description = this.GetString("goreman_description")
		PTask.GroupId, _ = this.GetInt("group_id")
		PTask.RetryTimes, _ = this.GetInt("goreman_retry_times")
		PTask.NotifyUsers  = this.GetString("goreman_notify_users", "")
		PTask.IntervalTime, _ = this.GetInt64("goreman_interval_time")
		PTask.Status, _ = this.GetInt8("goreman_status")
		PTask.OutputFile= this.GetString("goreman_output_file", "")
		PTask.Num, _ = this.GetInt("goreman_num")
		PTask.RunStatus = models.PTASK_STATUS_NOT_STARTED
		PTask.UpdateTime = time.Now().Unix()

		if PTask.Name == "" || PTask.Command == "" {
			this.ajaxMsg("脚本名称和命令不能为空", MSG_ERR)
		}
		if PTask.OutputFile != "" {
			err := this.checkAndCreateFile(PTask.OutputFile)
			if err != nil {
				this.ajaxMsg("创建任务输入文件时间错误"+err.Error(), MSG_ERR)
			}
		}
		maxProcNum := 20;
		if  maxProcNumConfig, error := beego.AppConfig.Int("ptask.max_proc_num"); error != nil {
			maxProcNum = maxProcNumConfig
		}
		if PTask.Num > maxProcNum || PTask.Num < 1 {
			this.ajaxMsg(fmt.Sprintf("单个任务最多设置%d个进程，最少1个", maxProcNum), MSG_ERR)
		}

		if _, err := PTask.Add(); err != nil {
			this.ajaxMsg(err.Error(), MSG_ERR)
		}

		if PTask.Status == 1 {
			goreman.PullProcs(PTask)
		}


		this.ajaxMsg("", MSG_OK)
	}

	groups, _ := models.TaskGroupGetList(1, 100)
	this.Data["groups"] = groups
	this.Data["pageTitle"] = "添加常驻任务"

	this.Data["pageNo"], _ = this.GetInt("pageNo", 1)
	this.Data["command"] = this.GetString("command", "")

	this.display()
}
func (this *PTaskController) Edit()  {
	id, _ := this.GetInt("id")
	pTaskModel, err := models.GetGoremanById(id)
	if err != nil {
		this.showMsg(err.Error())
	}

	if this.isPost() {
		if pTaskModel.RunStatus == models.PTASK_STATUS_RUNNING {
			this.ajaxMsg("任务正在运行中，不能修改", MSG_ERR)
		}

		pTaskModel.Name = this.GetString("goreman_name")
		pTaskModel.Command = this.GetString("goreman_command")
		pTaskModel.Description = this.GetString("goreman_description")
		pTaskModel.GroupId, _ = this.GetInt("group_id")
		pTaskModel.RetryTimes, _ = this.GetInt("goreman_retry_times")
		pTaskModel.IntervalTime, _ = this.GetInt64("goreman_interval_time")
		pTaskModel.Status, _ = this.GetInt8("goreman_status")
		pTaskModel.OutputFile= this.GetString("goreman_output_file", "")
		pTaskModel.NotifyUsers= this.GetString("goreman_notify_users", "")
		pTaskModel.Num, _ = this.GetInt("goreman_num")
		pTaskModel.UpdateTime = time.Now().Unix()

		if pTaskModel.Name == "" || pTaskModel.Command == "" {
			this.ajaxMsg("任务名称和命令不能为空", MSG_ERR)
		}
		if pTaskModel.OutputFile != "" {
			err = this.checkAndCreateFile(pTaskModel.OutputFile)
			if err != nil {
				this.ajaxMsg("创建任务输入文件时间错误"+err.Error(), MSG_ERR)
			}
		}

		maxProcNum := 20;
		if  maxProcNumConfig, error := beego.AppConfig.Int("ptask.max_proc_num"); error != nil {
			maxProcNum = maxProcNumConfig
		}
		if pTaskModel.Num > maxProcNum || pTaskModel.Num < 1 {
			this.ajaxMsg(fmt.Sprintf("单个任务最多设置%d个进程，最少1个", maxProcNum), MSG_ERR)
		}

		if err := pTaskModel.Update(); err != nil {
			this.ajaxMsg(err.Error(), MSG_ERR)
		}

		if pTaskModel.Status == 1 {
			goreman.PullProcs(pTaskModel)
		}
		//更新pros信息
		goreman.UpdateProc(pTaskModel)

		this.ajaxMsg("", MSG_OK)
	}

	groups, _ := models.TaskGroupGetList(1, 100)
	this.Data["groups"] = groups
	this.Data["pTask"] = pTaskModel
	this.Data["pageTitle"] = "编辑常驻任务"
	this.Data["pageNo"], _ = this.GetInt("pageNo", 1)
	this.Data["command"] = this.GetString("command", "")

	this.display()
}

func (this *PTaskController) checkAndCreateFile(filePath string) error {
	if !libs.Exist(filePath) {
		return libs.FileCreate(filePath)
	}

	return nil
}


func (this *PTaskController) List() {
	pageNo, _ := this.GetInt("page", 1)
	pageSize, _ := this.GetInt("pageSize", 10)
	command := this.GetString("command", "")

	filters := make([]interface{}, 0)
	if command != "" {
		filters = append(filters, "command", command)
	}

	groupId, _ := this.GetInt("groupid")
	if groupId > 0 {
		filters = append(filters, "group_id", groupId)
	}

	pTask := new(models.PTask)
	goremans, totalCount := pTask.GetList(pageNo, pageSize, filters...)

	list := make([]map[string]interface{}, len(goremans))
	procs := goreman.GetProcs()

	for k, v := range goremans {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["name"] = v.Name
		row["update_time"] = v.UpdateTime
		row["retry_times"] = v.RetryTimes
		row["interval_time"] = v.IntervalTime
		//row["notify_users"] = v.NotifyUsers
		row["command"] = v.Command
		row["description"] = v.Description
		row["status"] = v.Status
		row["run_status"] = v.RunStatus
		row["current_retry_times"] = "-"
		row["commandParam"] =   command
		row["pageNoParam"] =   pageNo
		_, ok := procs[v.Name]
		if ok {
			row["current_retry_times"] = procs[v.Name].FailureTimes
		}
		list[k] = row
	}

	// 分组列表
	groups, _ := models.TaskGroupGetList(1, 100)

	this.Data["pageTitle"] = "任务列表"
	this.Data["list"] = list
	this.Data["groups"] = groups
	this.Data["groupid"] = groupId
	this.Data["commandParam"] = command
	this.Data["pageNoParam"] = pageNo
	this.Data["pageBar"] = libs.NewPager(pageNo, int(totalCount), pageSize,
		beego.URLFor("PTaskController.List", "command", command), true).ToString()
	this.display()
}

func (this *PTaskController) Stop() {
	id, _ := this.GetInt("id", 0)
	if id <= 0 {
		this.showMsg("ID不能为空")
	}
	pTaskModel, err := models.GetGoremanById(id)

	if err != nil {
		this.showMsg(err.Error())
	}
	if pTaskModel.Status == 2 {
		this.showMsg("当前任务不可用状态")
	}
	if pTaskModel.RunStatus == models.PTASK_STATUS_NOT_STARTED  || pTaskModel.RunStatus == models.PTASK_STATUS_STOP{
		this.showMsg("当前未开启状态或已停止")
	}

	err = goreman.RunOpt("stop", pTaskModel.Name)
	if err != nil {
		this.showMsg(err.Error())
	}

	refer := this.Ctx.Request.Referer()
	if refer == "" {
		refer = beego.URLFor("PtaskController.List")
	}
	//异步停止任务，等待状态更新完成
	time.Sleep(time.Millisecond * 50)

	this.redirect(refer)
}

func (this *PTaskController) Run() {
	id, _ := this.GetInt("id", 0)
	if id <= 0 {
		this.showMsg("ID不能为空")
	}
	pTaskModel, err := models.GetGoremanById(id)

	if err != nil {
		this.showMsg(err.Error())
	}
	if pTaskModel.Status == 2 {
		this.showMsg("当前任务不可用状态")
	}
	if pTaskModel.RunStatus == models.PTASK_STATUS_RUNNING {
		this.showMsg("当前运行中状态1")
	}
	pTask := goreman.GetProcByName(pTaskModel.Name)
	if pTask == nil || pTask.RunStatus == models.PTASK_STATUS_RUNNING{
		this.showMsg("当前运行中状态2")
	}

	err = goreman.RunOpt("person-start", pTaskModel.Name)
	if err != nil {
		this.showMsg("启动任务失败：" + err.Error())
	}

	refer := this.Ctx.Request.Referer()
	if refer == "" {
		refer = beego.URLFor("PtaskController.List")
	}
	this.redirect(refer)
}

func (this *PTaskController) Batch() {
	action := this.GetString("action")
	ids := this.GetStrings("ids")
	if len(ids) < 1 {
		this.ajaxMsg("请选择要操作的项目", MSG_ERR)
	}

	switch action {
	case "start":
		var nameList []string
		for _, v := range ids {
			id, _ := strconv.Atoi(v)
			if id < 1 {
				continue
			}
			ptask, err := models.GetGoremanById(id)
			if err == nil {
				if goreman.GetProcByName(ptask.Name) != nil && goreman.GetProcByName(ptask.Name).RunStatus != models.PTASK_STATUS_RUNNING {
					nameList = append(nameList, ptask.Name)
				}
			}
		}

		goreman.RunOpt("start", nameList...)

	case "stop":
		var nameList []string
		for _, v := range ids {
			id, _ := strconv.Atoi(v)
			if id < 1 {
				continue
			}
			ptask, err := models.GetGoremanById(id)
			if err == nil {
				if goreman.GetProcByName(ptask.Name) != nil && goreman.GetProcByName(ptask.Name).RunStatus != models.PTASK_STATUS_STOP {
					nameList = append(nameList, ptask.Name)
				}
			}
		}
		goreman.RunOpt("stop", nameList...)

	case "delete":
		for _, v := range ids {
			id, _ := strconv.Atoi(v)
			if id < 1 {
				continue
			}
			models.PTaskDel(id)
		}
	}

	for _, v := range ids {
		id, _ := strconv.Atoi(v)
		if id < 1 {
			continue
		}

	}

	this.ajaxMsg("", MSG_OK)
}

// 任务执行日志列表
func (this *PTaskController) FailLogs() {
	ptaskId, _ := this.GetInt("id")
	page, _ := this.GetInt("page")
	if page < 1 {
		page = 1
	}

	ptask, err := models.GetGoremanById(ptaskId)
	if err != nil {
		this.showMsg(err.Error())
	}

	result, count := models.PTaskLogGetList(page, this.pageSize, "ptask_id", ptask.Id)

	list := make([]map[string]interface{}, len(result))
	for k, v := range result {
		row := make(map[string]interface{})
		row["id"] = v.Id
		row["start_time"] = beego.Date(time.Unix(v.CreateTime, 0), "Y-m-d H:i:s")
		row["error_info"] = v.Error
		list[k] = row
	}

	this.Data["pageNo"], _ = this.GetInt("pageNo", 1)
	this.Data["command"] = this.GetString("command", "")
	this.Data["pageTitle"] = "任务失败日志"
	this.Data["list"] = list
	this.Data["task"] = ptask
	this.Data["pageBar"] = libs.NewPager(page, int(count), this.pageSize, beego.URLFor("TaskController.FailLogs", "id", ptaskId), true).ToString()
	this.display()
}

// 批量操作日志
func (this *PTaskController) LogBatch() {
	action := this.GetString("action")
	ids := this.GetStrings("ids")
	if len(ids) < 1 {
		this.ajaxMsg("请选择要操作的项目", MSG_ERR)
	}
	for _, v := range ids {
		id, _ := strconv.Atoi(v)
		if id < 1 {
			continue
		}
		switch action {
		case "delete":
			models.PTaskLogDelById(id)
		}
	}

	if this.IsAjax() {
		 this.ajaxMsg("", MSG_OK)
		 return
	}
	refer := this.Ctx.Request.Referer()
	if refer == "" {
		refer = beego.URLFor("PtaskController.List")
	}
	this.redirect(refer)
}
