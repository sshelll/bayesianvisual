package model

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/sshelll/bayesianvisual/internal/styles"
)

// View 渲染视图
func (m Model) View() tea.View {
	var view tea.View
	if m.Quitting {
		return view
	}

	var content string
	switch m.State {
	case StateViewing:
		content = m.renderViewing()
	case StateMenu:
		content = m.renderMenu()
	case StateInputPriorA, StateInputLikelihoodA, StateInputLikelihoodNotA, StateInputDescA, StateInputDescB:
		content = m.renderInput()
	}

	view.SetContent(content)
	return view
}

func (m Model) renderViewing() string {
	diagram := m.renderBayesianDiagram()
	footer := styles.FooterStyle.Render("Press n/enter/space for new calculation • Press q to quit")
	return lipgloss.JoinVertical(lipgloss.Left, diagram, footer)
}

func (m Model) renderMenu() string {
	title := styles.TitleStyle.Render("📊 Choose Calculation Mode")

	menuItems := []string{
		"Iterative Calculation (use previous posterior as new prior)",
		"New Calculation (enter all values from scratch)",
		"Customize Descriptions (define what A and B represent)",
	}

	var items []string
	for i, item := range menuItems {
		cursor := " "
		if m.MenuCursor == i {
			cursor = "▶"
			item = styles.SelectedItemStyle.Render(item)
		} else {
			item = styles.NormalItemStyle.Render(item)
		}
		items = append(items, fmt.Sprintf("%s %s", cursor, item))
	}

	menu := styles.MenuStyle.Render(lipgloss.JoinVertical(lipgloss.Left, items...))
	footer := styles.FooterStyle.Render("↑/↓ or j/k to navigate • enter to select • esc to cancel")

	return lipgloss.JoinVertical(lipgloss.Left, title, menu, footer)
}

func (m Model) renderInput() string {
	var title, prompt string

	switch m.State {
	case StateInputPriorA:
		title = "📊 Enter Prior Probability"
		prompt = "P(A) - Prior probability:"
	case StateInputLikelihoodA:
		if m.IterativeMode {
			title = "📊 Iterative Calculation"
			prompt = fmt.Sprintf("Previous P(A|B) = %.4f (used as new prior)\nP(B|A) - Likelihood:", m.TempPriorA)
		} else {
			title = "📊 Enter Likelihood"
			prompt = "P(B|A) - Likelihood given A:"
		}
	case StateInputLikelihoodNotA:
		title = "📊 Enter Likelihood"
		prompt = "P(B|¬A) - Likelihood given not A:"
	case StateInputDescA:
		title = "📊 Customize Event A"
		prompt = "Describe what event A represents:"
	case StateInputDescB:
		title = "📊 Customize Event B"
		prompt = "Describe what event B represents:"
	}

	titleRendered := styles.TitleStyle.Render(title)
	promptRendered := styles.LabelStyle.Render(prompt)
	input := m.TextInput.View()

	var errorRendered string
	if m.ErrorMsg != "" {
		errorRendered = styles.ErrorStyle.Render(m.ErrorMsg)
	}

	footer := styles.FooterStyle.Render("enter to continue • esc to cancel")

	parts := []string{titleRendered, promptRendered, input}
	if errorRendered != "" {
		parts = append(parts, errorRendered)
	}
	parts = append(parts, footer)

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (m Model) renderBayesianDiagram() string {
	// 计算尺寸
	height := m.SquareSize
	width := m.SquareSize * 2

	// 计算分割位置
	leftWidth := int(float64(width) * m.PriorA)
	rightWidth := width - leftWidth

	// 设置最小高度，确保文字可读
	minHeight := 3

	leftTopHeight := int(float64(height) * (1 - m.LikelihoodA))
	leftBottomHeight := height - leftTopHeight

	rightTopHeight := int(float64(height) * (1 - m.LikelihoodNotA))
	rightBottomHeight := height - rightTopHeight

	// 如果概率区域太小，设置最小高度
	if leftBottomHeight > 0 && leftBottomHeight < minHeight {
		leftBottomHeight = minHeight
		leftTopHeight = max(height-leftBottomHeight, 0)
	}

	if rightBottomHeight > 0 && rightBottomHeight < minHeight {
		rightBottomHeight = minHeight
		rightTopHeight = max(height-rightBottomHeight, 0)
	}

	// 设置最小宽度，确保文字可读
	minWidth := 8
	if leftWidth > 0 && leftWidth < minWidth {
		leftWidth = minWidth
		rightWidth = max(width-leftWidth, 0)
	}

	if rightWidth > 0 && rightWidth < minWidth {
		rightWidth = minWidth
		leftWidth = max(width-rightWidth, 0)
	}

	// 创建四个区域并组合左右两侧
	// 跳过高度为 0 的区域，避免 lipgloss 渲染空字符串时产生多余的行
	var leftParts, rightParts []string

	if leftTopHeight > 0 {
		leftParts = append(leftParts, m.renderBox(leftWidth, leftTopHeight, false, ""))
	}
	if leftBottomHeight > 0 {
		leftParts = append(leftParts, m.renderBox(leftWidth, leftBottomHeight, true,
			fmt.Sprintf("P(B|A)\n%.0f%%", m.LikelihoodA*100)))
	}

	if rightTopHeight > 0 {
		rightParts = append(rightParts, m.renderBox(rightWidth, rightTopHeight, false, ""))
	}
	if rightBottomHeight > 0 {
		rightParts = append(rightParts, m.renderBox(rightWidth, rightBottomHeight, true,
			fmt.Sprintf("P(B|¬A)\n%.0f%%", m.LikelihoodNotA*100)))
	}

	leftSide := lipgloss.JoinVertical(lipgloss.Left, leftParts...)
	rightSide := lipgloss.JoinVertical(lipgloss.Left, rightParts...)

	// 添加竖线分隔符
	divider := m.renderVerticalDivider(height)

	// 组合完整的图表
	diagram := lipgloss.JoinHorizontal(lipgloss.Top, leftSide, divider, rightSide)

	// 添加容器边框
	boxed := styles.ContainerStyle.Render(diagram)

	// 创建底部标签（在边框外侧）
	// 每个标签对齐到对应矩形的下方
	leftLabel := fmt.Sprintf("P(A)=%.0f%%", m.PriorA*100)
	rightLabel := fmt.Sprintf("P(¬A)=%.0f%%", (1-m.PriorA)*100)

	// 计算标签需要的宽度，确保标签在矩形下方居中
	// leftWidth 对应左侧矩形，rightWidth 对应右侧矩形
	// 加上边框和分隔符的宽度：左边框(1) + leftWidth + 分隔符(1) + rightWidth + 右边框(1)
	leftLabelWidth := leftWidth + 1  // 包含左边框
	rightLabelWidth := rightWidth + 2 // 包含分隔符和右边框

	leftLabelStyled := lipgloss.NewStyle().
		Width(leftLabelWidth).
		Align(lipgloss.Center).
		Foreground(styles.TextColor).
		Render(leftLabel)

	rightLabelStyled := lipgloss.NewStyle().
		Width(rightLabelWidth).
		Align(lipgloss.Center).
		Foreground(styles.TextColor).
		Render(rightLabel)

	bottomLabels := lipgloss.JoinHorizontal(lipgloss.Top, leftLabelStyled, rightLabelStyled)

	// 组合边框和底部标签
	boxedWithLabels := lipgloss.JoinVertical(lipgloss.Left, boxed, bottomLabels)

	// 创建信息面板
	info := m.renderInfoPanel()

	// 将图表和信息面板并排放置（信息面板垂直居中）
	mainContent := lipgloss.JoinHorizontal(lipgloss.Center, boxedWithLabels, info)

	// 创建标题
	title := styles.TitleStyle.Render("📊 Bayesian Theorem Visualization")

	// 组合所有部分
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		mainContent,
	)

	return content
}

