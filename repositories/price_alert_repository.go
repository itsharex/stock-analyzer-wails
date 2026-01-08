package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"stock-analyzer-wails/models"

	"go.uber.org/zap"
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
	db *sql.DB
}

// NewPriceAlertRepository 创建价格预警Repository
func NewPriceAlertRepository(db *sql.DB) *PriceAlertRepository {
	return &PriceAlertRepository{db: db}
}

// CreateAlert 创建价格预警
func (r *PriceAlertRepository) CreateAlert(alert *PriceThresholdAlert) error {
	conditionsJSON, err := json.Marshal(alert.Conditions)
	if err != nil {
		return fmt.Errorf("序列化条件失败: %w", err)
	}

	result, err := r.db.Exec(`
		INSERT INTO price_threshold_alerts (
			stock_code, stock_name, alert_type, conditions, is_active,
			sensitivity, cooldown_hours, post_trigger_action, enable_sound,
			enable_desktop, template_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, alert.StockCode, alert.StockName, alert.AlertType, conditionsJSON,
		alert.IsActive, alert.Sensitivity, alert.CooldownHours, alert.PostTriggerAction,
		alert.EnableSound, alert.EnableDesktop, alert.TemplateID)

	if err != nil {
		return fmt.Errorf("创建预警失败: %w", err)
	}

	id, _ := result.LastInsertId()
	alert.ID = id
	return nil
}

// GetAlertByID 根据ID获取预警
func (r *PriceAlertRepository) GetAlertByID(id int64) (*PriceThresholdAlert, error) {
	alert := &PriceThresholdAlert{}

	var conditions sql.NullString
	var lastTriggeredAt sql.NullTime

	err := r.db.QueryRow(`
		SELECT id, stock_code, stock_name, alert_type, conditions, is_active,
			   sensitivity, cooldown_hours, post_trigger_action, enable_sound,
			   enable_desktop, template_id, created_at, updated_at, last_triggered_at
		FROM price_threshold_alerts WHERE id = ?
	`, id).Scan(
		&alert.ID, &alert.StockCode, &alert.StockName, &alert.AlertType,
		&conditions, &alert.IsActive, &alert.Sensitivity, &alert.CooldownHours,
		&alert.PostTriggerAction, &alert.EnableSound, &alert.EnableDesktop,
		&alert.TemplateID, &alert.CreatedAt, &alert.UpdatedAt, &lastTriggeredAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("预警不存在")
	}
	if err != nil {
		return nil, fmt.Errorf("查询预警失败: %w", err)
	}

	if conditions.Valid {
		alert.Conditions = conditions.String
	}

	if lastTriggeredAt.Valid {
		alert.LastTriggeredAt = lastTriggeredAt.Time
	}

	return alert, nil
}

// GetAllAlerts 获取所有预警
func (r *PriceAlertRepository) GetAllAlerts() ([]*PriceThresholdAlert, error) {
	query := `
		SELECT id, stock_code, stock_name, alert_type, conditions, is_active,
			   sensitivity, cooldown_hours, post_trigger_action, enable_sound,
			   enable_desktop, template_id, created_at, updated_at, last_triggered_at
		FROM price_threshold_alerts ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询所有预警失败: %w", err)
	}
	defer rows.Close()

	var alerts []*PriceThresholdAlert
	for rows.Next() {
		alert, err := r.scanAlert(rows)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}

	return alerts, nil
}

