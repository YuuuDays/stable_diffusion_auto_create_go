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
func checkSDConnection() error {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("http://127.0.0.1:7860/sdapi/v1/sd-models")
	if err != nil {
		return fmt.Errorf("SD WebUI APIに接続できません。SD WebUIをAPIモードで起動してください: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("SD APIが正常に応答しません (ステータス: %d)", resp.StatusCode)
	}
	fmt.Println("✅ SD WebUI APIに接続できました")
	return nil
}

func main() {
	// SD接続確認
	if err := checkSDConnection(); err != nil {
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
		http.Post("http://127.0.0.1:7860/sdapi/v1/interrupt", "application/json", nil)
	}()

	// 設定読み込み
	cfg, err := config.LoadSDConfig()
	if err != nil {
		fmt.Println("❌ 設定読み込みエラー:", err)
		os.Exit(1)
	}

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
