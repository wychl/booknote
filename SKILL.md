---
name: booknote
description: 基于豆瓣书籍信息生成读书笔记卡片，支持多种写作风格和视觉主题，输出 HTML/PNG 及结构化数据，适配抖音/微信等平台。
version: 1.4.0
author: 山海作品推荐
---

# booknote Skill

## 描述

`booknote` 可以：

- 读取豆瓣书籍信息（书名、作者、评分、简介等），调用 DeepSeek 大语言模型生成笔记正文、推荐标签和背景音乐，并渲染为 **1080×1920** 的精美读书笔记卡片。
- 直接生成金句卡片（单句名言+出处），同样支持所有主题和截图。

适用于读书博主、知识分享者生成抖音、小红书、微信等平台的竖版图文卡片。

## 触发场景

当用户要求为某本书生成：

- 读书笔记 / 读书卡片 / 读书分享
- 金句摘抄 + 个人感悟
- 抖音/小红书风格的竖版图文卡片

或要求生成单句名言卡片时，应使用此技能。

## 调用模式

- **执行器**: `command`
- **命令**:
  - `booknote card [书籍ID或书名] [选项]` – 生成读书笔记卡片
  - `booknote quote [金句内容] [选项]` – 生成金句卡片
- **输入**: 通过命令行参数传递书籍信息和选项
- **输出**:
  - **stdout**: JSON 对象，包含生成结果（成功/失败、文件路径、元数据等）
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

**如何获取 DOUBAN_COOKIE**：登录 <https://book.douban.com> → F12 → Network → 刷新 → 点任意请求 → 复制 Request Headers 中的 Cookie 完整字符串。

### 3. 安全读取 .env 的步骤（必须照做）

**步骤 1**：检查文件存在

```bash
if [ ! -f "$HOME/.openclaw/workspace/skills/booknote/.env" ]; then
  echo "❌ .env 不存在" >&2; exit 1
fi
```

**步骤 2**：用 `grep + sed` 安全读取（兼容内部双引号）

```bash
export DEEPSEEK_API_KEY="$(grep '^DEEPSEEK_API_KEY=' "$HOME/.openclaw/workspace/skills/booknote/.env" \
  | sed 's/^DEEPSEEK_API_KEY=//; s/^"//; s/"$//')"
export DOUBAN_COOKIE="$(grep '^DOUBAN_COOKIE=' "$HOME/.openclaw/workspace/skills/booknote/.env" \
  | sed 's/^DOUBAN_COOKIE=//; s/^"//; s/"$//')"
```

> **不要使用 `source .env`**，因为 DOUBAN_COOKIE 内部双引号会导致 shell 解析错误。

**步骤 3**：验证变量非空 + 格式检查

```bash
if [ -z "$DEEPSEEK_API_KEY" ] || [ -z "$DOUBAN_COOKIE" ]; then
  echo "❌ .env 中缺少必要变量" >&2; exit 1
fi
if [[ "$DEEPSEEK_API_KEY" != sk-* ]]; then
  echo "⚠️ DEEPSEEK_API_KEY 格式异常" >&2
fi
```

**步骤 4**：执行命令时显式传递

```bash
DEEPSEEK_API_KEY="$DEEPSEEK_API_KEY" DOUBAN_COOKIE="$DOUBAN_COOKIE" \
booknote card --book="活着" --style="深刻解读" --theme="gothic-dark" --output="/path/to/output"
```

### 4. 输出目录

- 建议使用绝对路径：`$HOME/.openclaw/workspace/skills/booknote/output/`
- `booknote` 的行为：当 `--output=path/to/prefix` 时，实际创建 `path/to/prefix/<bookid>.html`，因此需要先创建 `path/to/prefix/` 目录。
- 推荐：`mkdir -p "$OUTPUT_DIR"`

### 5. 安装 booknote

```bash
go install github.com/wychl/booknote/cmd/booknote@latest
```

确保 `$GOPATH/bin` 在 `PATH` 中。

---

## 参数映射

