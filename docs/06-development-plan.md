# 06 · 详细开发计划

> 基于 [01 可行性分析](01-feasibility.md)、[02 控件映射](02-widget-mapping.md)、[03 缺口控件](03-gap-widgets.md)、[04 主题系统](04-theme-system.md)、[05 紧凑密度](05-compact-density.md) 的落地计划。

## 1. 计划总览

### 1.1 目标

交付 `github.com/lishuangquan1987/mygo-antdesign` v0.1.0：

1. 一套 Ant Design 风格的 MyGo 原生 UI 控件库（换肤 + 缺口补齐）。
2. 全部控件共用主题，支持运行时切换亮/暗、品牌色、密度。
3. 提供标准与紧凑两套密度，覆盖上位机 / IDE 场景。
4. **每个控件配可运行示例**，展示不同主题/尺寸/密度/状态，代码与效果对应（见 07 文档）。
5. 具备 Gallery 展示应用、截图回归、静态一致性检查。

### 1.2 范围边界

| 纳入 | 不纳入 |
| --- | --- |
| MyGo 原生控件的 antd 换肤 | Fork / 修改 MyGo 源码 |
| antd 有而 mygo 无的控件补齐 | Web 框架层概念（App/ConfigProvider/Util 同名控件） |
| 主题 Token 与运行时切换 | 独立设计一套全新视觉语言 |
| 图标数据、Gallery、文档 | 商业组件 / Pro 组件 |
| 紧凑密度与 IDE 示例 | 移动端适配（MyGo 已原生支持，后续版本） |

### 1.3 总体量化目标

- 换肤控件：**≥ 55 个**（覆盖 MyGo 原生控件的绝大多数）。
- 缺口控件：**≥ 12 个**（P1/P2 全部 + P3 关键项）。
- 主题预设：4 套（亮/暗 × 标准/紧凑）+ 品牌色可配。
- **示例覆盖率：100%**（每个控件 ≥ 1 个示例页，高频控件覆盖全部变体/状态/尺寸）。
- 视觉回归截图：每个示例 × 4 种主题组合。

## 2. 技术基座与前置条件

| 项 | 要求 | 备注 |
| --- | --- | --- |
| Go | 1.27.3+ | MyGo native UI 的硬性要求 |
| MyGo | 锁定版本区间，如 `v0.7.x` | 0.7 阶段 API 可能变动，随版本跟进 |
| 依赖 | 仅纯 Go 库 | 二维码、上传等引入的库需无 cgo |
| 工具 | `go tool mygo dev/build` | Gallery 用无窗口渲染 PNG 做回归 |
| CI | GitHub Actions | `go test ./...` + 静态检查 + 截图生成 |

**前置动作**：`go mod init github.com/lishuangquan1987/mygo-antdesign`，在 `go.mod` 中固定 MyGo 版本；建立 `docs/decisions/` 记录关键取舍（ADR）。

## 3. 仓库与包结构（目标态）

```
mygo-antdesign/
├── go.mod / go.sum
├── theme/                  # 主题层：唯一视觉真源
│   ├── token.go            # AntToken 定义
│   ├── theme.go            # Option / Density / Apply / Sub
│   ├── bridge.go           # AntToken → ui.Theme 映射
│   ├── preset.go           # Light / Dark / Compact 预设
│   └── theme_test.go
├── internal/
│   ├── variant.go          # 状态 → 颜色/边框 集中映射（hover/pressed/disabled/focus）
│   ├── metrics.go          # 间距/字号 helper（统一 t.Space/t.Rem 用法）
│   └── check.go            # 静态检查辅助
├── samples/                # 控件示例：真实可编译代码，go:embed 提取片段
│   ├── registry.go         # 注册表与分类
│   ├── extract.go          # 代码片段提取工具
│   ├── button/  input/ ... # 每控件一个示例目录
│   └── golden/             # 示例截图基线
├── widgets/                # 换肤控件：与 MyGo 原生一一对应
│   ├── button.go  input.go  select.go  table.go ...
│   └── golden/             # 截图基线
├── extra/                  # 缺口控件
│   ├── pagination.go  card.go  cascader.go  upload.go ...
├── icon/
│   ├── icons.go            # 生成的 SVG 路径常量
│   └── generate.go         # 从 antd SVG 资源生成（go:generate）
├── gallery/                # 示例展示应用
│   ├── main.go             # 侧栏 + 示例区 + 控制条
│   ├── sidebar.go          # 控件目录与搜索
│   ├── controls.go         # 主题/密度/尺寸/并排对比切换器
│   └── render_test.go      # 无窗口渲染 PNG 回归
├── docs/                   # 本文档集
└── .github/workflows/ci.yml
```

