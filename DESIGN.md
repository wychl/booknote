# Open Design 卡片设计系统

*Version 1.0 | 基于语义令牌的主题驱动卡片系统*

---

## 设计哲学

### 核心原则

1. **主题驱动** — 所有视觉样式通过 CSS 变量定义，模板不硬编码任何颜色、字体、间距值。
2. **克制优于装饰** — 一个卡片只使用一种装饰手法，不叠加多层渐变、纹理和噪点。
3. **变体可组合** — 字体尺寸、间距、圆角、阴影等结构变量在所有主题中固定不变，仅颜色和品牌相关变量可调整，确保卡片尺寸统一。
4. **回退可读** — 每个 `var()` 引用必须提供硬编码回退值，确保无主题文件时仍可渲染。

### 分层架构

```
┌─────────────────────────────────────────┐
│  body (固定画布 1080×1920)              │
│  ┌───────────────────────────────────┐  │
│  │  od-[type]-card (卡片容器)        │  │
│  │  ┌─────────────────────────────┐  │  │
│  │  │  card-deco (装饰层, 可选)   │  │  │
│  │  └─────────────────────────────┘  │  │
│  │  ┌─────────────────────────────┐  │  │
│  │  │  safe-area (安全区, z:2)    │  │  │
│  │  │  ┌───────────────────────┐  │  │  │
│  │  │  │  card-content (内容)  │  │  │  │
│  │  │  └───────────────────────┘  │  │  │
│  │  └─────────────────────────────┘  │  │
│  └───────────────────────────────────┘  │
└─────────────────────────────────────────┘
```

### 各层职责

| 层级 | 必需 | CSS 定位 | 职责 |
|------|------|----------|------|
| `body` | ✅ | static | 固定画布尺寸、整体居中、页面外围边距 |
| `od-[type]-card` | ✅ | `position: relative` | 卡片容器，承载滚动、字体族（不设圆角/阴影/背景色） |
| `card-deco` | ❌ 可选 | `position: absolute; inset: 0` | 装饰层，单层径向渐变，`pointer-events: none` |
| `safe-area` | ✅ | `position: relative; z-index: 2` | 内容安全区，垂直居中，padding 定义内容与边缘距离 |
| `card-content` | ✅ | static flex (`gap`) | 内容区块容器，使用 `gap` 定义区块间距 |

---

## 主题系统

### 主题文件结构

```json
{
  "id": "theme-id",
  "name": "主题名称",
  "description": "简短描述，不超过50字",
  "variables": {
    // … 所有 CSS 变量
  }
}
```

- 文件存放位置：`internal/theme/themes/<id>.json`
- `id`：小写字母数字连字符
- `name`：中文，2~6 字

### 变量体系

#### 颜色（17个必填）

| 变量 | 用途 | 深色主题示例 | 浅色主题示例 |
|------|------|-------------|-------------|
| `--od-color-bg` | 卡片外围背景 | `#0a0a0a` | `#f0f4f8` |
| `--od-color-surface` | 卡片主体背景 | `#141414` | `#ffffff` |
| `--od-color-primary` | 主品牌色 | `#d4af37` | `#7c3aed` |
| `--od-color-primary-dark` | 主色暗变体 | `#b8922a` | `#5b21b6` |
| `--od-color-secondary` | 辅助色（同色系） | `#e4c668` | `#a78bfa` |
| `--od-color-hover` | 交互悬停色 | `#e4c668` | `#7c3aed` |
| `--od-color-link` | 链接色 | `#d4af37` | `#7c3aed` |
| `--od-color-divider` | 分割线/边框色 | `#2c2c2c` | `#e2e8f0` |
| `--od-color-text-high` | 高强调文字 | `#ffffff` | `#0f172a` |
| `--od-color-text-medium` | 中强调文字 | `#e0e0e0` | `#334155` |
| `--od-color-text-low` | 弱化文字 | `#a0a0a0` | `#94a3b8` |
| `--od-color-surface-accent` | 强调区块背景 | `color-mix(in srgb, var(--od-color-primary) 12%, var(--od-color-surface))` | 同左公式 |
| `--od-color-highlight` | 高亮色（星星等） | `#fbbf24` | `#f59e0b` |
| `--od-color-success` | 成功状态 | `#9ece6a` | `#22c55e` |
| `--od-color-error` | 错误状态 | `#f7768e` | `#ef4444` |
| `--od-color-warning` | 警告状态 | `#e0af68` | `#f97316` |

#### 字体尺寸（5个，固定值，不修改）

