package ui

import (
	"fmt"

	"sd-auto-new/common"
	"sd-auto-new/situation"
)

// GenerationSettings は生成設定を表す
type GenerationSettings struct {
	SituationRepeats map[string]int   // シチュエーションごとの繰り返し回数
	SituationSeeds   map[string]int64 // シチュエーションごとのシード値
	CategoryRepeats  int              // カテゴリ全体の繰り返し回数
}

// selectCharacter はキャラクターを選択
func selectCharacter(characters []common.PromptItem) *common.PromptItem {
	fmt.Println("👤 キャラクター選択")

	if len(characters) == 0 {
		fmt.Println("❌ キャラクターが読み込まれていません")
		return nil
	}

	fmt.Println("\nキャラクター一覧:")
	for i, char := range characters {
		fmt.Printf("  %d: %s\n", i, char.Ja)
	}

	fmt.Print("\nキャラクター番号を選択 (-1でランダム) >> ")
	idx := ReadInt()

	if idx == -1 {
		// ランダム選択
		randomIdx := 0 // 簡易的に最初のものを選択（本来はrand使用）
		fmt.Printf("✅ ランダム選択: %s\n", characters[randomIdx].Ja)
		fmt.Println(sep)
		return &characters[randomIdx]
	}

	if idx < 0 || idx >= len(characters) {
		fmt.Println("❌ 無効な選択です")
		return nil
	}

	fmt.Printf("✅ 選択: %s\n", characters[idx].Ja)
	fmt.Println(sep)
	return &characters[idx]
}

// selectSituationAndSettings はシチュエーションを選択し、設定を入力
func selectSituationAndSettings(categories []situation.SituationCategory) (*situation.SituationCategory, *GenerationSettings) {
	fmt.Println("📁 シチュエーションカテゴリ選択")

	if len(categories) == 0 {
		fmt.Println("❌ シチュエーションが読み込まれていません")
		return nil, nil
	}

	fmt.Println("\nカテゴリ一覧:")
	for i, cat := range categories {
		fmt.Printf("  %d: %s (%d個のシチュエーション)\n", i, cat.Name, len(cat.Situations))
	}

	fmt.Print("\nカテゴリ番号を選択 >> ")
	catIdx := ReadInt()

	if catIdx < 0 || catIdx >= len(categories) {
		fmt.Println("❌ 無効な選択です")
		return nil, nil
	}

	selectedCategory := &categories[catIdx]
	fmt.Printf("✅ 選択: %s\n", selectedCategory.Name)

	// 各シチュエーションのシードと繰り返し回数を入力
	settings := &GenerationSettings{
		SituationRepeats: make(map[string]int),
		SituationSeeds:   make(map[string]int64),
	}

	fmt.Println("\n各シチュエーションのシード指定と繰り返し回数を設定:")
	for _, sit := range selectedCategory.Situations {
		seed := int64(-1)
		fmt.Printf("  %s のシード指定を行いますか？ (1=はい 0=いいえ) >> ", sit.Name)
		if ReadInt() == 1 {
			fmt.Printf("  %s のシード値 >> ", sit.Name)
			seedInput := ReadInt()
			if seedInput >= 0 {
				seed = int64(seedInput)
			}
		}
		settings.SituationSeeds[sit.FileName] = seed

		fmt.Printf("  %s の回数 (0=スキップ) >> ", sit.Name)
		count := ReadInt()
		if count < 0 {
			count = 0
		}
		settings.SituationRepeats[sit.FileName] = count
	}

	// カテゴリ全体の繰り返し回数を入力
	fmt.Print("\nこのカテゴリ全体を何回繰り返しますか？ >> ")
	settings.CategoryRepeats = ReadInt()
	if settings.CategoryRepeats < 1 {
		settings.CategoryRepeats = 1
	}

	return selectedCategory, settings
}

// confirmGeneration は生成条件を確認
func confirmGeneration(char *common.PromptItem, category *situation.SituationCategory, settings *GenerationSettings) bool {
	fmt.Println("生成条件確認")
	fmt.Printf("キャラクター     : %s\n", char.Ja)
	fmt.Printf("カテゴリ         : %s\n", category.Name)
	fmt.Printf("カテゴリ繰り返し : %d回\n", settings.CategoryRepeats)
	fmt.Println("シチュエーション詳細:")
	totalImages := 0
	for _, sit := range category.Situations {
		repeats := settings.SituationRepeats[sit.FileName]
		seed := settings.SituationSeeds[sit.FileName]
		images := repeats * settings.CategoryRepeats
		totalImages += images
		if repeats == 0 {
			fmt.Printf("  %s: スキップ\n", sit.Name)
		} else if seed >= 0 {
			fmt.Printf("  %s: %d枚 (seed=%d)\n", sit.Name, images, seed)
		} else {
			fmt.Printf("  %s: %d枚\n", sit.Name, images)
		}
	}
	fmt.Printf("合計生成枚数     : %d枚\n", totalImages)

	fmt.Print("\nこの条件で生成しますか？ (1=はい 0=いいえ) >> ")
	answer := ReadString()
	if answer == "1" {
		fmt.Println("✅ 生成を実行します")
		return true
	} else {
		fmt.Println("❌ 生成をキャンセルしました")
		return false
	}
}
