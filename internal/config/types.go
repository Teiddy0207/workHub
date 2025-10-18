package config

import "workHub/pkg/jwt"

type Config struct {
	AppName          string     `yaml:"app_name" mapstructure:"app_name"`
	WebHttp          *WebHttp   `yaml:"webHttp" mapstructure:"webHttp"`
	WebSocket        *WebSocket `yaml:"webSocket" mapstructure:"webSocket"`
	Mysql            *Mysql     `yaml:"mysql" mapstructure:"mysql"`
	NodeEnv          string     `yaml:"node_env" mapstructure:"node_env"`
	Logger           Logger     `yaml:"logger" mapstructure:"logger"`
	Postgres         *Postgres  `yaml:"postgres" mapstructure:"postgres"`
	Redis            *Redis     `yaml:"redis" mapstructure:"redis"`
	Jwt              jwt.Config `yaml:"jwt" mapstructure:"jwt"`
	SessionSecretKey string     `yaml:"session_secret_key" mapstructure:"session_secret_key"`
	CookieSecretKey  string     `yaml:"cookie_secret_key" mapstructure:"cookie_secret_key"`
	CSRFSecretKey    string     `yaml:"csrf_secret_key" mapstructure:"csrf_secret_key"`
	Email            Email      `yaml:"email" mapstructure:"email"`
	AWS              AWS        `yaml:"aws" mapstructure:"aws"`
}

type Mysql struct {
	Host            string `yaml:"host" mapstructure:"host"`
	Port            int    `yaml:"port" mapstructure:"port"`
	Username        string `yaml:"username" mapstructure:"username"`
	Password        string `yaml:"password" mapstructure:"password"`
	Db_name         string `yaml:"db_name" mapstructure:"db_name"`
	MaxIdleConns    int    `yaml:"max_idle_conns" mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns" mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime" mapstructure:"conn_max_lifetime"` // seconds

}
type Logger struct {
	Encoding          string `yaml:"encoding" mapstructure:"encoding"`
	Level             string `yaml:"level" mapstructure:"level"`
	ZapType           string `yaml:"zap_type" mapstructure:"zap_type"`
	DisableCaller     bool   `yaml:"disable_caller" mapstructure:"disable_caller"`
	DisableStacktrace bool   `yaml:"disable_stacktrace" mapstructure:"disable_stacktrace"`
	LogFile           bool   `yaml:"log_file" mapstructure:"log_file"`
	Payload           bool   `yaml:"payload" mapstructure:"payload"`
}
type Postgres struct {
	Host            string `yaml:"host" mapstructure:"host"`
	Port            int    `yaml:"port" mapstructure:"port"`
	Username        string `yaml:"username" mapstructure:"username"`
	Password        string `yaml:"password" mapstructure:"password"`
	Db_name         string `yaml:"db_name" mapstructure:"db_name"`
	MaxIdleConns    int    `yaml:"max_idle_conns" mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns" mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime" mapstructure:"conn_max_lifetime"` // seconds
}

type Redis struct {
	Host         string `yaml:"host" mapstructure:"host"`
	Port         int    `yaml:"port" mapstructure:"port"`
	Password     string `yaml:"password" mapstructure:"password"`
	DB           int    `yaml:"db" mapstructure:"db"`
	PoolSize     int    `yaml:"pool_size" mapstructure:"pool_size"`
	ReadTimeout  int    `yaml:"read_timeout" mapstructure:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout" mapstructure:"write_timeout"`
	DialTimeout  int    `yaml:"dial_timeout" mapstructure:"dial_timeout"`
	Timeout      int    `yaml:"timeout" mapstructure:"timeout"`
}
type WebHttp struct {
	HttpHost     string `yaml:"http_host" mapstructure:"http_host"`
	HttpAddress  int    `yaml:"http_address" mapstructure:"http_address"`
}
type WebSocket struct {
	WsHost     string `yaml:"ws_host" mapstructure:"ws_host"`
	WsAddress  int    `yaml:"ws_address" mapstructure:"ws_address"`
	WsAllow    bool   `yaml:"ws_allow" mapstructure:"ws_allow"`
}

type Email struct {
	SMTPHost string `yaml:"smtp_host" mapstructure:"smtp_host"`
	SMTPPort int    `yaml:"smtp_port" mapstructure:"smtp_port"`
	Username string `yaml:"username" mapstructure:"username"`
	Password string `yaml:"password" mapstructure:"password"`
	FromName string `yaml:"from_name" mapstructure:"from_name"`
}

type AWS struct {
	Region          string `yaml:"region" mapstructure:"region"`
	Bucket          string `yaml:"bucket" mapstructure:"bucket"`
	AccessKeyID     string `yaml:"access_key_id" mapstructure:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key" mapstructure:"secret_access_key"`
	PublicEndpoint  string `yaml:"public_endpoint" mapstructure:"public_endpoint"`
	ApiEndpoint     string `yaml:"api_endpoint" mapstructure:"api_endpoint"`
	Endpoint        string `yaml:"endpoint" mapstructure:"endpoint"`
}
