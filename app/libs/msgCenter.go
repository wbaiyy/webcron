package libs

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"errors"
	"io/ioutil"
	"log"
)

const SERVICE_URL  = "/api-source/index"

type msgCenterParamContent struct {
	Channel, To, Title, Content string
}

type msgCenterResponse struct {
	Code string `json:"code"`
	Msg string `json:"msg"`
	data interface{}
}

type msgCenterParams struct {
	Account string `json:"account"`
	Password string  `json:"password"`
	ApiKey string `json:"api_key"`
	Data []msgCenterParamContent `json:"data"`
}

type MsgCenterService struct {
	Getaway, Password, Account, ApiKey string
	ConnectTimeout, Timeout int
}

func NewMsgCenterService() *MsgCenterService{
	connectTimeout, err := beego.AppConfig.Int("msgCenterService.connectTimeout ")
	if err!= nil {
		connectTimeout = 3
	}
	timeout,err :=beego.AppConfig.Int("msgCenterService.timeout")
	if err != nil {
		timeout = 3
	}

	return &MsgCenterService{
		Getaway: beego.AppConfig.String("msgCenterService.getaway"),
		Account:  beego.AppConfig.String("msgCenterService.account"),
		Password:  beego.AppConfig.String("msgCenterService.password"),
		ApiKey: beego.AppConfig.String("msgCenterService.apiKey"),
		ConnectTimeout: connectTimeout,
		Timeout: timeout,
	}
}

func (mcs *MsgCenterService) Send(title, content string, sendMapper map[string]string) error {
	var serivceUrl = fmt.Sprintf("%s%s", mcs.Getaway, SERVICE_URL)
	//var serivceUrl = "http://10.37.3.155/test.php"
	req := httplib.Post(serivceUrl)
	req.Header("Content-Type", "application/x-www-form-urlencoded")
	req.Param("account", mcs.Account)
	req.Param("password", mcs.Password)
	req.Param("api_key", mcs.ApiKey)

	var data []msgCenterParamContent

	var i = 0
	for channel, users := range sendMapper {
		data = append(data, msgCenterParamContent{
			channel, users, title, content,
		})
		req.Param(fmt.Sprintf("data[%d][channel]", i), channel)
		req.Param(fmt.Sprintf("data[%d][to]", i), users)
		req.Param(fmt.Sprintf("data[%d][title]", i), title)
		req.Param(fmt.Sprintf("data[%d][content]", i), content)
		i++
	}
	s := &msgCenterParams{
		Account: mcs.Account,
		Password: mcs.Password,
		ApiKey:  mcs.ApiKey,
		Data: data,
	}
	b, _ := json.Marshal(s)
	log.Println("【MsgCenterService】发送消息内容:", string(b))

	resp, _ := req.Response()
	//resp, _ := http.Post(serivceUrl, "application/x-www-form-urlencoded", strings.NewReader(string(b)))
	defer resp.Body.Close()
	//io.Reader
	body, err:= ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("【MsgCenterService】读取响应消息失败,原因：%s", err.Error()))
	}
	log.Println("【MsgCenterService】响应消息内容:", string(body))

	response := &msgCenterResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return errors.New(fmt.Sprintf("【MsgCenterService】响应内容解析失败,原因：%s", err))
	}

	if response.Code != "0" {
		return  errors.New(fmt.Sprintf("【MsgCenterService】发送消息失败,原因：%s", response.Msg))
	}

	return nil
}