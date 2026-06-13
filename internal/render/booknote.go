package render

import (
	"github.com/wychl/booknote/internal/datasource"
	"github.com/wychl/booknote/internal/theme"
)

// BookNoteData 包含渲染卡片所需的所有数据
type BookNoteData struct {
	Title        string
	Author       string
	Rating       string
	Stars        string
	NoteMainText string
	CSSVars      string
}

func (b *BookNoteData) Template() string {
	return "booknote"
}

// NewBookNoteData 从书籍信息、AI 笔记内容和主题构造 BookNoteData
func NewBookNoteData(book *datasource.BookDetail, note string, theme *theme.Theme) *BookNoteData {
	stars := generateStarString(book.Rating)
	return &BookNoteData{
		Title:        book.Title,
		Author:       book.Author,
		Rating:       book.Rating,
		Stars:        stars,
		NoteMainText: note,
		CSSVars:      theme.CSSVars(),
	}
}

func WithFormatBookNote() DataOption {
	return func(data any) {
		if v, ok := data.(*BookNoteData); ok {
			v.NoteMainText = formatStrongText(v.NoteMainText)
			v.NoteMainText = formatChineseQuotesStrong(v.NoteMainText)
		}
	}
}
