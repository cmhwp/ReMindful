package email

import (
	"ReMindful/internal/config"
	"crypto/tls"
	"fmt"
	"net/smtp"
)

type EmailSender struct {
	config config.EmailConfig
}

func NewEmailSender(config config.EmailConfig) *EmailSender {
	return &EmailSender{config: config}
}

func (s *EmailSender) SendVerificationCode(to, code string) error {
	subject := "ReMindful 验证码"
	body := fmt.Sprintf("您的验证码是: %s\n验证码有效期为5分钟。", code)

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n"+
		"\r\n"+
		"%s", s.config.From, to, subject, body)

	// 配置 TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         s.config.Host,
	}

	// 连接到 SMTP 服务器
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to create TLS connection: %v", err)
	}
	defer conn.Close()

	// 创建 SMTP 客户端
	c, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %v", err)
	}
	defer c.Close()

	// 认证
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	if err = c.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %v", err)
	}

	// 设置发件人和收件人
	if err = c.Mail(s.config.From); err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}
	if err = c.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %v", err)
	}

	// 发送邮件内容
	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("failed to create data writer: %v", err)
	}
	defer w.Close()

	_, err = w.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}

	return nil
}
