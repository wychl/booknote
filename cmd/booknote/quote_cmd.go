package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wychl/booknote/internal/render"
)

type quoteOptions struct {
	Text   string `json:"text,omitempty"`
	Source string `json:"source,omitempty"`
	Theme  string `json:"theme,omitempty"`
	Output string `json:"output,omitempty"`
}

func newQuoteCmd(lc *components) *cobra.Command {
	var opts quoteOptions

	cmd := &cobra.Command{
		Use:   "quote [金句内容]",
		Short: "生成金句卡片（单句名言+出处）",
		Long: `根据提供的金句和出处，生成一张 1080x1920 的竖版卡片。
支持多主题，可输出 HTML 和 PNG。`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuoteCmd(cmd.Context(), lc, opts, args)
		},
	}

	cmd.Flags().StringVar(&opts.Text, "text", "", "金句内容（支持空格，可用引号包裹；也可作为位置参数）")
	cmd.Flags().StringVar(&opts.Source, "source", "", "出处，例如 《追忆似水年华》普鲁斯特")
	cmd.Flags().StringVar(&opts.Theme, "theme", "black-gold", "视觉主题 ID")
	cmd.Flags().StringVarP(&opts.Output, "output", "o", "quote", "输出文件前缀（不含扩展名）")
	_ = cmd.MarkFlagRequired("text")

	return cmd
}

func runQuoteCmd(ctx context.Context, lc *components, opts quoteOptions, args []string) error {
	// 1. 确定金句文本
	text := opts.Text
	if text == "" && len(args) > 0 {
		text = args[0]
	}
	if text == "" {
		return fmt.Errorf("请提供金句内容 (--text 或位置参数)")
	}
	opts.Source = strings.TrimPrefix(opts.Source, "——")

	// 2. 准备输出目录
	outputDir := filepath.Dir(opts.Output)
	if outputDir == "." || outputDir == "" {
		outputDir = "."
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}
	baseName := filepath.Base(opts.Output)

	// 3. 初始化结果（输出 JSON）
	result := &Result{Success: false}
	defer result.output()

	// 4. 加载主题
	themeLoader, err := lc.themeLoader.Get()
	if err != nil {
		return fmt.Errorf("获取主题加载器失败: %w", err)
	}
	defaultTheme := "classic-dark"
	theme, ok := themeLoader.Get(opts.Theme)
	if !ok {
		opts.Theme = defaultTheme
		fmt.Fprintf(os.Stderr, "警告: 主题 '%s' 不存在，已使用默认主题 '%s'\n", opts.Theme, opts.Theme)
		theme, ok = themeLoader.Get(defaultTheme)
		if !ok {
			return fmt.Errorf("主题 '%s' 不存在", defaultTheme)
		}
	}

	// 5. 准备渲染数据（使用 render.QuoteData）
	quoteData := render.NewQuoteData(text, opts.Source, &theme)

	// 6. 获取渲染引擎并生成 HTML
	renderEngine, err := lc.renderEngine.Get()
	if err != nil {
		result.Error = err.Error()
		return nil
	}
	html, err := renderEngine.Render(quoteData) // 假设引擎支持 Render
	if err != nil {
		result.Error = fmt.Sprintf("渲染 HTML 失败: %v", err)
		return nil
	}

	htmlPath := filepath.Join(outputDir, baseName+".html")
	if err := os.WriteFile(htmlPath, []byte(html), 0644); err != nil {
		result.Error = fmt.Sprintf("保存 HTML 失败: %v", err)
		return nil
	}
	fmt.Fprintf(os.Stderr, "✅ 金句卡片 HTML 已保存: %s\n", htmlPath)
	result.HTML = htmlPath

	// 7. 截图
	imgPath := filepath.Join(outputDir, baseName+".png")
	fmt.Fprintln(os.Stderr, "📸 正在生成图片...")
	screenshotSvc, err := lc.screenshotSvc.Get()
	if err != nil {
		result.Error = fmt.Sprintf("获取截图服务失败: %v", err)
		return nil
	}
	if err := screenshotSvc.Capture(ctx, htmlPath, imgPath); err != nil {
		result.Error = fmt.Sprintf("图片生成失败: %v", err)
		return nil
	} else {
		result.Image = imgPath
		fmt.Fprintf(os.Stderr, "🖼️ 图片已保存: %s\n", imgPath)
	}

	result.Success = true
	result.Data = opts
	return nil
}