func (m Model) renderBox(width, height int, filled bool, label string) string {
	style := styles.EmptyStyle
	if filled {
		style = styles.ProbabilityStyle
	}

	return style.
		Width(width).
		Height(height).
		Render(label)
}

func (m Model) renderVerticalDivider(height int) string {
	divider := ""
	for range height {
		divider += "│\n"
	}
	return styles.DividerStyle.Render(divider[:len(divider)-1])
}

func (m Model) renderInfoPanel() string {
	var parts []string

	// 如果有自定义描述，显示描述性句子
	if m.DescA != "" && m.DescB != "" {
		posterior := m.CalculatePosterior()
		descStyle := lipgloss.NewStyle().
			Foreground(styles.AccentColor).
			Italic(true).
			Width(40)

		sentence := fmt.Sprintf("Given that %s, the probability that %s is %.2f%%",
			m.DescB, m.DescA, posterior*100)
		descRendered := descStyle.Render(sentence)
		parts = append(parts, descRendered)
		parts = append(parts, "") // 空行分隔
	}

	// 左侧先验概率
	leftLabel := styles.LabelStyle.Render("Left P(A):")
	leftValue := styles.ValueStyle.Render(fmt.Sprintf("%.1f%%", m.PriorA*100))
	leftInfo := lipgloss.JoinHorizontal(lipgloss.Left, leftLabel, " ", leftValue)

	// 右侧先验概率
	rightLabel := styles.LabelStyle.Render("Right P(¬A):")
	rightValue := styles.ValueStyle.Render(fmt.Sprintf("%.1f%%", (1-m.PriorA)*100))
	rightInfo := lipgloss.JoinHorizontal(lipgloss.Left, rightLabel, " ", rightValue)

	// 组合先验概率信息
	priorInfo := lipgloss.JoinHorizontal(lipgloss.Left, leftInfo, "  •  ", rightInfo)

	// 似然概率信息
	likelihoodLabel := styles.LabelStyle.Render("Likelihood:")
	likelihoodLeft := styles.ValueStyle.Render(fmt.Sprintf("P(B|A)=%.1f%%", m.LikelihoodA*100))
	likelihoodRight := styles.ValueStyle.Render(fmt.Sprintf("P(B|¬A)=%.1f%%", m.LikelihoodNotA*100))
	likelihoodInfo := lipgloss.JoinHorizontal(lipgloss.Left,
		likelihoodLabel, " ", likelihoodLeft, "  ", likelihoodRight)

	// 后验概率信息
	posterior := m.CalculatePosterior()
	posteriorLabel := styles.LabelStyle.Render("Posterior:")
	posteriorValue := styles.ValueStyle.Render(fmt.Sprintf("P(A|B) = %.2f%%", posterior*100))
	posteriorInfo := lipgloss.JoinHorizontal(lipgloss.Left, posteriorLabel, " ", posteriorValue)

	// 组合所有信息
	parts = append(parts, priorInfo, likelihoodInfo, posteriorInfo)
	allInfo := lipgloss.JoinVertical(lipgloss.Left, parts...)

	return styles.InfoPanelStyle.Render(allInfo)
}