| 变量 | 值 | 用途 |
|------|-----|------|
| `--od-font-size-1` | `84px` | 超大标题（书籍标题、金句正文） |
| `--od-font-size-2` | `40px` | 大标题/星星 |
| `--od-font-size-3` | `36px` | 正文（笔记内容） |
| `--od-font-size-4` | `34px` | 副文本/作者/来源 |
| `--od-font-size-5` | `32px` | 注释/标签 |

#### 间距（5个，固定值，不修改）

| 变量 | 值 | 用途 |
|------|-----|------|
| `--od-space-xs` | `12px` | 极小间距 |
| `--od-space-sm` | `24px` | 小间距 |
| `--od-space-md` | `40px` | 中间距 |
| `--od-space-lg` | `56px` | 大间距（金句与来源之间） |
| `--od-space-xl` | `72px` | 超大间距 |

#### 圆角（4个，固定值，不修改）

| 变量 | 值 | 用途 |
|------|-----|------|
| `--od-radius-sm` | `16px` | 小圆角 |
| `--od-radius-md` | `36px` | 中圆角 |
| `--od-radius-lg` | `48px` | 大圆角 |
| `--od-radius-full` | `999px` | 完全圆角（胶囊/标签/滚动条） |

#### 字体相关（固定值，不修改）

| 变量 | 值 |
|------|-----|
| `--od-font-family` | `system-ui, -apple-system, 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', 'Helvetica Neue', sans-serif` |
| `--od-font-weight-regular` | `400` |
| `--od-font-weight-medium` | `500` |
| `--od-font-weight-bold` | `700` |
| `--od-line-height-body` | `1.65` |
| `--od-letter-spacing-body` | `0.02em` |

#### 阴影（可选）

| 变量 | 示例值 |
|------|--------|
| `--od-shadow-card` | `0 24px 40px -14px rgba(0, 0, 0, 0.85)` |

> ⚠️ 卡片模板不预设阴影。阴影由上层容器按需添加。

---

## 模板规范

### 金句卡片 (`od-quote-card`)

```html
<div class="od-quote-card">
  <div class="card-deco"></div>
  <div class="safe-area">
    <div class="card-content">
      <div class="quote-text">{{.Quote}}</div>
      <div class="quote-source">{{.Source}}</div>
    </div>
  </div>
</div>
```

- `quote-text`：字号 `var(--od-font-size-1)`，引号通过 `::before`/`::after` 伪元素添加，颜色取 `var(--od-color-primary)`
- `quote-source`：字号 `var(--od-font-size-4)`，上方 `border-top: 2px solid var(--od-color-divider)`
- 区块间距：`gap: var(--od-space-lg, 56px)`

### 读书笔记卡片 (`od-note-card`)

```html
<div class="od-note-card">
  <div class="safe-area">
    <div class="card-content">
      <div class="book-header">
        <div class="book-title">{{.Title}}</div>
        <div class="author-rating">
          <span class="book-author">{{.Author}}</span>
          <span class="stars">{{.Stars}}</span>
          <span class="score-badge">豆瓣 {{.Rating}}</span>
        </div>
      </div>
      <div class="reflection">{{.NoteMainText | safeHTML}}</div>
    </div>
  </div>
</div>
```

- `book-title`：字号 `var(--od-font-size-1)`
- `book-author`：字号 `var(--od-font-size-5)`，颜色 `var(--od-color-primary)`
- `stars`：字号 `var(--od-font-size-5)`，颜色 `var(--od-color-highlight)`，字距 `0.3em`
- `score-badge`：背景 `var(--od-color-surface-accent)`，圆角 `var(--od-radius-full)`
- `reflection`：字号 `var(--od-font-size-3)`，`<strong>` 高亮颜色 `var(--od-color-primary)`
- 区块间距：`gap: var(--od-space-md, 40px)`

---

## 主题示例

### 黑金尊享 (`black-gold`)

