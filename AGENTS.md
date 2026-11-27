# AGENTS.md

> 面向 AI 编码助手（Copilot / CodeBuddy / Claude Code 等）的项目指南。
> 修改本仓库代码前请先读完本文件与 `docs/06-development-plan.md` 的第 4 节；写控件时另读 `docs/07-samples-gallery.md`。

## 1. 项目是什么

`github.com/lishuangquan1987/mygo-antdesign` 是 [MyGo](https://github.com/egoist/mygo) 原生 UI（`package ui`）的 **Ant Design 风格控件库**。目标：

1. 为 MyGo 原生控件提供 Ant Design 视觉与交互（换肤）。
2. 用 MyGo 的 Base / Element 补齐 antd 有而 mygo 无的控件。
3. 所有控件共用主题，支持运行时切换亮/暗、品牌色、密度。
4. 提供标准与紧凑两套密度（上位机 / IDE 场景）。
5. 每个控件配可运行示例（不同主题/尺寸/密度/状态），由 Gallery 展示。

当前状态：**v0.1.0 已发布，v0.2 已推送待打 tag，v0.3 开发中**——theme/、widgets/（包 ant）、extra/、icon/、samples/、gallery/ 均已落地并持续扩展（144 个控件构造函数、13 个示例分组 / 83 个示例片段）。先读 `AGENTS.md` 本文件、`docs/06-development-plan.md`（含 WBS 状态与进度章节）与现有控件实现风格再动手；文档中的「状态」列必须与代码一致，不要写成已完成。

## 2. 必须遵守的硬规则

### 2.1 样式铁律（违反会导致主题/密度切换失效）

1. **颜色只能来自主题**：控件内部一律用 `c.Theme()` 的字段（`Text`、`TextMuted`、`Surface`、`Border`、`Accent`、`Danger`、`Warning`、`Success`、`Focus`、`Inverse` 等），或 `internal/variant.go` 的映射函数。禁止在 `widgets/`、`extra/` 中写裸 `ui.Hex("#...")`。图标数据（`icon/`）是唯一豁免。
2. **间距只能走主题度量**：所有 `Padding`/`Gap`/`Margin`/`Size` 用 `t.Space(n)` 或 `t.Rem(n)`，禁止裸像素字面量。1px 描边、焦点环宽度等视觉基准例外。
3. **状态样式集中管理**：hover / pressed / disabled / selected / focus 的颜色映射统一放 `internal/variant.go`，不要在每个控件里各写一套分支。
4. **主题单一真源**：控件只读 `c.Theme()`，不直接读 `AntToken`，不自行 `SetTheme`（主题切换由应用/theme 层负责）。

### 2.2 架构边界

- `widgets/`：与 MyGo 原生控件一一对应的换肤实现，只依赖 `theme/` 与 MyGo `ui`。
- `extra/`：缺口控件（antd 有 mygo 无），同样只依赖 `theme/` 与 `ui`。
- `theme/`：`AntToken` → `ui.Theme` 的唯一桥接层，含 `Apply`（全局）/`Sub`（局部 `ui.Themed`）。
- `internal/`：共享私有 helper（状态映射、度量、静态检查）。
- `icon/`：从 antd SVG 生成的 Go 路径常量，带 `go:generate`。
- `samples/`：按分类分组的示例源码（真实可编译代码，`go:embed` 提取片段），在 `registry.go` 定义、由 `gallery/main.go` 显式注册。
- `gallery/`：示例展示应用，遍历注册表渲染；**不硬编码任何控件示例清单**。
- **禁止**跨控件硬编码颜色/间距，禁止控件之间循环依赖。

### 2.3 API 约定

- 函数式选项：`type Option func(*config)`，如 `ant.Button(c, "保存", ant.Primary(), ant.Size(ant.Small))`。
- 控件返回 `ui.Element`，与 MyGo 组合方式一致。
- 有状态项用 `*T` 绑定 + `ui.Local`，与 MyGo 原生一致。
- 签名贴近 MyGo 原生（如 `Button(c, label)`），降低迁移成本。

## 3. 常用命令

```bash
# 运行 Gallery 示例应用（开发预览）
go run ./gallery

# 构建 Gallery 可执行文件
go build -trimpath -ldflags "-X main.version=dev" -o bin/gallery ./gallery

# 构建全部包
go build ./...

# 测试（含 Gallery 与逐示例截图基线）
go test -count=1 ./...

# 静态一致性检查（CI 门禁）
go run ./internal/check

# 重新生成图标
go generate ./icon
```

> 可选：仓库根目录提供 `Makefile`（`make run` / `make build` / `make check` / `make baselines`）；Windows 可用 Git Bash/WSL 或直接执行上面的 go 命令。
>
> 前置：Go 1.27.4+；MyGo 版本在 `go.mod` 中锁定（当前 v0.10.0），勿随意升级，升级需跑全量视觉回归。

## 4. 开发工作流

1. 对照 `docs/02-widget-mapping.md` 或 `docs/03-gap-widgets.md` 确认控件归属与复杂度。
2. 按 `docs/06-development-plan.md` 第 4.3 节 DoD 检查表逐项实现。
3. **立即为该控件写示例**：在 `samples/<分组>/demo.go` 下写真实可编译示例（基础 + 主要变体 + 禁用 + 三档尺寸，见 `docs/07-samples-gallery.md` 第 4 节矩阵），并在 `gallery/main.go` 调用 `samples.Register(...)` 注册；不要等到后期补写。
4. 示例本身也遵守样式铁律（2.1），代码片段用 `// sample:xxx` 标记包起来供 `go:embed` 提取。
5. 生成或更新截图基线 `samples/golden/`，确保 3 种主题组合（亮标准、暗标准、亮紧凑）通过。
6. **按第 5 节完成自测 → 提交 → 推送**；未通过自测不得提交，提交后必须推送。

## 5. 任务收尾：自测 → 提交 → 推送（强制）

> 每完成一个任务（无论大小），必须走完以下流程。这是硬要求，不是可选步骤。
> 判定标准：一个任务 = WBS 中的一个任务 ID（如 B1）或一个独立修复/文档变更。

### 5.1 自测（全部通过才算完成）

按顺序执行，任一失败先修复再继续：

```bash
gofmt -l .                    # 输出为空，否则先 gofmt -w
go vet ./...                  # 无报错
go test ./...                 # 全部通过
go run ./internal/check       # 样式铁律静态检查通过
go test ./gallery -run TestGalleryBaselines # 示例截图 diff 通过（改控件/主题/示例时必跑）
```

补充人工确认（按需）：

- [ ] 改了控件：在 Gallery 里目视该控件，确认侧栏可见、示例可切换主题/密度/尺寸。
- [ ] 改了主题或度量：跑一遍 3 种主题组合（亮标准、暗标准、亮紧凑）截图，确认无遗漏旧色。
- [ ] 新增控件：确认 `samples/<分组>/demo.go` 示例已写，并在 `gallery/main.go` 注册。
- [ ] 新增依赖：确认为纯 Go、无 cgo，`go.sum` 已更新。

### 5.2 提交

- **一个任务一个提交**，不要把多个无关任务堆进一次提交；也不要一个任务拆成多次碎片提交。
- 先 `git add -A`，再提交；提交前用 `git status` 确认没有误纳入的文件（临时文件、截图大文件等）。
- 提交信息格式：`<type>: <中文描述>`，type 取 `feat` / `fix` / `docs` / `test` / `refactor` / `chore`。
  - 示例：`feat(widgets): 实现 Button 六种变体及示例` / `fix(theme): 修正紧凑档焦点环偏移`。

### 5.3 推送

```bash
git push origin main
```

- 提交后**必须立即推送**到远端，不要留到本地积压多个提交。
- 推送失败（如临时 TLS/网络错误）先重试一次；仍失败则如实报告错误，**不要声称已推送成功**。
- 推送后用 `git status` 确认工作区干净、`origin/main` 与本地一致。

### 5.4 收尾顺序

```
实现 → 自测全绿 → git add → git commit → git push → 确认干净 → 报告结果
```

> 未完成自测就提交、或提交后不推送，都视为任务未完成。

## 6. 不要做的事

- 不要 Fork 或修改 MyGo 源码；本库是上层包装。
- 不要引入 cgo 依赖；保持纯 Go 交叉编译能力。
- 不要写死任何颜色或间距（见 2.1）。
- 不要伪造纯 Web 概念的同名控件（App / ConfigProvider / Util）。
- 不要跳过截图回归与静态检查直接提交控件代码。
- 不要未自测就提交，或提交后不推送（见第 5 节）。
- 不要在文档中用占位图片 URL（如 example.com）；没有真实资源就用文字说明。

## 7. 关键文档索引

| 文件 | 用途 |
| --- | --- |
| `docs/01-feasibility.md` | 可行性、架构分层、风险 |
| `docs/02-widget-mapping.md` | 原生控件 → antd 换肤映射表 |
| `docs/03-gap-widgets.md` | 缺口控件实现路径与优先级 |
| `docs/04-theme-system.md` | Token 映射、运行时切换设计 |
| `docs/05-compact-density.md` | 双密度规范与数值表 |
| `docs/06-development-plan.md` | WBS、里程碑、时间线、DoD |
| `docs/07-samples-gallery.md` | 每控件示例的四维覆盖、Gallery 结构、代码片段与截图联动 |

## 8. 代码风格

- 遵循标准 Go 风格（`gofmt` / `go vet`）。
- 中文文档，代码注释用英文或中文均可，保持文件内一致。
- 提交信息用中文，格式：`docs: ...` / `feat(widgets): ...` / `fix(theme): ...`。
