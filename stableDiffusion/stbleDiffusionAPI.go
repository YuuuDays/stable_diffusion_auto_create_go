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

	"sd-auto-new/config"
)

type SDResponse struct {
	Images []string `json:"images"`
}

// GenerateImage は指定されたプロンプトで画像を生成し、指定ディレクトリに保存
func GenerateImage(ctx context.Context, prompt string, cfg *config.SDConfig, outputDir, fileName string) error {
	// ペイロード作成
	payload := map[string]interface{}{
		"prompt":          prompt,
		"negative_prompt": cfg.NegativePrompt,
		"steps":           cfg.Steps,
		"cfg_scale":       cfg.CfgScale,
		"width":           cfg.Width,
		"height":          cfg.Height,
		"sampler_name":    cfg.SamplerName,
		"seed":            cfg.Seed,
	}

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
	if cfg.SendImages {
		payload["send_images"] = cfg.SendImages
	}
	if cfg.SaveImages {
		payload["save_images"] = cfg.SaveImages
	}
	if cfg.AlwaysonScripts != nil {
		payload["alwayson_scripts"] = cfg.AlwaysonScripts
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("ペイロード作成エラー: %v", err)
	}

	// HTTPリクエスト作成
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		cfg.APIURL+"/sdapi/v1/txt2img",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("リクエスト作成エラー: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// リクエスト送信
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return fmt.Errorf("リクエストキャンセル")
		}
		return fmt.Errorf("APIリクエストエラー: %v", err)
	}
	defer resp.Body.Close()

	// レスポンス読み込み
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("レスポンス読み込みエラー: %v", err)
	}

	var result SDResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return fmt.Errorf("JSONパースエラー: %v", err)
	}

	if len(result.Images) == 0 {
		return fmt.Errorf("画像データが返されませんでした")
	}

	// Base64デコード
	imageBytes, err := base64.StdEncoding.DecodeString(result.Images[0])
	if err != nil {
		return fmt.Errorf("Base64デコードエラー: %v", err)
	}

	// 保存
	fullPath := filepath.Join(outputDir, fileName)
	err = os.WriteFile(fullPath, imageBytes, 0644)
	if err != nil {
		return fmt.Errorf("ファイル保存エラー: %v", err)
	}

	return nil
}