```json
{
  "id": "black-gold",
  "name": "黑金尊享",
  "description": "深邃黑背景，金色点缀，华丽沉稳，适合高端卡片、品牌展示",
  "variables": {
    "--od-color-bg": "#0a0a0a",
    "--od-color-surface": "#141414",
    "--od-color-primary": "#d4af37",
    "--od-color-primary-dark": "#b8922a",
    "--od-color-secondary": "#e4c668",
    "--od-color-hover": "#e4c668",
    "--od-color-link": "#d4af37",
    "--od-color-divider": "#2c2c2c",
    "--od-color-text-high": "#ffffff",
    "--od-color-text-medium": "#e5e5e5",
    "--od-color-text-low": "#a3a3a3",
    "--od-color-surface-accent": "color-mix(in srgb, var(--od-color-primary) 12%, var(--od-color-surface))",
    "--od-color-highlight": "#fbbf24",
    "--od-color-success": "#8bc34a",
    "--od-color-error": "#e57373",
    "--od-color-warning": "#ffb74d",
    "--od-space-xs": "12px",
    "--od-space-sm": "24px",
    "--od-space-md": "40px",
    "--od-space-lg": "56px",
    "--od-space-xl": "72px",
    "--od-radius-sm": "16px",
    "--od-radius-md": "36px",
    "--od-radius-lg": "48px",
    "--od-radius-full": "999px",
    "--od-font-family": "system-ui, -apple-system, 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', 'Helvetica Neue', sans-serif",
    "--od-font-weight-regular": "400",
    "--od-font-weight-medium": "500",
    "--od-font-weight-bold": "700",
    "--od-font-size-1": "84px",
    "--od-font-size-2": "40px",
    "--od-font-size-3": "36px",
    "--od-font-size-4": "34px",
    "--od-font-size-5": "32px",
    "--od-line-height-body": "1.65",
    "--od-letter-spacing-body": "0.02em",
    "--od-shadow-card": "0 24px 40px -14px rgba(0, 0, 0, 0.85)"
  }
}
```

---

## 创建新主题

### 步骤

1. **复制模板**：从 `black-gold.json` 或 `classic-dark.json` 复制一份。
2. **修改元数据**：更改 `id`、`name`、`description`。
3. **调整颜色变量**：
   - 背景与表面保持适宜对比度（深色主题：bg 接近 `#000`，surface 稍亮）。
   - 主色选取主题色调，辅助色取同色系较亮/暗版本，避免高反差撞色。
   - 其他功能色可沿用默认推荐值。
4. **保留结构变量不变**：字体尺寸、间距、圆角、字体族、行高、字距**不修改**，确保各主题卡片尺寸统一。
5. **设置阴影**：根据明暗选择透明度（深色 0.8+，浅色 0.1~0.2）。
6. **测试**：生成预览卡片确认渲染效果。

---

## 禁止清单

| 违规模式 | 后果 | 正确做法 |
|----------|------|----------|
| 卡片容器预设圆角和阴影 | 限制复用场景 | 由上层容器或主题按需添加 |
| 引号用内联 `<span>` 而非 `::before/after` | 主题无法控制引号样式 | 使用伪元素，颜色取 `var(--od-color-primary)` |
| 分割线用独立 `<div>` 而非 `border-top` | 增加无用 DOM | 使用 `border-top: 2px solid var(--od-color-divider)` |
| 多层装饰叠加（光晕+纹理+噪点） | 视觉杂乱 | 只使用一种装饰手法 |
| 使用 `gap: 0` + 手动 margin | 累积边距导致偏移 | 使用语义化 `gap` 变量 |
| 缺少 CSS 变量 fallback | 无主题时页面白屏 | 每个 `var()` 提供硬编码回退值 |
| 在 `body` 上设背景色让卡片透明 | 出现黑边 | 背景色交给主题或上层容器 |
| 结构变量（字号/间距/圆角）做个性化调整 | 卡片尺寸不统一 | 保持结构变量为推荐值不变 |

---

## 版本记录

| 版本 | 日期 | 变更内容 |
|------|------|----------|
| v1.0 | 2026-06-14 | 初始版本：整合主题规范 v2.0 + 模板规范 v1.1，含两个主题示例、禁止清单 |
````

### 设计系统文档结构说明

这份 `DESIGN.md` 涵盖了你在整轮对话中形成的所有设计资产：

| 章节 | 来源 |
|------|------|
| **设计哲学** | 从你强调的「克制」、「回退可读」、「卡片不设背景色」等决策提炼 |
| **分层架构** | 模板规范 v1.1 的骨架要求 + 各层职责表 |
| **主题系统** | 主题设计规范 v2.0，含变量体系的完整表格 |
| **模板规范** | 金句卡片和读书笔记卡片的完整 HTML 骨架与 CSS 模式 |
| **主题示例** | 黑金尊享（精修版 #0a0a0a）|
| **创建新主题** | 从规范中提取的 6 步流程 |
| **禁止清单** | 整合了你明确指出的所有违规模式（圆角/阴影/黑边/引号等） |