package model

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/sshelll/bayesianvisual/internal/bayesian"
)

// ViewState 视图状态
type ViewState int

const (
	StateViewing ViewState = iota
	StateMenu
	StateInputPriorA
	StateInputLikelihoodA
	StateInputLikelihoodNotA
	StateInputDescA
	StateInputDescB
)

// Model 应用模型
type Model struct {
	Width    int
	Height   int
	Quitting bool
	State    ViewState
	// 贝叶斯参数
	PriorA         float64 // P(A) 先验概率
	LikelihoodA    float64 // P(B|A) 似然概率
	LikelihoodNotA float64 // P(B|¬A) 似然概率
	// 正方形大小（可调整）
	SquareSize int
	// 菜单选择
	MenuCursor int
	// 输入框
	TextInput textinput.Model
	// 是否迭代模式
	IterativeMode bool
	// 临时存储输入值
	TempPriorA         float64
	TempLikelihoodA    float64
	TempLikelihoodNotA float64
	// 事件描述
	DescA string // A 事件的描述
	DescB string // B 事件的描述
	// 错误信息
	ErrorMsg string
}

// Init 初始化
func (m Model) Init() tea.Cmd {
	return nil
}

// CalculatePosterior 计算后验概率 P(A|B)
func (m Model) CalculatePosterior() float64 {
	calc := bayesian.Calculator{
		PriorA:         m.PriorA,
		LikelihoodA:    m.LikelihoodA,
		LikelihoodNotA: m.LikelihoodNotA,
	}
	return calc.CalculatePosterior()
}
