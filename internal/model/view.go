package model

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/shopspring/decimal"
	"github.com/sshelll/bayesianvisual/internal/styles"
)

// formatPercent 格式化百分比，最多保留 4 位小数，自动去掉尾部的 0
// 对于极小值，保留至少 2 位有效数字，避免显示为 0
func formatPercent(value decimal.Decimal) string {
	// 转换为百分比
	hundred := decimal.NewFromInt(100)
	percent := value.Mul(hundred)

	// 四舍五入到 4 位小数
	rounded := percent.Round(4)

	// 如果四舍五入后变成 0 或非常接近 0，说明原始值太小
	// 需要保留更多位数以显示至少 2 位有效数字
	if rounded.IsZero() || rounded.Abs().LessThan(decimal.NewFromFloat(0.0001)) {
		// 对于极小值，找到前 2 位有效数字的位置
		// 使用字符串处理来实现
		str := percent.String()

		// 如果是负数，保留负号
		if percent.IsNegative() {
			str = str[1:]
		}

		// 找到第一个非零数字的位置
		significantDigits := 0

		for i, ch := range str {
			if ch == '.' {
				continue
			}
			if ch != '0' {
				significantDigits++
				// 保留 2 位有效数字
				if significantDigits >= 2 {
					// 返回到当前位置的字符串
					result := str[:i+1]
					if percent.IsNegative() {
						result = "-" + result
					}
					return result
				}
			}
		}

		// 如果只有 1 位或 0 位有效数字，返回原始字符串
		if percent.IsNegative() {
			return "-" + str
		}
		return str
	}

	// 正常情况：使用 String() 自动去掉尾部的 0
	return rounded.String()
}

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
	case StateInputPriorA, StateInputLikelihoodA, StateInputLikelihoodNotA, StateInputDescA, StateInputDescB, StateInputExportPath:
		content = m.renderInput()
	case StateIterationDescChoice:
		content = m.renderIterationDescChoice()
	case StateNewCalculationDescChoice:
		content = m.renderNewCalculationDescChoice()
	}

	view.SetContent(content)
	return view
}

func (m Model) renderViewing() string {
	diagram := m.renderBayesianDiagram()

	// 如果有历史记录，在下方显示
	var historyPanel string
	if len(m.IterationHistory) > 0 {
		historyPanel = m.renderHistoryPanel()
	}

	footer := styles.FooterStyle.Render("Press n/enter/space for new calculation • Press q to quit")

	if historyPanel != "" {
		return lipgloss.JoinVertical(lipgloss.Left, diagram, historyPanel, footer)
	}
	return lipgloss.JoinVertical(lipgloss.Left, diagram, footer)
}

