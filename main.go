package main

import (
	"beego_cron_job/models"
	_ "beego_cron_job/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"time"
)

func init() {
	//当前时间
	models.Init(time.Now().Unix())
	//beego带的日志包，打印全部日志输出。
	logs.SetLevel(logs.LevelDebug)
}

func main() {
	beego.Run()
}

