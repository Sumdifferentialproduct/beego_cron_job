package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

const (
	TASK_SUCESS  = 0
	TASK_ERROR   = -1
	TSAK_TIMEOUT = -2
)

type Task struct {
	Id int
	GroupId int
	ServerId int
	TaskName string
	Description string
	CronSpec string
	Concurrent int
	Command string
	Timeout int
	ExecuteTimes int
	PrevTime int64
	IsNotify int
	NotifyType int
	NotifyUserIds string
	Status int
	CreateTime int64
	CreateId int
	UpdateTime int64
	UpdateId int
}

func (task *Task) TableName() string {
	return TableName("task")
}

func (task *Task) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(task, fields...); err != nil {
		return err
	}
	return nil
}

func TaskGetList(page, pageSize int, filters ...interface{}) ([]*Task, int64) {
	query := orm.NewOrm().QueryTable(TableName("task"))
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	total, _ := query.Count()
	list := make([]*Task, 0)
	offset := (page - 1) * pageSize
	_, err := query.OrderBy("-id").Limit(pageSize, offset).All(&list)
	if err != nil {
		fmt.Println(err)
	}
	return list, total
}

func TaskGetById(id int) (*Task, error) {
	task := &Task{
		Id: id,
	}
	err := orm.NewOrm().Read(task)
	if err != nil {
		return nil, err
	}
	return task, nil
}
