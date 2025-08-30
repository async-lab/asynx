package config

// 邮件配置
type ConfigEmail struct {
	Host     string `env:"SMTP_HOST,required"`
	Port     int    `env:"SMTP_PORT,required"`
	Username string `env:"SMTP_USERNAME,required"`
	Password string `env:"SMTP_PASSWORD,required"`
	From     string `env:"SMTP_FROM,required"`
	ReplyTo  string `env:"SMTP_REPLY_TO"`
}
