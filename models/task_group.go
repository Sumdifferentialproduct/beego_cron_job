package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type Group struct {
	Id int
	GroupName string
	Description string
	CreateId int
	CreateTime int64
	UpdateId int
	UpdateTime int64
	Status int
}

func (group *Group) TableName() string {
	return TableName("task_group")
}



func GroupGetList(page, pageSize int, filters ...interface{}) ([]*Group, int64) {
	query := orm.NewOrm().QueryTable(TableName("task_group"))
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	total, _ := query.Count()
	list := make([]*Group, 0)
	offset := (page - 1) * pageSize
	_, err := query.OrderBy("-id").Limit(pageSize, offset).All(&list)
	if err != nil {
		fmt.Println(err)
	}
	return list, total
}

func (group *Group) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(group, fields...); err != nil {
		return err
	}
	return nil

}

func TaskGroupGetById(id int) (*Group, error) {
	obj := &Group{
		Id: id,
	}
	err := orm.NewOrm().Read(obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}
