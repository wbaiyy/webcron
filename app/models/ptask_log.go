package models

import (
	"github.com/astaxie/beego/orm"
)

type PTaskLog struct {
	Id          int
	PtaskId      int
	Output      string
	Error       string
	Status      int
	ProcessTime int
	CreateTime  int64
}

func (t *PTaskLog) TableName() string {
	return TableName("ptask_log")
}

func PTaskLogAdd(t *PTaskLog) (int64, error) {
	return orm.NewOrm().Insert(t)
}

func PTaskLogGetList(page, pageSize int, filters ...interface{}) ([]*PTaskLog, int64) {
	offset := (page - 1) * pageSize

	logs := make([]*PTaskLog, 0)

	query := orm.NewOrm().QueryTable(TableName("ptask_log"))
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}

	total, _ := query.Count()
	query.OrderBy("-id").Limit(pageSize, offset).All(&logs)

	return logs, total
}

func PTaskLogGetById(id int) (*PTaskLog, error) {
	obj := &PTaskLog{
		Id: id,
	}

	err := orm.NewOrm().Read(obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func PTaskLogDelById(id int) error {
	_, err := orm.NewOrm().QueryTable(TableName("ptask_log")).Filter("id", id).Delete()
	return err
}

func PTaskLogDelByTaskId(taskId int) (int64, error) {
	return orm.NewOrm().QueryTable(TableName("ptask_log")).Filter("ptask_id", taskId).Delete()
}

