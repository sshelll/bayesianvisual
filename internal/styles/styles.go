package styles

import "charm.land/lipgloss/v2"

// 现代化颜色主题
var (
	ProbabilityColor = lipgloss.Color("62")  // 青蓝色
	BorderColor      = lipgloss.Color("252") // 浅灰色（边框更显眼）
	AccentColor      = lipgloss.Color("212") // 粉色
	TextColor        = lipgloss.Color("252") // 浅灰色
	DimTextColor     = lipgloss.Color("241") // 暗灰色（footer 等次要文本）
	HistoryTextColor = lipgloss.Color("249") // 中灰色（历史记录文本，比 DimText 亮）
	HighlightColor   = lipgloss.Color("99")  // 紫色
	SuccessColor     = lipgloss.Color("114") // 绿色（用于后验概率结果）
	ErrorColor       = lipgloss.Color("196") // 红色
)

// 标题样式
var TitleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(AccentColor).
	MarginBottom(1).
	Padding(0, 1)

// 主容器样式
var ContainerStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(BorderColor).
	Padding(0)

// 概率区域样式
var ProbabilityStyle = lipgloss.NewStyle().
	Background(ProbabilityColor).
	Foreground(lipgloss.Color("230")).
	Align(lipgloss.Center, lipgloss.Center).
	Bold(true)

// 空白区域样式
var EmptyStyle = lipgloss.NewStyle().
	Align(lipgloss.Center, lipgloss.Center)

// 分隔线样式
var DividerStyle = lipgloss.NewStyle().
	Foreground(BorderColor)

// 信息标签样式
var LabelStyle = lipgloss.NewStyle().
	Foreground(TextColor)

// 数值样式
var ValueStyle = lipgloss.NewStyle().
	Foreground(HighlightColor).
	Bold(true)

// 后验概率结果样式（加下划线强调）
var PosteriorValueStyle = lipgloss.NewStyle().
	Foreground(SuccessColor).
	Bold(true).
	Underline(true)

// Footer 样式
var FooterStyle = lipgloss.NewStyle().
	Foreground(DimTextColor).
	MarginTop(1).
	Italic(true)

// 信息面板样式
var InfoPanelStyle = lipgloss.NewStyle().
	Padding(1, 2).
	MarginLeft(2)

// 历史面板样式
var HistoryPanelStyle = lipgloss.NewStyle().
	Padding(1, 2)

// 水平分隔线样式
var HorizontalDividerStyle = lipgloss.NewStyle().
	Foreground(BorderColor).
	MarginTop(1)

// 菜单样式
var MenuStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(BorderColor).
	Padding(1, 2).
	MarginTop(1)

var SelectedItemStyle = lipgloss.NewStyle().
	Foreground(AccentColor).
	Bold(true)

var NormalItemStyle = lipgloss.NewStyle().
	Foreground(TextColor)

// 错误样式
var ErrorStyle = lipgloss.NewStyle().
	Foreground(ErrorColor).
	Bold(true).
	MarginTop(1)

// 历史记录迭代编号样式
var IterLabelStyle = lipgloss.NewStyle().
	Foreground(HighlightColor).
	Bold(true)

// 历史记录主文本样式
var HistoryMainTextStyle = lipgloss.NewStyle().
	Foreground(HistoryTextColor)

// 历史记录描述文本样式（带描述时使用斜体）
var HistoryDescTextStyle = lipgloss.NewStyle().
	Foreground(HistoryTextColor).
	Italic(true).
	Width(50)

// 历史记录 likelihood 详情样式
var HistoryDetailStyle = lipgloss.NewStyle().
	Foreground(DimTextColor)

// 信息面板小标题样式
var InfoSectionTitle = lipgloss.NewStyle().
	Foreground(AccentColor).
	Bold(true)

// 描述性句子样式
var DescSentenceStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("219")). // 浅粉色，比 AccentColor 柔和
	Italic(true).
	Width(50)
