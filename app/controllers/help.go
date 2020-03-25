package controllers

import (
	"webcron-source/app/libs"
)

type HelpController struct {
	BaseController
}

func (this *HelpController) Index() {

	this.Data["pageTitle"] = "使用帮助"
	this.display()
}

func (this *HelpController) Test() {
	MsgCenterService := libs.NewMsgCenterService()
	sendWappers := make(map[string]string)
	sendWappers["vv"] = "607084"
	MsgCenterService.Send("测试消息1", "收到的消息内容", sendWappers)
}
