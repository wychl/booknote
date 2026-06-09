---
name: booknote
description: 基于豆瓣书籍信息生成读书笔记卡片，支持多种写作风格和视觉主题，输出 HTML/PNG 及结构化数据，适配抖音/微信等平台。
version: 1.2.0
author: wychl
---

# booknote Skill

## 描述

`booknote` 读取豆瓣书籍信息（书名、作者、评分、简介等），调用 DeepSeek 大语言模型生成笔记正文、推荐标签和背景音乐，并渲染为 **1080×1920** 的精美卡片。**自动输出 HTML 和 PNG 图片**（通过 chromedp 无头浏览器截图），适用于读书博主生成分享卡片，适配抖音、小红书、微信等平台。

## 触发场景

当用户要求为某本书生成：

- 读书笔记 / 读书卡片 / 读书分享
- 金句摘抄 + 个人感悟
- 抖音/小红书风格的竖版图文卡片

时应使用此技能。

## 调用模式

- **执行器**: `command`
- **命令**: `booknote card [书籍ID或书名] [选项]`
- **输入**: 通过命令行参数传递书籍信息和选项
- **输出**:
  - **stdout**: JSON 对象，包含完整的生成结果（成功/失败、文件路径、书籍元数据、笔记内容等）
  - **stderr**: 人类可读的日志信息（进度、错误等）

---

## ⚠️ 前置条件

### 1. 所需环境变量

| 变量名 | 用途 | 来源 |
|--------|------|------|
| `DEEPSEEK_API_KEY` | AI 生成笔记正文 | DeepSeek 账户 |
| `DOUBAN_COOKIE` | 获取豆瓣书籍详情 | 浏览器 Cookie |

这两个变量**不会自动从系统环境继承**，必须从 `.env` 文件读取并显式注入。

### 2. `.env` 文件

- **路径**: `$HOME/.openclaw/workspace/skills/booknote/.env`
- **格式**: 标准 `KEY="VALUE"`，每行一个

```bash
DEEPSEEK_API_KEY="sk-xxxxxxxxxxxxxxxxxxxx"
DOUBAN_COOKIE="bid=xxx; ll=\"108296\"; dbcl2=\"xxxx\"; ..."
```

> ⚠️ **关键陷阱**：DOUBAN_COOKIE 内部包含双引号（如 `ll="108296"`、`dbcl2="..."`），`.env` 文件必须用**外层双引号+转义内部双引号**。**不要外层套单引号**，否则 `source` 会报错。

**如何获取 DOUBAN_COOKIE**：登录 <https://book.douban.com> → F12 → Network → 刷新 → 点任意请求 → 复制 Request Headers 中的 Cookie 完整字符串。

### 3. 安全读取 .env 的步骤（必须照做）

**步骤 1**：检查文件存在

```bash
if [ ! -f "$HOME/.openclaw/workspace/skills/booknote/.env" ]; then
  echo "❌ .env 不存在" >&2
  exit 1
fi
```

**步骤 2**：用 `grep + sed` 安全读取（兼容内部双引号）

```bash
# 提取外层引号内的值，去掉 KEY 前缀和首尾引号
export DEEPSEEK_API_KEY="$(grep '^DEEPSEEK_API_KEY=' "$HOME/.openclaw/workspace/skills/booknote/.env" \
  | sed 's/^DEEPSEEK_API_KEY=//; s/^"//; s/"$//')"
export DOUBAN_COOKIE="$(grep '^DOUBAN_COOKIE=' "$HOME/.openclaw/workspace/skills/booknote/.env" \
  | sed 's/^DOUBAN_COOKIE=//; s/^"//; s/"$//')"
```

> **不要使用 `source .env`**，因为 DOUBAN_COOKIE 内部双引号会导致 shell 解析错误。

**步骤 3**：验证变量非空 + 格式检查

```bash
# 空值检查
if [ -z "$DEEPSEEK_API_KEY" ] || [ -z "$DOUBAN_COOKIE" ]; then
  echo "❌ .env 中缺少 DEEPSEEK_API_KEY 或 DOUBAN_COOKIE" >&2
  exit 1
fi

# Key 前缀检查（防止截断错误）
if [[ "$DEEPSEEK_API_KEY" != sk-* ]]; then
  echo "⚠️ DEEPSEEK_API_KEY 格式异常（不以 sk- 开头）" >&2
  echo "   实际值前20字符: ${DEEPSEEK_API_KEY:0:20}" >&2
fi

# DOUBAN_COOKIE 长度检查（正常应在 800-1500 字符）
if [ "${#DOUBAN_COOKIE}" -lt 200 ]; then
  echo "⚠️ DOUBAN_COOKIE 过短（${#DOUBAN_COOKIE} 字符），可能不完整" >&2
fi
```

