
# 📖 读书笔记卡片生成器

一个基于 Go + DeepSeek API 的智能读书笔记卡片生成工具，支持多种风格主题，自动生成结构化的读书笔记、标签和推荐 BGM，并渲染为适配抖音等平台的 1080×1920 高清卡片。

## ✨ 功能特性

- **📝 AI 自动生成笔记**：调用 DeepSeek API（兼容 OpenAI 风格），根据书籍信息生成高质量读书笔记。
- **🎨 多主题配色**：内置多种视觉主题（极简白、复古纸、暖木调、黑金、黑紫、哥特暗黑、变革蓝绿等），满足不同书籍气质和平台调性。
- **🏷️ 智能标签**：自动生成 `#类型 #主题 #风格` 三类标签，便于内容分类和推荐。
- **🎵 BGM 推荐**：根据书籍情感基调，推荐 3 首契合的背景音乐（歌名 - 歌手名）。
- **📱 平台优化**：卡片尺寸固定 1080×1920，内置安全边距，适配抖音、小红书等短视频平台展示。
- **🎭 多种写作风格**：支持**感性推荐**、**深度解读**、**极简金句**、**故事叙事** 四种提示词风格，可灵活扩展。
- **🖼️ 富文本支持**：笔记正文自动将 `**关键词**` 转换为 `<strong>` 高亮，段落自动分行为 `<p>` 标签。
- **🚀 易于集成**：提供完整的 Go 模块化代码，支持嵌入提示词文件，开箱即用。

## 📦 项目结构

```plaintext
.
├── cmd/                # 命令行入口（示例）
├── internal/
│   ├── datasource/     # 书籍数据模型
│   ├── llm/            # LLM 客户端、提示词管理
│   ├── templates/      # HTML 模板（card.tmpl）
│   └── theme/          # 主题 JSON 配置
│   └── prompts/          # 四种风格的提示词 .md 文件
├── go.mod
└── README.md
```

## 快速开始

### 1. 环境要求

- Go 1.21+
- DeepSeek API Key（或其他 OpenAI 兼容 API）
- Douban Cookie（用于获取书籍详情）

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 配置环境变量

```bash
export DEEPSEEK_API_KEY="your-deepseek-api-key"
export DOUBAN_COOKIE="your-douban-cookie"
```

## 主题配置

所有主题位于 `themes/` 目录（JSON 格式），例如 `classic-dark-gold.json`：

```json
{
  "id": "classic-dark-gold",
  "name": "黑金质感",
  "description": "纯黑背景配金色点缀，高端醒目",
  "variables": {
    "--od-color-bg-stage": "#000000",
    "--od-color-surface": "#0d0d0d",
    "--od-color-primary": "#d4af37",
    ...
  }
}
```

使用时将 `{{.CSSVars}}` 替换为主题对应的 CSS 变量字符串。项目内置了以下主题：

| 主题 ID | 名称 | 适用场景 |
|---------|------|----------|
| pure-white | 极简白 | 干货、书评 |
| vintage-paper | 复古纸 | 经典、历史 |
| warm-wood | 暖木调 | 文学、散文 |
| classic-dark-gold | 黑金质感 | 金句、励志 |
| dark-purple | 黑紫幻夜 | 科幻、创意 |
| gothic-dark | 哥特暗黑 | 悬疑、奇幻 |
| transformative-teal | 变革蓝绿 | 心理、成长 |

## 📝 提示词风格

内置三种写作风格，位于 `prompts/` 目录：

- `感性推荐.md` – 温暖走心，侧重情感共鸣
- `深度解读.md` – 理性分析，挖掘核心思想
- `故事叙事.md` – 讲故事，引人入胜

## 扩展开发

### 添加新主题

1. 在 `themes/` 目录下新建 JSON 文件，遵循同一结构。
2. 在 `card.tmpl` 中通过 `{{.CSSVars}}` 注入变量。

### 添加新写作风格

1. 在 `prompts/` 目录下新建 `.md` 文件，内容为提示词模板（包含 `{{title}}` `{{author}}` `{{intro}}`）。
2. 提示词末尾需规定输出 JSON 格式，包含 `main_text`, `theme`, `tags`, `bgm`。

### 更换 LLM 提供商

实现 `Generator` 接口（`GenerateNote` 方法）即可轻松替换为 OpenAI、Ollama 等。

## 📄 License

MIT © [wychl](https://github.com/wychl)