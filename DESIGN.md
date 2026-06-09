# OpenDesign 设计规范

## 主题规范 · Theme Specification

*Version 1.0 | 基于 OpenDesign 品牌指南*

---

### 1. 主题定义

主题（Theme）是一组**CSS自定义属性**（CSS Variables）的集合，用于统一控制读书笔记卡片的外观。每个主题通过覆盖全局变量来实现独特的视觉风格。

主题仅影响**视觉表现**（颜色、阴影、间距等），不影响卡片的内容结构和布局逻辑。

---

### 2. 主题文件格式

每个主题以独立的 JSON 文件存储，文件名为 `<theme-id>.json`。

**JSON Schema**：

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "required": ["id", "name", "description", "variables"],
  "properties": {
    "id": {
      "type": "string",
      "pattern": "^[a-z][a-z0-9-]*$",
      "description": "主题唯一标识，小写字母、数字、连字符"
    },
    "name": {
      "type": "string",
      "maxLength": 20,
      "description": "主题显示名称"
    },
    "description": {
      "type": "string",
      "maxLength": 100,
      "description": "主题简短描述（适用场景）"
    },
    "variables": {
      "type": "object",
      "description": "CSS 变量键值对，必须包含所有必填变量",
      "additionalProperties": true
    }
  }
}
```

---

### 3. 必填变量（Mandatory Variables）

每个主题**必须**定义以下全部变量，无默认回退。这些变量覆盖颜色、间距、圆角、字体、阴影等所有可控视觉属性。

#### 3.1 颜色变量

| 变量名 | 用途 |
|--------|------|
| `--od-color-bg-stage` | 卡片外围背景（整体界面底色） |
| `--od-color-surface` | 卡片主体背景 |
| `--od-color-primary` | 主色调（强调色） |
| `--od-color-primary-dark` | 主色调暗色（悬停、压深） |
| `--od-color-secondary` | 辅助色（渐变、次强调） |
| `--od-color-hover` | 悬停/交互色 |
| `--od-color-link` | 链接颜色 |
| `--od-color-divider` | 分割线颜色 |
| `--od-color-text-title` | 标题文字颜色 |
| `--od-color-text-body` | 正文文字颜色 |
| `--od-color-text-muted` | 辅助/弱化文字颜色 |
| `--od-color-accent-paper` | 强调区块背景（金句卡片等） |
| `--od-color-reflection-bg` | 个人感悟区块背景 |
| `--od-color-star` | 评分星星颜色（点亮状态） |
| `--od-color-success` | 成功状态色 |
| `--od-color-error` | 错误状态色 |
| `--od-color-warning` | 警告状态色 |

#### 3.2 间距变量（基于 4px 基础单位）

| 变量名 | 推荐值 | 用途 |
|--------|--------|------|
| `--od-space-xs` | `12px` | 极小间距 |
| `--od-space-sm` | `24px` | 小间距 |
| `--od-space-md` | `40px` | 中间距 |
| `--od-space-lg` | `56px` | 大间距 |
| `--od-space-xl` | `72px` | 超大间距 |

> **约束**：各主题的间距值**应保持一致**（品牌统一性原则），除非有特殊设计理由。

#### 3.3 圆角变量

| 变量名 | 推荐值 | 用途 |
|--------|--------|------|
| `--od-radius-sm` | `16px` | 小圆角 |
| `--od-radius-md` | `36px` | 中圆角 |
| `--od-radius-lg` | `48px` | 大圆角 |
| `--od-radius-full` | `999px` | 完全圆角（胶囊） |

> **约束**：所有主题的圆角值应保持一致，确保视觉风格统一。

#### 3.4 字体变量

| 变量名 | 推荐值 | 用途 |
|--------|--------|------|
| `--od-font-family` | `system-ui, -apple-system, ...` | 字体栈（需包含中文回退） |
| `--od-font-weight-regular` | `400` | 常规字重 |
| `--od-font-weight-medium` | `500` | 中等字重 |
| `--od-font-weight-bold` | `700` | 粗体字重 |

> **约束**：字体栈和字重值在所有主题中**必须一致**，品牌指南不允许多种字体混用。


### 4. 可选变量（Optional Variables）

主题可以额外定义以下变量以精细化控制，若不定义则使用全局默认值（如果存在）。但为了主题自包含，建议全部显式定义。

| 变量名 | 说明 |
|--------|------|
| `--od-shadow-paper` | 内部卡片投影 |
| `--od-font-size-1` | 标题字号（默认 1.8rem） |
| `--od-font-size-base` | 正文字号（默认 0.95rem） |
| `--od-font-size-small` | 小字号（默认 0.8rem） |
| `--od-font-size-micro` | 微字号（默认 0.7rem） |

---

### 5. 主题命名规范

- **ID**：小写字母开头，仅包含小写字母、数字、连字符 `-`，例如 `pure-white`、`dark-purple`、`classic-dark-gold`
- **Name**：中文名称，2-6 个字，如 `极简白`、`黑金尊享`
- **Description**：一句话说明适用场景，例如“清爽纯净，适合干货、商业类书籍”

---

### 6. 内置主题清单

根据品牌指南，OpenDesign 内置以下 7 个主题（已全部实现）：

| 主题 ID | 名称 | 色调特征 |
|---------|------|----------|
| `pure-dawn` | 极简晨光 | 纯白+蓝，高亮 |
| `vintage-parchment` | 复古羊皮 | 暖米+褐，怀旧 |
| `warm-autumn-wood` | 暖木秋色 | 木色+琥珀，温馨 |
| `black-gold` | 黑金尊享 | 黑+金，奢华 |
| `midnight-purple` | 紫夜幻境 | 深紫+霓虹，科幻 |
| `gothic-shadow` | 哥特暗影 | 深灰+暗红，神秘 |
| `transformative-teal` | 变革碧涛 | 蓝绿+白，清新 |

---

### 7. 创建新主题的规则

1. **复制一个现有主题 JSON 文件**作为模板。
2. **修改 `id`、`name`、`description`**。
3. **调整颜色变量**：确保对比度满足 WCAG AA（正文与背景对比度 ≥4.5:1，标题与背景 ≥3:1）。
4. **微调阴影**：亮底浅阴影，暗底深阴影。
5. **保持间距、圆角、字体、字重不变**（除非品牌指南未来修订）。
6. **测试**：在 `booknote` 中实际渲染验证可读性。

---

### 8. 主题加载与合并

- 加载时，主题的 `variables` 直接合并到 `:root` 或卡片容器。
- 若某个变量未在主题中定义，**系统不应提供默认值**，而是要求主题显式定义所有必填变量（避免不一致）。
- 主题 JSON 文件放置在 `themes/` 目录，由 `theme.Loader` 动态加载。

---

### 9. 示例：最小合法主题（省略部分颜色但结构正确）

```json
{
  "id": "sample-theme",
  "name": "示例主题",
  "description": "演示主题规范",
  "variables": {
    "--od-color-bg-stage": "#f5f5f5",
    "--od-color-surface": "#ffffff",
    "--od-color-primary": "#3b82f6",
    "--od-color-primary-dark": "#2563eb",
    "--od-color-secondary": "#60a5fa",
    "--od-color-hover": "#60a5fa",
    "--od-color-link": "#3b82f6",
    "--od-color-divider": "#e5e7eb",
    "--od-color-text-title": "#111827",
    "--od-color-text-body": "#374151",
    "--od-color-text-muted": "#6b7280",
    "--od-color-accent-paper": "#f9fafb",
    "--od-color-reflection-bg": "#ffffff",
    "--od-color-star": "#fbbf24",
    "--od-color-success": "#10b981",
    "--od-color-error": "#ef4444",
    "--od-color-warning": "#f59e0b",
    "--od-space-xs": "12px",
    "--od-space-sm": "24px",
    "--od-space-md": "40px",
    "--od-space-lg": "56px",
    "--od-space-xl": "72px",
    "--od-radius-sm": "16px",
    "--od-radius-md": "36px",
    "--od-radius-lg": "48px",
    "--od-radius-full": "999px",
    "--od-font-family": "system-ui, -apple-system, 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif",
    "--od-font-weight-regular": "400",
    "--od-font-weight-medium": "500",
    "--od-font-weight-bold": "700"
  }
}
```

---

**本规范是 OpenDesign 品牌指南的一部分，用于指导主题的设计与实现。**
