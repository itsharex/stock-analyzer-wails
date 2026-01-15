package repositories

import (
	"fmt"
	"stock-analyzer-wails/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ConfigRepository 定义配置持久化接口
type ConfigRepository interface {
	GetConfigValue(key string) (string, error)
	SetConfigValue(key string, value string) error
}

// SQLiteConfigRepository 基于 SQLite 的实现
type SQLiteConfigRepository struct {
	db *gorm.DB
}

// NewSQLiteConfigRepository 构造函数
func NewSQLiteConfigRepository(db *gorm.DB) *SQLiteConfigRepository {
	return &SQLiteConfigRepository{db: db}
}

// GetConfigValue 从数据库中读取配置值
func (r *SQLiteConfigRepository) GetConfigValue(key string) (string, error) {
	var entity models.ConfigEntity
	if err := r.db.First(&entity, "key = ?", key).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil // 配置项不存在
		}
		return "", fmt.Errorf("查询配置项 %s 失败: %w", key, err)
	}
	return entity.Value, nil
}

// SetConfigValue 向数据库中写入配置值
func (r *SQLiteConfigRepository) SetConfigValue(key string, value string) error {
	entity := models.ConfigEntity{
		Key:   key,
		Value: value,
	}

	if err := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&entity).Error; err != nil {
		return fmt.Errorf("保存配置项 %s 失败: %w", key, err)
	}
	return nil
}
