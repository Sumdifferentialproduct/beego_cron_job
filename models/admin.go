package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type Admin struct {
	Id int
	LoginName string
	RealName string
	Password string
	RoleIds string
	Phone string
	Email string
	Salt string
	LastLogin int64
	LastIp string
	Status int
	CreateId int
	UpdateId int
	CreateTime int64
	UpdateTime int64
}

func (admin *Admin) TableName() string {
	return TableName("uc_admin")
}




func AdminGetList(page, pageSize int, filters ...interface{}) ([]*Admin, int64) {
	query := orm.NewOrm().QueryTable(TableName("uc_admin"))
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	total, _ := query.Count()
	list := make([]*Admin, 0)
	offset := (page - 1) * pageSize
	_, err := query.OrderBy("-id").
		Limit(pageSize, offset).All(&list)
	if err != nil {
		fmt.Println(err)
	}
	return list, total
}

func AdminGetByName(loginName string) (*Admin, error) {
	admin := new(Admin)
	err := orm.NewOrm().QueryTable(TableName("uc_admin")).
		Filter("login_name", loginName).One(admin)
	if err != nil {
		return nil, err
	}
	return admin, nil
}

func AdminGetById(id int) (*Admin, error) {
	admin := new(Admin)
	err := orm.NewOrm().QueryTable(TableName("uc_admin")).
		Filter("id", id).One(admin)
	if err != nil {
		return nil, err
	}
	return admin, nil
}
