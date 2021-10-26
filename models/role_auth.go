package models

import (
	"bytes"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
)

type RoleAuth struct {
	AuthId int `orm:"pk"`
	RoleId int64
}

func (roleauth *RoleAuth) TableName() string {
	return TableName("uc_role_auth")
}

func RoleAuthGetByIds(RoleIds string) (Authids string, err error) {
	query := orm.NewOrm().QueryTable(TableName("uc_role_auth"))
	ids := strings.Split(RoleIds, ",")
	list := make([]*RoleAuth, 0)
	_, err = query.Filter("role_id__in", ids).All(&list, "auth_id")
	if err != nil {
		return "", err
	}
	b := bytes.Buffer{}
	for _, v := range list {
		if v.AuthId != 0 && v.AuthId != 1 {
			b.WriteString(strconv.Itoa(v.AuthId))
			b.WriteString(",")
		}
	}
	Authids = strings.TrimRight(b.String(), ",")
	return Authids, nil
}

func RoleAuthBatchAdd(ras *[]RoleAuth) (int64, error) {
	return orm.NewOrm().InsertMulti(len(*ras), ras)
}

func RoleAuthGetById(id int) ([]*RoleAuth, error) {
	list := make([]*RoleAuth, 0)
	_, err := orm.NewOrm().QueryTable(TableName("uc_role_auth")).Filter("role_id", id).All(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func RoleAuthDelete(id int) (int64, error) {
	return orm.NewOrm().QueryTable(TableName("uc_role_auth")).Filter("role_id", id).Delete()
}
