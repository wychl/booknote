package render

import (
	"github.com/wychl/booknote/internal/theme"
)

// QuoteData 包含渲染引用所需的所有数据
type QuoteData struct {
	Quote   string
	Source  string
	CSSVars string
}

func (b *QuoteData) Template() string {
	return "quote"
}

// NewQuoteData 从引用内容和主题构造 QuoteData
func NewQuoteData(quote, source string, theme *theme.Theme) *QuoteData {
	return &QuoteData{
		Quote:   quote,
		Source:  source,
		CSSVars: theme.CSSVars(),
	}
}
