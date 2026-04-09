package common

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PromptItem struct {
	En string `json:"en"`
	Ja string `json:"ja"`
}

// ReadInt は標準入力から整数を読み取る
func ReadInt() int {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if num, err := strconv.Atoi(input); err == nil {
			return num
		}
		fmt.Print("無効な入力です。整数を入力してください: ")
	}
	return 0
}
