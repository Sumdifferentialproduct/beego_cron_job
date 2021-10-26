package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type Ban struct {
	Id int
	Code string
	CreateTime int64
	UpdateTime int64
	Status int
}

func (ban *Ban) TableName() string {
	return TableName("task_ban")
}

func BanGetList(page, pageSize int, filters ...interface{}) ([]*Ban, int64) {
	query := orm.NewOrm().QueryTable(TableName("task_ban"))
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			query = query.Filter(filters[k].(string), filters[k+1])
		}
	}
	total, _ := query.Count()
	list := make([]*Ban, 0)
	offset := (page - 1) * pageSize
	_, err := query.OrderBy("-id").Limit(pageSize, offset).All(&list)
	if err != nil {
		fmt.Println(err)
	}
	return list, total
}
