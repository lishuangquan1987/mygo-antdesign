# AGENTS.md

> 面向 AI 编码助手（Copilot / CodeBuddy / Claude Code 等）的项目指南。
> 修改本仓库代码前请先读完本文件与 `docs/06-development-plan.md` 的第 4 节。

## 1. 项目是什么

`github.com/lishuangquan1987/mygo-antdesign` 是 [MyGo](https://github.com/egoist/mygo) 原生 UI（`package ui`）的 **Ant Design 风格控件库**。目标：

1. 为 MyGo 原生控件提供 Ant Design 视觉与交互（换肤）。
2. 用 MyGo 的 Base / Element 补齐 antd 有而 mygo 无的控件。
3. 所有控件共用主题，支持运行时切换亮/暗、品牌色、密度。
4. 提供标准与紧凑两套密度（上位机 / IDE 场景）。

当前状态：**设计阶段，尚未编写实现代码**。先读 `docs/README.md`（若有）与 `docs/01`–`docs/06` 再动手。

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
- **禁止**跨控件硬编码颜色/间距，禁止控件之间循环依赖。

### 2.3 API 约定

- 函数式选项：`type Option func(*config)`，如 `ant.Button(c, "保存", ant.Primary(), ant.Size(ant.Small))`。
- 控件返回 `ui.Element`，与 MyGo 组合方式一致。
- 有状态项用 `*T` 绑定 + `ui.Local`，与 MyGo 原生一致。
- 签名贴近 MyGo 原生（如 `Button(c, label)`），降低迁移成本。

## 3. 常用命令

```bash
# 开发（Gallery 预览）
go tool mygo dev

# 构建
go tool mygo build

# 测试
go test ./...

# 静态一致性检查（CI 门禁）
go run ./internal/check

# 重新生成图标
go generate ./icon
```

> 前置：Go 1.27.3+；MyGo 版本在 `go.mod` 中锁定，勿随意升级，升级需跑全量视觉回归。

## 4. 开发工作流

1. 对照 `docs/02-widget-mapping.md` 或 `docs/03-gap-widgets.md` 确认控件归属与复杂度。
2. 按 `docs/06-development-plan.md` 第 4.3 节 DoD 检查表逐项实现。
3. 在 `gallery/` 添加该控件条目（含标准/紧凑两档预览）。
4. 生成或更新截图基线 `widgets/golden/`，确保 4 种主题组合（亮/暗 × 标准/紧凑）通过。
5. 提交前确认 CI 全绿：`go test ./...` + 静态检查 + 截图 diff。

## 5. 不要做的事

- 不要 Fork 或修改 MyGo 源码；本库是上层包装。
- 不要引入 cgo 依赖；保持纯 Go 交叉编译能力。
- 不要写死任何颜色或间距（见 2.1）。
- 不要伪造纯 Web 概念的同名控件（App / ConfigProvider / Util）。
- 不要跳过截图回归与静态检查直接提交控件代码。
- 不要在文档中用占位图片 URL（如 example.com）；没有真实资源就用文字说明。

## 6. 关键文档索引

| 文件 | 用途 |
| --- | --- |
| `docs/01-feasibility.md` | 可行性、架构分层、风险 |
| `docs/02-widget-mapping.md` | 原生控件 → antd 换肤映射表 |
| `docs/03-gap-widgets.md` | 缺口控件实现路径与优先级 |
| `docs/04-theme-system.md` | Token 映射、运行时切换设计 |
| `docs/05-compact-density.md` | 双密度规范与数值表 |
| `docs/06-development-plan.md` | WBS、里程碑、时间线、DoD |

## 7. 代码风格

- 遵循标准 Go 风格（`gofmt` / `go vet`）。
- 中文文档，代码注释用英文或中文均可，保持文件内一致。
- 提交信息用中文，格式：`docs: ...` / `feat(widgets): ...` / `fix(theme): ...`。
