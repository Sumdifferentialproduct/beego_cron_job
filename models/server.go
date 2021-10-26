package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type TaskServer struct {
	Id int
	GroupId int
	ServerName string
	ServerAccount string
	ServerOuterIp string
	ServerIp string
	Port int
	Password string
	Type int
	Detail string
	CreateTime int64
	UpdateTime int64
	Status int
}

func (server *TaskServer) TableName() string {
	return TableName("task_server")
}

func TaskSeverGetById(id int) (*TaskServer, error) {
	obj := &TaskServer{
		Id: id,
	}
	err := orm.NewOrm().Read(obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func TaskServerGetList(page, pageSize int, filters ...interface{}) ([]*TaskServer, int64) {
	query := orm.NewOrm().QueryTable(TableName("task_server"))
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	total, _ := query.Count()
	list := make([]*TaskServer, 0)
	offset := (page - 1) * pageSize
	_, err := query.OrderBy("-id").Limit(pageSize, offset).All(&list)
	if err != nil {
		fmt.Println(err)
	}
	return list, total
}
