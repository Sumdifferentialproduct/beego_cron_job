package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"net/url"
)

//用于首页服务运行时间的计算
var  StartTime  int64


//db.host = 127.0.0.1
//db.user = root
//db.password = 123
//db.port = 3306
//db.name = cron_job
//db.prefix = pp_
//db.timezone = Asia/Shanghai
func Init(st int64) {
	StartTime =st
	//连接数据库做初始化
	dbhost := beego.AppConfig.String("db.host")
	dbuser := beego.AppConfig.String("db.user")
	dbpassword := beego.AppConfig.String("db.password")
	dbname := beego.AppConfig.String("db.name")
	dbport := beego.AppConfig.String("db.port")
	timezone := beego.AppConfig.String("db.timezone")
	if dbport == ""{
		dbport = "3306"
	}
	ds := dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname+"?charset=utf8"
	//Asia/Shanghai  编码后Asia%2FShanghai
	if  timezone !=""{
		//URL编码处理处理后“Asia%2FShanghai”，可以被识别
		ds = ds+"&loc="+url.QueryEscape(timezone)
		//ds = ds + "&loc=" + timezone
		fmt.Println("dssssss",ds)

	}
	//连接数据库
	err := orm.RegisterDataBase("default", "mysql", ds, 50, 30)
	if err != nil{
		fmt.Println(err)
		return
	}
	//注册model
	orm.RegisterModel(
		new(Admin),
		new(Auth),
		new(Ban),
		new(Role),
		new(RoleAuth),
		new(TaskServer),
		new(ServerGroup),
		new(Task),
		new(Group),
		new(TaskLog),
		)
	//调试阶段 ,可以设置执行sql语句
	orm.Debug  =  true

}
func TableName(name string) string {
return beego.AppConfig.String("db.prefix") + name
}

