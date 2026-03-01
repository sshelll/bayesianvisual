package bayesian

// Calculator 贝叶斯计算器
type Calculator struct {
	PriorA         float64 // P(A) 先验概率
	LikelihoodA    float64 // P(B|A) 似然概率
	LikelihoodNotA float64 // P(B|¬A) 似然概率
}

// CalculatePosterior 计算后验概率 P(A|B)
// 使用贝叶斯定理: P(A|B) = P(B|A) * P(A) / P(B)
// 其中 P(B) = P(B|A) * P(A) + P(B|¬A) * P(¬A)
func (c *Calculator) CalculatePosterior() float64 {
	// P(B) = P(B|A) * P(A) + P(B|¬A) * P(¬A)
	pB := c.LikelihoodA*c.PriorA + c.LikelihoodNotA*(1-c.PriorA)

	// 避免除以零
	if pB == 0 {
		return 0
	}

	// P(A|B) = P(B|A) * P(A) / P(B)
	return (c.LikelihoodA * c.PriorA) / pB
}
