package stablediffusion

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"sd-auto-new/config"
)

type SDResponse struct {
	Images []string `json:"images"`
	Info   string   `json:"info"`
}

type SDOptions struct {
	SDModelCheckpoint string `json:"sd_model_checkpoint"`
}

type SDResultInfo struct {
	Seed int64 `json:"seed"`
}

// GetCurrentModel はSD WebUIの現在のモデルを取得
func GetCurrentModel(apiURL string) (string, error) {
	resp, err := http.Get(apiURL + "/sdapi/v1/options")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var options SDOptions
	if err := json.NewDecoder(resp.Body).Decode(&options); err != nil {
		return "", err
	}

	return options.SDModelCheckpoint, nil
}

// GetLastUsedLora は最後に使用したLoRAを取得
func GetLastUsedLora() (string, error) {
	data, err := os.ReadFile("config/.last_lora")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// SaveLastUsedLora は最後に使用したLoRAを保存
func SaveLastUsedLora(lora string) error {
	return os.WriteFile("config/.last_lora", []byte(lora), 0644)
}

// GetEffectiveLora は有効なLoRAを取得（設定または最後に使用したもの）
func GetEffectiveLora(cfg *config.SDConfig) (string, error) {
	lora := cfg.Lora
	if lora == "" {
		lastLora, err := GetLastUsedLora()
		if err == nil && lastLora != "" {
			lora = lastLora
			// 表示はrunGenerationで行う
		}
	} else {
		// LoRAが指定された場合、保存
		SaveLastUsedLora(lora)
	}
	return lora, nil
}

// GenerateImage は指定されたプロンプトで画像を生成し、指定ディレクトリに保存
func GenerateImage(ctx context.Context, prompt string, cfg *config.SDConfig, outputDir, baseName string, seed int64) (int64, error) {
	// モデル取得（設定されていない場合、現在のモデルを使用）
	model := cfg.Model
	if model == "" {
		currentModel, err := GetCurrentModel(cfg.APIURL)
		if err != nil {
			return 0, fmt.Errorf("現在のモデル取得エラー: %v", err)
		}
		model = currentModel
		// 表示はrunGenerationで行う
	}

	// ペイロード作成
	payload := map[string]interface{}{
		"prompt":          prompt,
		"negative_prompt": cfg.NegativePrompt,
		"steps":           cfg.Steps,
		"cfg_scale":       cfg.CfgScale,
		"width":           cfg.Width,
		"height":          cfg.Height,
		"sampler_name":    cfg.SamplerName,
		"seed":            seed,
	}

	// モデル指定
	payload["model"] = model

	// Hires.fix 設定
	if cfg.EnableHiresFix {
		payload["enable_hr"] = true
		if cfg.HrScale != 0 {
			payload["hr_scale"] = cfg.HrScale
		}
		if cfg.HrUpscaler != "" {
			payload["hr_upscaler"] = cfg.HrUpscaler
		}
		if cfg.HrResizeX != 0 {
			payload["hr_resize_x"] = cfg.HrResizeX
		}
		if cfg.HrResizeY != 0 {
			payload["hr_resize_y"] = cfg.HrResizeY
		}
		if cfg.HrSecondPassSteps != 0 {
			payload["hr_second_pass_steps"] = cfg.HrSecondPassSteps
		}
	}

	// モデル指定
	payload["model"] = model

	// 追加パラメータ（値がある場合のみ）
	if len(cfg.Styles) > 0 {
		payload["styles"] = cfg.Styles
	}
	if cfg.Subseed != 0 {
		payload["subseed"] = cfg.Subseed
	}
	if cfg.SubseedStrength != 0 {
		payload["subseed_strength"] = cfg.SubseedStrength
	}
	if cfg.SeedResizeFromH != 0 {
		payload["seed_resize_from_h"] = cfg.SeedResizeFromH
	}
	if cfg.SeedResizeFromW != 0 {
		payload["seed_resize_from_w"] = cfg.SeedResizeFromW
	}
	if cfg.BatchSize != 0 {
		payload["batch_size"] = cfg.BatchSize
	}
	if cfg.NIter != 0 {
		payload["n_iter"] = cfg.NIter
	}
	if cfg.RestoreFaces {
		payload["restore_faces"] = cfg.RestoreFaces
	}
	if cfg.Tiling {
		payload["tiling"] = cfg.Tiling
	}
	if cfg.DoNotSaveSamples {
		payload["do_not_save_samples"] = cfg.DoNotSaveSamples
	}
	if cfg.DoNotSaveGrid {
		payload["do_not_save_grid"] = cfg.DoNotSaveGrid
	}
	if cfg.Eta != 0 {
		payload["eta"] = cfg.Eta
	}
	if cfg.DenoisingStrength != 0 {
		payload["denoising_strength"] = cfg.DenoisingStrength
	} else if cfg.EnableHiresFix {
		payload["denoising_strength"] = 0.5 // Hires.fix有効時のデフォルト値
	}
	if cfg.SMinUncond != 0 {
		payload["s_min_uncond"] = cfg.SMinUncond
	}
	if cfg.SChurn != 0 {
		payload["s_churn"] = cfg.SChurn
	}
	if cfg.STmax != 0 {
		payload["s_tmax"] = cfg.STmax
	}
	if cfg.STmin != 0 {
		payload["s_tmin"] = cfg.STmin
	}
	if cfg.SNoise != 0 {
		payload["s_noise"] = cfg.SNoise
	}
	if cfg.OverrideSettings != nil {
		payload["override_settings"] = cfg.OverrideSettings
	}
	if cfg.OverrideSettingsRestoreAfterwards {
		payload["override_settings_restore_afterwards"] = cfg.OverrideSettingsRestoreAfterwards
	}
	if cfg.RefinerCheckpoint != "" {
		payload["refiner_checkpoint"] = cfg.RefinerCheckpoint
	}
	if cfg.RefinerSwitchAt != 0 {
		payload["refiner_switch_at"] = cfg.RefinerSwitchAt
	}
	if cfg.DisableExtraNetworks {
		payload["disable_extra_networks"] = cfg.DisableExtraNetworks
	}
	if cfg.ResizeMode != 0 {
		payload["resize_mode"] = cfg.ResizeMode
	}
	if cfg.ImageCfgScale != 0 {
		payload["image_cfg_scale"] = cfg.ImageCfgScale
	}
	if cfg.MaskBlur != 0 {
		payload["mask_blur"] = cfg.MaskBlur
	}
	if cfg.InpaintingFill != 0 {
		payload["inpainting_fill"] = cfg.InpaintingFill
	}
	if cfg.InpaintFullRes {
		payload["inpaint_full_res"] = cfg.InpaintFullRes
	}
	if cfg.InpaintFullResPadding != 0 {
		payload["inpaint_full_res_padding"] = cfg.InpaintFullResPadding
	}
	if cfg.InpaintingMaskInvert != 0 {
		payload["inpainting_mask_invert"] = cfg.InpaintingMaskInvert
	}
	if cfg.InitialNoiseMultiplier != 0 {
		payload["initial_noise_multiplier"] = cfg.InitialNoiseMultiplier
	}
	if cfg.ForceTaskId != "" {
		payload["force_task_id"] = cfg.ForceTaskId
	}
	if cfg.SamplerIndex != "" {
		payload["sampler_index"] = cfg.SamplerIndex
	}
	if cfg.IncludeInitImages {
		payload["include_init_images"] = cfg.IncludeInitImages
	}
	if cfg.ScriptName != "" {
		payload["script_name"] = cfg.ScriptName
	}
	if cfg.ScriptArgs != nil {
		payload["script_args"] = cfg.ScriptArgs
	}
	payload["send_images"] = true // 常に画像データをレスポンスに含める
	if cfg.SaveImages {
		payload["save_images"] = false // プログラム側で保存するため、SD側は保存しない
	}
	if cfg.AlwaysonScripts != nil {
		payload["alwayson_scripts"] = cfg.AlwaysonScripts
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("ペイロード作成エラー: %v", err)
	}

	// HTTPリクエスト作成
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		cfg.APIURL+"/sdapi/v1/txt2img",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return 0, fmt.Errorf("リクエスト作成エラー: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// リクエスト送信
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return 0, fmt.Errorf("リクエストキャンセル")
		}
		return 0, fmt.Errorf("APIリクエストエラー: %v", err)
	}
	defer resp.Body.Close()

	// レスポンス読み込み
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("レスポンス読み込みエラー: %v", err)
	}

	var result SDResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, fmt.Errorf("JSONパースエラー: %v", err)
	}

	if len(result.Images) == 0 {
		return 0, fmt.Errorf("画像データが返されませんでした（send_imagesが無効な可能性）")
	}

	if result.Images[0] == "" {
		return 0, fmt.Errorf("画像データが空です（SDのレスポンスにimagesデータがありません）")
	}

	// 常にInfoからactualSeedを取得する（指定seedに関わらず）
	actualSeed := seed
	if result.Info != "" {
		var info SDResultInfo
		if err := json.Unmarshal([]byte(result.Info), &info); err == nil && info.Seed > 0 {
			actualSeed = info.Seed
		} else if err != nil {
			fmt.Printf("⚠️ Info JSONパースエラー: %v\n", err)
		}
	} else {
		fmt.Println("⚠️ Infoフィールドが空です")
	}

	if actualSeed <= 0 {
		actualSeed = seed
	}

	fileName := fmt.Sprintf("%s_seed_%d.png", baseName, actualSeed)

	// Base64デコード（標準およびURLセーフbase64に対応）
	imageBytes, err := base64.StdEncoding.DecodeString(result.Images[0])
	if err != nil {
		// URLセーフbase64でリトライ
		imageBytes, err = base64.URLEncoding.DecodeString(result.Images[0])
		if err != nil {
			return 0, fmt.Errorf("Base64デコードエラー: %v", err)
		}
	}

	if len(imageBytes) == 0 {
		return 0, fmt.Errorf("デコード後の画像データが0バイトです")
	}

	// 保存
	fullPath := filepath.Join(outputDir, fileName)
	err = os.WriteFile(fullPath, imageBytes, 0644)
	if err != nil {
		return 0, fmt.Errorf("ファイル保存エラー: %v", err)
	}

	return actualSeed, nil
}
