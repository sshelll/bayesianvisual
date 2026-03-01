package model

import (
	"strconv"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

// Update 更新模型
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.KeyPressMsg:
		switch m.State {
		case StateViewing:
			return m.updateViewing(msg)
		case StateMenu:
			return m.updateMenu(msg)
		case StateInputPriorA, StateInputLikelihoodA, StateInputLikelihoodNotA, StateInputDescA, StateInputDescB:
			return m.updateInput(msg)
		}
	default:
		// Forward other messages (e.g., paste/IME events, cursor blink) to textinput when in input states
		switch m.State {
		case StateInputPriorA, StateInputLikelihoodA, StateInputLikelihoodNotA, StateInputDescA, StateInputDescB:
			m.TextInput, cmd = m.TextInput.Update(msg)
			return m, cmd
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
			// 直接进入似然概率输入
			m.State = StateInputLikelihoodA
			m.TextInput = textinput.New()
			m.TextInput.Focus()
			m.TextInput.CharLimit = 10
			m.TextInput.SetWidth(20)
			m.TextInput.Placeholder = "0.0 - 1.0"
			m.ErrorMsg = ""
			return m, textinput.Blink
		case 1:
			// 重新绘图模式：从先验概率开始输入
			m.IterativeMode = false
			m.State = StateInputPriorA
			m.TextInput = textinput.New()
			m.TextInput.Focus()
			m.TextInput.CharLimit = 10
			m.TextInput.SetWidth(20)
			m.TextInput.Placeholder = "0.0 - 1.0"
			m.ErrorMsg = ""
			return m, textinput.Blink
		case 2:
			// 自定义描述模式
			m.State = StateInputDescA
			m.TextInput = textinput.New()
			m.TextInput.Focus()
			m.TextInput.CharLimit = -1
			m.TextInput.SetWidth(80)
			m.TextInput.Placeholder = "e.g., a thief came in"
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
				// 返回查看状态
				m.State = StateViewing
				m.TextInput.Blur()
				return m, nil
			}
		}

		// 验证并保存概率输入
		value, err := strconv.ParseFloat(m.TextInput.Value(), 64)
		if err != nil || value < 0 || value > 1 {
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
			// 返回查看状态
			m.State = StateViewing
			m.TextInput.Blur()
			return m, nil
		}
	}

	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}
