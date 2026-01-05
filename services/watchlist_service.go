package services

import (
	"encoding/json"
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
	// 使用跨平台的应用数据目录
	path := filepath.Join(GetAppDataDir(), "watchlist.json")
	
	// 迁移逻辑：如果当前目录下有旧文件，移动到新位置
	if _, err := os.Stat("watchlist.json"); err == nil {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.Rename("watchlist.json", path)
		}
	}

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

	for _, s := range stocks {
		if s.Code == stock.Code {
			return nil // 已存在则不重复添加
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
	for _, s := range stocks {
		if s.Code == code {
			continue
		}
		newStocks = append(newStocks, s)
	}

	return r.saveInternal(newStocks)
}

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
