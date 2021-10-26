package utils

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"net/smtp"
	"strings"
	"time"
)

// 根据app.conf定义配置的结构体
//email.host = smtp.126.com
//email.port = 25
//email.from = golang123456@126.com
//email.username = golang123456@126.com
//email.password = JZQDIGQFAVMDVMXM


type  EmailConfig struct {
	Host   	string
	Port   	string
	From    string
	User 	string
	Pwd 	string
}

//封装邮件结构体

type  Email struct {
	Subject  	string
	Body  		string
	//传给多个人用，约定用封号分割    xxx1;xxx2;xxx3
	To			string
	Format 		string
	Config 		*EmailConfig
}
//邮件配置协程全局变量
//以及把邮件管道放在全局，方便下面各个函数调用
var (
	config    *EmailConfig
	emailChan	chan *Email

)


func init() {
	host := beego.AppConfig.String("email.host")
	port := beego.AppConfig.String("email.port")
	from := beego.AppConfig.String("email.from")
	username := beego.AppConfig.String("email.username")
	password := beego.AppConfig.String("email.password")
	poolSize, _ := beego.AppConfig.Int("email.pool")
	config = &EmailConfig{
		Host: host,
		From: from,
		Port: port,
		User: username,
		Pwd:  password,
	}
	//创建邮件管道，poolSize为缓冲区大小
	emailChan = make(chan *Email,poolSize)
	//开子协程，循环监控邮件管道
	//若管道读到Email数据，则发送邮件
	go func() {
		for{
			select {
			//管道取出要发的邮件
			case  email ,ok := <-emailChan:
				if !ok{
					return
				}
				//设置发邮件的用户名  密码（授权码），  服务地址
				auth	:=	smtp.PlainAuth(
					"",
					email.Config.User,
					email.Config.Pwd,
					email.Config.Host)
				//处理接收者信息
				sendTo := strings.Split(email.To, ";")
				//设置邮件发送格式
				var  contentType  string
				if  email.Format==""{
					contentType = "Content-Type: text/plain;charset=UTF-8"
				}else{
					contentType = "Content-Type: "+email.Format + ";charset=UTF-8"
				}
				//拼接发送内容
				msg :=  []byte("To: " + email.To + "\r\nFrom: " + email.Config.User + ">\r\nSubject: " +
					email.Subject + "\r\n" + contentType + "\r\n\r\n" + email.Body)
				//严谨判断。端口是否为25
				var  err  error
				if  email.Config.Port == "25" {
					//发送邮件
					err = smtp.SendMail(email.Config.Host+":"+email.Config.Port, auth, email.Config.User, sendTo, msg)
				}else{
					err = errors.New("邮件端口不正确！邮件发送失败！")
				}
				if  err !=nil{
					logs.Error("SendEmail: ",err.Error())
				}
			}
		}
	}()
}

//提供外界调用的，发邮件的方法
//参数1 ：邮件发送给谁
//参数2 ：邮件主题
//参数3 ：邮件解析格式
//返回值： 返回是否发送成功
func  SendToChan(to,subject,body,mailtype string)bool{
	//1.组装邮件对象
	email := &Email{
		Config :config,
		Body :body,
		Subject: subject,
		Format: mailtype,
		To :to,
	}
	//2.将邮件发送到管道
	select {
	case  emailChan <- email:
		return  true
		case <-time.After(3 * time.Second):
			return false
	}

}




















