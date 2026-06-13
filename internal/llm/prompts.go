package llm

import (
	"embed"
	"fmt"
	"log"
	"strings"

	"github.com/wychl/booknote/internal/datasource"
	"github.com/wychl/booknote/internal/theme"
)

//go:embed prompts/*
var promptsFS embed.FS

// promptCache 缓存所有模板内容，key 为风格名称（不含扩展名）
var promptCache map[string]string

func init() {
	// 初始化缓存
	promptCache = make(map[string]string)

	// 读取 prompts 目录下所有文件
	entries, err := promptsFS.ReadDir("prompts")
	if err != nil {
		panic(fmt.Sprintf("读取 prompts 目录失败: %v", err))
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		// 文件名作为风格名称（去掉 .md 后缀）
		name := strings.TrimSuffix(entry.Name(), ".md")
		if name == entry.Name() {
			// 不是 .md 文件，跳过
			continue
		}
		// 读取文件内容
		data, err := promptsFS.ReadFile("prompts/" + entry.Name())
		if err != nil {
			panic(fmt.Sprintf("读取文件 %s 失败: %v", entry.Name(), err))
		}
		promptCache[name] = string(data)
	}

	if len(promptCache) == 0 {
		panic("prompts 目录下没有找到任何 .md 模板文件")
	}
}

// BuildNotePrompt 构建提示词
// styleName: 风格名称（对应文件名，不含 .md）
// book: 书籍详情
func BuildNotePrompt(promptStyle string, themes []theme.Theme, book *datasource.BookDetail) (string, error) {
	if book == nil {
		return "", fmt.Errorf("book 参数不能为空")
	}
	if themes == nil {
		return "", fmt.Errorf("themes 参数不能为空")
	}
	template, ok := promptCache[promptStyle]
	if !ok {
		// 列出可用风格
		styles := make([]string, 0, len(promptCache))
		for s := range promptCache {
			styles = append(styles, s)
		}
		return "", fmt.Errorf("未找到风格: %s，可用风格: %v", promptStyle, styles)
	}

	themeStr := ""
	for _, t := range themes {
		themeStr += fmt.Sprintf("{{%s:%s}}\n", t.ID, t.Description)
	}

	// 替换变量
	replacer := strings.NewReplacer(
		"{{title}}", book.Title,
		"{{author}}", book.Author,
		"{{intro}}", book.Intro,
		"{{themes}}", themeStr,
	)
	prompt := replacer.Replace(template)
	log.Printf("构建提示词: %s", prompt)
	return prompt, nil
}

// GetStyleNames 返回所有可用风格名称
func GetStyleNames() []string {
	names := make([]string, 0, len(promptCache))
	for name := range promptCache {
		names = append(names, name)
	}
	return names
}
