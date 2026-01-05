package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"stock-analyzer-wails/internal/logger"
	"stock-analyzer-wails/models"
	"sync"
	"time"

	"go.uber.org/zap"
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
	start := time.Now()
	r.mu.RLock()
	defer r.mu.RUnlock()

	logger.Debug("开始获取自选股列表",
		zap.String("module", "services.watchlist"),
		zap.String("op", "get_all_watchlist"),
		zap.String("file_path", r.filePath),
	)

	stocks, err := r.getAllInternal()
	if err != nil {
		logger.Error("获取自选股列表失败",
			zap.String("module", "services.watchlist"),
			zap.String("op", "get_all_watchlist"),
			zap.Error(err),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
		)
		return nil, fmt.Errorf("获取自选股列表失败: %w", err)
	}

	logger.Debug("成功获取自选股列表",
		zap.String("module", "services.watchlist"),
		zap.String("op", "get_all_watchlist"),
		zap.Int("stocks_count", len(stocks)),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
	return stocks, nil
}

func (r *FileWatchlistRepository) Add(stock *models.StockData) error {
	start := time.Now()
	r.mu.Lock()
	defer r.mu.Unlock()

	logger.Info("开始添加股票到自选股",
		zap.String("module", "services.watchlist"),
		zap.String("op", "add_to_watchlist"),
		zap.String("stock_code", stock.Code),
		zap.String("stock_name", stock.Name),
		zap.String("file_path", r.filePath),
	)

	stocks, err := r.getAllInternal()
	if err != nil {
		logger.Error("获取自选股列表失败",
			zap.String("module", "services.watchlist"),
			zap.String("op", "add_to_watchlist"),
			zap.String("stock_code", stock.Code),
			zap.Error(err),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
		)
		return fmt.Errorf("获取自选股列表失败: %w", err)
	}

	// 检查是否已存在
	for _, s := range stocks {
		if s.Code == stock.Code {
			logger.Info("股票已存在于自选股中，跳过添加",
				zap.String("module", "services.watchlist"),
				zap.String("op", "add_to_watchlist"),
				zap.String("stock_code", stock.Code),
				zap.Int64("duration_ms", time.Since(start).Milliseconds()),
			)
			return nil // 已存在则不重复添加
		}
	}

	stocks = append(stocks, stock)
	if err := r.saveInternal(stocks); err != nil {
		logger.Error("保存自选股列表失败",
			zap.String("module", "services.watchlist"),
			zap.String("op", "add_to_watchlist"),
			zap.String("stock_code", stock.Code),
			zap.Int("total_stocks", len(stocks)),
			zap.Error(err),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
		)
		return fmt.Errorf("保存自选股列表失败: %w", err)
	}

	logger.Info("成功添加股票到自选股",
		zap.String("module", "services.watchlist"),
		zap.String("op", "add_to_watchlist"),
		zap.String("stock_code", stock.Code),
		zap.String("stock_name", stock.Name),
		zap.Int("total_stocks", len(stocks)),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
	return nil
}

func (r *FileWatchlistRepository) Remove(code string) error {
	start := time.Now()
	r.mu.Lock()
	defer r.mu.Unlock()

	logger.Info("开始从自选股移除股票",
		zap.String("module", "services.watchlist"),
		zap.String("op", "remove_from_watchlist"),
		zap.String("stock_code", code),
		zap.String("file_path", r.filePath),
	)

	stocks, err := r.getAllInternal()
	if err != nil {
		logger.Error("获取自选股列表失败",
			zap.String("module", "services.watchlist"),
			zap.String("op", "remove_from_watchlist"),
			zap.String("stock_code", code),
			zap.Error(err),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
		)
		return fmt.Errorf("获取自选股列表失败: %w", err)
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
		logger.Warn("要移除的股票不存在于自选股中",
			zap.String("module", "services.watchlist"),
			zap.String("op", "remove_from_watchlist"),
			zap.String("stock_code", code),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
		)
		return nil // 不存在也算成功
	}

	if err := r.saveInternal(newStocks); err != nil {
		logger.Error("保存自选股列表失败",
			zap.String("module", "services.watchlist"),
			zap.String("op", "remove_from_watchlist"),
			zap.String("stock_code", code),
			zap.Int("remaining_stocks", len(newStocks)),
			zap.Error(err),
			zap.Int64("duration_ms", time.Since(start).Milliseconds()),
		)
		return fmt.Errorf("保存自选股列表失败: %w", err)
	}

	logger.Info("成功从自选股移除股票",
		zap.String("module", "services.watchlist"),
		zap.String("op", "remove_from_watchlist"),
		zap.String("stock_code", code),
		zap.Int("remaining_stocks", len(newStocks)),
		zap.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
	return nil
}

func (r *FileWatchlistRepository) getAllInternal() ([]*models.StockData, error) {
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		return []*models.StockData{}, nil
	}
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		logger.Error("读取自选股文件失败",
			zap.String("module", "services.watchlist"),
			zap.String("op", "read_file"),
			zap.String("file_path", r.filePath),
			zap.Error(err),
		)
		return nil, fmt.Errorf("读取自选股文件失败: %w", err)
	}

	var stocks []*models.StockData
	if err := json.Unmarshal(data, &stocks); err != nil {
		logger.Error("解析自选股文件失败，尝试备份并重置",
			zap.String("module", "services.watchlist"),
			zap.String("op", "parse_json"),
			zap.String("file_path", r.filePath),
			zap.Int("file_size", len(data)),
			zap.Error(err),
		)

		// 备份损坏的文件
		if backupErr := r.backupCorruptedFile(); backupErr != nil {
			logger.Warn("备份损坏文件失败",
				zap.String("module", "services.watchlist"),
				zap.Error(backupErr),
			)
		}

		// 返回空列表，让系统继续工作
		return []*models.StockData{}, nil
	}
	return stocks, nil
}