**依赖规则**：`widgets/` 和 `extra/` 只能依赖 `theme/` 与 MyGo `ui`，禁止跨控件硬编码颜色/间距；`internal/` 为两者共享的私有helper。

## 4. 统一开发规范（所有任务必须遵守）

### 4.1 API 约定

函数式选项 + 贴近 MyGo 的签名：

```go
// 目标 API 形态（设计示意，非已实现代码）
ant.Button(c, "保存", ant.Primary(), ant.Size(ant.Small), ant.Loading(app.saving))
ant.Table(c, &state, cols, rows, ant.Size(ant.Compact))
extra.Pagination(c, &app.page, app.total)
theme.Apply(c, theme.Light(theme.WithPrimary("#1677ff"), theme.WithDensity(theme.Compact)))
```

- 控件返回 `ui.Element`，与 MyGo 组合方式一致。
- 有状态项用 `*T` 绑定 + `ui.Local`，与 MyGo 原生一致。
- 选项类型统一为 `type Option func(*config)`。

### 4.2 样式铁律（密度/主题可切换的前提）

1. 颜色只能来自 `c.Theme()` 字段或 `internal/variant.go` 的映射函数。
2. 间距只能写 `t.Space(n)` / `t.Rem(n)`，禁止裸像素。
3. 1px 描边、焦点环宽度等视觉基准保持固定。
4. 状态样式（hover/pressed/disabled/selected/focus）集中在 `variant.go`，禁止各控件各写一套。

### 4.3 每个控件的完成检查表（DoD）

- [ ] 签名符合 4.1 API 约定
- [ ] 颜色/间距符合 4.2 铁律（静态检查通过）
- [ ] 标准与紧凑两档截图无异常
- [ ] 亮/暗主题截图无异常
- [ ] 键盘可达、有 Role/Label（Base 实现天然满足，纯组合需补）
- [ ] **示例已编写并注册**（`samples/<widget>/`，覆盖基础 + 主要变体 + 禁用 + 三档尺寸）
- [ ] Gallery 侧栏可见该控件
- [ ] 4 种主题组合的示例截图通过 diff
- [ ] 文档/注释齐全

## 5. WBS 详细任务分解

> 估算单位：人日（1 名熟悉 Go 的开发者）。任务 ID 用于跟踪。
> **进度**（2026-11-26）：A1–A8、F1–F5、B1 已完成并推送；详见文末「当前进度」。

### WS-A · 基础设施与主题层

| ID | 任务 | 估算 | 依赖 | 状态 |
| --- | --- | --- | --- | --- |
| A1 | `go.mod`、目录骨架、`.github/workflows/ci.yml`、ADR 模板 | 1 | — | ✅ |
| A2 | `theme/token.go`：`AntToken`、`Density`、Option 框架 | 2 | A1 | ✅ |
| A3 | `theme/bridge.go`：Token → `ui.Theme` 映射，含 `WithAccent` 色阶生成 | 2 | A2 | ✅ |
| A4 | `theme/preset.go`：亮/暗/紧凑预设、品牌色、密度组合 | 2 | A3 | ✅ |
| A5 | `theme/theme.go`：`Apply`（全局）/`Sub`（局部 `ui.Themed`） | 1 | A3 | ✅ |
| A6 | `internal/variant.go` + `metrics.go`：状态映射与度量 helper | 2 | A3 | ✅ 随 B1 |
| A7 | 静态检查（禁止 widgets/extra 中出现裸 `ui.Hex`/裸间距） | 1 | A1 | ✅ `internal/check` |
| A8 | `gallery/`：应用骨架（空壳启动、无窗口 PNG 渲染入口） | 2 | A4 | ✅ |
| | **小计** | **13** | | |

