package repositories

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"stock-analyzer-wails/models"

	"gorm.io/gorm"
)

// PriceAlertCondition 价格预警条件结构
type PriceAlertCondition struct {
	Field     string  `json:"field"`     // 字段名: price_change_percent, close_price, volume_ratio等
	Operator  string  `json:"operator"`  // 操作符: >, <, >=, <=, ==, !=
	Value     float64 `json:"value"`     // 比较值
	Reference string  `json:"reference"` // 引用值: historical_high, historical_low, ma5, ma20等
}

// PriceAlertConditions 预警条件列表
type PriceAlertConditions struct {
	Logic      string                `json:"logic"`      // 逻辑关系: AND, OR
	Conditions []PriceAlertCondition `json:"conditions"` // 条件列表
}

// PriceThresholdAlert 价格预警配置（对应数据库表）
type PriceThresholdAlert struct {
	ID                int64     `json:"id"`
	StockCode         string    `json:"stockCode"`
	StockName         string    `json:"stockName"`
	AlertType         string    `json:"alertType"`
	Conditions        string    `json:"conditions"` // JSON字符串
	IsActive          bool      `json:"isActive"`
	Sensitivity       float64   `json:"sensitivity"`
	CooldownHours     int       `json:"cooldownHours"`
	PostTriggerAction string    `json:"postTriggerAction"`
	EnableSound       bool      `json:"enableSound"`
	EnableDesktop     bool      `json:"enableDesktop"`
	TemplateID        string    `json:"templateId"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	LastTriggeredAt   time.Time `json:"lastTriggeredAt"`
}

// PriceAlertTemplate 预警模板
type PriceAlertTemplate struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AlertType   string    `json:"alertType"`
	Conditions  string    `json:"conditions"` // JSON字符串
	CreatedAt   time.Time `json:"createdAt"`
}

// PriceAlertTriggerHistory 预警触发历史
type PriceAlertTriggerHistory struct {
	ID             int64     `json:"id"`
	AlertID        int64     `json:"alertId"`
	StockCode      string    `json:"stockCode"`
	StockName      string    `json:"stockName"`
	AlertType      string    `json:"alertType"`
	TriggerPrice   float64   `json:"triggerPrice"`
	TriggerMessage string    `json:"triggerMessage"`
	TriggeredAt    time.Time `json:"triggeredAt"`
}

// PriceAlertRepository 价格预警数据访问层
type PriceAlertRepository struct {
	db *gorm.DB
}

// NewPriceAlertRepository 创建价格预警Repository
func NewPriceAlertRepository(db *gorm.DB) *PriceAlertRepository {
	return &PriceAlertRepository{db: db}
}

// CreateAlert 创建价格预警
func (r *PriceAlertRepository) CreateAlert(alert *PriceThresholdAlert) error {
	entity := models.PriceThresholdAlertEntity{
		StockCode:         alert.StockCode,
		StockName:         alert.StockName,
		AlertType:         alert.AlertType,
		Conditions:        alert.Conditions,
		IsActive:          alert.IsActive,
		Sensitivity:       alert.Sensitivity,
		CooldownHours:     alert.CooldownHours,
		PostTriggerAction: alert.PostTriggerAction,
		EnableSound:       alert.EnableSound,
		EnableDesktop:     alert.EnableDesktop,
		TemplateID:        alert.TemplateID,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := r.db.Create(&entity).Error; err != nil {
		return fmt.Errorf("创建预警失败: %w", err)
	}

	alert.ID = int64(entity.ID)
	return nil
}

// GetAlertByID 根据ID获取预警
func (r *PriceAlertRepository) GetAlertByID(id int64) (*PriceThresholdAlert, error) {
	var entity models.PriceThresholdAlertEntity
	if err := r.db.First(&entity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("预警不存在")
		}
		return nil, fmt.Errorf("查询预警失败: %w", err)
	}

	return r.entityToAlert(&entity), nil
}

// normalizePriceAlertConditions 兼容历史数据
func normalizePriceAlertConditions(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return raw
	}

	var cond PriceAlertConditions
	if json.Unmarshal([]byte(raw), &cond) == nil {
		return raw
	}

	var inner string
	if json.Unmarshal([]byte(raw), &inner) == nil {
		inner = strings.TrimSpace(inner)
		if inner != "" {
			if json.Unmarshal([]byte(inner), &cond) == nil {
				return inner
			}
			return inner
		}
	}

	return raw
}

// GetAllAlerts 获取所有预警
func (r *PriceAlertRepository) GetAllAlerts() ([]*PriceThresholdAlert, error) {
	var entities []models.PriceThresholdAlertEntity
	if err := r.db.Order("created_at DESC").Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("查询所有预警失败: %w", err)
	}

	var alerts []*PriceThresholdAlert
	for _, entity := range entities {
		alerts = append(alerts, r.entityToAlert(&entity))
	}
	return alerts, nil
}

// GetActiveAlerts 获取所有活跃的预警
func (r *PriceAlertRepository) GetActiveAlerts() ([]*PriceThresholdAlert, error) {
	var entities []models.PriceThresholdAlertEntity
	if err := r.db.Where("is_active = ?", true).Order("created_at DESC").Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("查询活跃预警失败: %w", err)
	}

	var alerts []*PriceThresholdAlert
	for _, entity := range entities {
		alerts = append(alerts, r.entityToAlert(&entity))
	}
	return alerts, nil
}

// GetAlertsByStockCode 根据股票代码获取预警
func (r *PriceAlertRepository) GetAlertsByStockCode(stockCode string) ([]*PriceThresholdAlert, error) {
	var entities []models.PriceThresholdAlertEntity
	if err := r.db.Where("stock_code = ?", stockCode).Order("created_at DESC").Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("查询股票预警失败: %w", err)
	}

	var alerts []*PriceThresholdAlert
	for _, entity := range entities {
		alerts = append(alerts, r.entityToAlert(&entity))
	}
	return alerts, nil
}

// UpdateAlert 更新预警
func (r *PriceAlertRepository) UpdateAlert(alert *PriceThresholdAlert) error {
	updates := map[string]interface{}{
		"stock_code":          alert.StockCode,
		"stock_name":          alert.StockName,
		"alert_type":          alert.AlertType,
		"conditions":          alert.Conditions,
		"is_active":           alert.IsActive,
		"sensitivity":         alert.Sensitivity,
		"cooldown_hours":      alert.CooldownHours,
		"post_trigger_action": alert.PostTriggerAction,
		"enable_sound":        alert.EnableSound,
		"enable_desktop":      alert.EnableDesktop,
		"template_id":         alert.TemplateID,
		"updated_at":          time.Now(),
	}

	if err := r.db.Model(&models.PriceThresholdAlertEntity{}).Where("id = ?", alert.ID).Updates(updates).Error; err != nil {
		return fmt.Errorf("更新预警失败: %w", err)
	}
	return nil
}

// DeleteAlert 删除预警
func (r *PriceAlertRepository) DeleteAlert(id int64) error {
	if err := r.db.Delete(&models.PriceThresholdAlertEntity{}, id).Error; err != nil {
		return fmt.Errorf("删除预警失败: %w", err)
	}
	return nil
}

// ToggleAlertStatus 切换预警状态
func (r *PriceAlertRepository) ToggleAlertStatus(id int64, isActive bool) error {
	if err := r.db.Model(&models.PriceThresholdAlertEntity{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_active":  isActive,
			"updated_at": time.Now(),
		}).Error; err != nil {
		return fmt.Errorf("切换预警状态失败: %w", err)
	}
	return nil
}

// UpdateLastTriggeredTime 更新最后触发时间
func (r *PriceAlertRepository) UpdateLastTriggeredTime(id int64) error {
	now := time.Now()
	if err := r.db.Model(&models.PriceThresholdAlertEntity{}).Where("id = ?", id).
		Update("last_triggered_at", now).Error; err != nil {
		return fmt.Errorf("更新最后触发时间失败: %w", err)
	}
	return nil
}

// IsInCooldown 检查是否在冷却时间内
func (r *PriceAlertRepository) IsInCooldown(id int64, cooldownHours int) (bool, error) {
	alert, err := r.GetAlertByID(id)
	if err != nil {
		return false, err
	}

	if alert.LastTriggeredAt.IsZero() {
		return false, nil
	}

	cooldownDuration := time.Duration(cooldownHours) * time.Hour
	elapsed := time.Since(alert.LastTriggeredAt)

	return elapsed < cooldownDuration, nil
}

// SaveTriggerHistory 保存预警触发历史
func (r *PriceAlertRepository) SaveTriggerHistory(history *PriceAlertTriggerHistory) error {
	entity := models.PriceAlertTriggerHistoryEntity{
		AlertID:        uint(history.AlertID),
		StockCode:      history.StockCode,
		StockName:      history.StockName,
		AlertType:      history.AlertType,
		TriggerPrice:   history.TriggerPrice,
		TriggerMessage: history.TriggerMessage,
		TriggeredAt:    time.Now(),
	}

	if err := r.db.Create(&entity).Error; err != nil {
		return fmt.Errorf("保存触发历史失败: %w", err)
	}
	return nil
}

// GetTriggerHistory 获取预警触发历史
func (r *PriceAlertRepository) GetTriggerHistory(stockCode string, limit int) ([]*PriceAlertTriggerHistory, error) {
	var entities []models.PriceAlertTriggerHistoryEntity
	query := r.db.Order("triggered_at DESC").Limit(limit)

	if stockCode != "" {
		query = query.Where("stock_code = ?", stockCode)
	}

	if err := query.Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("查询触发历史失败: %w", err)
	}

	var histories []*PriceAlertTriggerHistory
	for _, entity := range entities {
		histories = append(histories, &PriceAlertTriggerHistory{
			ID:             int64(entity.ID),
			AlertID:        int64(entity.AlertID),
			StockCode:      entity.StockCode,
			StockName:      entity.StockName,
			AlertType:      entity.AlertType,
			TriggerPrice:   entity.TriggerPrice,
			TriggerMessage: entity.TriggerMessage,
			TriggeredAt:    entity.TriggeredAt,
		})
	}

	return histories, nil
}

// GetAllTemplates 获取所有预警模板
func (r *PriceAlertRepository) GetAllTemplates() ([]*PriceAlertTemplate, error) {
	var entities []models.PriceAlertTemplateEntity
	if err := r.db.Order("created_at ASC").Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("查询预警模板失败: %w", err)
	}

	var templates []*PriceAlertTemplate
	for _, entity := range entities {
		templates = append(templates, &PriceAlertTemplate{
			ID:          entity.ID,
			Name:        entity.Name,
			Description: entity.Description,
			AlertType:   entity.AlertType,
			Conditions:  entity.Conditions,
			CreatedAt:   entity.CreatedAt,
		})
	}
	return templates, nil
}

// GetTemplateByID 根据ID获取预警模板
func (r *PriceAlertRepository) GetTemplateByID(id string) (*PriceAlertTemplate, error) {
	var entity models.PriceAlertTemplateEntity
	if err := r.db.First(&entity, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("模板不存在")
		}
		return nil, fmt.Errorf("查询模板失败: %w", err)
	}

	return &PriceAlertTemplate{
		ID:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
		AlertType:   entity.AlertType,
		Conditions:  entity.Conditions,
		CreatedAt:   entity.CreatedAt,
	}, nil
}

// CreateTemplate 创建预警模板
func (r *PriceAlertRepository) CreateTemplate(template *PriceAlertTemplate) error {
	entity := models.PriceAlertTemplateEntity{
		ID:          template.ID,
		Name:        template.Name,
		Description: template.Description,
		AlertType:   template.AlertType,
		Conditions:  template.Conditions,
		CreatedAt:   time.Now(),
	}

	if err := r.db.Create(&entity).Error; err != nil {
		return fmt.Errorf("创建模板失败: %w", err)
	}
	return nil
}

// DeleteTemplate 删除预警模板
func (r *PriceAlertRepository) DeleteTemplate(id string) error {
	if err := r.db.Delete(&models.PriceAlertTemplateEntity{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("删除模板失败: %w", err)
	}
	return nil
}

// entityToAlert 辅助方法：Entity转Alert
func (r *PriceAlertRepository) entityToAlert(entity *models.PriceThresholdAlertEntity) *PriceThresholdAlert {
	alert := &PriceThresholdAlert{
		ID:                int64(entity.ID),
		StockCode:         entity.StockCode,
		StockName:         entity.StockName,
		AlertType:         entity.AlertType,
		Conditions:        normalizePriceAlertConditions(entity.Conditions),
		IsActive:          entity.IsActive,
		Sensitivity:       entity.Sensitivity,
		CooldownHours:     entity.CooldownHours,
		PostTriggerAction: entity.PostTriggerAction,
		EnableSound:       entity.EnableSound,
		EnableDesktop:     entity.EnableDesktop,
		TemplateID:        entity.TemplateID,
		CreatedAt:         entity.CreatedAt,
		UpdatedAt:         entity.UpdatedAt,
		LastTriggeredAt:   entity.LastTriggeredAt,
	}

	return alert
}

// ToModelsPriceAlert 转换为models.PriceAlert（兼容旧的Alert类型）
func (r *PriceThresholdAlert) ToModelsPriceAlert() *models.PriceAlert {
	return &models.PriceAlert{
		StockCode:     r.StockCode,
		StockName:     r.StockName,
		Type:          r.AlertType,
		Price:         0, // 价格预警不支持单个价格
		Label:         r.AlertType,
		IsActive:      r.IsActive,
		LastTriggered: r.LastTriggeredAt,
	}
}
