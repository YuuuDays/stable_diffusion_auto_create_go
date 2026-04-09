package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// SDConfig はSD生成パラメータを管理
type SDConfig struct {
	NegativePrompt string
	Steps          int
	CfgScale       float64
	Width          int
	Height         int
	SamplerName    string
	Seed           int64
}

// LoadSDConfig は.envファイルからSD設定を読み込む
func LoadSDConfig() (*SDConfig, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cfg := &SDConfig{
		NegativePrompt: getEnv("NEGATIVE_PROMPT", "score_6, score_5, score_4"),
		Steps:          getEnvInt("STEPS", 20),
		CfgScale:       getEnvFloat("CFG_SCALE", 7.0),
		Width:          getEnvInt("WIDTH", 768),
		Height:         getEnvInt("HEIGHT", 1024),
		SamplerName:    getEnv("SAMPLER_NAME", "DPM++ 2M Karras"),
		Seed:           getEnvInt64("SEED", -1),
	}

	return cfg, nil
}

// SaveSDConfig はSD設定を.envファイルに保存
func SaveSDConfig(cfg *SDConfig) error {
	envContent := fmt.Sprintf(`NEGATIVE_PROMPT=%s
STEPS=%d
CFG_SCALE=%.1f
WIDTH=%d
HEIGHT=%d
SAMPLER_NAME=%s
SEED=%d
`, cfg.NegativePrompt, cfg.Steps, cfg.CfgScale, cfg.Width, cfg.Height, cfg.SamplerName, cfg.Seed)

	return os.WriteFile(".env", []byte(envContent), 0644)
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
