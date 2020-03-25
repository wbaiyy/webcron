package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

const (
	PTASK_STATUS_NOT_STARTED = 1
	PTASK_STATUS_RUNNING = 2
	PTASK_STATUS_STOP = 3
)
type PTask struct {
	BaseModel
	Id int
	Name string
	Command string
	RetryTimes int
	IntervalTime int64
	Description string
	GroupId int
	UpdateTime int64
	Status int8
	RunStatus int8
	OutputFile string
	NotifyUsers string
	Num int
}

func (g *PTask) TableName() string {
	return TableName("ptasks")
}

func (g *PTask) Add() (int64, error) {
	if g.UpdateTime == 0 {
		g.UpdateTime = time.Now().Unix()
	}
	return orm.NewOrm().Insert(g)
}

func (g *PTask) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(g, fields...); err != nil {
		return err
	}
	return nil
}

func GetGoremanById(id int) (*PTask, error) {
	goreman := &PTask{
		Id: id,
	}

	err := orm.NewOrm().Read(goreman)
	if err != nil {
		return nil, err
	}
	return goreman, nil
}

func GetGoremanByName(name string) (*PTask, error) {
	goreman := &PTask{
		Name: name,
	}

	err := orm.NewOrm().Read(goreman, "name")
	if err != nil {
		return nil, err
	}
	return goreman, nil
}

func(g *PTask) GetList(page, pageSize int, filters ...interface{}) ([]*PTask, int64) {
	offset := (page - 1) * pageSize
	goremans := make([]*PTask, 0)

	query := orm.NewOrm().QueryTable(g.TableName())

	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			if filters[k].(string) == "command" {
				query = query.Filter("command__icontains", filters[k+1])
			} else {
				query = query.Filter(filters[k].(string), filters[k+1])
			}
		}
	}
	total, _ := query.Count()
	query.OrderBy("-id").Limit(pageSize, offset).All(&goremans)

	return goremans, total
}

func GetAllGoreman() ([]*PTask, int64) {
	goremans := make([]*PTask, 0)

	query := orm.NewOrm().QueryTable(new(PTask).TableName())
	total, _ := query.Count()
	query.OrderBy("-id").All(&goremans)

	return goremans, total
}

func PTaskDel(id int) error {
	_, err := orm.NewOrm().QueryTable(TableName("ptasks")).Filter("id", id).Delete()
	return err
}

func SetPtaskRunningStatus(name string, runningStatus int8) (int64, error){
	o := orm.NewOrm()
	ptaskModel := &PTask{}
	result, err := o.Raw(fmt.Sprintf("UPDATE %s SET run_status=? WHERE name =?", ptaskModel.TableName()), runningStatus, name).Exec()

	if err == nil {
		num, _ := result.RowsAffected()
		return num, nil
	} else {
		return 0, err
	}
}
