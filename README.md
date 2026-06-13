# 📖 读书笔记卡片生成器

一个基于 Go + DeepSeek API 的智能读书笔记卡片生成工具，支持多种视觉主题，自动生成结构化的读书笔记、标签、推荐 BGM，并渲染为适配抖音等平台的 1080×1920 高清卡片。同时支持**金句卡片**生成（单句名言+出处）。

## ✨ 功能特性

- **📝 AI 自动生成笔记**：调用 DeepSeek API（兼容 OpenAI 风格），根据书籍信息生成高质量读书笔记。
- **💬 金句卡片**：快速生成单句名言+出处卡片，同样支持所有主题和截图。
- **🎨 多主题配色**：内置 13+ 视觉主题（极简白、复古纸、暖木调、黑金、黑紫、哥特暗黑、赛博朋克、电光蓝、霓虹紫等），满足不同书籍气质和平台调性。
- **🏷️ 智能标签**：自动生成 `#类型 #主题 #风格` 三类标签，便于内容分类和推荐。
- **🎵 BGM 推荐**：根据书籍情感基调，推荐 3 首契合的背景音乐（歌名 - 歌手名）。
- **📱 平台优化**：卡片尺寸固定 1080×1920，内置安全边距，适配抖音、小红书等短视频平台展示。
- **🎭 多种写作风格**：支持**感性推荐**、**深度解读**、**极简金句**、**故事叙事** 四种提示词风格。
- **🖼️ 富文本支持**：笔记正文自动将 `**关键词**` 转换为 `<strong>` 高亮，段落自动分行为 `<p>` 标签。
- **🚀 易于集成**：提供完整的 Go 模块化代码，支持嵌入模板和提示词，开箱即用。

## 📦 项目结构

```plaintext
.
├── cmd/                   # 命令行入口（booknote）
├── internal/
│   ├── config/            # 环境变量加载（支持 .env）
│   ├── datasource/        # 豆瓣书籍数据源
│   ├── llm/               # LLM 客户端、提示词管理
│   ├── render/            # HTML 渲染引擎（card.tmpl + quote.tmpl）
│   ├── theme/             # 主题加载器 + JSON 配置目录
│   ├── export/            # 截图服务（chromedp）
│   └── prompts/           # 四种写作风格的提示词 .md 文件
├── go.mod
└── README.md
```

## 🚀 快速开始

### 1. 环境要求

- Go 1.21+
- DeepSeek API Key（或其他 OpenAI 兼容 API）
- Douban Cookie（用于获取书籍详情）
- （可选）Google Chrome / Chromium，用于自动截图

### 2. 安装

```bash
go install github.com/wychl/booknote/cmd/booknote@latest
```

### 3. 配置环境变量

推荐使用 `.env` 文件，当前目录下的 `.env`。

```env
DEEPSEEK_API_KEY="sk-xxxxxx"
DOUBAN_COOKIE="bid=xxx; ll=\"108296\"; dbcl2=\"xxxx\"; ..."
```

### 4. 使用示例

#### 生成读书笔记卡片

```bash
booknote card --book="活着" --style="故事叙事" --theme="gothic-shadow"
```

#### 生成金句卡片

```bash
booknote quote "真正的发现之旅不在于寻找新风景，而在于拥有新的眼睛。" \
               --source "——《追忆似水年华》普鲁斯特" \
               --theme black-gold
```

#### 列出所有可用主题

```bash
booknote themes
```

输出示例：

```plaintext
ID               名称               描述
black-gold       黑金尊享          深邃黑背景，金色点缀，华丽沉稳...
classic-dark     经典深色          纯黑底深色主题，极致对比...
cyberpunk        赛博朋克          电光青与深蓝黑背景，霓虹高对比...
...
```

## 🎨 主题系统

所有主题位于 `internal/theme/themes/` 目录（JSON 格式），每个主题定义以下变量：

| 类别 | 变量示例 | 说明 |
|------|----------|------|
| 颜色 | `--od-color-bg`, `--od-color-surface`, `--od-color-primary`... | 16+ 语义颜色 |
| 字体尺寸 | `--od-font-size-1` (84px) ~ `--od-font-size-5` (32px) | 数字后缀，清晰层级 |
| 间距 | `--od-space-xs` ~ `--od-space-xl` | 12px ~ 72px |
| 圆角 | `--od-radius-sm` ~ `--od-radius-full` | 16px ~ 999px |
| 阴影 | `--od-shadow-card` | 卡片投影 |
| 行高/字间距 | `--od-line-height-body`, `--od-letter-spacing-body` | 正文行高1.65，字间距0.02em |

主题完全独立于业务逻辑，可用于任何卡片场景（不仅仅是读书笔记）。

## 📝 提示词风格

内置四种写作风格，位于 `internal/llm/prompts/` 目录：

- `感性推荐.md` – 温暖走心，侧重情感共鸣
- `深度解读.md` – 理性分析，挖掘核心思想
- `极简金句.md` – 只摘录书中金句 + 简短感悟
- `故事叙事.md` – 用讲故事的方式串联书籍内容

可通过 `--style` 参数选择。

## 🖼️ 自动截图

`--image` 标志会使用 `chromedp` 自动调用系统 Chrome/Chromium 生成 PNG 图片。如果失败，请检查：

- 是否安装 Chrome 并位于 PATH
- 或设置 `CHROME_PATH` 环境变量指向浏览器可执行文件

截图失败不影响 HTML 和 JSON 的生成。

## 🧩 扩展开发

### 添加新主题

1. 在 `internal/theme/themes/` 目录下新建 JSON 文件，遵循标准变量结构。
2. 运行 `booknote themes` 验证是否被加载。

### 添加新写作风格

1. 在 `internal/llm/prompts/` 目录下新建 `.md` 文件。
2. 内容为提示词模板，包含 `{{title}}` `{{author}}` `{{intro}}` 占位符。
3. 提示词末尾需规定输出 JSON 格式，包含 `main_text`, `theme`, `tags`, `bgm`。

### 更换 LLM 提供商

实现 `llm.Generator` 接口，并通过 `llm.NewGenerator` 工厂注册。

## 📄 License

MIT © [wychl](https://github.com/wychl)