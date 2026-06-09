package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/wychl/booknote/internal/config"
	"github.com/wychl/booknote/internal/datasource"
	"github.com/wychl/booknote/internal/export"
	"github.com/wychl/booknote/internal/llm"
	"github.com/wychl/booknote/internal/render"
	"github.com/wychl/booknote/internal/theme"
)

func main() {
	cfg := config.Load()

	if err := run(cfg); err != nil {
		outputError(err.Error())
		os.Exit(1)
	}
}

func run(cfg *config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 优雅退出处理
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	components, err := initComponents(cfg)
	if err != nil {
		return fmt.Errorf("初始化组件失败: %w", err)
	}

	rootCmd := newRootCmd(components)
	return rootCmd.ExecuteContext(ctx)
}

type components struct {
	datasourceClient datasource.BookSource
	aiGen            llm.Generator
	themeLoader      theme.Loader
	renderEngine     render.Engine
	screenshotSvc    export.Screenshot
}

func initComponents(cfg *config.Config) (*components, error) {
	datasourceClient := datasource.NewBookDoubanSource(cfg.DoubanCookie)

	aiGen, err := llm.NewGenerator(cfg.LLM.Provider, cfg.LLM.APIKey)
	if err != nil {
		return nil, fmt.Errorf("创建 AI 生成器失败: %w", err)
	}

	themeLoader, err := theme.NewThemeLoader()
	if err != nil {
		return nil, fmt.Errorf("加载主题失败: %w", err)
	}

	renderEngine, err := render.NewHTMLEngine()
	if err != nil {
		return nil, fmt.Errorf("初始化渲染引擎失败: %w", err)
	}

	return &components{
		datasourceClient: datasourceClient,
		aiGen:            aiGen,
		themeLoader:      themeLoader,
		renderEngine:     renderEngine,
		screenshotSvc:    export.NewBookNoteScreenshot(),
	}, nil
}

func newRootCmd(comp *components) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "booknote",
		Short: "读书笔记生成工具",
		Long: `一个基于命令行的工具，用于获取豆瓣书籍信息，并通过 AI 生成精美的读书笔记卡片。
支持多种视觉主题，可导出 HTML 文件及 PNG 图片。`,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	rootCmd.AddCommand(newCardCmd(comp))
	return rootCmd
}

func outputError(errMsg string) {
	resp := Result{Success: false, Error: errMsg}
	_ = json.NewEncoder(os.Stdout).Encode(resp)
}
