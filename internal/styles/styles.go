package styles

import "charm.land/lipgloss/v2"

// 现代化颜色主题
var (
	ProbabilityColor = lipgloss.Color("62")  // 青蓝色
	BorderColor      = lipgloss.Color("252") // 浅灰色（边框更显眼）
	AccentColor      = lipgloss.Color("212") // 粉色
	TextColor        = lipgloss.Color("252") // 浅灰色
	DimTextColor     = lipgloss.Color("241") // 暗灰色
	HighlightColor   = lipgloss.Color("99")  // 紫色
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

// Footer 样式
var FooterStyle = lipgloss.NewStyle().
	Foreground(DimTextColor).
	MarginTop(1).
	Italic(true)

// 信息面板样式
var InfoPanelStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(BorderColor).
	Padding(1, 2).
	MarginLeft(2)

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
