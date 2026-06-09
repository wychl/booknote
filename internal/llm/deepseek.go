package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-deepseek/deepseek"
	"github.com/go-deepseek/deepseek/request"
	"github.com/wychl/booknote/internal/datasource"
	"github.com/wychl/booknote/internal/theme"
)

type deepseekGenerator struct {
	apiKey string
	client deepseek.Client
}

func NewDeepSeekGenerator(apiKey string) (Generator, error) {
	client, err := deepseek.NewClient(apiKey)
	if err != nil {
		return nil, err
	}
	return &deepseekGenerator{
		apiKey: apiKey,
		client: client,
	}, nil
}

// GenerateNote 调用 DeepSeek 生成笔记并解析为结构化结果
func (g *deepseekGenerator) GenerateNote(ctx context.Context, styleName string, themes []theme.Theme, book *datasource.BookDetail) (*GenerateResult, error) {
	// 1. 构建提示词（假设 BuildNotePrompt 已按 JSON 格式要求构造）
	prompt, err := BuildNotePrompt(styleName, themes, book)
	if err != nil {
		return nil, fmt.Errorf("构建提示失败: %w", err)
	}

	// 2. 调用 DeepSeek API
	req := &request.ChatCompletionsRequest{
		Model: deepseek.DEEPSEEK_CHAT_MODEL,
		Messages: []*request.Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	resp, err := g.client.CallChatCompletionsChat(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("调用 DeepSeek API 失败: %w", err)
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("DeepSeek 返回空结果")
	}

	content := strings.TrimSpace(resp.Choices[0].Message.Content)

	// 3. 尝试解析 JSON
	var result GenerateResult
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		// 降级：如果解析失败，尝试从响应中提取 JSON（例如响应前后可能有额外文本）
		jsonStart := strings.Index(content, "{")
		jsonEnd := strings.LastIndex(content, "}")
		if jsonStart != -1 && jsonEnd != -1 && jsonEnd > jsonStart {
			jsonStr := content[jsonStart : jsonEnd+1]
			if err2 := json.Unmarshal([]byte(jsonStr), &result); err2 != nil {
				return nil, fmt.Errorf("解析 JSON 失败: 原始内容=%s, 错误=%w", content, err)
			}
		} else {
			return nil, fmt.Errorf("解析 JSON 失败: 原始内容=%s, 错误=%w", content, err)
		}
	}

	// 4. 基本校验（可选）
	if result.MainText == "" {
		return nil, fmt.Errorf("LLM 返回的 main_text 为空")
	}
	result.Tags = append([]string{"推荐好书", "读书感悟人生", "山海作品推荐", book.Title, book.Author}, result.Tags...)

	return &result, nil
}