### WS-B · 换肤控件（widgets/）

| ID | 任务包 | 控件 | 估算 | 依赖 |
| --- | --- | --- | --- | --- |
| B1 | 按钮与排版 | Button（含 primary/dashed/text/link/danger）✅（Typography/Icon/Divider/Space 包装待续） | 4 | A6 |
| B2 | 文本录入 | Input、TextArea、InputNumber、SearchField、Mentions 输入框部分 | 4 | A6 |
| B3 | 选择类 | Select、AutoComplete、Cascader 触发器部分、TreeSelect 触发器部分 | 5 | A6, B1 |
| B4 | 选择控件 | Checkbox、Radio、Switch/Toggle、Slider/RangeSlider、Rate | 4 | A6 |
| B5 | 日期时间 | DatePicker（Calendar+浮层）、TimePicker | 4 | B3 |
| B6 | 表单 | Form（标签/必填/校验信息） | 2 | B2 |
| B7 | 数据展示-基础 | Table、Tree、List、Card、Descriptions、Statistic | 6 | A6 |
| B8 | 数据展示-辅助 | Avatar、Badge、Tag、Progress、Timeline、Empty、Result、Skeleton | 5 | A6 |
| B9 | 导航 | Tabs、Segmented、Steps、Breadcrumb、Menu/Sidebar、Dropdown | 5 | A6, B1 |
| B10 | 折叠/浮层 | Collapse/Accordion、Popover、Tooltip | 3 | A6 |
| B11 | 反馈 | Modal、AlertDialog/Popconfirm、Message/Notification（Toast）、Spin、Watermark | 5 | B10 |
| B12 | 其他 | Splitter、Affix/Sticky、Drawer | 2 | A6 |
| | **小计** | | **44** | |

### WS-C · 缺口控件（extra/）

| ID | 任务包 | 控件 | 估算 | 依赖 |
| --- | --- | --- | --- | --- |
| C1 | P1 简单组合 | Card、Tag、Empty、Result、Divider、Statistic、Timeline（注：若 B 已做则合并） | 3 | A6 |
| C2 | P1 数据展示 | Descriptions、Skeleton | 2 | B8 |
| C3 | P2 分页与穿梭 | Pagination、Transfer | 3 | B1, B7 |
| C4 | P2 浮层 | Drawer、Popconfirm、Watermark、Tour | 6 | B10, B11 |
| C5 | P2 复杂选择 | Cascader、TreeSelect、Mentions | 8 | B3 |
| C6 | P3 依赖能力 | QRCode（引入纯 Go 编码库）、Upload（本地文件 dialog + 拖放）、Spin、FloatButton、BorderBeam、Notification、Anchor、Carousel、Masonry | 12 | B11 |
| | **小计** | | **34** | |

### WS-D · 图标与资源

| ID | 任务 | 估算 | 依赖 |
| --- | --- | --- | --- |
| D1 | `icon/generate.go`：从 antd SVG 资源生成 Go 路径常量 | 2 | A1 |
| D2 | 图标集生成与校验（尺寸随 `t.Rem`） | 3 | D1 |
| | **小计** | **5** | |

### WS-E · 文档、示例与发布

| ID | 任务 | 估算 | 依赖 |
| --- | --- | --- | --- |
| E1 | README 完善、迁移指南（MyGo → 本库） | 2 | M1 |
| E2 | 主题与密度使用文档、API 参考 | 3 | M1 |
| E3 | IDE / 上位机风格示例应用 | 4 | M3 |
| E4 | v0.1.0 发布：CHANGELOG、tag、release 说明 | 2 | M4 |
| | **小计** | **11** | |

### WS-F · 控件示例与 Gallery（见 07 文档）

> 示例不集中在一个阶段写完，而是**随每个控件完成时立即编写**（纳入 WS-B/WS-C 的 DoD）。下表是独立的框架与专项工作；示例编写人日已摊入各控件的实现成本。