// GetActiveAlerts 获取所有活跃的预警
func (r *PriceAlertRepository) GetActiveAlerts() ([]*PriceThresholdAlert, error) {
	query := `
		SELECT id, stock_code, stock_name, alert_type, conditions, is_active,
			   sensitivity, cooldown_hours, post_trigger_action, enable_sound,
			   enable_desktop, template_id, created_at, updated_at, last_triggered_at
		FROM price_threshold_alerts WHERE is_active = 1 ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询活跃预警失败: %w", err)
	}
	defer rows.Close()

	var alerts []*PriceThresholdAlert
	for rows.Next() {
		alert, err := r.scanAlert(rows)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}

	return alerts, nil
}

// GetAlertsByStockCode 根据股票代码获取预警
func (r *PriceAlertRepository) GetAlertsByStockCode(stockCode string) ([]*PriceThresholdAlert, error) {
	query := `
		SELECT id, stock_code, stock_name, alert_type, conditions, is_active,
			   sensitivity, cooldown_hours, post_trigger_action, enable_sound,
			   enable_desktop, template_id, created_at, updated_at, last_triggered_at
		FROM price_threshold_alerts WHERE stock_code = ? ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, stockCode)
	if err != nil {
		return nil, fmt.Errorf("查询股票预警失败: %w", err)
	}
	defer rows.Close()

	var alerts []*PriceThresholdAlert
	for rows.Next() {
		alert, err := r.scanAlert(rows)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}

	return alerts, nil
}

// UpdateAlert 更新预警
func (r *PriceAlertRepository) UpdateAlert(alert *PriceThresholdAlert) error {
	conditionsJSON, err := json.Marshal(alert.Conditions)
	if err != nil {
		return fmt.Errorf("序列化条件失败: %w", err)
	}

	_, err = r.db.Exec(`
		UPDATE price_threshold_alerts
		SET stock_code = ?, stock_name = ?, alert_type = ?, conditions = ?,
		    is_active = ?, sensitivity = ?, cooldown_hours = ?,
		    post_trigger_action = ?, enable_sound = ?, enable_desktop = ?,
		    template_id = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, alert.StockCode, alert.StockName, alert.AlertType, conditionsJSON,
		alert.IsActive, alert.Sensitivity, alert.CooldownHours, alert.PostTriggerAction,
		alert.EnableSound, alert.EnableDesktop, alert.TemplateID, alert.ID)

	if err != nil {
		return fmt.Errorf("更新预警失败: %w", err)
	}

	return nil
}

// DeleteAlert 删除预警
func (r *PriceAlertRepository) DeleteAlert(id int64) error {
	_, err := r.db.Exec(`DELETE FROM price_threshold_alerts WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("删除预警失败: %w", err)
	}
	return nil
}

// ToggleAlertStatus 切换预警状态
func (r *PriceAlertRepository) ToggleAlertStatus(id int64, isActive bool) error {
	_, err := r.db.Exec(`
		UPDATE price_threshold_alerts SET is_active = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, isActive, id)
	if err != nil {
		return fmt.Errorf("切换预警状态失败: %w", err)
	}
	return nil
}

// UpdateLastTriggeredTime 更新最后触发时间
func (r *PriceAlertRepository) UpdateLastTriggeredTime(id int64) error {
	_, err := r.db.Exec(`
		UPDATE price_threshold_alerts SET last_triggered_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, id)
	if err != nil {
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
		return false, nil // 从未触发过，不在冷却期
	}

	cooldownDuration := time.Duration(cooldownHours) * time.Hour
	elapsed := time.Since(alert.LastTriggeredAt)

	return elapsed < cooldownDuration, nil
}

// SaveTriggerHistory 保存预警触发历史
func (r *PriceAlertRepository) SaveTriggerHistory(history *PriceAlertTriggerHistory) error {
	_, err := r.db.Exec(`
		INSERT INTO price_alert_trigger_history (
			alert_id, stock_code, stock_name, alert_type,
			trigger_price, trigger_message, triggered_at
		) VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	`, history.AlertID, history.StockCode, history.StockName, history.AlertType,
		history.TriggerPrice, history.TriggerMessage)

	if err != nil {
		return fmt.Errorf("保存触发历史失败: %w", err)
	}

	return nil
}

