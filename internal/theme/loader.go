package theme

import (
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
)

//go:embed themes/*.json
var themeFS embed.FS

// Loader 主题加载器接口
type Loader interface {
	LoadAll() error
	Get(name string) (Theme, bool)
	GetThemes() []Theme
}

type themeLoader struct {
	themes map[string]Theme
	mu     sync.RWMutex
	loaded bool
}

// NewThemeLoader 创建从嵌入文件系统加载主题的加载器
func NewThemeLoader() (Loader, error) {
	loader := &themeLoader{}
	if err := loader.LoadAll(); err != nil {
		return nil, err
	}
	return loader, nil
}

// LoadAll 读取 themes 目录下所有 JSON 文件并解析为主题
func (l *themeLoader) LoadAll() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	entries, err := themeFS.ReadDir("themes")
	if err != nil {
		return fmt.Errorf("读取主题目录失败: %w", err)
	}

	themes := make(map[string]Theme)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		data, err := themeFS.ReadFile(filepath.Join("themes", entry.Name()))
		if err != nil {
			return fmt.Errorf("读取主题文件 %s 失败: %w", entry.Name(), err)
		}
		var theme Theme
		if err := json.Unmarshal(data, &theme); err != nil {
			return fmt.Errorf("解析主题文件 %s 失败: %w", entry.Name(), err)
		}
		if err := theme.Validate(); err != nil {
			return fmt.Errorf("主题 %s 无效: %w", entry.Name(), err)
		}
		themes[theme.ID] = theme
	}

	if len(themes) == 0 {
		return fmt.Errorf("未找到任何有效主题")
	}
	l.themes = themes
	l.loaded = true
	return nil
}

func (l *themeLoader) GetThemes() []Theme {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if !l.loaded {
		return make([]Theme, 0)
	}
	themes := make([]Theme, 0, len(l.themes))
	for _, theme := range l.themes {
		themes = append(themes, theme)
	}
	return themes
}

// Register 注册主题到加载器中
func (l *themeLoader) Register(name string, theme Theme) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, ok := l.themes[name]; ok {
		return fmt.Errorf("主题 %s 已存在", name)
	}
	l.themes[name] = theme
	return nil
}

// Get 根据主题 ID 获取主题
func (l *themeLoader) Get(name string) (Theme, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if !l.loaded {
		return Theme{}, false
	}
	theme, ok := l.themes[name]
	return theme, ok
}