**步骤 4**：执行命令时显式传递

```bash
DEEPSEEK_API_KEY="$DEEPSEEK_API_KEY" DOUBAN_COOKIE="$DOUBAN_COOKIE" \
booknote card --book="活着" --style="故事叙事" --theme="gothic-dark" \
  --output="$HOME/.openclaw/workspace/skills/booknote/output/huozhe"
```

### 4. 输出目录

- 建议使用绝对路径：`$HOME/.openclaw/workspace/skills/booknote/output/`
- `booknote` 的行为：当 `--output=path/to/prefix` 时，实际创建 `path/to/prefix/<bookid>.html`
  - 所以需要预先确保 `path/to/prefix/` 目录**存在**
  - 推荐 mkdir 后再执行：`mkdir -p "$HOME/.openclaw/workspace/skills/booknote/output/huozhe"`

### 5. 自动截图（chromedp）

`booknote` 使用 `chromedp` 自动截图生成 PNG，需满足：

- 系统存在可用的 Chrome/Chromium 实例（自动查找）
- 若失败，可设置 `CHROME_PATH` 指定浏览器路径
- 截图失败不影响 HTML/JSON 生成

### 6. 安装 booknote

```bash
go install github.com/wychl/booknote/cmd/booknote@latest
```

确保 `$GOPATH/bin` 或 `$GOBIN` 在 `PATH` 中。

---

## 参数映射

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--book` 或位置参数 | string | **必填** | 书籍名称或豆瓣 ID（推荐用 ID） |
| `--style` | string | `"感性推荐"` | 写作风格：`感性推荐` `深度解读` `极简金句` `故事叙事` |
| `--theme` | string | 自动推荐 | 视觉主题 ID |
| `--output` | string | 书籍 ID | 输出文件前缀（不含扩展名） |
| `--from-json` | string | `""` | 从已有 JSON 加载（跳过 AI 调用） |

**注意**：
- **推荐用豆瓣 ID**：`booknote card 4913064`，避免同名搜索歧义
- 书名含空格需引号：`--book="活着"`
- **输出实际为 `--output/<bookid>.{html,png,json}`**，需先创建 `--output` 目录

---

## 完整可用主题

以 `booknote card --help` 输出为准。当前已知主题：

| 主题 ID | 名称 | 适用场景 |
|---------|------|----------|
| `pure-white` | 极简白 | 干货、书评、实用类 |
| `vintage-paper` | 复古纸 | 经典文学、历史 |
| `warm-wood` | 暖木调 | 文学、散文、情感 |
| `classic-dark-gold` | 黑金质感 | 金句、励志、商业 |
| `dark-purple` | 黑紫幻夜 | 科幻、奇幻、创意 |
| `gothic-dark` | 哥特暗黑 | 悬疑、恐怖、暗黑系 |
| `transformative-teal` | 变革蓝绿 | 心理学、个人成长 |

---

## 命令模板

### 完整执行流程（推荐）

```bash
# 1. 检查 .env
if [ ! -f "$HOME/.openclaw/workspace/skills/booknote/.env" ]; then
  echo "❌ .env 不存在" >&2; exit 1
fi

# 2. 读取变量
export DEEPSEEK_API_KEY="$(grep '^DEEPSEEK_API_KEY=' "$HOME/.openclaw/workspace/skills/booknote/.env" | sed 's/^DEEPSEEK_API_KEY=//; s/^"//; s/"$//')"
export DOUBAN_COOKIE="$(grep '^DOUBAN_COOKIE=' "$HOME/.openclaw/workspace/skills/booknote/.env" | sed 's/^DOUBAN_COOKIE=//; s/^"//; s/"$//')"

# 3. 验证
[ -z "$DEEPSEEK_API_KEY" ] && { echo "❌ KEY 为空" >&2; exit 1; }
[ -z "$DOUBAN_COOKIE" ] && { echo "❌ COOKIE 为空" >&2; exit 1; }

# 4. 创建输出目录（booknote 会写为 OUTPUT_DIR/<bookid>.html）
OUTPUT_DIR="$HOME/.openclaw/workspace/skills/booknote/output/huozhe"
mkdir -p "$OUTPUT_DIR"

