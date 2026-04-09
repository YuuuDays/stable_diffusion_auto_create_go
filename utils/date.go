package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// GetTodayDateString は今日の日付をYYYY-MM-DD形式で返す
func GetTodayDateString() string {
	return time.Now().Format("2006-01-02")
}

// CreateOutputDir は今日の日付フォルダを作成し、パスを返す
func CreateOutputDir(baseDir string) (string, error) {
	dateStr := GetTodayDateString()
	outputDir := filepath.Join(baseDir, dateStr)

	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	return outputDir, nil
}

// GetNextFileNumber は指定ディレクトリ内で指定パターンにマッチするファイルの次の連番を取得
func GetNextFileNumber(dir, pattern string) (int, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return 1, nil // ディレクトリが存在しない場合は1から
	}

	maxNumber := 0

	for _, file := range files {
		name := file.Name()
		if matched, _ := filepath.Match(pattern, name); matched {
			// パターンから番号を抽出（例: "character_situation_001.png" → 1）
			var num int
			fmt.Sscanf(name, pattern, &num)
			if num > maxNumber {
				maxNumber = num
			}
		}
	}

	return maxNumber + 1, nil
}
