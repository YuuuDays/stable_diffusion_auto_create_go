package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// SDConfig はSD生成パラメータを管理
type SDConfig struct {
	APIURL         string
	NegativePrompt string
	Steps          int
	CfgScale       float64
	Width          int
	Height         int
	SamplerName    string
	Seed           int64
}

// LoadSDConfig は.envファイルからAPI_URLを読み込み、sd_config.txtからSD設定を読み込む
func LoadSDConfig() (*SDConfig, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	apiURL := getEnv("API_URL", "http://127.0.0.1:7860")

	// sd_config.txtから設定を読み込み
	sdConfig, err := loadSDConfigFromFile()
	if err != nil {
		return nil, err
	}

	cfg := &SDConfig{
		APIURL:         apiURL,
		NegativePrompt: sdConfig["NEGATIVE_PROMPT"],
		Steps:          parseInt(sdConfig["STEPS"], 20),
		CfgScale:       parseFloat(sdConfig["CFG_SCALE"], 7.0),
		Width:          parseInt(sdConfig["WIDTH"], 768),
		Height:         parseInt(sdConfig["HEIGHT"], 1024),
		SamplerName:    sdConfig["SAMPLER_NAME"],
		Seed:           parseInt64(sdConfig["SEED"], -1),
	}

	return cfg, nil
}

// SaveSDConfig はSD設定をsd_config.txtファイルに保存
func SaveSDConfig(cfg *SDConfig) error {
	content := fmt.Sprintf(`NEGATIVE_PROMPT=%s
STEPS=%d
CFG_SCALE=%.1f
WIDTH=%d
HEIGHT=%d
SAMPLER_NAME=%s
SEED=%d
`, cfg.NegativePrompt, cfg.Steps, cfg.CfgScale, cfg.Width, cfg.Height, cfg.SamplerName, cfg.Seed)

	return os.WriteFile("config/sd_config.txt", []byte(content), 0644)
}

// loadSDConfigFromFile はsd_config.txtから設定を読み込み、マップとして返す
func loadSDConfigFromFile() (map[string]string, error) {
	data, err := os.ReadFile("config/sd_config.txt")
	if err != nil {
		return nil, err
	}

	config := make(map[string]string)
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			config[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return config, nil
}

func parseInt(s string, defaultValue int) int {
	if value, err := strconv.Atoi(s); err == nil {
		return value
	}
	return defaultValue
}

func parseFloat(s string, defaultValue float64) float64 {
	if value, err := strconv.ParseFloat(s, 64); err == nil {
		return value
	}
	return defaultValue
}

func parseInt64(s string, defaultValue int64) int64 {
	if value, err := strconv.ParseInt(s, 10, 64); err == nil {
		return value
	}
	return defaultValue
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}
