package controllers

import (
	"fmt"
	"strconv"
	"strings"
	"webcron-source/app/goreman"
	"webcron-source/app/models"
)

type ApiController struct {
	BaseController
}

func (this *ApiController) getProcNamesByIds(ids string, needProcStatus int8) []string {
	idList := strings.Split(ids, ",")

	var names []string
	for _, v := range idList {
		id, _ := strconv.Atoi(v)
		if id < 1 {
			continue
		}
		ptask, err := models.GetGoremanById(id)
		if err == nil {
			if ptask.RunStatus != needProcStatus{
				continue
			}
			names = append(names, ptask.Name)
		}
	}

	return names
}

func (this *ApiController) Restart() {
	ids := this.GetString("ids")
	names := this.getProcNamesByIds(ids, goreman.PROC_STATUS_RUNNING)
	err := goreman.RunOpt("restart", names...)
	if err != nil {
		this.ajaxMsg(err.Error(), MSG_ERR)
	}

	this.ajaxMsg(fmt.Sprintf("重启%v成功!!", names), MSG_OK)
}

func (this *ApiController) Stop() {
	ids := this.GetString("ids")
	names := this.getProcNamesByIds(ids, goreman.PROC_STATUS_RUNNING)

	err := goreman.RunOpt("stop", names...)
	if err != nil {
		this.ajaxMsg(err.Error(), MSG_ERR)
	}

	this.ajaxMsg(fmt.Sprintf("停止%v成功!!", names), MSG_OK)
}

func (this *ApiController) Start() {
	ids := this.GetString("ids")
	idList := strings.Split(ids, ",")

	var names []string
	for _, v := range idList {
		id, _ := strconv.Atoi(v)
		if id < 1 {
			continue
		}
		ptask, err := models.GetGoremanById(id)
		if err == nil {
			if ptask.RunStatus == goreman.PROC_STATUS_RUNNING{
				continue
			}
			names = append(names, ptask.Name)
		}
	}

	err := goreman.RunOpt("start", names...)
	if err != nil {
		this.ajaxMsg(err.Error(), MSG_ERR)
	}

	this.ajaxMsg(fmt.Sprintf("启动%v成功!!", names), MSG_OK)
}

func (this *ApiController) Batch0pt()  {
	cmd := this.GetString("cmd")

	err := goreman.RunOpt(cmd, "")
	switch cmd {
	case "start-all":
		if err != nil {
			this.ajaxMsg("启动所有失败" + err.Error(), MSG_ERR)
		}
		this.ajaxMsg(fmt.Sprintf("启动所有成功!!"), MSG_OK)
	case "stop-all":
		if err != nil {
			this.ajaxMsg("停止所有失败" + err.Error(), MSG_ERR)
		}
		this.ajaxMsg(fmt.Sprintf("停止所有成功!!"), MSG_OK)
	case "restart-all":
		if err != nil {
			this.ajaxMsg("重启所有失败" + err.Error(), MSG_ERR)
		}
		this.ajaxMsg(fmt.Sprintf("重启所有成功!!"), MSG_OK)
	}

	this.ajaxMsg(fmt.Sprintf("未知的命令：%s", cmd), MSG_OK)
}