package config

import (
	"os"
	"time"
)

type Config struct {
	APP         string
	Environment string
	LogLevel    string
	RPCPort     string
	Context     struct {
		Timeout string
	}
	UserUrl string
	Token   struct {
		Secret     string
		AccessTTL  time.Duration
		RefreshTTL time.Duration
	}
	RedisURL string
	DB       struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
	}
	NotificationUrl string

	MessageBrokerUses struct {
		URL          string
		Topic        string
		TopicBooking string
		TopicIncome  string
		Keys         struct {
			Create         []byte
			Update         []byte
			Delete         []byte
			UpdateEmail    []byte
			UpdatePassword []byte
		}
		IncomeCreate   string
		ExpensesCreate string
	}
	BudgetServiceUrl string
	IncomeServiceUrl string
	ReportUcl        string
}

func Token() string {
	c := Config{}
	c.Token.Secret = getEnv("TOKEN_SECRET", "token_secret")
	return c.Token.Secret
}

func New() *Config {
	var config Config

	config.APP = getEnv("APP", "app")
	config.Environment = getEnv("ENVIRONMENT", "develop")
	config.LogLevel = getEnv("LOG_LEVEL", "local")
	config.RPCPort = getEnv("RPC_PORT", "api_gateway:9006")
	config.Context.Timeout = getEnv("CONTEXT_TIMEOUT", "30s")

	config.Token.Secret = getEnv("TOKEN_SECRET", "D1YORTOP4EEK")
	accessTTl, err := time.ParseDuration(getEnv("TOKEN_ACCESS_TTL", "1h"))
	if err != nil {
		return nil
	}
	refreshTTL, err := time.ParseDuration(getEnv("TOKEN_REFRESH_TTL", "24h"))
	if err != nil {
		return nil
	}
	config.Token.AccessTTL = accessTTl
	config.Token.RefreshTTL = refreshTTL

	config.UserUrl = getEnv("User_URL", "user_service:9000")
	config.NotificationUrl = getEnv("Notification_URL", "notification_service:9001")
	config.RedisURL = getEnv("REDIS_URL", "redis:6379")

	config.MessageBrokerUses.URL = getEnv("KAFKA_URL", "broker:29092")
	config.MessageBrokerUses.Topic = getEnv("MESSAGE_BROKER_USE_TOKEN", "USER_SERVICE")
	config.MessageBrokerUses.Keys.Create = []byte(getEnv("MESSAGE_BROKER_USE_KEY", "CREATE"))
	config.MessageBrokerUses.Keys.Delete = []byte(getEnv("MESSAGE_BROKER_USE_KEY", "DELETE"))
	config.MessageBrokerUses.Keys.Update = []byte(getEnv("MESSAGE_BROKER_USE_KEY", "UPDATE"))
	config.MessageBrokerUses.Keys.UpdateEmail = []byte(getEnv("MESSAGE_BROKER_USE_KEY", "UPDATE_EMAIL"))
	config.MessageBrokerUses.Keys.UpdatePassword = []byte(getEnv("MESSAGE_BROKER_USE_KEY", "UPDATE_PASSWORD"))
	config.MessageBrokerUses.TopicBooking = getEnv("MESSAGE_BROKER_USE_TOKEN", "USER_SERVICE")

	config.MessageBrokerUses.TopicIncome = getEnv("MESSAGE_BROKER_USE_TOKEN", "incomeexpenses17")
	config.MessageBrokerUses.IncomeCreate = getEnv("MESSAGE_BROKER_USE_KEY", "incomecreate")
	config.MessageBrokerUses.ExpensesCreate = getEnv("MESSAGE_BROKER_USE_KEY", "expensescreate")

	config.ReportUcl = getEnv("REPORT_UCL", "report-service:8000")
	config.BudgetServiceUrl = getEnv("BUDGET_SERVICE_URL", "byudjet-service:8888")
	config.IncomeServiceUrl = getEnv("INCOME_SERVICE_URL", "income-expenses_container:8080")

	return &config
}

func getEnv(key string, defaultVaule string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultVaule
}
