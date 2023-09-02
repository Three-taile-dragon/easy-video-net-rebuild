package email

import (
	"dragonsss.cn/evn_user/config"
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
	"strconv"
)

func SendMail(mailTo []string, subject string, body string) error {
	// 设置邮箱主体
	mailConn := map[string]string{
		"user": config.C.Email.User,
		"pass": config.C.Email.Pass,
		"host": config.C.Email.Host,
		"port": config.C.Email.Port,
	}

	port, _ := strconv.Atoi(mailConn["port"])
	m := gomail.NewMessage()
	// 添加别名
	m.SetHeader("From", m.FormatAddress(mailConn["user"], "发验证码滴!!!"))
	// 发送给用户(可以多个)
	m.SetHeader("To", mailTo...)
	// 设置邮件主题
	m.SetHeader("Subject", subject)
	// 设置邮件正文
	m.SetBody("text/html", body)
	// 设置邮件正文
	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])
	err := d.DialAndSend(m)
	if err != nil {
		zap.L().Error(fmt.Sprintf("发送至%s邮箱验证码失败,内容 ： %s 错误原因: ", mailTo, body), zap.Error(err))
	}
	return err
}
