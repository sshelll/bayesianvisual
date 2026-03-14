package model

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

// BuildTextView 以纯文本方式输出迭代历史的计算步骤和结果
func (m Model) BuildTextView() string {
	if len(m.IterationHistory) == 0 {
		return "No iteration history.\n"
	}

	var b strings.Builder

	// 标题
	if m.DescA != "" && m.DescA != "A" {
		b.WriteString(fmt.Sprintf("Bayesian Iteration History  (A = \"%s\")\n", m.DescA))
	} else {
		b.WriteString("Bayesian Iteration History\n")
	}
	b.WriteString(strings.Repeat("=", 60) + "\n\n")

	one := decimal.NewFromInt(1)

	for i, r := range m.IterationHistory {
		// 迭代编号
		b.WriteString(fmt.Sprintf("#%d", i+1))
		if r.DescB != "" && r.DescB != "B" {
			b.WriteString(fmt.Sprintf("  B = \"%s\"", r.DescB))
		}
		b.WriteString("\n")
		b.WriteString(strings.Repeat("-", 40) + "\n")

		// 参数
		b.WriteString(fmt.Sprintf("  P(A)   = %s%%\n", formatPercent(r.PriorA)))
		b.WriteString(fmt.Sprintf("  P(B|A) = %s%%,  P(B|¬A) = %s%%\n",
			formatPercent(r.LikelihoodA), formatPercent(r.LikelihoodNotA)))

		// 计算过程
		pB := r.LikelihoodA.Mul(r.PriorA).Add(r.LikelihoodNotA.Mul(one.Sub(r.PriorA)))
		b.WriteString(fmt.Sprintf("  P(B)   = P(B|A)·P(A) + P(B|¬A)·P(¬A) = %s%%\n", formatPercent(pB)))
		b.WriteString(fmt.Sprintf("  P(A|B) = P(B|A)·P(A) / P(B) = %s%%\n", formatPercent(r.Posterior)))

		b.WriteString("\n")
	}

	// 总结
	last := m.IterationHistory[len(m.IterationHistory)-1]
	b.WriteString(strings.Repeat("=", 60) + "\n")
	if m.DescA != "" && m.DescA != "A" {
		b.WriteString(fmt.Sprintf("After %d iteration(s), the probability of \"%s\" is %s%%.\n",
			len(m.IterationHistory), m.DescA, formatPercent(last.Posterior)))
	} else {
		b.WriteString(fmt.Sprintf("After %d iteration(s), P(A) = %s%%.\n",
			len(m.IterationHistory), formatPercent(last.Posterior)))
	}

	return b.String()
}
