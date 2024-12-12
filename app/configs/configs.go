package configs

import (
	"sync"

	"log"

	"github.com/spf13/viper"
)

var (
	onceConf sync.Once
	Conf     Config
)

type Config struct {
	App         App      `mapstructure:",squash"`
	Database    Database `mapstructure:",squash"`
	SBA         SBA      `mapstructure:",squash"`
	Client      Client   `mapstructure:",squash"`
	PrefixRefID string   `mapstructure:"PREFIX_REF_ID"`
	Logger      Logger   `mapstructure:",squash"`
}

type SBA struct {
	PID                string `mapstructure:"SBA_PID"`
	EncryptionKey      string `mapstructure:"SBA_KEY"`
	Hashword           string `mapstructure:"SBA_HASHWORD"`
	MACAddress         string `mapstructure:"SBA_MACADDRESS"`
	MessageKey         string `mapstructure:"SBA_MSG_KEY"`
	OperationTimeStart string `mapstructure:"SBA_OPERATION_TIME_START"`
	OperationTimeEnd   string `mapstructure:"SBA_OPERATION_TIME_END"`
}

type App struct {
	Env  string `mapstructure:"APP_ENV"`
	Port int    `mapstructure:"APP_PORT"`
	Name string `mapstructure:"APP_NAME"`
}

type Client struct {
	CreditUrl     string `mapstructure:"CREDIT_URL"`
	AccountUrl    string `mapstructure:"ACCOUNT_URL"`
	AccountApiKey string `mapstructure:"ACCOUNT_API_KEY"`
}

type Database struct {
	INVXCreditLimitDatabase  INVXCreditLimitDatabase `mapstructure:",squash"`
	INVXCenterDatabase       INVXCenterDatabase      `mapstructure:",squash"`
	RedisDatabase            RedisDatabase           `mapstructure:",squash"`
	SSLMode                  string                  `mapstructure:"DB_SSL_MODE"`
	DBMaxConnections         int                     `mapstructure:"DB_MAX_CONNECTIONS"`
	DBMaxIdleConnections     int                     `mapstructure:"DB_MAX_IDLE_CONNECTIONS"`
	DBMaxLifetimeConnections int                     `mapstructure:"DB_MAX_LIFETIME_CONNECTIONS"`
}

type INVXCreditLimitDatabase struct {
	Host     string `mapstructure:"DB_INVX_CREDIT_HOST"`
	Port     int    `mapstructure:"DB_INVX_CREDIT_PORT"`
	User     string `mapstructure:"DB_INVX_CREDIT_USER"`
	Name     string `mapstructure:"DB_INVX_CREDIT_NAME"`
	Password string `mapstructure:"DB_INVX_CREDIT_PASSWORD"`
}

type INVXCenterDatabase struct {
	Host     string `mapstructure:"DB_INVX_CENTER_HOST"`
	Port     int    `mapstructure:"DB_INVX_CENTER_PORT"`
	User     string `mapstructure:"DB_INVX_CENTER_USER"`
	Name     string `mapstructure:"DB_INVX_CENTER_NAME"`
	Password string `mapstructure:"DB_INVX_CENTER_PASSWORD"`
}

type RedisDatabase struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     string `mapstructure:"REDIS_PORT"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	SSLMode  bool   `mapstructure:"REDIS_SSL_MODE"`
	DB       int    `mapstructure:"REDIS_DB"`
}

type Logger struct {
	Feed  bool   `mapstructure:"LOG_FEED"`
	Level string `mapstructure:"LOG_LEVEL"`
	Url   string `mapstructure:"LOG_URL"`
}

func InitConfig() Config {
	onceConf.Do(func() {
		viper.SetConfigFile(".env")
		viper.SetConfigType("env")
		viper.AutomaticEnv()
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file: %s", err)
		}
		err := viper.Unmarshal(&Conf)
		if err != nil {
			log.Fatalf("unable to decode into struct , %v", err)
		}
	})
	return Conf
}