func (m Model) renderMenu() string {
	title := styles.TitleStyle.Render("📊 Choose Calculation Mode")

	menuItems := []string{
		"Iterative Calculation (use previous posterior as new prior)",
		"New Calculation (start from scratch)",
		"Export Iteration History (save to JSON file)",
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

func (m Model) renderIterationDescChoice() string {
	title := styles.TitleStyle.Render("📊 Iterative Calculation - Event Descriptions")

	// 显示当前描述
	currentDesc := styles.DescSentenceStyle.Render(
		fmt.Sprintf("Current: A=\"%s\", B=\"%s\"", m.DescA, m.DescB))

	// 提示信息
	hint := styles.HistoryDetailStyle.Italic(true).
		Render("Note: Event A remains the same in iteration, only B changes")

	menuItems := []string{
		"Use the same event descriptions",
		"Enter new B description (A stays the same)",
	}

	var items []string
	for i, item := range menuItems {
		cursor := " "
		if m.IterationDescCursor == i {
			cursor = "▶"
			item = styles.SelectedItemStyle.Render(item)
		} else {
			item = styles.NormalItemStyle.Render(item)
		}
		items = append(items, fmt.Sprintf("%s %s", cursor, item))
	}

	menu := styles.MenuStyle.Render(lipgloss.JoinVertical(lipgloss.Left, items...))
	footer := styles.FooterStyle.Render("↑/↓ or j/k to navigate • enter to select • esc to cancel")

	return lipgloss.JoinVertical(lipgloss.Left, title, currentDesc, hint, "", menu, footer)
}

func (m Model) renderNewCalculationDescChoice() string {
	title := styles.TitleStyle.Render("📊 New Calculation - Event Descriptions")

	hint := styles.HistoryDetailStyle.Italic(true).
		Render("Choose how to describe events A and B")

	menuItems := []string{
		"Use default (A and B)",
		"Enter custom descriptions",
	}

	var items []string
	for i, item := range menuItems {
		cursor := " "
		if m.NewCalcDescCursor == i {
			cursor = "▶"
			item = styles.SelectedItemStyle.Render(item)
		} else {
			item = styles.NormalItemStyle.Render(item)
		}
		items = append(items, fmt.Sprintf("%s %s", cursor, item))
	}

	menu := styles.MenuStyle.Render(lipgloss.JoinVertical(lipgloss.Left, items...))
	footer := styles.FooterStyle.Render("↑/↓ or j/k to navigate • enter to select • esc to cancel")

	return lipgloss.JoinVertical(lipgloss.Left, title, hint, "", menu, footer)
}

func (m Model) renderInput() string {
	var title, prompt string

	switch m.State {
	case StateInputPriorA:
		title = "📊 Enter Prior Probability"
		if m.DescA != "" && m.DescA != "A" {
			prompt = fmt.Sprintf("P(A) - Prior probability (probability of \"%s\"):", m.DescA)
		} else {
			prompt = "P(A) - Prior probability (probability of A):"
		}
	case StateInputLikelihoodA:
		if m.IterativeMode {
			title = "📊 Iterative Calculation"
			var likelihoodExplanation string
			// 如果有自定义描述（不是默认的 A/B），显示详细说明
			if m.hasCustomDesc() {
				likelihoodExplanation = fmt.Sprintf("probability of \"%s\" given \"%s\"", m.DescB, m.DescA)
			} else {
				likelihoodExplanation = "probability of B given A"
			}
			prompt = fmt.Sprintf("Previous P(A|B) = %s%% (used as new prior)\nP(B|A) - Likelihood (%s):",
				formatPercent(m.TempPriorA), likelihoodExplanation)
		} else {
			title = "📊 Enter Likelihood"
			// 如果有自定义描述（不是默认的 A/B），显示详细说明
			if m.hasCustomDesc() {
				prompt = fmt.Sprintf("P(B|A) - Likelihood (probability of \"%s\" given \"%s\"):", m.DescB, m.DescA)
			} else {
				prompt = "P(B|A) - Likelihood (probability of B given A):"
			}
		}
	case StateInputLikelihoodNotA:
		title = "📊 Enter Likelihood"
		// 如果有自定义描述（不是默认的 A/B），显示详细说明
		if m.hasCustomDesc() {
			prompt = fmt.Sprintf("P(B|¬A) - Likelihood (probability of \"%s\" given NOT \"%s\"):", m.DescB, m.DescA)
		} else {
			prompt = "P(B|¬A) - Likelihood (probability of B given NOT A):"
		}
	case StateInputDescA:
		title = "📊 Customize Event A"
		prompt = "Describe what event A represents:"
	case StateInputDescB:
		title = "📊 Customize Event B"
		prompt = "Describe what event B represents:"
	case StateInputExportPath:
		title = "📊 Export Iteration History"
		prompt = "Enter file path to save (format: JSON):"
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

// hasCustomDesc 判断是否有自定义事件描述（非默认 A/B）
func (m Model) hasCustomDesc() bool {
	return m.DescA != "" && m.DescB != "" && !(m.DescA == "A" && m.DescB == "B")
}

func (m Model) renderBayesianDiagram() string {
	// 计算尺寸
	height := m.SquareSize
	width := m.SquareSize * 2

	// 计算分割位置（将 decimal 转换为 float64 用于计算）
	priorAFloat, _ := m.PriorA.Float64()
	likelihoodAFloat, _ := m.LikelihoodA.Float64()
	likelihoodNotAFloat, _ := m.LikelihoodNotA.Float64()

	leftWidth := int(float64(width) * priorAFloat)
	rightWidth := width - leftWidth

	// 设置最小高度，确保文字可读
	minHeight := 3

	leftTopHeight := int(float64(height) * (1 - likelihoodAFloat))
	leftBottomHeight := height - leftTopHeight

	rightTopHeight := int(float64(height) * (1 - likelihoodNotAFloat))
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
			fmt.Sprintf("P(B|A)\n%s%%", formatPercent(m.LikelihoodA))))
	}

	if rightTopHeight > 0 {
		rightParts = append(rightParts, m.renderBox(rightWidth, rightTopHeight, false, ""))
	}
	if rightBottomHeight > 0 {
		rightParts = append(rightParts, m.renderBox(rightWidth, rightBottomHeight, true,
			fmt.Sprintf("P(B|¬A)\n%s%%", formatPercent(m.LikelihoodNotA))))
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
	one := decimal.NewFromInt(1)
	leftLabel := fmt.Sprintf("P(A)=%s%%", formatPercent(m.PriorA))
	rightLabel := fmt.Sprintf("P(¬A)=%s%%", formatPercent(one.Sub(m.PriorA)))

	// 计算标签需要的宽度，确保标签在矩形下方居中
	// leftWidth 对应左侧矩形，rightWidth 对应右侧矩形
	// 加上边框和分隔符的宽度：左边框(1) + leftWidth + 分隔符(1) + rightWidth + 右边框(1)
	leftLabelWidth := leftWidth + 1   // 包含左边框
	rightLabelWidth := rightWidth + 2 // 包含分隔符和右边框

	// 只在宽度足够时才设置固定宽度并居中，避免文本换行
	leftLabelStyle := lipgloss.NewStyle().Foreground(styles.TextColor)
	if len(leftLabel) <= leftLabelWidth {
		leftLabelStyle = leftLabelStyle.Width(leftLabelWidth).Align(lipgloss.Center)
	}
	leftLabelStyled := leftLabelStyle.Render(leftLabel)

	rightLabelStyle := lipgloss.NewStyle().Foreground(styles.TextColor)
	if len(rightLabel) <= rightLabelWidth {
		rightLabelStyle = rightLabelStyle.Width(rightLabelWidth).Align(lipgloss.Center)
	}
	rightLabelStyled := rightLabelStyle.Render(rightLabel)

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
	var b strings.Builder
	for i := range height {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteRune('│')
	}
	return styles.DividerStyle.Render(b.String())
}

func (m Model) renderInfoPanel() string {
	var parts []string

	// 标题：当前计算
	currentTitle := styles.InfoSectionTitle.Render("Current Calculation")
	parts = append(parts, currentTitle)

	// 如果有自定义描述，显示描述性句子
	if m.DescA != "" && m.DescB != "" {
		posterior := m.CalculatePosterior()
		sentence := fmt.Sprintf("Given that \"%s\", the probability that \"%s\" is %s%%",
			m.DescB, m.DescA, formatPercent(posterior))
		parts = append(parts, styles.DescSentenceStyle.Render(sentence))
	}

	parts = append(parts, "") // 空行分隔

	// 先验概率
	leftLabel := styles.LabelStyle.Render("Prior P(A):")
	leftValue := styles.ValueStyle.Render(fmt.Sprintf("%s%%", formatPercent(m.PriorA)))
	leftInfo := lipgloss.JoinHorizontal(lipgloss.Left, leftLabel, " ", leftValue)

	// 似然概率
	likelihoodLabel := styles.LabelStyle.Render("Likelihood:")
	likelihoodLeft := styles.ValueStyle.Render(fmt.Sprintf("P(B|A)=%s%%", formatPercent(m.LikelihoodA)))
	likelihoodRight := styles.ValueStyle.Render(fmt.Sprintf("P(B|¬A)=%s%%", formatPercent(m.LikelihoodNotA)))
	likelihoodInfo := lipgloss.JoinHorizontal(lipgloss.Left,
		likelihoodLabel, " ", likelihoodLeft, "  ", likelihoodRight)

	parts = append(parts, leftInfo, likelihoodInfo)

	parts = append(parts, "") // 空行分隔

	// 后验概率（用醒目的样式强调结果）
	posterior := m.CalculatePosterior()
	posteriorLabel := styles.LabelStyle.Render("Posterior:")
	posteriorValue := styles.PosteriorValueStyle.Render(fmt.Sprintf("P(A|B) = %s%%", formatPercent(posterior)))
	posteriorInfo := lipgloss.JoinHorizontal(lipgloss.Left, posteriorLabel, " ", posteriorValue)
	parts = append(parts, posteriorInfo)

	allInfo := lipgloss.JoinVertical(lipgloss.Left, parts...)

	return styles.InfoPanelStyle.Render(allInfo)
}

// buildHistoryContent 构建历史记录内容字符串（不含标题）
func (m Model) buildHistoryContent() string {
	var parts []string

	// 渲染所有历史记录（包含最新的一条）
	endIdx := len(m.IterationHistory)

	for i := range endIdx {
		record := m.IterationHistory[i]
		iterNum := i + 1

		iterLabel := styles.IterLabelStyle.Render(fmt.Sprintf("#%d", iterNum))

		// 主要信息行
		var mainText string
		if record.DescA != "" && record.DescB != "" {
			mainText = styles.HistoryDescTextStyle.Render(
				fmt.Sprintf("Given \"%s\" → \"%s\": %s%%",
					record.DescB, record.DescA, formatPercent(record.Posterior)))
		} else {
			mainText = styles.HistoryMainTextStyle.Render(
				fmt.Sprintf("P(A)=%s%% → P(A|B)=%s%%",
					formatPercent(record.PriorA), formatPercent(record.Posterior)))
		}
		parts = append(parts, lipgloss.JoinHorizontal(lipgloss.Left, iterLabel, " ", mainText))

		// Likelihood 详情
		likelihoodText := styles.HistoryDetailStyle.Render(
			fmt.Sprintf("  P(B|A)=%s%%, P(B|¬A)=%s%%, P(A)=%s%%",
				formatPercent(record.LikelihoodA),
				formatPercent(record.LikelihoodNotA),
				formatPercent(record.PriorA),
			),
		)
		parts = append(parts, likelihoodText)
	}

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

// renderHistoryPanel 渲染历史记录面板（独立显示在下方）
func (m Model) renderHistoryPanel() string {
	// 水平分隔线
	divider := styles.HorizontalDividerStyle.Render(strings.Repeat("─", 60))

	// 标题始终置顶，不随 viewport 滚动
	historyTitle := styles.InfoSectionTitle.Render(
		fmt.Sprintf("Iteration History (%d)", len(m.IterationHistory)))
	scrollHint := lipgloss.NewStyle().
		Foreground(styles.DimTextColor).
		Italic(true).
		Render("  ↑/↓ to scroll")

	titleLine := lipgloss.JoinHorizontal(lipgloss.Bottom, historyTitle, scrollHint)

	// 渲染折线图
	chart := m.renderProbabilityChart()

	content := lipgloss.JoinVertical(lipgloss.Left, divider, titleLine, chart, "", m.HistoryViewport.View())
	return styles.HistoryPanelStyle.Render(content)
}

// renderProbabilityChart 渲染后验概率散点图
func (m Model) renderProbabilityChart() string {
	if len(m.IterationHistory) <= 1 {
		return ""
	}

	// 只展示最后 10 条历史记录
	history := m.IterationHistory
	startIdx := 0
	if len(history) > 10 {
		startIdx = len(history) - 10
		history = history[startIdx:]
	}

	// 收集后验概率百分比数据
	n := len(history)
	values := make([]float64, n)
	labels := make([]string, n)
	minVal, maxVal := 100.0, 0.0
	for i, r := range history {
		f, _ := r.Posterior.Float64()
		values[i] = f * 100
		pct, _ := r.Posterior.Mul(decimal.NewFromInt(100)).Float64()
		labels[i] = fmt.Sprintf("%.2f", pct)
		if values[i] < minVal {
			minVal = values[i]
		}
		if values[i] > maxVal {
			maxVal = values[i]
		}
	}

	// 图表参数
	chartHeight := 8
	chartWidth := n // 每个数据点占一列

	// 计算 Y 轴范围，留一点边距
	yRange := maxVal - minVal
	if yRange < 1 {
		yRange = 1
	}
	yMin := minVal - yRange*0.1
	yMax := maxVal + yRange*0.1
	if yMin < 0 {
		yMin = 0
	}
	if yMax > 100 {
		yMax = 100
	}
	yRange = yMax - yMin

	// Y 轴标签宽度
	labelWidth := 8

	// 构建图表网格
	var lines []string

	dotStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("114"))
	probStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("249"))
	axisStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("249"))

	// 预计算每个数据点所在的行
	dataRows := make([]int, chartWidth)
	for col := range chartWidth {
		dataRows[col] = int((yMax - values[col]) / yRange * float64(chartHeight))
		dataRows[col] = max(0, min(dataRows[col], chartHeight-1))
	}

	// 预计算每列的显示宽度：" ●(xx.xx%)" 的可见字符宽度
	// 点列格式: " ●(xx%%)" ，空列需要用同等宽度的空格填充
	colWidths := make([]int, chartWidth)
	for col := range chartWidth {
		// " ●" = 2 可见字符宽度 + "(xx%%)" 的长度
		probText := fmt.Sprintf("(%s%%)", labels[col])
		colWidths[col] = 2 + len(probText) // " ●" 占 2 宽度（空格+●）
	}

	for row := range chartHeight {
		// 当前行对应的 Y 值（从上到下递减）
		yVal := yMax - (float64(row)+0.5)*yRange/float64(chartHeight)

		// Y 轴标签
		label := fmt.Sprintf("%6.1f%%", yVal)
		label = labelStyle.Render(label)

		// 绘制该行
		var rowBuf strings.Builder
		rowBuf.WriteString(label)
		rowBuf.WriteString(axisStyle.Render(" │"))

		for col := range chartWidth {
			if dataRows[col] == row {
				probText := fmt.Sprintf("(%s%%)", labels[col])
				rowBuf.WriteString(dotStyle.Render(" ●"))
				rowBuf.WriteString(probStyle.Render(probText))
			} else {
				rowBuf.WriteString(strings.Repeat(" ", colWidths[col]))
			}
		}

		lines = append(lines, rowBuf.String())
	}

	// X 轴线
	var xAxisBuf strings.Builder
	xAxisBuf.WriteString(strings.Repeat(" ", labelWidth))
	xAxisBuf.WriteString(axisStyle.Render(" └"))
	for col := range chartWidth {
		xAxisBuf.WriteString(axisStyle.Render(strings.Repeat("─", colWidths[col])))
	}
	lines = append(lines, xAxisBuf.String())

	// X 轴标签（迭代编号）
	var xLabelBuf strings.Builder
	xLabelBuf.WriteString(strings.Repeat(" ", labelWidth+2))
	for col := range chartWidth {
		numStr := fmt.Sprintf("#%d", startIdx+col+1)
		padding := colWidths[col] - len(numStr)
		left := padding / 2
		right := padding - left
		xLabelBuf.WriteString(strings.Repeat(" ", left))
		xLabelBuf.WriteString(labelStyle.Render(numStr))
		xLabelBuf.WriteString(strings.Repeat(" ", right))
	}
	lines = append(lines, xLabelBuf.String())

	// Caption
	captionStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
	caption := captionStyle.Render("P(A|B) Posterior Trend (%)")
	lines = append(lines, strings.Repeat(" ", labelWidth+2)+caption)

	graph := strings.Join(lines, "\n")

	// 用圆角边框包裹
	bordered := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("238")).
		Padding(0, 1).
		Render(graph)

	return bordered
}
