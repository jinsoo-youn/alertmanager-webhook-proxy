package config

import (
	"flag"
	"os"
	"strings"
)

type Config struct {
	ListenAddress string
	LogLevel      string
	Stage         string
	Region        string
	TemplateDir   string
	WardConfig    *WardConfig
	DoorayConfig  *DoorayConfig
}

type WardConfig struct {
	Enable   bool
	EventURL string
	Env      string
	Region   string
	Actor    string
}

type DoorayConfig struct {
	Enable     bool
	WebhookURL string
	Stage      string
	Region     string
}

func getEnvOrDefault(key string, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func getBoolEnvOrDefault(key string, fallback bool) bool {
	if val, ok := os.LookupEnv(key); ok {
		lowerVal := strings.ToLower(val)
		if lowerVal == "true" || lowerVal == "1" || lowerVal == "yes" {
			return true
		}
		return false
	}
	return fallback
}

func LoadConfig() (*Config, error) {
	listen := flag.String("listen", getEnvOrDefault("LISTEN_ADDRESS", "0.0.0.0:8080"), "HTTP listen address")
	logLevel := flag.String("log-level", getEnvOrDefault("LOG_LEVEL", "info"), "Log level: debug | info | warn | error")
	stage := flag.String("stage", getEnvOrDefault("STAGE", "beta"), "App stage: alpha | beta | public | gov")
	region := flag.String("region", getEnvOrDefault("REGION", "kr2"), "Deployment region: kr1 | kr2 | kr3 | p01")
	TemplateDir := flag.String("template-dir", getEnvOrDefault("TEMPLATE_DIR", "./templates"), "Template directory")

	wardEnable := flag.Bool("ward-enable", getBoolEnvOrDefault("WARD_ENABLE", false), "WARD enable flag")
	wardURL := flag.String("ward-url", getEnvOrDefault("WARD_EVENT_URL", "https://ward_url/event"), "WARD Event URL")
	wardActor := flag.String("ward-actor", getEnvOrDefault("WARD_ACTOR", "buoy"), "WARD Actor")

	doorayEnable := flag.Bool("dooray-enable", getBoolEnvOrDefault("DOORAY_ENABLE", false), "Dooray enable flag")
	doorayURL := flag.String("dooray-webhook", getEnvOrDefault("DOORAY_WEBHOOK_URL", "https://nhnent.dooray.com/services/xxxx/yyyy"), "Dooray Webhook URL")

	flag.Parse()

	return &Config{
		ListenAddress: *listen,
		LogLevel:      *logLevel,
		Stage:         *stage,
		Region:        *region,
		TemplateDir:   *TemplateDir,
		WardConfig: &WardConfig{
			EventURL: *wardURL,
			Env:      *stage,
			Region:   *region,
			Actor:    *wardActor,
			Enable:   *wardEnable,
		},
		DoorayConfig: &DoorayConfig{
			WebhookURL: *doorayURL,
			Stage:      *stage,
			Region:     *region,
			Enable:     *doorayEnable,
		},
	}, nil
}
