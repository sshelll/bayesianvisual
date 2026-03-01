package model

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// ExportData 导出数据结构
type ExportData struct {
	Version          string            `json:"version"`
	IterationHistory []IterationRecord `json:"iteration_history"`
}

// expandPath 展开路径中的 ~ 为用户 home 目录
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}

// ExportToJSON 导出迭代历史到 JSON 文件
func (m *Model) ExportToJSON(path string) error {
	data := ExportData{
		Version:          "1.0",
		IterationHistory: m.IterationHistory,
	}

	file, err := os.Create(expandPath(path))
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// LoadFromJSON 从 JSON 文件加载迭代历史
func (m *Model) LoadFromJSON(path string) error {
	file, err := os.Open(expandPath(path))
	if err != nil {
		return err
	}
	defer file.Close()

	var data ExportData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return err
	}

	// 加载历史记录
	m.IterationHistory = data.IterationHistory

	// 如果有历史记录，使用最后一条记录的数据作为当前状态
	if len(m.IterationHistory) > 0 {
		lastRecord := m.IterationHistory[len(m.IterationHistory)-1]
		m.PriorA = lastRecord.PriorA
		m.LikelihoodA = lastRecord.LikelihoodA
		m.LikelihoodNotA = lastRecord.LikelihoodNotA
		m.DescA = lastRecord.DescA
		m.DescB = lastRecord.DescB
	}

	return nil
}
