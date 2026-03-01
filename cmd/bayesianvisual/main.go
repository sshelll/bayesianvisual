package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/shopspring/decimal"
	"github.com/sshelll/bayesianvisual/internal/model"
)

func main() {
	// 初始化模型，设置贝叶斯参数示例值
	m := model.Model{
		PriorA:         decimal.NewFromFloat(0.05), // P(A) = 5%
		LikelihoodA:    decimal.NewFromFloat(0.95), // P(B|A) = 95%
		LikelihoodNotA: decimal.NewFromFloat(0.2),  // P(B|¬A) = 20%
		SquareSize:     15,                         // 固定正方形大小（可调整）
		State:          model.StateViewing,
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
