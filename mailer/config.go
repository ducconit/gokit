package mailer

// Config represents the configuration for the mailer.
type Config struct {
	SMTP SMTPConfig `json:"smtp" mapstructure:"smtp"`
}

// SMTPConfig represents the SMTP configuration.
type SMTPConfig struct {
	Host     string `json:"host" mapstructure:"host"`
	Port     int    `json:"port" mapstructure:"port"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
	From     string `json:"from" mapstructure:"from"`
	TLS      bool   `json:"tls" mapstructure:"tls"`
}
