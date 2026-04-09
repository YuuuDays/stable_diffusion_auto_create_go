package situation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"sd-auto-new/common"
)

// SituationCategory はシチュエーションカテゴリを表す
type SituationCategory struct {
	Name       string      // カテゴリ名（フォルダ名）
	Situations []Situation // シチュエーション一覧
}

// Situation は個別のシチュエーションを表す
type Situation struct {
	FileName string              // ファイル名（01_school.txt）
	Name     string              // 日本語名
	Prompts  []common.PromptItem // プロンプト一覧
}

// LoadAll はsituationフォルダ内の全シチュエーションを読み込む
func LoadAll(dirPath string) ([]SituationCategory, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var categories []SituationCategory

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		categoryName := entry.Name()
		categoryPath := filepath.Join(dirPath, categoryName)

		situations, err := loadSituationsFromCategory(categoryPath)
		if err != nil {
			fmt.Printf("⚠️ カテゴリ %s の読み込みエラー: %v\n", categoryName, err)
			continue
		}

		categories = append(categories, SituationCategory{
			Name:       categoryName,
			Situations: situations,
		})
	}

	return categories, nil
}

// loadSituationsFromCategory はカテゴリフォルダ内のシチュエーションを読み込む
func loadSituationsFromCategory(categoryPath string) ([]Situation, error) {
	files, err := filepath.Glob(filepath.Join(categoryPath, "*.txt"))
	if err != nil {
		return nil, err
	}

	sort.Strings(files)

	var situations []Situation

	for _, file := range files {
		prompts, err := loadPromptsFromFile(file)
		if err != nil {
			fmt.Printf("⚠️ ファイル %s の読み込みエラー: %v\n", file, err)
			continue
		}

		if len(prompts) == 0 {
			continue
		}

		fileName := filepath.Base(file)
		name := prompts[0].Ja // 最初のプロンプトの日本語名を使用

		situations = append(situations, Situation{
			FileName: fileName,
			Name:     name,
			Prompts:  prompts,
		})
	}

	return situations, nil
}

// loadPromptsFromFile はシチュエーションファイルからプロンプトを読み込む
func loadPromptsFromFile(filename string) ([]common.PromptItem, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var prompts []common.PromptItem
	err = json.Unmarshal(file, &prompts)
	if err != nil {
		return nil, err
	}

	return prompts, nil
}
