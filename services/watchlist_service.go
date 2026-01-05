package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"stock-analyzer-wails/models"
	"sync"
)

// WatchlistRepository 定义自选股持久化接口
type WatchlistRepository interface {
	Add(stock *models.StockData) error
	Remove(code string) error
	GetAll() ([]*models.StockData, error)
}

// FileWatchlistRepository 基于 JSON 文件的实现
type FileWatchlistRepository struct {
	filePath string
	mu       sync.RWMutex
}

func NewFileWatchlistRepository() (*FileWatchlistRepository, error) {
	cwd, _ := os.Getwd()
	path := filepath.Join(cwd, "watchlist.json")
	return &FileWatchlistRepository{
		filePath: path,
	}, nil
}

func (r *FileWatchlistRepository) GetAll() ([]*models.StockData, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		return []*models.StockData{}, nil
	}

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}

	var stocks []*models.StockData
	if err := json.Unmarshal(data, &stocks); err != nil {
		return nil, err
	}

	return stocks, nil
}

func (r *FileWatchlistRepository) Add(stock *models.StockData) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	stocks, err := r.getAllInternal()
	if err != nil {
		return err
	}

	// 检查是否已存在
	for _, s := range stocks {
		if s.Code == stock.Code {
			return fmt.Errorf("股票 %s 已在自选股中", stock.Code)
		}
	}

	stocks = append(stocks, stock)
	return r.saveInternal(stocks)
}

func (r *FileWatchlistRepository) Remove(code string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	stocks, err := r.getAllInternal()
	if err != nil {
		return err
	}

	newStocks := make([]*models.StockData, 0)
	found := false
	for _, s := range stocks {
		if s.Code == code {
			found = true
			continue
		}
		newStocks = append(newStocks, s)
	}

	if !found {
		return fmt.Errorf("未找到股票代码: %s", code)
	}

	return r.saveInternal(newStocks)
}

// 内部辅助方法（不带锁）
func (r *FileWatchlistRepository) getAllInternal() ([]*models.StockData, error) {
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		return []*models.StockData{}, nil
	}
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}
	var stocks []*models.StockData
	json.Unmarshal(data, &stocks)
	return stocks, nil
}

func (r *FileWatchlistRepository) saveInternal(stocks []*models.StockData) error {
	data, err := json.MarshalIndent(stocks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath, data, 0644)
}

// WatchlistService 业务逻辑层
type WatchlistService struct {
	repo WatchlistRepository
}

func NewWatchlistService() (*WatchlistService, error) {
	// 目前默认使用文件存储，后续可轻松切换为 MongoDB 实现
	repo, err := NewFileWatchlistRepository()
	if err != nil {
		return nil, err
	}
	return &WatchlistService{repo: repo}, nil
}

func (s *WatchlistService) AddToWatchlist(stock *models.StockData) error {
	return s.repo.Add(stock)
}

func (s *WatchlistService) RemoveFromWatchlist(code string) error {
	return s.repo.Remove(code)
}

func (s *WatchlistService) GetWatchlist() ([]*models.StockData, error) {
	return s.repo.GetAll()
}
