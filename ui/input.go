package ui

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// ReadInt は整数入力を読み取る
func ReadInt() int {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	num, err := strconv.Atoi(input)
	if err != nil {
		return -1
	}
	return num
}

// ReadFloat は浮動小数点入力を読み取る
func ReadFloat() float64 {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	num, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return -1.0
	}
	return num
}

// ReadString は文字列入力を読み取る
func ReadString() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
