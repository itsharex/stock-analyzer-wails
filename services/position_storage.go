package services

import (
	"encoding/json"
	"os"
	"path/filepath"
	"stock-analyzer-wails/models"
	"sync"
	"time"

	"stock-analyzer-wails/internal/logger"
	"go.uber.org/zap"
)

type PositionStorageService struct {
	storageDir string
	filePath   string
	mu         sync.RWMutex
}

func NewPositionStorageService() *PositionStorageService {
	home, _ := os.UserHomeDir()
	storageDir := filepath.Join(home, ".stock-analyzer", "positions")
	os.MkdirAll(storageDir, 0755)

	return &PositionStorageService{
		storageDir: storageDir,
		filePath:   filepath.Join(storageDir, "active_positions.json"),
	}
}

func (s *PositionStorageService) SavePosition(pos *models.Position) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	positions, err := s.loadPositions()
	if err != nil {
		positions = make(map[string]*models.Position)
	}

	pos.UpdatedAt = time.Now()
	positions[pos.StockCode] = pos

	return s.saveToFile(positions)
}

func (s *PositionStorageService) GetPositions() (map[string]*models.Position, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.loadPositions()
}

func (s *PositionStorageService) RemovePosition(code string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	positions, err := s.loadPositions()
	if err != nil {
		return err
	}

	delete(positions, code)
	return s.saveToFile(positions)
}

func (s *PositionStorageService) loadPositions() (map[string]*models.Position, error) {
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		return make(map[string]*models.Position), nil
	}

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, err
	}

	var positions map[string]*models.Position
	if err := json.Unmarshal(data, &positions); err != nil {
		return nil, err
	}

	return positions, nil
}

func (s *PositionStorageService) saveToFile(positions map[string]*models.Position) error {
	data, err := json.MarshalIndent(positions, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(s.filePath, data, 0644)
	if err != nil {
		logger.Error("保存持仓数据失败", zap.Error(err))
		return err
	}

	return nil
}
