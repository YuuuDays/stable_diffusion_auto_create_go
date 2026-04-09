package ui

import (
	"context"
	"fmt"
	"time"

	"sd-auto-new/common"
	"sd-auto-new/config"
	"sd-auto-new/situation"
	stablediffusion "sd-auto-new/stableDiffusion"
	"sd-auto-new/utils"
)

// runGeneration は画像生成を実行
func runGeneration(ctx context.Context, char *common.PromptItem, category *situation.SituationCategory, settings *GenerationSettings, cfg *config.SDConfig) {
	start := time.Now()

	// 出力ディレクトリ作成
	outputDir, err := utils.CreateOutputDir("output")
	if err != nil {
		fmt.Println("❌ 出力ディレクトリ作成エラー:", err)
		return
	}
	fmt.Printf("📁 出力先: %s\n", outputDir)

	// 使用モデル取得と表示
	model := cfg.Model
	if model == "" {
		currentModel, err := stablediffusion.GetCurrentModel(cfg.APIURL)
		if err != nil {
			fmt.Printf("⚠️ 現在のモデル取得エラー: %v\n", err)
		} else {
			model = currentModel
			fmt.Printf("📋 使用モデル: %s\n", model)
		}
	} else {
		fmt.Printf("📋 使用モデル: %s\n", model)
	}

	// 使用LoRA取得と表示
	currentLora, err := stablediffusion.GetEffectiveLora(cfg)
	if err != nil {
		fmt.Printf("⚠️ LoRA取得エラー: %v\n", err)
	}
	if currentLora != "" {
		fmt.Printf("🔗 使用LoRA: %s\n", currentLora)
	}

	// Hires.fix表示
	if cfg.EnableHiresFix {
		fmt.Println("✨ Hires.fix が有効です")
	}

	// 総生成数を計算
	totalImages := 0
	for _, sit := range category.Situations {
		repeats := settings.SituationRepeats[sit.FileName]
		totalImages += repeats * settings.CategoryRepeats
	}

	fmt.Printf("🎨 合計 %d 枚の画像を生成します\n", totalImages)

	currentImage := 0

	// カテゴリ繰り返し
	for categoryRound := 0; categoryRound < settings.CategoryRepeats; categoryRound++ {
		fmt.Printf("\n🔄 カテゴリ %d/%d 回目\n", categoryRound+1, settings.CategoryRepeats)

		// 各シチュエーション
		for _, sit := range category.Situations {
			sitRepeats := settings.SituationRepeats[sit.FileName]

			for sitRound := 0; sitRound < sitRepeats; sitRound++ {
				select {
				case <-ctx.Done():
					fmt.Println("\n🛑 生成をキャンセルしました")
					return
				default:
				}

				currentImage++
				fmt.Printf("📸 生成中 %d/%d - キャラ:%s, シチュ:%s\n",
					currentImage, totalImages, char.Ja, sit.Name)

				// プロンプト生成（キャラクター + シチュエーション）
				prompt := char.En + ", " + sit.Prompts[0].En // 簡易的に最初のプロンプトを使用

				// LoRA追加（設定されている場合）
				if currentLora != "" {
					prompt = currentLora + ", " + prompt
				}

				seed := settings.SituationSeeds[sit.FileName]
				if seed == 0 || seed == -1 {
					seed = cfg.Seed
				}

				baseName := fmt.Sprintf("%s_%s_%03d", char.Ja, sit.Name, currentImage)

				// 生成実行
				actualSeed, err := stablediffusion.GenerateImage(ctx, prompt, cfg, outputDir, baseName, seed)
				if err != nil {
					fmt.Printf("❌ 生成エラー: %v\n", err)
					continue
				}
				fmt.Printf("✅ seed=%d で保存しました\n", actualSeed)
			}
		}
	}

	fmt.Println("\n✅ すべての生成が完了しました！")
	elapsed := time.Since(start)
	fmt.Printf("⏱️ 所要時間: %.2f秒\n", elapsed.Seconds())
}
