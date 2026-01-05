package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"stock-analyzer-wails/models"
	"strings"
	"sync"
	"time"
)

type AlertStorage struct {
	baseDir string
	mu      sync.RWMutex
}

func NewAlertStorage() (*AlertStorage, error) {
	// 获取用户主目录下的存储路径
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	baseDir := filepath.Join(home, ".stock-analyzer", "alerts")
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}

	return &AlertStorage{
		baseDir: baseDir,
	}, nil
}

// SaveAlert 保存告警记录到当月文件
func (s *AlertStorage) SaveAlert(alert *models.PriceAlert, advice string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	fileName := fmt.Sprintf("alerts_%d_%02d.jsonl", now.Year(), now.Month())
	filePath := filepath.Join(s.baseDir, fileName)

	record := map[string]interface{}{
		"timestamp": now.Format(time.RFC3339),
		"stockCode": alert.StockCode,
		"stockName": alert.StockName,
		"type":      alert.Type,
		"price":     alert.Price,
		"label":     alert.Label,
		"role":      alert.Role,
		"advice":    advice,
	}

	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(append(data, '\n'))
	return err
}

// SaveActiveAlerts 保存当前活跃的预警订阅
func (s *AlertStorage) SaveActiveAlerts(alerts []*models.PriceAlert) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	filePath := filepath.Join(s.baseDir, "active_alerts.json")
	data, err := json.MarshalIndent(alerts, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

// LoadActiveAlerts 加载保存的活跃预警订阅
func (s *AlertStorage) LoadActiveAlerts() ([]*models.PriceAlert, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	filePath := filepath.Join(s.baseDir, "active_alerts.json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []*models.PriceAlert{}, nil
	}
	
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	
	var alerts []*models.PriceAlert
	if err := json.Unmarshal(data, &alerts); err != nil {
		return nil, err
	}
	return alerts, nil
}

// GetAlertHistory 获取告警历史，支持分页和股票代码筛选
func (s *AlertStorage) GetAlertHistory(stockCode string, limit int) ([]map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	files, err := filepath.Glob(filepath.Join(s.baseDir, "alerts_*.jsonl"))
	if err != nil {
		return nil, err
	}

	// 按文件名倒序排列（从最近的月份开始读）
	sort.Sort(sort.Reverse(sort.StringSlice(files)))

	var history []map[string]interface{}
	for _, file := range files {
		if len(history) >= limit {
			break
		}

		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		lines := strings.Split(string(data), "\n")
		// 从后往前读，获取最新记录
		for i := len(lines) - 1; i >= 0; i-- {
			line := strings.TrimSpace(lines[i])
			if line == "" {
				continue
			}

			var record map[string]interface{}
			if err := json.Unmarshal([]byte(line), &record); err != nil {
				continue
			}

			if stockCode == "" || record["stockCode"] == stockCode {
				history = append(history, record)
				if len(history) >= limit {
					break
				}
			}
		}
	}

	return history, nil
}
