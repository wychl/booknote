package main

import (
	"github.com/wychl/booknote/internal/config"
	"github.com/wychl/booknote/internal/datasource"
	"github.com/wychl/booknote/internal/export"
	"github.com/wychl/booknote/internal/llm"
	"github.com/wychl/booknote/internal/render"
	"github.com/wychl/booknote/internal/theme"
	"github.com/wychl/booknote/pkg/lazy"
)

type components struct {
	cfg *config.Config

	datasource    *lazy.Lazy[datasource.BookSource]
	aiGen         *lazy.Lazy[llm.Generator]
	themeLoader   *lazy.Lazy[theme.Loader]
	renderEngine  *lazy.Lazy[render.Engine]
	screenshotSvc *lazy.Lazy[export.Screenshot]
}

func initComponents(cfg *config.Config) *components {
	lc := &components{cfg: cfg}
	lc.datasource = lazy.NewLazy(func() (datasource.BookSource, error) {
		return datasource.NewBookDoubanSource(cfg.DoubanCookie), nil
	})
	lc.aiGen = lazy.NewLazy(func() (llm.Generator, error) {
		return llm.NewGenerator(cfg.LLM.Provider, cfg.LLM.APIKey)
	})
	lc.themeLoader = lazy.NewLazy(func() (theme.Loader, error) {
		return theme.NewThemeLoader()
	})
	lc.renderEngine = lazy.NewLazy(func() (render.Engine, error) {
		return render.NewHTMLEngine()
	})
	lc.screenshotSvc = lazy.NewLazy(func() (export.Screenshot, error) {
		return export.NewChromedpScreenshot(), nil
	})
	return lc
}
