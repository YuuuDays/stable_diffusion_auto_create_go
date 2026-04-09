package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// SDConfig はSD生成パラメータを管理
type SDConfig struct {
	APIURL                            string
	NegativePrompt                    string
	Steps                             int
	CfgScale                          float64
	Width                             int
	Height                            int
	SamplerName                       string
	Seed                              int64
	Styles                            []string
	Subseed                           int
	SubseedStrength                   float64
	SeedResizeFromH                   int
	SeedResizeFromW                   int
	BatchSize                         int
	NIter                             int
	RestoreFaces                      bool
	Tiling                            bool
	DoNotSaveSamples                  bool
	DoNotSaveGrid                     bool
	Eta                               float64
	DenoisingStrength                 float64
	SMinUncond                        float64
	SChurn                            float64
	STmax                             float64
	STmin                             float64
	SNoise                            float64
	OverrideSettings                  map[string]interface{}
	OverrideSettingsRestoreAfterwards bool
	RefinerCheckpoint                 string
	RefinerSwitchAt                   float64
	DisableExtraNetworks              bool
	ResizeMode                        int
	ImageCfgScale                     float64
	MaskBlur                          int
	InpaintingFill                    int
	InpaintFullRes                    bool
	InpaintFullResPadding             int
	InpaintingMaskInvert              int
	InitialNoiseMultiplier            float64
	ForceTaskId                       string
	SamplerIndex                      string
	IncludeInitImages                 bool
	ScriptName                        string
	ScriptArgs                        []interface{}
	SendImages                        bool
	SaveImages                        bool
	AlwaysonScripts                   map[string]interface{}
}

// LoadSDConfig はsd_config.txtから全ての設定を読み込む
func LoadSDConfig() (*SDConfig, error) {
	// sd_config.txtから設定を読み込み
	sdConfig, err := loadSDConfigFromFile()
	if err != nil {
		return nil, err
	}

	cfg := &SDConfig{
		APIURL:                            sdConfig["API_URL"],
		NegativePrompt:                    sdConfig["NEGATIVE_PROMPT"],
		Steps:                             parseInt(sdConfig["STEPS"], 20),
		CfgScale:                          parseFloat(sdConfig["CFG_SCALE"], 7.0),
		Width:                             parseInt(sdConfig["WIDTH"], 768),
		Height:                            parseInt(sdConfig["HEIGHT"], 1024),
		SamplerName:                       sdConfig["SAMPLER_NAME"],
		Seed:                              parseInt64(sdConfig["SEED"], -1),
		Styles:                            parseStringArray(sdConfig["STYLES"]),
		Subseed:                           parseInt(sdConfig["SUBSEED"], -1),
		SubseedStrength:                   parseFloat(sdConfig["SUBSEED_STRENGTH"], 0.0),
		SeedResizeFromH:                   parseInt(sdConfig["SEED_RESIZE_FROM_H"], -1),
		SeedResizeFromW:                   parseInt(sdConfig["SEED_RESIZE_FROM_W"], -1),
		BatchSize:                         parseInt(sdConfig["BATCH_SIZE"], 1),
		NIter:                             parseInt(sdConfig["N_ITER"], 1),
		RestoreFaces:                      parseBool(sdConfig["RESTORE_FACES"], false),
		Tiling:                            parseBool(sdConfig["TILING"], false),
		DoNotSaveSamples:                  parseBool(sdConfig["DO_NOT_SAVE_SAMPLES"], false),
		DoNotSaveGrid:                     parseBool(sdConfig["DO_NOT_SAVE_GRID"], false),
		Eta:                               parseFloat(sdConfig["ETA"], 0.0),
		DenoisingStrength:                 parseFloat(sdConfig["DENOISING_STRENGTH"], 0.0),
		SMinUncond:                        parseFloat(sdConfig["S_MIN_UNCOND"], 0.0),
		SChurn:                            parseFloat(sdConfig["S_CHURN"], 0.0),
		STmax:                             parseFloat(sdConfig["S_TMAX"], 0.0),
		STmin:                             parseFloat(sdConfig["S_TMIN"], 0.0),
		SNoise:                            parseFloat(sdConfig["S_NOISE"], 1.0),
		OverrideSettings:                  parseObject(sdConfig["OVERRIDE_SETTINGS"]),
		OverrideSettingsRestoreAfterwards: parseBool(sdConfig["OVERRIDE_SETTINGS_RESTORE_AFTERWARDS"], false),
		RefinerCheckpoint:                 sdConfig["REFINER_CHECKPOINT"],
		RefinerSwitchAt:                   parseFloat(sdConfig["REFINER_SWITCH_AT"], 0.0),
		DisableExtraNetworks:              parseBool(sdConfig["DISABLE_EXTRA_NETWORKS"], false),
		ResizeMode:                        parseInt(sdConfig["RESIZE_MODE"], 0),
		ImageCfgScale:                     parseFloat(sdConfig["IMAGE_CFG_SCALE"], 0.0),
		MaskBlur:                          parseInt(sdConfig["MASK_BLUR"], 4),
		InpaintingFill:                    parseInt(sdConfig["INPAINTING_FILL"], 0),
		InpaintFullRes:                    parseBool(sdConfig["INPAINT_FULL_RES"], true),
		InpaintFullResPadding:             parseInt(sdConfig["INPAINT_FULL_RES_PADDING"], 0),
		InpaintingMaskInvert:              parseInt(sdConfig["INPAINTING_MASK_INVERT"], 0),
		InitialNoiseMultiplier:            parseFloat(sdConfig["INITIAL_NOISE_MULTIPLIER"], 1.0),
		ForceTaskId:                       sdConfig["FORCE_TASK_ID"],
		SamplerIndex:                      sdConfig["SAMPLER_INDEX"],
		IncludeInitImages:                 parseBool(sdConfig["INCLUDE_INIT_IMAGES"], false),
		ScriptName:                        sdConfig["SCRIPT_NAME"],
		ScriptArgs:                        parseArray(sdConfig["SCRIPT_ARGS"]),
		SendImages:                        parseBool(sdConfig["SEND_IMAGES"], true),
		SaveImages:                        parseBool(sdConfig["SAVE_IMAGES"], true),
		AlwaysonScripts:                   parseObject(sdConfig["ALWAYSON_SCRIPTS"]),
	}

	return cfg, nil
}

