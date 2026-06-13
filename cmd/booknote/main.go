package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/wychl/booknote/internal/config"
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

	components := initComponents(cfg)

	rootCmd := newRootCmd(components)
	return rootCmd.ExecuteContext(ctx)
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
	rootCmd.AddCommand(newBookNoteCmd(comp))
	rootCmd.AddCommand(newQuoteCmd(comp))
	rootCmd.AddCommand(newThemesCmd(comp))
	return rootCmd
}

func outputError(errMsg string) {
	resp := Result{Success: false, Error: errMsg}
	_ = json.NewEncoder(os.Stdout).Encode(resp)
}
