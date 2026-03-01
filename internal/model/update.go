package model

import (
	"fmt"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/shopspring/decimal"
)

// newProbabilityInput 创建概率输入框
func newProbabilityInput() textinput.Model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 10
	ti.SetWidth(20)
	ti.Placeholder = "0.0 - 1.0"
	return ti
}

// newTextInput 创建文本输入框
func newTextInput(placeholder string) textinput.Model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = -1
	ti.SetWidth(80)
	ti.Placeholder = placeholder
	return ti
}

// Update 更新模型
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		if !m.HistoryReady {
			m.HistoryViewport.SetWidth(80)
			m.HistoryViewport.SetHeight(8)
			m.HistoryReady = true
			// 同步已有的历史记录（例如从文件导入的）
			if len(m.IterationHistory) > 1 {
				m.SyncHistoryViewport()
			}
		}
		return m, nil

	case tea.KeyPressMsg:
		switch m.State {
		case StateViewing:
			return m.updateViewing(msg)
		case StateMenu:
			return m.updateMenu(msg)
		case StateInputPriorA, StateInputLikelihoodA, StateInputLikelihoodNotA, StateInputDescA, StateInputDescB, StateInputExportPath:
			return m.updateInput(msg)
		case StateIterationDescChoice:
			return m.updateIterationDescChoice(msg)
		case StateNewCalculationDescChoice:
			return m.updateNewCalculationDescChoice(msg)
		}
	default:
		// Forward other messages (e.g., paste/IME events, cursor blink) to textinput when in input states
		switch m.State {
		case StateInputPriorA, StateInputLikelihoodA, StateInputLikelihoodNotA, StateInputDescA, StateInputDescB, StateInputExportPath:
			m.TextInput, cmd = m.TextInput.Update(msg)
			return m, cmd
		case StateViewing:
			// Forward mouse scroll events to history viewport
			if len(m.IterationHistory) > 1 {
				m.HistoryViewport, cmd = m.HistoryViewport.Update(msg)
				return m, cmd
			}
		}
	}

	return m, cmd
}