### 读书笔记卡片 (`card`)

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--book` 或位置参数 | string | **必填** | 书籍名称或豆瓣 ID（推荐用 ID） |
| `--style` | string | `"感性推荐"` | 写作风格：`感性推荐` `深度解读` `极简金句` `故事叙事` |
| `--theme` | string | 自动推荐 | 视觉主题 ID |
| `--output` | string | 书籍 ID | 输出文件前缀（不含扩展名） |
| `--from-json` | string | `""` | 从已有 JSON 加载（跳过 AI 调用） |

### 金句卡片 (`quote`)

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `text` 或位置参数 | string | **必填** | 金句内容 |
| `--source` | string | `""` | 出处，例如 `《追忆似水年华》普鲁斯特` |
| `--theme` | string | `"black-gold"` | 视觉主题 ID |
| `--output` | string | `"quote"` | 输出文件前缀（不含扩展名） |

---

## 完整可用主题

执行以下命令可获取当前版本支持的所有主题（实时更新）：

```bash
booknote themes
```

当前全部主题（共 13 个）：

| ID | 名称 | 描述 | 适用场景 |
|----|------|------|----------|
| `black-gold` | 黑金尊享 | 深邃黑背景，暖金点缀，华丽沉稳 | 高端卡片、金句、品牌展示 |
| `classic-dark` | 经典深色·精修版 | 纯黑底深色主题，极致对比 | 深夜阅读、长时间使用 |
| `midnight-purple` | 紫夜幻境 | 深邃紫黑背景，霓虹紫点缀，充满未来感 | 科幻、奇幻、创意类 |
| `neon-purple` | 霓虹紫 | 深紫底配亮紫，炫酷夜店风 | 潮流、音乐、亚文化 |
| `gothic-shadow` | 哥特暗影 | 深灰背景，暗红点缀，神秘压抑 | 悬疑、恐怖、暗黑系 |
| `cyberpunk` | 赛博朋克 | 电光青与深蓝黑背景，霓虹高对比 | 赛博朋克、反乌托邦、科幻 |
| `electric-blue` | 电光蓝 | 深蓝黑底，亮蓝闪电感，动感刺激 | 运动、电竞、科技类 |
| `ocean-deep` | 深海沉蓝 | 深邃海底暗蓝色调，荧光青绿点缀 | 科幻、海洋主题 |
| `serene-blue` | 静谧蓝 | 深蓝调安静深邃 | 哲学、心理类 |
| `transformative-teal` | 变革碧涛 | 清透蓝绿色背景，自然宁静 | 心理学、个人成长 |
| `pure-dawn` | 极简晨光 | 明亮纯白背景，天空蓝强调色 | 干货、商业类 |
| `vintage-parchment` | 复古羊皮 | 仿古羊皮纸色调，暖黄与深褐 | 经典文学、历史 |
| `warm-autumn-wood` | 暖木秋色 | 温暖木调背景，琥珀色细节 | 文学、散文、情感 |

主题预览图片位于：`~/.openclaw/workspace/skills/booknote/image/<theme-id>.png`

> 💡 主题会持续新增，建议始终使用 `booknote themes` 获取最新列表。

---

## 命令模板

### 生成读书笔记卡片（推荐流程）

```bash
# 1. 读取 .env
export DEEPSEEK_API_KEY="$(grep '^DEEPSEEK_API_KEY=' "$HOME/.openclaw/workspace/skills/booknote/.env" | sed 's/^DEEPSEEK_API_KEY=//; s/^"//; s/"$//')"
export DOUBAN_COOKIE="$(grep '^DOUBAN_COOKIE=' "$HOME/.openclaw/workspace/skills/booknote/.env" | sed 's/^DOUBAN_COOKIE=//; s/^"//; s/"$//')"

# 2. 创建输出目录
OUTPUT_DIR="$HOME/.openclaw/workspace/skills/booknote/output/huozhe"
mkdir -p "$OUTPUT_DIR"

# 3. 执行
DEEPSEEK_API_KEY="$DEEPSEEK_API_KEY" DOUBAN_COOKIE="$DOUBAN_COOKIE" \
booknote card --book="活着" \
  --style "深刻解读" \
  --theme "gothic-dark" \
  --output "$OUTPUT_DIR"
```

### 生成金句卡片（无需环境变量）

```bash
# 直接执行，无需 .env 环境变量
OUTPUT_DIR="$HOME/.openclaw/workspace/skills/booknote/output/quote"
mkdir -p "$OUTPUT_DIR"

booknote quote --text "凡有所相，皆是虚妄。" \
  --source "《金刚经》" \
  --theme "black-gold" \
  --output "$OUTPUT_DIR/diamond-sutra"