func (r *FileWatchlistRepository) saveInternal(stocks []*models.StockData) error {
	data, err := json.MarshalIndent(stocks, "", "  ")
	if err != nil {
		logger.Error("序列化自选股数据失败",
			zap.String("module", "services.watchlist"),
			zap.String("op", "marshal_json"),
			zap.Int("stocks_count", len(stocks)),
			zap.Error(err),
		)
		return fmt.Errorf("序列化自选股数据失败: %w", err)
	}

	// 检查磁盘空间和文件权限
	if err := r.checkFileWritePermission(); err != nil {
		logger.Error("文件写入权限检查失败",
			zap.String("module", "services.watchlist"),
			zap.String("op", "check_permission"),
			zap.String("file_path", r.filePath),
			zap.Error(err),
		)
		return fmt.Errorf("文件写入权限不足: %w", err)
	}

	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		logger.Error("写入自选股文件失败",
			zap.String("module", "services.watchlist"),
			zap.String("op", "write_file"),
			zap.String("file_path", r.filePath),
			zap.Int("data_size", len(data)),
			zap.Error(err),
		)
		return fmt.Errorf("写入自选股文件失败: %w", err)
	}

	return nil
}

// backupCorruptedFile 备份损坏的自选股文件
func (r *FileWatchlistRepository) backupCorruptedFile() error {
	backupPath := r.filePath + ".backup." + time.Now().Format("20060102_150405")
	if err := os.Rename(r.filePath, backupPath); err != nil {
		return fmt.Errorf("备份损坏文件失败: %w", err)
	}

	logger.Info("已备份损坏的自选股文件",
		zap.String("module", "services.watchlist"),
		zap.String("op", "backup_corrupted_file"),
		zap.String("original_path", r.filePath),
		zap.String("backup_path", backupPath),
	)
	return nil
}

// checkFileWritePermission 检查文件写入权限和磁盘空间
func (r *FileWatchlistRepository) checkFileWritePermission() error {
	// 检查目录是否存在，不存在则创建
	dir := filepath.Dir(r.filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录失败: %w", err)
		}
	}

	// 尝试创建临时文件来检查权限
	tempFile := r.filePath + ".tmp"
	if err := os.WriteFile(tempFile, []byte("test"), 0644); err != nil {
		return fmt.Errorf("写入权限检查失败: %w", err)
	}

	// 清理临时文件
	os.Remove(tempFile)
	return nil
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
