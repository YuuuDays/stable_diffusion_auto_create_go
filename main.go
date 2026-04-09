package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sd-auto-new/config"
	"sd-auto-new/prompt"
	"sd-auto-new/situation"
	"sd-auto-new/ui"
)

// SD API接続確認
func checkSDConnection(cfg *config.SDConfig) error {
	client := &http.Client{Timeout: 5 * time.Second}
	urls := []string{
		cfg.APIURL + "/sdapi/v1/options",
		cfg.APIURL + "/sdapi/v1/sd-models",
	}
	var lastErr error
	for _, url := range urls {
		resp, err := client.Get(url)
		if err != nil {
			lastErr = err
			continue
		}
		resp.Body.Close()
		if resp.StatusCode == 200 {
			fmt.Println("✅ SD WebUI APIに接続できました")
			return nil
		}
		lastErr = fmt.Errorf("ステータス: %d", resp.StatusCode)
	}
	return fmt.Errorf("SD WebUI APIが正常に応答しません。API URL と SD WebUI の起動モードを確認してください: %v", lastErr)
}

func main() {
	// 設定読み込み
	cfg, err := config.LoadSDConfig()
	if err != nil {
		fmt.Println("❌ 設定読み込みエラー:", err)
		os.Exit(1)
	}

	// SD接続確認
	if err := checkSDConnection(cfg); err != nil {
		fmt.Println("❌", err)
		os.Exit(1)
	}

	// Ctrl+C対応 context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n🛑 Ctrl+C 検知 生成停止します...")
		cancel()
		// SD側も強制停止
		http.Post(cfg.APIURL+"/sdapi/v1/interrupt", "application/json", nil)
	}()

	// キャラクター読み込み
	characters, err := prompt.LoadCharacters("src/character.txt")
	if err != nil {
		fmt.Println("❌ キャラクター読み込みエラー:", err)
		os.Exit(1)
	}

	// シチュエーション読み込み
	situations, err := situation.LoadAll("situation")
	if err != nil {
		fmt.Println("❌ シチュエーション読み込みエラー:", err)
		os.Exit(1)
	}

	// UI開始
	ui.Run(ctx, characters, situations, cfg)
}
