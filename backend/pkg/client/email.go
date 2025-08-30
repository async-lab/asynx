package client

import (
	"html/template"
	"strings"

	"asynclab.club/asynx/backend/pkg/config"
	"gopkg.in/gomail.v2"
)

type EmailClient struct {
	cfg  *config.ConfigEmail
	tmpl []byte
}

func NewEmailClient(cfg *config.ConfigEmail, tmpl []byte) (*EmailClient, error) {
	return &EmailClient{
		cfg:  cfg,
		tmpl: tmpl,
	}, nil
}

func (c *EmailClient) send(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", c.cfg.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetHeader("Reply-To", c.cfg.ReplyTo)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(c.cfg.Host, c.cfg.Port, c.cfg.Username, c.cfg.Password)

	return d.DialAndSend(m)
}

// 加载嵌入的HTML模板
func (c *EmailClient) loadTemplate(name string) (*template.Template, error) {
	// content, err := templateFS.ReadFile("templates/" + name)
	// if err != nil {
	// 	return nil, err
	// }
	return template.New(name).Parse(string(c.tmpl))
}

// 发送邮件处理器
func (c *EmailClient) SendMail(to string, subject string, body any) error {
	// 加载邮件模板
	tmpl, err := c.loadTemplate("email.html")
	if err != nil {
		return err
	}

	var htmlBody strings.Builder
	err = tmpl.Execute(&htmlBody, body)
	if err != nil {
		return err
	}

	err = c.send(to, subject, htmlBody.String())
	if err != nil {
		return err
	}

	return nil
}