| ID | 任务 | 估算 | 依赖 |
| --- | --- | --- | --- |
| F1 | Gallery 界面：侧栏分类/搜索、示例区、顶部控制条 | 4 | A8 |
| F2 | 并排对比模式（`ui.Themed` 亮/暗、标准/紧凑同屏） | 2 | F1, A5 |
| F3 | 代码视图：只读 `ui.CodeEditor` + 复制按钮 | 2 | F1 |
| F4 | 示例注册表 `samples/registry.go` + 分类 | 1 | A1 |
| F5 | `go:embed` 代码片段提取工具 `samples/extract.go` | 2 | F4 |
| F6 | 截图生成与文档嵌图联动（`render_test.go`） | 3 | F1, A8 |
| | **小计** | **14** | |

**总计估算：约 121 人日**（含缓冲前；其中示例编写约 21 人日摊入各控件实现）；单人约 6 个月，2 人并行约 3.5 个月。若示例仅做「基础+主要变体」最小覆盖，可压缩约 8 人日。

## 6. 里程碑与验收标准

| 里程碑 | 周次 | 范围 | 退出标准（Exit Criteria） |
| --- | --- | --- | --- |
| **M0 基建** | W1 | A1–A8 | 仓库骨架、CI 通过、主题层 API 可用、Gallery 空壳能渲染 PNG |
| **M1 主题闭环** | W2–W5 | A 全部 + B1–B6 + F1–F5 | 主题/密度/品牌色切换闭环；基础录入控件 ≥ 20 个且**配好示例**；Gallery 可展示并切换；截图回归通过 |
| **M2 常用控件** | W6–W10 | B7–B12 | 数据展示与反馈控件完成并**配好示例**；覆盖常见后台/工具应用 80% 场景 |
| **M3 密度适配** | W11–W12 | 全控件紧凑审计 + F2/F6 + E3 | 亮/暗 × 标准/紧凑 4 组合**示例截图**全通过；并排对比可用；IDE 示例可运行 |
| **M4 缺口与发布** | W13–W18 | C 全部 + D + E4 | 缺口控件 ≥ 12 个且**配好示例**；示例覆盖率 100%；v0.1.0 发布 |

## 7. 时间线（单人 18 周基线）

```
W1    M0 基建 + 主题层 + Gallery 空壳
W2-5  M1 基础控件（按钮/录入/选择/表单）+ 每控件示例 + Gallery 闭环
W6-10 M2 数据展示/导航/浮层/反馈 + 同步配示例
W11-12 M3 紧凑密度审计 + 并排对比 + 全量示例截图 + IDE 示例
W13-18 M4 缺口控件 + 缺口示例 + 图标 + 文档嵌图 + 发布
```

**并行建议**：若 2 人，A+WS-B 由一人负责，WS-C+WS-D 由另一人在 M1 完成后启动，可压缩至约 12–14 周。

## 8. 测试策略

1. **单元测试**：主题映射正确性（Token→Theme 各字段）、Option 组合、密度数值换算。
2. **静态一致性检查**：CI 中扫描 `widgets/`、`extra/`，禁止裸 `ui.Hex("#...")` 与裸间距字面量（图标数据目录豁免）。
3. **视觉回归**：MyGo 无窗口渲染 PNG → 与 `golden/` 基线对比；4 种主题组合 × 每个控件。
4. **切换测试**：运行时切换主题/密度后截图，断言无遗漏旧色控件。
5. **交互测试**：用 MyGo 的测试工具（tester）验证键盘可达、点击、焦点环。
6. **CI 门禁**：`go test ./...` + 静态检查 + 截图 diff 全部通过才允许合并。

## 9. 风险登记册

