// card.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/wychl/booknote/internal/datasource"
	"github.com/wychl/booknote/internal/render"
)

type cardOptions struct {
	id       string
	theme    string
	output   string
	book     string
	fromJSON string
	style    string
}

func newCardCmd(comp *components) *cobra.Command {
	var opts cardOptions

	cmd := &cobra.Command{
		Use:   "card [书籍ID]",
		Short: "生成书籍笔记卡片",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCard(cmd.Context(), comp, opts, args)
		},
	}

	cmd.Flags().StringVarP(&opts.theme, "theme", "t", "", "视觉主题")
	cmd.Flags().StringVarP(&opts.output, "output", "o", "", "输出目录")
	cmd.Flags().StringVarP(&opts.book, "book", "b", "", "书籍名称关键词，自动选择第一个搜索结果")
	cmd.Flags().StringVar(&opts.fromJSON, "from-json", "", "从 JSON 文件加载数据（跳过 AI 调用和搜索）")
	cmd.Flags().StringVar(&opts.style, "style", "深度解读", "笔记风格 (深度解读, 感性推荐等)")

	return cmd
}

func runCard(ctx context.Context, comp *components, opts cardOptions, args []string) error {
	outputDir := getCardOutputDir(opts)

	result := &Result{Success: false}
	defer result.output()

	bookNote, err := loadBookData(ctx, comp, opts, args)
	if err != nil {
		result.Error = err.Error()
		return nil
	}
	result.BookNote = *bookNote

	if err := applyTheme(comp, bookNote); err != nil {
		result.Error = err.Error()
		return nil
	}

	htmlPath, err := renderAndSaveHTML(comp, outputDir, bookNote)
	if err != nil {
		result.Error = err.Error()
		return nil
	}
	result.HTML = htmlPath

	imagePath := exportCardImage(ctx, comp, htmlPath, outputDir, bookNote)
	result.Image = imagePath

	jsonPath := saveJSONFile(outputDir, bookNote)
	result.JSONFile = jsonPath

	result.Success = true
	return nil
}

func getCardOutputDir(opts cardOptions) string {
	if opts.output != "" {
		return opts.output
	}
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}

func loadBookData(ctx context.Context, comp *components, opts cardOptions, args []string) (*BookNote, error) {
	if opts.fromJSON != "" {
		bn, err := loadBookNoteFromJSON(opts.fromJSON)
		if err != nil {
			return nil, fmt.Errorf("从 JSON 加载失败: %w", err)
		}
		fmt.Fprintf(os.Stderr, "📖 已从本地文件加载: %s\n", bn.Book.Title)
		if opts.theme != "" {
			bn.Theme = opts.theme
		}
		return bn, nil
	}

	bookID, err := resolveBookID(ctx, comp.datasourceClient, opts, args)
	if err != nil {
		return nil, err
	}
	opts.id = bookID

	book, err := comp.datasourceClient.GetBookDetail(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("获取书籍详情失败: %w", err)
	}
	fmt.Fprintf(os.Stderr, "✅ 已获取书籍: %s\n", book.Title)

	themes := comp.themeLoader.GetThemes()
	genResult, err := comp.aiGen.GenerateNote(ctx, opts.style, themes, book)
	if err != nil {
		return nil, fmt.Errorf("AI 生成笔记失败: %w", err)
	}

	bookNote := &BookNote{
		Book:  *book,
		Note:  genResult.MainText,
		Tags:  genResult.Tags,
		Bgm:   genResult.Bgm,
		Style: opts.style,
		Theme: genResult.Theme,
	}
	if opts.theme != "" {
		bookNote.Theme = opts.theme
	}
	return bookNote, nil
}

func resolveBookID(ctx context.Context, client datasource.BookSource, opts cardOptions, args []string) (string, error) {
	if opts.book != "" {
		summaries, err := client.Search(ctx, opts.book, 10)
		if err != nil {
			return "", fmt.Errorf("搜索失败: %w", err)
		}
		if len(summaries) == 0 {
			return "", fmt.Errorf("未找到相关书籍: %s", opts.book)
		}
		fmt.Fprintf(os.Stderr, "📚 已自动选择: %s\n", summaries[0].Title)
		return summaries[0].ID, nil
	}
	if len(args) == 1 {
		return args[0], nil
	}
	return "", fmt.Errorf("请提供书籍ID、使用 --book 或 --from-json")
}

func applyTheme(comp *components, bookNote *BookNote) error {
	tm, ok := comp.themeLoader.Get(bookNote.Theme)
	if !ok {
		defaultTheme := "classic-dark"
		tm, ok = comp.themeLoader.Get(defaultTheme)
		if !ok {
			return fmt.Errorf("默认主题 '%s' 不存在", defaultTheme)
		}
		fmt.Fprintf(os.Stderr, "警告: 主题 '%s' 不存在，已使用默认主题 '%s'\n", bookNote.Theme, defaultTheme)
		bookNote.Theme = defaultTheme
	}
	_ = tm
	return nil
}

func renderAndSaveHTML(comp *components, outputDir string, bookNote *BookNote) (string, error) {
	tm, _ := comp.themeLoader.Get(bookNote.Theme)
	cardData := render.NewCardData(&bookNote.Book, bookNote.Note, &tm)

	html, err := comp.renderEngine.RenderCard(cardData, render.WithFormatNoteMainText())
	if err != nil {
		return "", fmt.Errorf("渲染 HTML 失败: %w", err)
	}

	outputPrefix := getCardOutputPrefix(bookNote)
	htmlPath := filepath.Join(outputDir, outputPrefix+".html")

	if err := os.WriteFile(htmlPath, []byte(html), 0644); err != nil {
		return "", fmt.Errorf("保存 HTML 失败: %w", err)
	}
	fmt.Fprintf(os.Stderr, "✅ HTML 卡片已保存: %s\n", htmlPath)
	return htmlPath, nil
}

func exportCardImage(ctx context.Context, comp *components, htmlPath string, outputDir string, bookNote *BookNote) string {
	outputPrefix := getCardOutputPrefix(bookNote)
	imagePath := filepath.Join(outputDir, outputPrefix+".png")

	fmt.Fprintln(os.Stderr, "📸 正在生成图片（需要 Chrome）...")
	if err := comp.screenshotSvc.Capture(ctx, htmlPath, imagePath); err != nil {
		fmt.Fprintf(os.Stderr, "图片生成失败: %v\n", err)
		return ""
	}
	fmt.Fprintf(os.Stderr, "🖼️ 图片已保存: %s\n", imagePath)
	return imagePath
}

func saveJSONFile(outputDir string, bookNote *BookNote) string {
	outputPrefix := getCardOutputPrefix(bookNote)
	jsonPath := filepath.Join(outputDir, outputPrefix+".json")

	if err := saveBookNoteToJSON(bookNote, jsonPath); err != nil {
		fmt.Fprintf(os.Stderr, "警告: 保存 JSON 文件失败: %v\n", err)
		return ""
	}
	fmt.Fprintf(os.Stderr, "💾 书籍数据已保存至: %s\n", jsonPath)
	return jsonPath
}

func getCardOutputPrefix(bookNote *BookNote) string {
	if bookNote != nil && bookNote.Book.ID != "" {
		return bookNote.Book.ID
	}
	return "booknote"
}

func saveBookNoteToJSON(bn *BookNote, filePath string) error {
	data, err := json.MarshalIndent(bn, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

func loadBookNoteFromJSON(filePath string) (*BookNote, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var bn BookNote
	if err := json.Unmarshal(data, &bn); err != nil {
		return nil, err
	}
	return &bn, nil
}
