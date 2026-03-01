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
	StateIterationDescChoice      // 迭代时选择是否使用相同描述
	StateNewCalculationDescChoice // 新运算时选择是否自定义描述
)

// IterationRecord 迭代记录
type IterationRecord struct {
	PriorA         float64
	LikelihoodA    float64
	LikelihoodNotA float64
	Posterior      float64
	DescA          string
	DescB          string
}

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
	// 迭代历史记录
	IterationHistory []IterationRecord
	// 迭代描述选择
	IterationDescCursor int // 0: 使用相同描述, 1: 输入新描述
	// 新运算描述选择
	NewCalcDescCursor int // 0: 使用默认 A/B, 1: 输入自定义描述
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

// AddIterationRecord 添加迭代记录
func (m *Model) AddIterationRecord() {
	record := IterationRecord{
		PriorA:         m.PriorA,
		LikelihoodA:    m.LikelihoodA,
		LikelihoodNotA: m.LikelihoodNotA,
		Posterior:      m.CalculatePosterior(),
		DescA:          m.DescA,
		DescB:          m.DescB,
	}
	m.IterationHistory = append(m.IterationHistory, record)
}
