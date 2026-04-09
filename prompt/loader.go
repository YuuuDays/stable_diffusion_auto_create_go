package prompt

import (
	"encoding/json"
	"os"

	"sd-auto-new/common"
)

// LoadCharacters はキャラクター定義ファイルを読み込む
func LoadCharacters(filename string) ([]common.PromptItem, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var characters []common.PromptItem
	err = json.Unmarshal(file, &characters)
	if err != nil {
		return nil, err
	}

	return characters, nil
}
