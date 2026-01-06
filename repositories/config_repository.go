package repositories

import (
	"database/sql"
	"fmt"
)

// ConfigRepository 定义配置持久化接口
type ConfigRepository interface {
	GetConfigValue(key string) (string, error)
	SetConfigValue(key string, value string) error
}

// SQLiteConfigRepository 基于 SQLite 的实现
type SQLiteConfigRepository struct {
	db *sql.DB
}

// NewSQLiteConfigRepository 构造函数
func NewSQLiteConfigRepository(db *sql.DB) *SQLiteConfigRepository {
	return &SQLiteConfigRepository{db: db}
}

// GetConfigValue 从数据库中读取配置值
func (r *SQLiteConfigRepository) GetConfigValue(key string) (string, error) {
	var value string
	err := r.db.QueryRow("SELECT value FROM config WHERE key = ?", key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // 配置项不存在
		}
		return "", fmt.Errorf("查询配置项 %s 失败: %w", key, err)
	}
	return value, nil
}

// SetConfigValue 向数据库中写入配置值
func (r *SQLiteConfigRepository) SetConfigValue(key string, value string) error {
	query := `
		INSERT OR REPLACE INTO config (key, value) VALUES (?, ?)
	`
	_, err := r.db.Exec(query, key, value)
	if err != nil {
		return fmt.Errorf("保存配置项 %s 失败: %w", key, err)
	}
	return nil
}
