package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"stock-analyzer-wails/models"
	"sync"
	"time"
)

type AnalysisCacheService struct {
	cachePath string
	cache     models.AnalysisCache
	mutex     sync.RWMutex
}

func NewAnalysisCacheService() (*AnalysisCacheService, error) {
	homeDir, _ := os.UserHomeDir()
	cacheDir := filepath.Join(homeDir, ".stock-analyzer")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, err
	}

	cachePath := filepath.Join(cacheDir, "analysis_cache.json")
	s := &AnalysisCacheService{
		cachePath: cachePath,
		cache:     models.AnalysisCache{Entries: make(map[string]models.AnalysisCacheEntry)},
	}

	s.load()
	return s, nil
}

func (s *AnalysisCacheService) load() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	data, err := os.ReadFile(s.cachePath)
	if err != nil {
		return
	}

	json.Unmarshal(data, &s.cache)
}

func (s *AnalysisCacheService) save() error {
	data, err := json.MarshalIndent(s.cache, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.cachePath, data, 0644)
}

func (s *AnalysisCacheService) Get(code, role, period string) (*models.TechnicalAnalysisResult, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	key := fmt.Sprintf("%s_%s_%s", code, role, period)
	entry, ok := s.cache.Entries[key]
	if !ok {
		return nil, false
	}

	// 缓存有效期：4小时
	if time.Since(entry.Timestamp) > 4*time.Hour {
		return nil, false
	}

	return &entry.Result, true
}

func (s *AnalysisCacheService) Set(code, role, period string, result models.TechnicalAnalysisResult) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	key := fmt.Sprintf("%s_%s_%s", code, role, period)
	s.cache.Entries[key] = models.AnalysisCacheEntry{
		StockCode: code,
		Role:      role,
		Period:    period,
		Result:    result,
		Timestamp: time.Now(),
	}

	return s.save()
}
