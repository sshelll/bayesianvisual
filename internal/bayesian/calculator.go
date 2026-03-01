package bayesian

import "github.com/shopspring/decimal"

// Calculator 贝叶斯计算器
type Calculator struct {
	PriorA         decimal.Decimal // P(A) 先验概率
	LikelihoodA    decimal.Decimal // P(B|A) 似然概率
	LikelihoodNotA decimal.Decimal // P(B|¬A) 似然概率
}

// CalculatePosterior 计算后验概率 P(A|B)
// 使用贝叶斯定理: P(A|B) = P(B|A) * P(A) / P(B)
// 其中 P(B) = P(B|A) * P(A) + P(B|¬A) * P(¬A)
func (c *Calculator) CalculatePosterior() decimal.Decimal {
	one := decimal.NewFromInt(1)

	// P(B) = P(B|A) * P(A) + P(B|¬A) * P(¬A)
	pB := c.LikelihoodA.Mul(c.PriorA).Add(c.LikelihoodNotA.Mul(one.Sub(c.PriorA)))

	// 避免除以零
	if pB.IsZero() {
		return decimal.Zero
	}

	// P(A|B) = P(B|A) * P(A) / P(B)
	return c.LikelihoodA.Mul(c.PriorA).Div(pB)
}
