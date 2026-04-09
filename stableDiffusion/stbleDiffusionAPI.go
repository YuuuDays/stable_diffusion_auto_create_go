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
		"width":           cfg.Width,
		"height":          cfg.Height,
		"cfg_scale":       cfg.CfgScale,
		"steps":           cfg.Steps,
		"sampler_name":    cfg.SamplerName,
		"seed":            cfg.Seed,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("ペイロード作成エラー: %v", err)
	}

	// HTTPリクエスト作成
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"http://127.0.0.1:7860/sdapi/v1/txt2img",
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
