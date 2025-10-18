package logger

type Config struct {
	Mode              string `yaml:"mode" mapstructure:"mode"`
	DisableCaller     bool   `yaml:"disable_caller" mapstructure:"disable_caller"`
	DisableStacktrace bool   `yaml:"disable_stacktrace" mapstructure:"disable_stacktrace"`
	LogFile           bool   `yaml:"log_file" mapstructure:"log_file"`
	Payload           bool   `yaml:"payload" mapstructure:"payload"`
	Encoding          string `yaml:"encoding" mapstructure:"encoding"`
	Level             string `yaml:"level" mapstructure:"level"`
	ZapType           string `yaml:"zap_type" mapstructure:"zap_type"`
}
