package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                   string
	MongoURI               string
	DBName                 string
	AdminEmail             string
	AdminPass              string
	JWTSecret              string
	EmailFrom              string
	EmailPassword          string
	SMTPHost               string
	SMTPPort               string
	TwilioAccountSID       string
	TwilioAuthToken        string
	TwilioVerifyServiceSID string
	TwilioPhoneNumber      string
	Fast2SMSAPIKey         string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading from environment variables")
	}

	return &Config{
		Port:                   getEnv("PORT", "8080"),
		MongoURI:               getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DBName:                 getEnv("DB_NAME", "shopdb"),
		AdminEmail:             getEnv("ADMIN_EMAIL", "admin@shop.com"),
		AdminPass:              getEnv("ADMIN_PASS", "admin123"),
		JWTSecret:              getEnv("JWT_SECRET", "mysecretkey"),
		EmailFrom:              getEnv("EMAIL_FROM", ""),
		EmailPassword:          getEnv("EMAIL_PASSWORD", ""),
		SMTPHost:               getEnv("SMTP_HOST", ""),
		SMTPPort:               getEnv("SMTP_PORT", ""),
		TwilioAccountSID:       getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:        getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioVerifyServiceSID: getEnv("TWILIO_VERIFY_SERVICE_SID", ""),
		TwilioPhoneNumber:      getEnv("TWILIO_PHONE_NUMBER", ""),
		Fast2SMSAPIKey:         getEnv("FAST2SMS_API_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
