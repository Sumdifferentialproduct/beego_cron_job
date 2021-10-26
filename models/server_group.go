package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type ServerGroup struct {
	Id int
	GroupName string
	Description string
	Status int
	CreateTime int64
	UpdateTime int64
	CreateId int
	UpdateId int
}

func (servergroup *ServerGroup) TableName() string {
	return TableName("task_server_group")
}

func ServerGroupGetList(page, pageSize int, filters ...interface{}) ([]*ServerGroup, int64) {
	query := orm.NewOrm().QueryTable(TableName("task_server_group"))
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	total, _ := query.Count()
	list := make([]*ServerGroup, 0)
	offset := (page - 1) * pageSize
	_, err := query.OrderBy("-id").Limit(pageSize, offset).All(&list)
	if err != nil {
		fmt.Println(err)
	}
	return list, total
}

func ServerGroupGetById(id int) (*ServerGroup, error) {
	obj := &ServerGroup{
		Id: id,
	}
	err := orm.NewOrm().Read(obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}
