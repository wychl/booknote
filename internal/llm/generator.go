package llm

import (
	"context"
	"fmt"

	"github.com/wychl/booknote/internal/datasource"
	"github.com/wychl/booknote/internal/theme"
)

type GenerateResult struct {
	MainText string   `json:"main_text"`
	Tags     []string `json:"tags"`
	Bgm      []string `json:"bgm"`
	Theme    string   `json:"theme"`
}

// Generator 笔记生成接口
type Generator interface {
	GenerateNote(ctx context.Context, styleName string, themes []theme.Theme, book *datasource.BookDetail) (*GenerateResult, error)
}

func NewGenerator(provider, apiKey string) (Generator, error) {
	if provider == "deepseek" {
		return NewDeepSeekGenerator(apiKey)
	}
	return nil, fmt.Errorf("不支持的模型: %s", provider)
}
