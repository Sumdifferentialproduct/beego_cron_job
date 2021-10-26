package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type Role struct {
	Id int
	RoleName string
	Detail string
	ServerGroupIds string
	TaskGroupIds string
	CreateId int
	UpdateId int
	Status int
	CreateTime int64
	UpdateTime int64
}

func (role *Role) TableName() string {
	return TableName("uc_role")
}

func RoleGetList(page, pageSize int, filters ...interface{}) ([]*Role, int64) {
	query := orm.NewOrm().QueryTable(TableName("uc_role"))
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	total, _ := query.Count()
	list := make([]*Role, 0)
	offset := (page - 1) * pageSize
	_, err := query.OrderBy("-id").Limit(pageSize, offset).All(&list)
	if err != nil {
		fmt.Println(err)
	}
	return list, total
}

func RoleAdd(role *Role) (int64, error) {
	id, err := orm.NewOrm().Insert(role)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func RoleGetById(id int) (*Role, error) {
	r := new(Role)
	err := orm.NewOrm().QueryTable(TableName("uc_role")).Filter("id", id).One(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (role *Role) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(role, fields...); err != nil {
		return err
	}
	return nil
}