# 5. 执行
DEEPSEEK_API_KEY="$DEEPSEEK_API_KEY" DOUBAN_COOKIE="$DOUBAN_COOKIE" \
booknote card --book="活着" \
  --style "故事叙事" \
  --theme "gothic-dark" \
  --output "$OUTPUT_DIR"
```

### 从 JSON 恢复（跳过 AI 调用）

```bash
DEEPSEEK_API_KEY="$DEEPSEEK_API_KEY" DOUBAN_COOKIE="$DOUBAN_COOKIE" \
booknote card --from-json "$OUTPUT_DIR/4913064.json" --output "$OUTPUT_DIR"
```

---

## 输出解析

### 成功时

```json
{
  "success": true,
  "error": "",
  "html": "/absolute/path/to/prefix/4913064.html",
  "image": "/absolute/path/to/prefix/4913064.png",
  "json_file": "/absolute/path/to/prefix/4913064.json",
  "book": { "id": "4913064", "title": "活着", "author": "余华", "rating": "9.4", ... },
  "note": "笔记正文（含 ****高亮**** 标记）...",
  "theme": "gothic-dark",
  "style": "故事叙事",
  "tags": ["#小说", "#苦难与生存", "#沉重感人"],
  "bgm": ["Moby - Porcelain", "李宗盛 - 山丘", "Leonard Cohen - Famous Blue Raincoat"]
}
```

### 失败时

```json
{
  "success": false,
  "error": "错误描述",
  "html": "", "image": "", "json_file": "",
  "book": {}, "note": "", "theme": "", "style": "",
  "tags": null, "bgm": null
}
```

---

## 处理指南

### 成功时

1. **告知用户**：列出生成的文件路径（HTML、PNG、JSON）
2. **展示关键信息**：tags、bgm、note 前 100 字
3. **修复加粗标记**（按需）：
   - booknote 在 HTML 中把 `****内容****` 渲染为 `<strong></strong>内容<strong></strong>`
   - 检查生成的 HTML 是否包含空 `<strong></strong>`，如有则修复：

```bash
# macOS
sed -i '' 's/<strong><\/strong>\(.*\)<strong><\/strong>/<strong>\1<\/strong>/g' /path/to/output.html

# Linux
sed -i 's/<strong><\/strong>\(.*\)<strong><\/strong>/<strong>\1<\/strong>/g' /path/to/output.html
```

4. 提供后续操作建议（发微信/小红书等）

### 失败时

| 错误 | 原因 | 处理 |
|------|------|------|
| `api key should not be blank` | `DEEPSEEK_API_KEY` 未设置 | 检查 `.env` |
| `Authentication Fails` | Key 无效/过期 | 检查 Key 是否完整、是否被截断 |
| `DOUBAN_COOKIE required` | Cookie 未设置 | 在 `.env` 中添加 |
| `HTTP 404` | 书籍详情获取失败 | 改用豆瓣 ID |
| `no such file or directory` | `--output` 目录不存在 | 先 `mkdir -p` |
| `command not found` | 未安装 booknote | `go install ...` |
| `failed to take screenshot` | Chrome 不可用 | 使用 HTML 手动截图 |

---

## 常见问题

### Key 被截断怎么办？

检查 `.env` 中 key 的实际字符数与 DeepSeek 后台显示的字符数一致。用 `echo ${#DEEPSEEK_API_KEY}` 查看长度。

### 为什么 output 目录写成了子目录？

booknote 行为：`--output=path/to/foo` → 输出为 `path/to/foo/<bookid>.{html,png,json}`。所以先 `mkdir -p path/to/foo`。

### 是否可以不使用 .env？

可以手动传入，但 Cookie 很长且易错，不推荐。

### 图片生成失败怎么办？

```bash
"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" \
  --headless --disable-gpu --screenshot=output.png --window-size=1080,1920 \
  "file:///path/to/output.html"
```

---

## 示例对话

**用户**: 帮我做一张《活着》的读书笔记卡片，风格用故事叙事，主题用哥特暗黑。

**AI 处理步骤**:

1. 检查 `.env` 是否存在，grep/sed 读取变量
2. 验证 KEY 以 `sk-` 开头、COOKIE 长度 > 200
3. 创建输出目录 `mkdir -p output/huozhe`
4. 执行命令并传环境变量
5. 解析 stdout JSON，检查 `success`
6. 检查 HTML 中空 `<strong>`，必要时修复
7. 返回结果给用户
