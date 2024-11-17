package config

import (
	"os"
	"time"
)

type Config struct {
	Database struct {
		User     string
		Password string
		Host     string
		Port     string
		Dbname   string
	}
	APP         string
	Environment string
	LogLevel    string
	RPCPort     string
	Context     struct {
		Timeout string
	}
	Email struct {
		From     string
		Password string
		SmtHost  string
		SmtPort  string
	}
	UserUrl string
	Token   struct {
		Secret     string
		AccessTTL  time.Duration
		RefreshTTL time.Duration
	}
	RedisURL string
	DB       struct {
		Host           string
		Port           string
		Name           string
		User           string
		Password       string
		CollectionName string
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

	config.Database.User = getEnv("User", "postgres")
	config.Database.Password = getEnv("Password", "2005")
	config.Database.Host = getEnv("Host", "localhost")
	config.Database.Port = getEnv("Port", "5432")
	config.Database.Dbname = getEnv("Dbname", "udevs")
	config.DB.CollectionName = "notifications"
	config.DB.Host = getEnv("DB_HOST", "localhost")
	config.DB.Port = getEnv("DB_PORT", ":27017")
	config.DB.User = getEnv("DB_USER", "postgres")
	config.DB.Password = getEnv("DB_PASSWORD", "postgres")
	config.DB.Name = getEnv("DB_NAME", "notification")
	config.APP = getEnv("APP", "app")
	config.Environment = getEnv("ENVIRONMENT", "develop")
	config.LogLevel = getEnv("LOG_LEVEL", "local")
	config.RPCPort = getEnv("RPC_PORT", "localhost:9006")
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
	config.RedisURL = getEnv("REDIS_URL", "localhost:6379")
	config.RedisURL = getEnv("REDIS_URL", "localhost:6379")

	config.MessageBrokerUses.TopicIncome = getEnv("MESSAGE_BROKER_USE_TOKEN", "incomeexpenses17")
	config.MessageBrokerUses.IncomeCreate = getEnv("MESSAGE_BROKER_USE_KEY", "incomecreate")
	config.MessageBrokerUses.ExpensesCreate = getEnv("MESSAGE_BROKER_USE_KEY", "expensescreate")

	config.ReportUcl = getEnv("REPORT_UCL", "report-service:8000")
	config.BudgetServiceUrl = getEnv("BUDGET_SERVICE_URL", "byudjet-service:8888")
	config.IncomeServiceUrl = getEnv("INCOME_SERVICE_URL", "income-expenses_container:8080")

	config.Email.SmtHost = getEnv("SMT_HOST", "smtp.gmail.com")
	config.Email.SmtPort = getEnv("SMTP_PORT", "587")
	config.Email.From = getEnv("EMAIL_FROM", "diyordev3@gmail.com")
	config.Email.Password = getEnv("EMAIL_PASSWORD", "ueus bord hbep ttam")
	return &config
}

func getEnv(key string, defaultVaule string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultVaule
}
