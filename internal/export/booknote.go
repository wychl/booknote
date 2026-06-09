package export

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
)

type bookNoteScreenshot struct {
	opts []chromedp.ExecAllocatorOption
}

func NewBookNoteScreenshot() Screenshot {
	return &bookNoteScreenshot{
		opts: []chromedp.ExecAllocatorOption{
			chromedp.NoFirstRun,
			chromedp.NoDefaultBrowserCheck,
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("headless", true),
		},
	}
}

func (s *bookNoteScreenshot) Capture(ctx context.Context, htmlPath, outputPath string) error {
	// 检查 HTML 文件是否存在
	if _, err := os.Stat(htmlPath); err != nil {
		return fmt.Errorf("HTML 文件不存在: %w", err)
	}

	absPath, err := filepath.Abs(htmlPath)
	if err != nil {
		return fmt.Errorf("获取绝对路径失败: %w", err)
	}
	url := "file://" + absPath

	// 确保输出目录存在
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 目标尺寸
	width := 1080
	height := 1920

	// 创建带窗口大小的分配器上下文
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx,
		append(s.opts, chromedp.WindowSize(width, height))...,
	)
	defer cancelAlloc()

	// 创建浏览器上下文
	chromeCtx, cancelCtx := chromedp.NewContext(allocCtx)
	defer cancelCtx()

	// 设置超时（最长 30 秒）
	timeoutCtx, cancelTimeout := context.WithTimeout(chromeCtx, 30*time.Second)
	defer cancelTimeout()

	var buf []byte
	err = chromedp.Run(timeoutCtx,
		chromedp.EmulateViewport(1080, 1920), // 关键：视口 = 1080x1920
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		// 强制设置视口大小为 1080x1920，缩放比为 1
		chromedp.EmulateViewport(int64(width), int64(height), chromedp.EmulateScale(1)),
		// 等待字体和样式渲染完成
		chromedp.Sleep(500*time.Millisecond),
		// 全页截图（此时页面的 body 尺寸已是 1080x1920）
		chromedp.FullScreenshot(&buf, 100),
	)
	if err != nil {
		return fmt.Errorf("chromedp 执行失败: %w", err)
	}

	return os.WriteFile(outputPath, buf, 0644)
}