```

### 从 JSON 恢复卡片（跳过 AI）

```bash
booknote card --from-json "$OUTPUT_DIR/4913064.json" --output "$OUTPUT_DIR"
```

---

## 关键区别

| 命令 | 是否需要 DeepSeek Key | 是否需要豆瓣 Cookie | 用途 |
|------|-----------------------|---------------------|------|
| `booknote card` | ✅ 需要 | ✅ 需要 | 生成读书笔记，含 AI 评语、标签、BGM |
| `booknote quote` | ❌ 不需要 | ❌ 不需要 | 生成金句卡片，纯渲染文本+出处 |

---

## 输出解析

### 读书笔记成功时

```json
{
  "success": true,
  "error": "",
  "html": "/absolute/path/to/4913064.html",
  "image": "/absolute/path/to/4913064.png",
  "json_file": "/absolute/path/to/4913064.json",
  "book": { "id": "4913064", "title": "活着", "author": "余华", "rating": "9.4", ... },
  "note": "笔记正文（含 ****高亮**** 标记）...",
  "theme": "gothic-dark",
  "style": "故事叙事",
  "tags": ["#小说", "#苦难与生存", "#沉重感人"],
  "bgm": ["Moby - Porcelain", "李宗盛 - 山丘", "Leonard Cohen - Famous Blue Raincoat"]
}
```

### 金句卡片成功时

```json
{
  "success": true,
  "error": "",
  "html": "/absolute/path/to/quote.html",
  "image": "/absolute/path/to/quote.png",
  "theme": "black-gold",
  "tags": null,
  "bgm": null
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

1. **告知用户**：列出生成的文件路径（HTML、PNG、JSON）。
2. **展示关键信息**：读书笔记卡片需展示 tags、bgm、note 前 100 字；金句卡片展示原文和出处。
3. **修复加粗标记**（仅读书笔记）：`booknote` 生成 HTML 中 `****内容****` 可能渲染为 `<strong></strong>内容<strong></strong>`，需手动修复：

```bash
# macOS
sed -i '' 's/<strong><\/strong>\(.*\)<strong><\/strong>/<strong>\1<\/strong>/g' output.html
# Linux
sed -i 's/<strong><\/strong>\(.*\)<strong><\/strong>/<strong>\1<\/strong>/g' output.html
```

1. **提供后续操作**：建议用户将图片发布到抖音/小红书等平台。

### 失败时

| 错误 | 原因 | 处理 |
|------|------|------|
| `api key should not be blank` | `DEEPSEEK_API_KEY` 未设置 | 检查 `.env` |
| `Authentication Fails` | Key 无效/过期 | 检查 Key 完整性和额度 |
| `DOUBAN_COOKIE required` | Cookie 未设置 | 在 `.env` 中添加 |
| `HTTP 404` | 书籍详情获取失败 | 改用豆瓣 ID |
| `no such file or directory` | `--output` 目录不存在 | 先 `mkdir -p` |
| `command not found` | 未安装 booknote | `go install ...` |
| `failed to take screenshot` | Chrome 不可用 | 仅使用 HTML，或手动截图 |

---

## 常见问题

### Q: 为什么 output 目录写成了子目录？

booknote 行为：`--output=path/to/foo` → 输出为 `path/to/foo/<bookid>.{html,png,json}`。务必先创建该目录。

### Q: 是否可以不使用 .env？

可以手动传入环境变量，但 Cookie 很长且易错，不推荐。

### Q: 图片生成失败怎么办？

```bash
# 手动截图
"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" \
  --headless --disable-gpu --screenshot=output.png --window-size=1080,1920 \
  "file:///path/to/output.html"
```

---

## 示例对话

**用户**: 帮我做一张《活着》的读书笔记卡片，风格用故事叙事，主题用哥特暗黑。

**AI 处理步骤**:

1. 检查 `.env` 文件，用 `grep`+`sed` 读取 `DEEPSEEK_API_KEY` 和 `DOUBAN_COOKIE`。
2. 创建输出目录 `mkdir -p output/huozhe`。
3. 执行命令：

   ```bash
   DEEPSEEK_API_KEY=... DOUBAN_COOKIE=... booknote card --book="活着" --style="故事叙事" --theme="gothic-dark" --output="output/huozhe" --image
   ```

4. 解析 stdout JSON，检查 `success`。
5. 若成功，输出文件路径、标签、BGM；并检查 HTML 中空 `<strong>` 问题（如有则用 `sed` 修复）。
6. 回复用户：“《活着》的故事叙事卡片已生成！HTML: ..., PNG: ..., 标签: ..., BGM: ...”

**用户**: 生成一张金句卡片：“凡有所相，皆是虚妄。”出自《金刚经》，用黑金主题。

**AI 处理步骤**:

1. 创建输出目录 `mkdir -p output/jinju`。
2. 直接执行（无需 `.env` 环境变量）：

```bash
mkdir -p output/jinju
booknote quote --text "凡有所相，皆是虚妄。" \
  --source "《金刚经》" \
  --theme "black-gold" \
  --output output/jinju/blackgold-quote
```

> ⚠️ `--source` 书名号、引号直接传入即可，无需转义。

3. 解析 stdout JSON，检查 `success`。
4. 告知用户图片路径，建议在访达中查看。
