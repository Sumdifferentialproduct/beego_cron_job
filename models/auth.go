package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type Auth struct {
	Id int
	Pid int
	AuthName string
	AuthUrl string
	Sort int
	Icon string
	IsShow int
	UserId int
	CreateId int
	UpdateId int
	Status int
	CreateTime int64
	UpdateTime int64
}

func (auth *Auth) TableName() string {
	return TableName("uc_auth")
}

func AuthGetList(page, pageSize int, filters ...interface{}) ([]*Auth, int64) {
	query := orm.NewOrm().QueryTable(TableName("uc_auth"))
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	total, _ := query.Count()
	list := make([]*Auth, 0)
	offset := (page - 1) * pageSize
	_, err := query.OrderBy("pid", "sort").
		Limit(pageSize, offset).All(&list)
	if err != nil {
		fmt.Println(err)
	}
	return list, total
}

func AuthAdd(auth *Auth) (int64, error) {
	return orm.NewOrm().Insert(auth)
}

func AuthGetById(id int) (*Auth, error) {
	auth := &Auth{
		Id: id,
	}
	err := orm.NewOrm().Read(auth)
	if err != nil {
		return nil, err
	}
	return auth, nil
}

func (auth *Auth) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(auth, fields...); err != nil {
		return err
	}
	return nil
}
