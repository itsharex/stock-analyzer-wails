package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/go-resty/resty/v2"
)

// CookieManager 负责维护和自动刷新 Cookie
type CookieManager struct {
	client *resty.Client
}

func NewCookieManager() *CookieManager {
	return &CookieManager{
		client: resty.New().SetRedirectPolicy(resty.FlexibleRedirectPolicy(5)),
	}
}

// GetAutoCookie 获取并组装 Cookie 字符串
func (m *CookieManager) GetAutoCookie() (string, error) {
	// 1. 模拟浏览器访问东财首页或行情中心
	resp, err := m.client.R().
		SetHeaders(map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			"Accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8",
		}).
		Get("https://quote.eastmoney.com/center/gridlist.html") // 访问行情中心更容易拿到有效的 QGQT 等关键 Cookie

	if err != nil {
		return "", err
	}

	// 2. 从响应中提取所有 Set-Cookie
	cookies := resp.Cookies()
	if len(cookies) == 0 {
		return "", fmt.Errorf("未能从服务器获取到 Cookie")
	}

	var cookiePairs []string
	for _, c := range cookies {
		cookiePairs = append(cookiePairs, fmt.Sprintf("%s=%s", c.Name, c.Value))
	}

	// 3. 将 Cookie 对象列表转换为 Header 需要的字符串格式
	fullCookie := strings.Join(cookiePairs, "; ")
	return fullCookie, nil
}

// GetStockCookie 模拟真实浏览器行为，获取东方财富关键 Cookie
// targetURL 建议传入个股详情页，因为个股页触发的鉴权请求最全
func (m *CookieManager) GetStockCookie(targetURL string) (string, error) {
	// 1. 禁用 GPU 和沙盒模式，提高在 Linux/Docker 环境下的兼容性
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.NoSandbox,
		chromedp.Flag("headless", true), // 无头模式
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// 2. 创建上下文
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// 3. 设置执行超时（防止死锁）
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var cookiesStr string

	// 4. 执行浏览器任务
	err := chromedp.Run(ctx,
		// 监听网络事件，允许访问网络库
		network.Enable(),
		// 访问页面
		chromedp.Navigate(targetURL),
		// 核心：必须等待页面加载完成，且预留时间让 JS 执行异步设置 Cookie
		chromedp.Sleep(5*time.Second),
		// 提取所有 Cookie
		chromedp.ActionFunc(func(ctx context.Context) error {
			cookies, err := network.GetCookies().Do(ctx)
			if err != nil {
				return err
			}

			var cookieParts []string
			for _, c := range cookies {
				// 打印日志方便调试（可选）
				// log.Printf("找到 Cookie: %s = %s", c.Name, c.Value)
				cookieParts = append(cookieParts, fmt.Sprintf("%s=%s", c.Name, c.Value))
			}
			cookiesStr = strings.Join(cookieParts, "; ")
			return nil
		}),
	)

	if err != nil {
		return "", fmt.Errorf("浏览器抓取 Cookie 失败: %v", err)
	}

	if cookiesStr == "" {
		return "", fmt.Errorf("未能提取到任何 Cookie")
	}

	return cookiesStr, nil
}
