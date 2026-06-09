package render

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/wychl/booknote/internal/datasource"
	"github.com/wychl/booknote/internal/theme"
)

// CardData 包含渲染卡片所需的所有数据
type CardData struct {
	Title        string
	Author       string
	Rating       string
	Stars        string
	NoteMainText string
	CSSVars      string
}

// NewCardData 从书籍信息、AI 笔记内容和主题构造 CardData
func NewCardData(book *datasource.BookDetail, note string, theme *theme.Theme) *CardData {
	stars := generateStarString(book.Rating)
	return &CardData{
		Title:        book.Title,
		Author:       book.Author,
		Rating:       book.Rating,
		Stars:        stars,
		NoteMainText: note,
		CSSVars:      theme.CSSVars(),
	}
}

// generateStarString 根据豆瓣评分（如 "9.4"）生成星级字符串（"★★★★★" 或 "★★★★☆"）
func generateStarString(ratingStr string) string {
	var rating float64
	fmt.Sscanf(ratingStr, "%f", &rating)
	fullStars := int(math.Round(rating / 2))
	if fullStars > 5 {
		fullStars = 5
	}
	if fullStars < 0 {
		fullStars = 0
	}
	return strings.Repeat("★", fullStars) + strings.Repeat("☆", 5-fullStars)
}

func WithFormatNoteMainText() DataOption {
	return func(data any) {
		if v, ok := data.(*CardData); ok {
			v.NoteMainText = formatStrongText(v.NoteMainText)
			v.NoteMainText = formatChineseQuotesStrong(v.NoteMainText)
		}
	}
}

func formatStrongText(raw string) string {
	// 将 **xxx** 替换为 <strong>xxx</strong>
	re := regexp.MustCompile(`\*\*(.*?)\*\*`)
	html := re.ReplaceAllString(raw, `<strong>$1</strong>`)
	// 将换行符转为 <p> 段落（可选）
	paragraphs := strings.Split(html, "\n")
	var builder strings.Builder
	for _, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		builder.WriteString("<p>")
		builder.WriteString(p)
		builder.WriteString("</p>")
	}
	return builder.String()
}

// formatChineseQuoteStrong 将中文引号包裹的内容转为 <strong>
// 支持：『xxx』、‘xxx’、“xxx” （常见中文引号对）
// 示例：'重要' -> <strong>'重要'</strong>
func formatChineseQuotesStrong(raw string) string {
	// 匹配从左引号到右引号之间的内容（包括引号本身）
	// 使用非贪婪匹配，确保最短成对匹配
	re := regexp.MustCompile(`[‘“].*?[’”]`)
	html := re.ReplaceAllString(raw, "<strong>$0</strong>")

	// 段落处理：按换行分割并包裹 <p> 标签
	paragraphs := strings.Split(html, "\n")
	var builder strings.Builder
	for _, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		builder.WriteString("<p>")
		builder.WriteString(p)
		builder.WriteString("</p>")
	}
	return builder.String()
}
