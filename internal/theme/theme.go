package theme

import (
	"fmt"
	"strings"
)

// Theme 表示一个完整的视觉主题
type Theme struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Variables   map[string]string `json:"variables"`
}

// CSSVars 生成 CSS 自定义属性字符串，每行格式 "    --key: value;"
func (t *Theme) CSSVars() string {
	var b strings.Builder
	for k, v := range t.Variables {
		fmt.Fprintf(&b, "    %s: %s;\n", k, v)
	}
	return b.String()
}

// Validate 检查主题是否包含必要字段
func (t *Theme) Validate() error {
	if t.ID == "" {
		return fmt.Errorf("主题 ID 不能为空")
	}
	if t.Name == "" {
		return fmt.Errorf("主题名称不能为空")
	}
	if len(t.Variables) == 0 {
		return fmt.Errorf("主题变量不能为空")
	}
	return nil
}

func (t *Theme) Desc() string {
	return fmt.Sprintf("%s: %s", t.ID, t.Description)
}