func (m Model) updateViewing(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c", "esc":
		m.Quitting = true
		return m, tea.Quit
	case "n", "enter", " ":
		// 进入菜单
		m.State = StateMenu
		m.MenuCursor = 0
		m.ErrorMsg = ""
		return m, nil
	case "up", "k", "down", "j":
		// 滚动历史记录视口
		if len(m.IterationHistory) > 1 {
			var cmd tea.Cmd
			m.HistoryViewport, cmd = m.HistoryViewport.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func (m Model) updateMenu(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c", "esc":
		m.State = StateViewing
		return m, nil
	case "up", "k":
		if m.MenuCursor > 0 {
			m.MenuCursor--
		}
	case "down", "j":
		if m.MenuCursor < 2 {
			m.MenuCursor++
		}
	case "enter", " ":
		switch m.MenuCursor {
		case 0:
			// 迭代模式：使用上次的后验概率作为新的先验
			m.IterativeMode = true
			posterior := m.CalculatePosterior()
			m.TempPriorA = posterior

			// 只有在有迭代历史记录且有描述时，才询问是否使用相同描述
			if len(m.IterationHistory) > 0 && m.DescA != "" && m.DescB != "" {
				m.State = StateIterationDescChoice
				m.IterationDescCursor = 0
				m.ErrorMsg = ""
				return m, nil
			}

			// 没有历史记录或没有描述，直接进入似然概率输入
			m.State = StateInputLikelihoodA
			m.TextInput = newProbabilityInput()
			m.ErrorMsg = ""
			return m, textinput.Blink
		case 1:
			// 新运算模式：先询问是否自定义描述，稍后再清空历史和描述
			m.IterativeMode = false
			m.State = StateNewCalculationDescChoice
			m.NewCalcDescCursor = 0
			m.ErrorMsg = ""
			return m, nil
		case 2:
			// 导出迭代历史
			if len(m.IterationHistory) == 0 {
				m.ErrorMsg = "No iteration history to export."
				return m, nil
			}
			m.State = StateInputExportPath
			m.TextInput = newTextInput("e.g., ./history.json")
			m.ErrorMsg = ""
			return m, textinput.Blink
		}
	}
	return m, nil
}

func (m Model) updateInput(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "ctrl+c":
		m.Quitting = true
		return m, tea.Quit
	case "esc":
		m.State = StateMenu
		m.ErrorMsg = ""
		return m, nil
	case "enter":
		// 处理导出路径输入
		if m.State == StateInputExportPath {
			filepath := m.TextInput.Value()
			if filepath == "" {
				m.ErrorMsg = "Path cannot be empty."
				return m, nil
			}

			// 导出到文件
			if err := m.ExportToJSON(filepath); err != nil {
				m.ErrorMsg = fmt.Sprintf("Export failed: %v", err)
				return m, nil
			}

			// 导出成功，返回查看状态
			m.State = StateViewing
			m.TextInput.Blur()
			m.ErrorMsg = ""
			return m, nil
		}

		// 处理描述输入
		if m.State == StateInputDescA || m.State == StateInputDescB {
			if m.TextInput.Value() == "" {
				m.ErrorMsg = "Description cannot be empty."
				return m, nil
			}

			m.ErrorMsg = ""

			if m.State == StateInputDescA {
				m.DescA = m.TextInput.Value()
				m.State = StateInputDescB
				m.TextInput.SetValue("")
				m.TextInput.Placeholder = "e.g., the dog barked"
				return m, nil
			} else {
				m.DescB = m.TextInput.Value()

				// 如果是迭代模式，继续输入似然概率
				if m.IterativeMode {
					m.State = StateInputLikelihoodA
					m.TextInput = newProbabilityInput()
					return m, nil
				}

				// 非迭代模式（New Calculation），继续输入先验概率
				m.State = StateInputPriorA
				m.TextInput = newProbabilityInput()
				return m, nil
			}
		}

		// 验证并保存概率输入
		value, err := decimal.NewFromString(m.TextInput.Value())
		zero := decimal.Zero
		one := decimal.NewFromInt(1)
		if err != nil || value.LessThan(zero) || value.GreaterThan(one) {
			m.ErrorMsg = "Invalid input. Please enter a value between 0 and 1."
			m.TextInput.SetValue("")
			return m, nil
		}

		m.ErrorMsg = ""

		switch m.State {
		case StateInputPriorA:
			m.TempPriorA = value
			m.State = StateInputLikelihoodA
			m.TextInput.SetValue("")
			return m, nil

		case StateInputLikelihoodA:
			m.TempLikelihoodA = value
			m.State = StateInputLikelihoodNotA
			m.TextInput.SetValue("")
			return m, nil

		case StateInputLikelihoodNotA:
			m.TempLikelihoodNotA = value
			// 更新模型参数
			m.PriorA = m.TempPriorA
			m.LikelihoodA = m.TempLikelihoodA
			m.LikelihoodNotA = m.TempLikelihoodNotA
			// 添加迭代记录
			m.AddIterationRecord()
			// 返回查看状态
			m.State = StateViewing
			m.TextInput.Blur()
			return m, nil
		}
	}

	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

func (m Model) updateIterationDescChoice(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c", "esc":
		m.State = StateMenu
		return m, nil
	case "up", "k":
		if m.IterationDescCursor > 0 {
			m.IterationDescCursor--
		}
	case "down", "j":
		if m.IterationDescCursor < 1 {
			m.IterationDescCursor++
		}
	case "enter", " ":
		switch m.IterationDescCursor {
		case 0:
			// 使用相同描述，直接进入似然概率输入
			m.State = StateInputLikelihoodA
			m.TextInput = newProbabilityInput()
			m.ErrorMsg = ""
			return m, textinput.Blink
		case 1:
			// 输入新描述（迭代模式下，A 描述保持不变，只输入新的 B 描述）
			m.State = StateInputDescB
			m.TextInput = newTextInput("e.g., the dog barked")
			m.ErrorMsg = ""
			return m, textinput.Blink
		}
	}
	return m, nil
}

func (m Model) updateNewCalculationDescChoice(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c", "esc":
		m.State = StateMenu
		return m, nil
	case "up", "k":
		if m.NewCalcDescCursor > 0 {
			m.NewCalcDescCursor--
		}
	case "down", "j":
		if m.NewCalcDescCursor < 1 {
			m.NewCalcDescCursor++
		}
	case "enter", " ":
		switch m.NewCalcDescCursor {
		case 0:
			// 使用默认 A/B 描述
			m.IterationHistory = nil
			m.DescA = "A"
			m.DescB = "B"
			m.State = StateInputPriorA
			m.TextInput = newProbabilityInput()
			m.ErrorMsg = ""
			return m, textinput.Blink
		case 1:
			// 输入自定义描述
			m.IterationHistory = nil
			m.DescA = ""
			m.DescB = ""
			m.State = StateInputDescA
			m.TextInput = newTextInput("e.g., a thief came in")
			m.ErrorMsg = ""
			return m, textinput.Blink
		}
	}
	return m, nil
}
