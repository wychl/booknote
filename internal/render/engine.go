package render

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
)

type DataOption func(any)

//go:embed templates/*.tmpl
var templateFS embed.FS

// Engine 定义模板渲染接口
type Engine interface {
	RenderCard(data any, opts ...DataOption) (string, error)
}

// htmlEngine 实现基于 Go 模板的渲染
type htmlEngine struct {
	templates map[string]*template.Template
	mu        sync.RWMutex
}

// NewHTMLEngine 创建一个新的 HTML 渲染引擎
func NewHTMLEngine() (Engine, error) {
	engine := &htmlEngine{
		templates: make(map[string]*template.Template),
	}
	if err := engine.loadTemplates(); err != nil {
		return nil, fmt.Errorf("加载模板失败: %w", err)
	}
	return engine, nil
}

// loadTemplates 加载嵌入的模板文件
func (e *htmlEngine) loadTemplates() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// 遍历 templates 目录下的所有 .tmpl 文件
	err := fs.WalkDir(templateFS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".tmpl" {
			return nil
		}
		// 读取模板内容
		data, err := templateFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("读取模板文件 %s 失败: %w", path, err)
		}
		tmpl := template.New(filepath.Base(path)).Funcs(template.FuncMap{
			"safeHTML": safeHTML,
			"safeCSS":  safeCSS,
		})
		tmpl, err = tmpl.Parse(string(data))
		if err != nil {
			return fmt.Errorf("解析模板 %s 失败: %w", path, err)
		}
		name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)) // 去掉 .tmpl 后缀
		e.templates[name] = tmpl
		return nil
	})
	if err != nil {
		return err
	}
	if _, ok := e.templates["card"]; !ok {
		return fmt.Errorf("未找到 card 模板")
	}
	return nil
}

// RenderCard 渲染读书笔记卡片
func (e *htmlEngine) RenderCard(data any, opts ...DataOption) (string, error) {
	e.mu.RLock()
	tmpl, ok := e.templates["card"]
	e.mu.RUnlock()
	if !ok {
		return "", fmt.Errorf("card 模板未加载")
	}

	// 应用选项
	for _, opt := range opts {
		opt(data)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("模板执行失败: %w", err)
	}
	return buf.String(), nil
}

func safeHTML(raw string) template.HTML {
	return template.HTML(raw)
}

func safeCSS(raw string) template.CSS {
	return template.CSS(raw)
}