// SaveSDConfig はSD設定をsd_config.txtファイルに保存
func SaveSDConfig(cfg *SDConfig) error {
	var lines []string

	if cfg.NegativePrompt != "" {
		lines = append(lines, fmt.Sprintf("NEGATIVE_PROMPT=%s", cfg.NegativePrompt))
	}
	if cfg.Steps != 0 {
		lines = append(lines, fmt.Sprintf("STEPS=%d", cfg.Steps))
	}
	if cfg.CfgScale != 0 {
		lines = append(lines, fmt.Sprintf("CFG_SCALE=%.1f", cfg.CfgScale))
	}
	if cfg.Width != 0 {
		lines = append(lines, fmt.Sprintf("WIDTH=%d", cfg.Width))
	}
	if cfg.Height != 0 {
		lines = append(lines, fmt.Sprintf("HEIGHT=%d", cfg.Height))
	}
	if cfg.SamplerName != "" {
		lines = append(lines, fmt.Sprintf("SAMPLER_NAME=%s", cfg.SamplerName))
	}
	if cfg.Seed != 0 {
		lines = append(lines, fmt.Sprintf("SEED=%d", cfg.Seed))
	}
	// 他のフィールドも同様に、値がある場合のみ追加
	// 簡略化のため、主要なものだけ。必要に応じて追加。

	content := strings.Join(lines, "\n") + "\n"
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

func parseBool(s string, defaultValue bool) bool {
	if value, err := strconv.ParseBool(s); err == nil {
		return value
	}
	return defaultValue
}

func parseStringArray(s string) []string {
	if s == "" {
		return nil
	}
	// シンプルにカンマ区切り
	return strings.Split(s, ",")
}

func parseArray(s string) []interface{} {
	if s == "" {
		return nil
	}
	// JSONとしてパース
	var arr []interface{}
	if err := json.Unmarshal([]byte(s), &arr); err == nil {
		return arr
	}
	return nil
}

func parseObject(s string) map[string]interface{} {
	if s == "" {
		return nil
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(s), &obj); err == nil {
		return obj
	}
	return nil
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