// GetTriggerHistory 获取预警触发历史
func (r *PriceAlertRepository) GetTriggerHistory(stockCode string, limit int) ([]*PriceAlertTriggerHistory, error) {
	var rows *sql.Rows
	var err error

	if stockCode != "" {
		rows, err = r.db.Query(`
			SELECT id, alert_id, stock_code, stock_name, alert_type,
				   trigger_price, trigger_message, triggered_at
			FROM price_alert_trigger_history
			WHERE stock_code = ?
			ORDER BY triggered_at DESC LIMIT ?
		`, stockCode, limit)
	} else {
		rows, err = r.db.Query(`
			SELECT id, alert_id, stock_code, stock_name, alert_type,
				   trigger_price, trigger_message, triggered_at
			FROM price_alert_trigger_history
			ORDER BY triggered_at DESC LIMIT ?
		`, limit)
	}

	if err != nil {
		return nil, fmt.Errorf("查询触发历史失败: %w", err)
	}
	defer rows.Close()

	var histories []*PriceAlertTriggerHistory
	for rows.Next() {
		history := &PriceAlertTriggerHistory{}
		err := rows.Scan(
			&history.ID, &history.AlertID, &history.StockCode, &history.StockName,
			&history.AlertType, &history.TriggerPrice, &history.TriggerMessage, &history.TriggeredAt,
		)
		if err != nil {
			return nil, err
		}
		histories = append(histories, history)
	}

	return histories, nil
}

// GetAllTemplates 获取所有预警模板
func (r *PriceAlertRepository) GetAllTemplates() ([]*PriceAlertTemplate, error) {
	query := `
		SELECT id, name, description, alert_type, conditions, created_at
		FROM price_alert_templates ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询预警模板失败: %w", err)
	}
	defer rows.Close()

	var templates []*PriceAlertTemplate
	for rows.Next() {
		template := &PriceAlertTemplate{}
		err := rows.Scan(
			&template.ID, &template.Name, &template.Description, &template.AlertType,
			&template.Conditions, &template.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		templates = append(templates, template)
	}

	return templates, nil
}

// GetTemplateByID 根据ID获取预警模板
func (r *PriceAlertRepository) GetTemplateByID(id string) (*PriceAlertTemplate, error) {
	template := &PriceAlertTemplate{}

	err := r.db.QueryRow(`
		SELECT id, name, description, alert_type, conditions, created_at
		FROM price_alert_templates WHERE id = ?
	`, id).Scan(
		&template.ID, &template.Name, &template.Description, &template.AlertType,
		&template.Conditions, &template.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("模板不存在")
	}
	if err != nil {
		return nil, fmt.Errorf("查询模板失败: %w", err)
	}

	return template, nil
}

// CreateTemplate 创建预警模板
func (r *PriceAlertRepository) CreateTemplate(template *PriceAlertTemplate) error {
	_, err := r.db.Exec(`
		INSERT INTO price_alert_templates (id, name, description, alert_type, conditions)
		VALUES (?, ?, ?, ?, ?)
	`, template.ID, template.Name, template.Description, template.AlertType, template.Conditions)

	if err != nil {
		return fmt.Errorf("创建模板失败: %w", err)
	}

	return nil
}

// DeleteTemplate 删除预警模板
func (r *PriceAlertRepository) DeleteTemplate(id string) error {
	_, err := r.db.Exec(`DELETE FROM price_alert_templates WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("删除模板失败: %w", err)
	}
	return nil
}

// scanAlert 扫描一行数据到Alert结构体
func (r *PriceAlertRepository) scanAlert(rows *sql.Rows) (*PriceThresholdAlert, error) {
	alert := &PriceThresholdAlert{}

	var conditions sql.NullString
	var lastTriggeredAt sql.NullTime

	err := rows.Scan(
		&alert.ID, &alert.StockCode, &alert.StockName, &alert.AlertType,
		&conditions, &alert.IsActive, &alert.Sensitivity, &alert.CooldownHours,
		&alert.PostTriggerAction, &alert.EnableSound, &alert.EnableDesktop,
		&alert.TemplateID, &alert.CreatedAt, &alert.UpdatedAt, &lastTriggeredAt,
	)

	if err != nil {
		return nil, err
	}

	if conditions.Valid {
		alert.Conditions = conditions.String
	}

	if lastTriggeredAt.Valid {
		alert.LastTriggeredAt = lastTriggeredAt.Time
	}

	return alert, nil
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