| 风险 | 概率/影响 | 应对 |
| --- | --- | --- |
| MyGo API 变动 | 中/高 | 锁版本区间；换肤层做薄；每跟进一个版本跑全量回归 |
| Token 无法 1:1 映射 | 高/中 | 语义收敛的 `AntToken` 子集；细节用 Element 局部样式表达（见 04） |
| 紧凑档遗漏硬编码间距 | 中/中 | 静态检查 + 4 组合截图 + DoD 检查表 |
| 缺口控件工作量超估 | 中/中 | 按 P1→P3 分批评审，必要时推迟 C5/C6 到 v0.2 |
| 图标版权/一致性 | 低/中 | 用 antd 官方 MIT 许可 SVG，`go:generate` 可追溯来源 |
| 单人周期长 | 高/中 | 优先 M1 闭环验证价值；社区认领控件；里程碑交付可用版本 |

## 10. 立即行动（Next Actions）

1. 执行 A1：初始化 `go.mod`，固定 MyGo 版本，建 CI 与目录骨架。
2. 执行 A2–A5：完成主题层，写第一个单测验证 Token→Theme 映射与 `Apply`/`Sub`。
3. 执行 B1：实现 Button 全套变体，**并写好 Button 示例（变体/尺寸/状态/并排对比）**，接入 Gallery，跑通「换肤 + 主题切换 + 密度切换 + 示例展示」闭环。
4. 执行 F1–F5：完成 Gallery 侧栏、控制条、并排对比、代码视图与片段提取工具。
5. 闭环通过后按 WBS 顺序推进，**每个控件完成时立即配示例**，每个里程碑结束时更新本文档进度。

> 完成定义（全局）：v0.1.0 发布时，WS-A/B/C/D/E/F 全部任务满足第 4.3 节控件 DoD，示例覆盖率 100%，且全部里程碑退出标准达成。

## 11. 当前进度（2026-11-26）

### 已完成（全部含自测、提交与推送）

| 范围 | 交付物 | 验证 |
| --- | --- | --- |
| A1 骨架 | `go.mod`（mygo v0.10.0）、CI、ADR 0001/0002 | CI 全绿 |
| A2–A5 主题层 | `theme/`：AntToken、antd v5 亮暗预设、桥接（WithAccent 色阶）、Apply/Sub | 单测 + NewTester 闭环 |
| A6 helpers | `internal/metrics.go`（三档尺寸×密度）、`internal/variant.go`（按钮状态色集中映射） | 随 Button 测试 |
| A7 静态检查 | `internal/check`：禁裸颜色/裸间距（0 值豁免），CI 门禁 | 工具自验证 |
| B1 Button | `widgets/`（包 ant）：五型 × danger × 三档尺寸 × Loading/Disabled/Block/OnClick | NewTester：变体/点击/禁用/4 主题组合 |
| F4/F5 示例框架 | `samples/`：注册表、`// sample:` 标记提取（go:embed，代码即源码） | 提取与渲染测试 |
| Button 示例 | 基础/类型/尺寸/状态四组，真实可编译、可交互 | 点击反馈测试 |
| F1–F3+A8 Gallery | 侧栏 + 卡片流 + 控制条（主题/密度/品牌色/亮暗并排对比）+ 代码视图 | 交互测试 + 截图基线 |
| 截图基线 | `samples/golden/`：亮/暗 × 标准/紧凑 4 张 PNG | 二次运行比对稳定 |

### 关键实现决策（实现期新增）

- 包名 `ant`（目录 `widgets/`），导入别名按需；示例目录 `samples/<widget>`、包名 `sample`（目录 `samples/`）。
- MyGo v0.10.0 无 `OnClick`/`ui.Local` 文档化的 main 分支 API：`OnClick` 用 `ui.Local` 上升沿实现（`Button` 内置）。
- `Wrap()` 无参（与设计文档 `Wrap(true)` 写法不同）。
- 品牌色与系统偏好冲突：`Apply` 默认不跟随系统强调色（ADR 0002）。

### 下一步（按 WBS 顺序）

1. B1 余项：Typography、Icon 封装、Divider、Space/Flex/Grid 包装（各带示例）。
2. B2 文本录入：Input、TextArea、InputNumber、SearchField。
3. B4 选择控件：Checkbox、Radio、Switch、Slider、Rate。
4. 每个控件完成即配示例（DoD 第 4.3 节）。
