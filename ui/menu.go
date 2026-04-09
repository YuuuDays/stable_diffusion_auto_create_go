package ui

import (
	"context"
	"fmt"

	"sd-auto-new/common"
	"sd-auto-new/config"
	"sd-auto-new/situation"
)

// Run はメインUIを実行
func Run(ctx context.Context, characters []common.PromptItem, situations []situation.SituationCategory, cfg *config.SDConfig) {
	for {
		fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("🎨 SD Auto Generation Tool")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("1. 生成モード")
		fmt.Println("0. 終了")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Print("選択 >> ")

		choice := ReadInt()

		switch choice {
		case 1:
			runGenerationMode(ctx, characters, situations, cfg)
		case 0:
			fmt.Println("終了します")
			return
		default:
			fmt.Println("❌ 無効な選択です")
		}
	}
}

// runGenerationMode は生成モードを実行
func runGenerationMode(ctx context.Context, characters []common.PromptItem, situations []situation.SituationCategory, cfg *config.SDConfig) {
	fmt.Println("\n🎲 生成モード")

	// キャラクター選択
	selectedChar := selectCharacter(characters)
	if selectedChar == nil {
		return // キャンセル
	}

	// シチュエーション選択
	selectedCategory, generationSettings := selectSituationAndSettings(situations)
	if selectedCategory == nil {
		return // キャンセル
	}

	// 確認
	if !confirmGeneration(selectedChar, selectedCategory, generationSettings) {
		return
	}

	// 生成実行
	runGeneration(ctx, selectedChar, selectedCategory, generationSettings, cfg)
}
