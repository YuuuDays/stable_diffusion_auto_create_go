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
		fmt.Println("2. 設定モード")
		fmt.Println("0. 終了")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Print("選択 >> ")

		choice := ReadInt()

		switch choice {
		case 1:
			runGenerationMode(ctx, characters, situations, cfg)
		case 2:
			runSettingsMode(cfg)
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

// runSettingsMode は設定モードを実行
func runSettingsMode(cfg *config.SDConfig) {
	fmt.Println("\n🔧 設定モード")

	for {
		displayCurrentSettings(cfg)

		fmt.Println("\n変更したい項目を選択:")
		fmt.Println("1. Negative Prompt")
		fmt.Println("2. Steps")
		fmt.Println("3. CFG Scale")
		fmt.Println("4. Width")
		fmt.Println("5. Height")
		fmt.Println("6. Sampler Name")
		fmt.Println("7. Seed")
		fmt.Println("0. 戻る")
		fmt.Print("選択 >> ")

		choice := ReadInt()

		switch choice {
		case 0:
			return
		case 1:
			fmt.Print("新しいNegative Prompt >> ")
			cfg.NegativePrompt = ReadString()
		case 2:
			fmt.Print("新しいSteps (1-100) >> ")
			newSteps := ReadInt()
			if newSteps < 1 || newSteps > 100 {
				fmt.Println("❌ Stepsは1-100の範囲で入力してください")
				continue
			}
			cfg.Steps = newSteps
		case 3:
			fmt.Print("新しいCFG Scale (1.0-20.0) >> ")
			newCfg := ReadFloat()
			if newCfg < 1.0 || newCfg > 20.0 {
				fmt.Println("❌ CFG Scaleは1.0-20.0の範囲で入力してください")
				continue
			}
			cfg.CfgScale = newCfg
		case 4:
			fmt.Print("新しいWidth (64-2048, 64の倍数) >> ")
			newWidth := ReadInt()
			if newWidth < 64 || newWidth > 2048 || newWidth%64 != 0 {
				fmt.Println("❌ Widthは64-2048の範囲で64の倍数で入力してください")
				continue
			}
			cfg.Width = newWidth
		case 5:
			fmt.Print("新しいHeight (64-2048, 64の倍数) >> ")
			newHeight := ReadInt()
			if newHeight < 64 || newHeight > 2048 || newHeight%64 != 0 {
				fmt.Println("❌ Heightは64-2048の範囲で64の倍数で入力してください")
				continue
			}
			cfg.Height = newHeight
		case 6:
			fmt.Print("新しいSampler Name >> ")
			cfg.SamplerName = ReadString()
		case 7:
			fmt.Print("新しいSeed (-1でランダム) >> ")
			cfg.Seed = int64(ReadInt())
		default:
			fmt.Println("❌ 無効な選択です")
			continue
		}

		// 設定保存
		if err := config.SaveSDConfig(cfg); err != nil {
			fmt.Println("❌ 設定保存エラー:", err)
		} else {
			fmt.Println("✅ 設定を保存しました")
		}
	}
}

// displayCurrentSettings は現在の設定を表示
func displayCurrentSettings(cfg *config.SDConfig) {
	fmt.Println("\n📋 現在の設定:")
	fmt.Printf("  Negative Prompt: %s\n", cfg.NegativePrompt)
	fmt.Printf("  Steps: %d\n", cfg.Steps)
	fmt.Printf("  CFG Scale: %.1f\n", cfg.CfgScale)
	fmt.Printf("  Size: %dx%d\n", cfg.Width, cfg.Height)
	fmt.Printf("  Sampler: %s\n", cfg.SamplerName)
	fmt.Printf("  Seed: %d\n", cfg.Seed)
}
