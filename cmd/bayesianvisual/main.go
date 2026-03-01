package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
	"github.com/sshelll/bayesianvisual/internal/model"
)

var (
	importFile string
)

var rootCmd = &cobra.Command{
	Use:   "bayesianvisual",
	Short: "A terminal-based interactive visualization tool for Bayes' Theorem",
	Long:  `Bayesian Visual is a TUI application that helps you visualize and understand Bayesian probability calculations with interactive diagrams and iteration history.`,
	Run:   run,
}

func init() {
	rootCmd.Flags().StringVarP(&importFile, "import", "i", "", "Import iteration history from JSON file")
}

func run(cmd *cobra.Command, args []string) {
	// 初始化模型，设置贝叶斯参数示例值
	m := model.Model{
		PriorA:         decimal.NewFromFloat(0.05), // P(A) = 5%
		LikelihoodA:    decimal.NewFromFloat(0.95), // P(B|A) = 95%
		LikelihoodNotA: decimal.NewFromFloat(0.2),  // P(B|¬A) = 20%
		SquareSize:     18,                         // 固定正方形大小（可调整）
		State:          model.StateViewing,
	}

	// 如果指定了导入文件，加载历史记录
	if importFile != "" {
		if err := m.LoadFromJSON(importFile); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to import history: %v\n", err)
			os.Exit(1)
		}
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
