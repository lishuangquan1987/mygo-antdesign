# 06 · 详细开发计划

> 基于 [01 可行性分析](01-feasibility.md)、[02 控件映射](02-widget-mapping.md)、[03 缺口控件](03-gap-widgets.md)、[04 主题系统](04-theme-system.md)、[05 紧凑密度](05-compact-density.md) 的落地计划。
>
> 本文的「状态列」与代码同步，核对日期 **2026-11-27**；基线 MyGo **v0.10.0**、Go **1.27.4**。

## 1. 计划总览

### 1.1 目标与当前版本

交付 `github.com/lishuangquan1987/mygo-antdesign`：

1. 一套 Ant Design 风格的 MyGo 原生 UI 控件库（换肤 + 缺口补齐）。
2. 全部控件共用主题，支持运行时切换亮/暗、品牌色、密度。
3. 提供标准与紧凑两套密度，覆盖上位机 / IDE 场景。
4. 每个控件配可运行示例（当前按分类成组，逐控件示例是后续目标）。
5. 具备 Gallery 展示应用、截图回归、静态一致性检查。

**版本状态**：

| 版本 | 状态 | 内容 |
| --- | --- | --- |
| v0.1.0 | ✅ 已打 tag | 主题层 + 9 组换肤控件 + 首批缺口控件（Tree / Upload / Transfer / Carousel / QRCode）+ 44 图标 + Gallery + 文档 |
| v0.2 | ✅ 已推送、待打 tag | README 快速开始重写、Release 流水线、Cascader / TreeSelect |
| v0.3 | 🚧 开发中 | P1 缺口控件：Empty / Result / Skeleton / Statistic / Descriptions / Timeline / Spin |

### 1.2 范围边界

| 纳入 | 不纳入 |
| --- | --- |
| MyGo 原生控件的 antd 换肤 | Fork / 修改 MyGo 源码 |
| antd 有而 mygo 无的控件补齐 | Web 框架层概念（App/ConfigProvider/Util 同名控件） |
| 主题 Token 与运行时切换 | 独立设计一套全新视觉语言 |
| 图标数据、Gallery、文档 | 商业组件 / Pro 组件 |
| 紧凑密度与 IDE 示例（待做） | 移动端适配（MyGo 已原生支持，后续版本） |

### 1.3 量化目标与现状

| 指标 | 目标 | 现状（2026-11-27） |
| --- | --- | --- |
| 换肤 + 缺口控件构造函数 | ≥ 55 | **106**（`widgets/` 99 + `extra/` 7，另有 `icon.Icon`） |
| 缺口控件 | ≥ 12 | `extra/` **7** + `widgets/` 内 **27**（Pagination / Steps / Alert / Card / Tag / Badge / Empty / Result / Skeleton / Statistic / Descriptions / Timeline / List / Spin / Confirm / Tooltip / Popconfirm / Collapse / Form / Drawer / Image / ColorPicker / ColorPanel / 环形 Progress / Toggle / Affix 等） |
| 主题预设 | 4 套（亮/暗 × 标准/紧凑） | ✅ 通过 `Light/Dark × Comfortable/Compact` 组合得到 |
| 示例覆盖率 | 每个控件 ≥ 1 个示例页 | **13 个示例分组 / 84 个示例片段**，按 docs/02 §九 全覆盖矩阵继续补齐 |
| 视觉回归 | 每个示例 × 3 种主题组合 | **3 张 Gallery 级基线**（`samples/golden/gallery-*.png`，按平台生效），尚未逐示例 |

> 现状与目标之间的差距是明确的待办项，见 §10，不再宣称「100% 覆盖」。

## 2. 技术基座与前置条件

| 项 | 要求 | 备注 |
| --- | --- | --- |
| Go | **1.27.4**（`go.mod` 与 CI 一致） | MyGo native UI 的硬性要求 |
| MyGo | **v0.10.0**（`go.mod` 锁定） | 升级需跑全量回归；早期文档写的 `v0.7.x` 已过时 |
| 依赖 | 仅纯 Go 库、无 cgo | 当前额外依赖 `skip2/go-qrcode`（二维码） |
| 工具 | `go run ./gallery`、`go build ./...`、`scripts/visualcheck`、`scripts/package` | Gallery 用 `ui.NewTester` 无窗口渲染 PNG；发布打包为纯 Go 实现 |
| CI | GitHub Actions | `gofmt` + `go vet` + `go test` + `internal/check` |

**已完成的前置动作**：`go.mod`、`docs/decisions/`（ADR 0001/0002）、`.github/workflows/`（CI + Release）。

## 3. 仓库与包结构（实际）

```
mygo-antdesign/
├── go.mod / go.sum                 # Go 1.27.4，mygo v0.10.0，skip2/go-qrcode
├── theme/                          # 主题层：唯一视觉真源
│   ├── token.go                    # AntToken / Density / Option
│   ├── preset.go                   # Light() / Dark() 与 antd 色板
│   ├── bridge.go                   # AntToken → ui.Theme
│   ├── theme.go                    # Apply / Sub
│   └── theme_test.go
├── internal/
│   ├── metrics.go                  # SizeKind / ControlHeight / ControlPaddingX / ControlFont
│   ├── variant.go                  # 按钮状态色集中映射
│   └── check/main.go               # 样式铁律静态检查（可执行）
├── widgets/                        # 包名 ant：换肤控件
│   ├── button.go        typography.go   divider.go    layout.go
│   ├── input.go         select.go       picker.go     table.go
│   ├── display.go       feedback.go     nav.go        data.go
│   └── *_test.go                   # 每个源文件对应一套测试
├── extra/                          # 包名 extra：缺口控件
│   ├── tree.go          treeselect.go   cascader.go   trigger.go
│   ├── upload.go        transfer.go     carousel.go    qrcode.go
│   └── extra_test.go
├── icon/
│   ├── icon.go                     # Icon 渲染
│   ├── icons_gen.go                # 44 个 outlined SVG 常量（生成）
│   └── icon_test.go
├── scripts/genicons/main.go        # 图标生成器（go:generate 入口在 icon/）
├── samples/                        # 示例源码：真实可编译
│   ├── registry.go  extract.go     # 注册表 + // sample: 标记提取
│   ├── button/ typography/ divider/ layout/ input/ selectdemo/
│   ├── display/ table/ feedback/ nav/ extrademo/   # 11 个示例分组
│   └── golden/                     # gallery-*.png 截图基线（3 张）+ manifest.json（生成平台）
├── gallery/
│   ├── main.go                     # 入口 + 显式注册示例
│   ├── app.go                      # 侧栏 + 卡片流 + 控制条
│   └── app_test.go                 # 交互测试 + 截图基线比对
├── docs/                           # 本文档集 + decisions/ + images/
└── .github/workflows/{ci.yml,release.yml}
```

**依赖规则**：`widgets/` 和 `extra/` 只能依赖 `theme/`、`internal/` 与 MyGo `ui`，禁止跨控件硬编码颜色/间距；`gallery/` 是演示应用，允许少量展示用 `ui.Hex`（不在静态检查范围）。

## 4. 统一开发规范（所有任务必须遵守）

### 4.1 API 约定

函数式选项 + 贴近 MyGo 的签名：

```go
// 已实现的 API 形态。
ant.Button(c, "保存", ant.Primary(), ant.Size(ant.Small), ant.Loading(app.saving))
ant.Table(c, columns, rows, ant.TableSelection(&sel), ant.TableOnSort(app.sort))
ant.Pagination(c, &app.page, app.total, ant.OnPageChange(app.load))
ant.Steps(c, &app.step, []string{"填写", "确认", "完成"})
extra.Cascader(c, &app.path, nodes, extra.CascaderOnChange(app.apply))
theme.Apply(c, theme.Light(theme.WithPrimary(ui.Hex("#1677ff")), theme.WithDensity(theme.Compact)))
```

- 控件返回 `ui.Element`，与 MyGo 组合方式一致。
- 有状态项用 `*T` 绑定 + `ui.Local`，与 MyGo 原生一致。
- 选项类型统一为 `type Option func(*config)`；控件级尺寸用 `ant.Size(ant.Large)` 强类型枚举。
- `OnClick` 等 React 风格回调用 `ui.Local` 上升沿实现（MyGo v0.10.0 的 Element 没有 `OnClick`）。

### 4.2 样式铁律（密度/主题可切换的前提）

1. 颜色只能来自 `c.Theme()` 字段或 `internal/variant.go` 的映射函数。
2. 间距只能写 `t.Space(n)` / `t.Rem(n)`，禁止裸像素（0 与 1px 描边豁免）。
3. 1px 描边、焦点环宽度等视觉基准保持固定。
4. 状态样式（hover/pressed/disabled/selected/focus）集中在 `variant.go`，禁止各控件各写一套。

### 4.3 每个控件的完成检查表（DoD）

- [ ] 签名符合 4.1 API 约定
- [ ] 颜色/间距符合 4.2 铁律（`go run ./internal/check` 通过）
- [ ] 标准与紧凑两档截图无异常
- [ ] 亮/暗主题截图无异常
- [ ] 键盘可达、有 Role/Label（Base 实现天然满足，纯组合需补）
- [ ] **示例已编写并注册**（`samples/<分组>/demo.go`，覆盖基础 + 主要变体 + 禁用 + 三档尺寸）
- [ ] Gallery 侧栏可见该控件
- [ ] 截图基线通过（当前是 Gallery 级 3 张；逐控件基线为后续目标）
- [ ] 文档/注释齐全

## 5. WBS 详细任务分解与状态

> 估算单位：人日（1 名熟悉 Go 的开发者），是立项时的原始估算，未随实现重估。
> 状态：✅ 完成 / ◐ 部分完成 / ◻ 未开始 / — 不适用。

### WS-A · 基础设施与主题层

| ID | 任务 | 估算 | 状态 | 交付物 |
| --- | --- | --- | --- | --- |
| A1 | `go.mod`、目录骨架、CI、ADR 模板 | 1 | ✅ | `go.mod`、`ci.yml`、`docs/decisions/` |
| A2 | `theme/token.go`：`AntToken`、`Density`、`Option` | 2 | ✅ | `token.go` |
| A3 | `theme/bridge.go`：Token → `ui.Theme`，含 `WithAccent` 色阶 | 2 | ✅ | `bridge.go` |
| A4 | `theme/preset.go`：亮/暗/紧凑预设、品牌色 | 2 | ✅ | `preset.go` |
| A5 | `theme/theme.go`：`Apply` / `Sub` | 1 | ✅ | `theme.go` |
| A6 | `internal/variant.go` + `metrics.go` | 2 | ✅ | 状态色与三档尺寸 |
| A7 | 静态检查 | 1 | ✅ | `internal/check/main.go` |
| A8 | `gallery/` 应用骨架 + 无窗口 PNG 入口 | 2 | ✅ | `gallery/app.go` |
| | **小计** | **13** | | |

### WS-B · 换肤控件（widgets/）

| ID | 任务包 | 覆盖控件 | 估算 | 状态 | 备注 |
| --- | --- | --- | --- | --- | --- |
| B1 | 按钮与排版 | Button、Typography、Divider、Space/Flex/Grid | 4 | ✅ | 五型 × danger × 三档尺寸 |
| B2 | 文本录入 | Input、TextArea、SearchInput、NumberInputV、DecimalInput、SearchField、TokenField、EditableText | 4 | ✅ | Mentions 归 P2（C5） |
| B3 | 选择类 | Select、Combobox、AutoComplete | 5 | ✅ | Cascader/TreeSelect 在 extra |
| B4 | 选择控件 | Checkbox、CheckboxAll、RadioGroup、Switch、Slider、RangeSlider、Rate、Segmented | 4 | ✅ | Toggle/ToggleGroup 已包装（v0.3） |
| B5 | 日期时间 | DatePicker、TimePicker、Calendar | 4 | ✅ | 基于 `ui.DateInput` / `ui.TimeInput` / `ui.Calendar` |
| B6 | 表单 | Form（标签/必填/校验） | 2 | ✅ | `ant.Form` / `FormItem` / `FormGroup`（v0.3） |
| B7 | 数据展示-基础 | Table、Tree、Card、List、Descriptions、Statistic | 6 | ✅ | List 已包装（v0.3） |
| B8 | 数据展示-辅助 | Avatar、Badge、Tag、Progress、Timeline、Empty、Result、Skeleton | 5 | ✅ | 全部完成（Timeline/Empty/Result/Skeleton 在 v0.3） |
| B9 | 导航 | Tabs、Segmented、Steps、Breadcrumb、Dropdown、Sidebar | 5 | ✅ | Sidebar 已包装；Menu 复用 Dropdown |
| B10 | 折叠/浮层 | Collapse/Accordion、Popover、Tooltip | 3 | ✅ | v0.3 全部完成 |
| B11 | 反馈 | Modal、Confirm、Message/Notify、Spin、Watermark | 5 | ✅ | v0.3 全部完成 |
| B12 | 其他 | Splitter、Affix、Drawer | 2 | ✅ | 全部完成（v0.3） |
| | **小计** | | **44** | | |

### WS-C · 缺口控件（extra/ 与 widgets/ 内的自建控件）

| ID | 任务包 | 覆盖控件 | 估算 | 状态 | 备注 |
| --- | --- | --- | --- | --- | --- |
| C1 | P1 简单组合 | Card、Tag、Empty、Result、Divider、Statistic、Timeline | 3 | ✅ | v0.3 补齐 Empty/Result/Statistic/Timeline |
| C2 | P1 数据展示 | Descriptions、Skeleton | 2 | ✅ | v0.3 完成 |
| C3 | P2 分页与穿梭 | Pagination、Transfer | 3 | ✅ | Pagination 在 widgets，Transfer 在 extra |
| C4 | P2 浮层 | Drawer、Popconfirm、Watermark、Tour | 6 | ✅ | v0.3 全部完成 |
| C5 | P2 复杂选择 | Cascader、TreeSelect、Mentions | 8 | ✅ | 全部完成 |
| C6 | P3 依赖能力 | QRCode、Upload、Spin、FloatButton、BorderBeam、Notification、Anchor、Carousel、Masonry | 12 | ✅ | v0.3 全部完成 |
| | **小计** | | **34** | | 实际分批交付，见 §11 |

### WS-D · 图标与资源

| ID | 任务 | 估算 | 状态 | 交付物 |
| --- | --- | --- | --- | --- |
| D1 | 从 antd SVG 生成 Go 路径常量 | 2 | ✅ | `scripts/genicons` + `icon/icons_gen.go` |
| D2 | 图标集与渲染（尺寸随 `t.Rem`） | 3 | ✅ | 44 个 outlined 图标 + `icon.Icon` |
| | **小计** | **5** | | |

### WS-E · 文档、示例与发布

| ID | 任务 | 估算 | 状态 | 备注 |
| --- | --- | --- | --- | --- |
| E1 | README、迁移指南 | 2 | ✅ | README + docs/08 |
| E2 | 主题与密度使用文档、API 参考 | 3 | ✅ | docs/04、05 已对齐实现 |
| E3 | IDE / 上位机风格示例应用 | 4 | ✅ | `samples/ide` 紧凑型 IDE 布局（侧栏/标签页/代码区/状态栏） |
| E4 | v0.1.0 发布（CHANGELOG、tag、release） | 2 | ✅ | v0.1.0/v0.2.0/v0.3.0 tag 已推送，Release 工作流自动构建安装包 |
| | **小计** | **11** | | |

### WS-F · 控件示例与 Gallery

| ID | 任务 | 估算 | 状态 | 备注 |
| --- | --- | --- | --- | --- |
| F1 | Gallery 骨架：侧栏 + 示例区 + 顶部控制条 | 4 | ✅ | 侧栏分类分组 + 搜索 + 每组数量 |
| F2 | 并排对比模式（亮/暗） | 2 | ✅ | `theme.Sub` + 左右两栏 |
| F3 | 代码视图（真实源码 + 只读展示） | 2 | ✅ | `ui/editor` 只读代码视图（语法高亮 + 行号）+ 复制代码按钮 |
| F4 | 示例注册表 + 分类 | 1 | ✅ | `samples/registry.go` |
| F5 | `go:embed` 片段提取 | 2 | ✅ | `samples/extract.go` |
| F6 | 截图生成与回归 | 3 | ✅ | 3 张 Gallery 级 + 324 张逐示例基线（4 组合）；README 嵌图手工维护 |
| | **小计** | **14** | | |

**总计估算：约 121 人日**（立项估算；示例编写约 21 人日摊入各控件）。单人约 6 个月，2 人并行约 3.5 个月。

## 6. 里程碑与验收标准

| 里程碑 | 计划周次 | 范围 | 状态 | 实际/差距 |
| --- | --- | --- | --- | --- |
| **M0 基建** | W1 | A1–A8 | ✅ | 仓库、CI、主题层、Gallery 可渲染 |
| **M1 主题闭环** | W2–W5 | A + B1–B6 + F1–F5 | ✅ | Button/录入/选择/日期与 11 组示例卡片；主题/密度/品牌色可切换 |
| **M2 常用控件** | W6–W10 | B7–B12 | ✅ | v0.3 补齐 Form/Collapse/Tooltip/Sidebar/Drawer/List 等 |
| **M3 密度适配** | W11–W12 | 全控件紧凑审计 + F2/F6 + E3 | ✅ | 3 组合基线 + IDE 风格示例 + 密度并排对比 |
| **M4 缺口与发布** | W13–W18 | C 全部 + D + E4 | ✅ | 缺口控件、Gallery/IDE 示例、CHANGELOG 与三个 tag 全部完成 |

## 7. 时间线（历史基线）

```
W1    M0 基建 + 主题层 + Gallery 空壳
W2-5  M1 基础控件（按钮/录入/选择/表单）+ 示例 + Gallery 闭环
W6-10 M2 数据展示/导航/浮层/反馈 + 同步配示例
W11-12 M3 紧凑密度审计 + 并排对比 + 全量示例截图 + IDE 示例
W13-18 M4 缺口控件 + 缺口示例 + 图标 + 文档嵌图 + 发布
```

> 这是立项时的单人 18 周基线；实际以 WBS 状态为准。

**并行建议**：若 2 人，A+WS-B 由一人负责，WS-C+WS-D 由另一人在 M1 完成后启动，可压缩至约 12–14 周。

## 8. 测试策略

1. **单元测试**：每个控件文件对应 `_test.go`，覆盖变体渲染、点击/键入、禁用、主题与密度组合。
2. **静态一致性检查**：`go run ./internal/check` 扫描 `widgets/`、`extra/`，禁止裸 `ui.Hex/RGB/RGBA/Oklch` 与裸间距字面量（0 与 1px 豁免）。
3. **视觉回归**：`gallery` 包用 `ui.NewTester` 渲染 **3 种组合**（亮标准、暗标准、亮紧凑），与 `samples/golden/gallery-*.png` 逐字节比对；基线带平台清单（`manifest.json`），跨平台只渲染不比对（字体栈不同），详见 [07 §8](07-samples-gallery.md#8-截图生成与文档联动实际)。
4. **交互测试**：`TestGalleryInteractions` / `TestGallerySampleInteraction` 验证切换与点击反馈。
5. **CI 门禁**：`gofmt` + `go vet` + `go test ./...` + 静态检查全部通过才允许合并。

> 尚未实现：逐控件/逐示例截图基线、文档嵌图自动化、宏基准（帧耗时）。

## 9. 风险登记册

| 风险 | 概率/影响 | 应对 |
| --- | --- | --- |
| MyGo API 变动 | 中/高 | 锁 `v0.10.0`；换肤层做薄；升级时跑全量回归 |
| Token 无法 1:1 映射 | 高/中 | 语义收敛的 `AntToken` 子集；细节用 Element 局部样式表达（见 04） |
| 紧凑档遗漏硬编码间距 | 中/中 | 静态检查 + 3 组合截图 + DoD 检查表 |
| 缺口控件工作量超估 | 中/中 | 按 P1→P3 分批评审，必要时推迟 P3 到 v0.3 |
| 图标版权/一致性 | 低/中 | 用 antd 官方 MIT 许可 SVG，`scripts/genicons` 可追溯来源 |
| 单人周期长 | 高/中 | 优先高频控件；里程碑交付可用版本；社区认领 |
| 文档与实现漂移 | 高/中 | 每个任务收尾同步状态列与路径；核对日期写在文首 |

## 10. 后续（v0.x 范围外）

WBS A–F 全部完成（见 §11），v0.3.0 已发布。以下为后续可选增强，不阻塞 v0.x：

1. **系统亮暗自动跟随**：MyGo `ui` 暂未暴露系统外观 API（ADR 0002 记录为上游限制，本库无法完成）。
2. ~~Gallery 全局尺寸切换~~：**已完成（第十二轮）**——顶部「尺寸 小/中/大」通过 `theme.WithSpacing/WithRadius/WithFontSize` 统一缩放。
3. ~~`ui/editor` 代码视图~~：**已完成（第十二轮）**——Gallery 代码区改为 `ui/editor` 只读视图（语法高亮 + 行号 + 复制按钮）。
4. ~~文档嵌图流水线~~：**已完成（第十二轮）**——`TestWriteDocsScreenshots` 生成 README 截图，`scripts/visualcheck` 按分组批量出图，均基于无窗口 `ui.NewTester`。
5. **Pro 组件与更多 antd 高级组件**：超出 v0.x 范围。

## 11. 当前进度（2026-11-27）

### 11.1 v0.1.0（已打 tag）

| 范围 | 交付物 | 验证 |
| --- | --- | --- |
| A1 骨架 | `go.mod`（mygo v0.10.0）、CI、ADR 0001/0002 | CI 全绿 |
| A2–A5 主题层 | `theme/`：AntToken、antd v5 亮暗预设、桥接（WithAccent 色阶）、Apply/Sub | 单测 + NewTester 闭环 |
| A6 helpers | `internal/metrics.go`（三档尺寸 × 密度）、`internal/variant.go`（按钮状态色） | 随 Button 测试 |
| A7 静态检查 | `internal/check`：禁裸颜色/裸间距（0 值豁免） | 工具自验证 |
| B1–B5、B7–B9、B11 主体 | Button/Typography/Divider/Layout/Input/Select/Checkbox/Radio/Switch/Slider/Rate/Segmented/Select/Combobox/DatePicker/TimePicker/Table/Tag/Badge/Avatar/Card/Alert/Modal/Message/Notify/Progress/Confirm/Breadcrumb/Tabs/Dropdown/Pagination/Steps（部分控件仍有未做项，见 §5） | 各控件单测 + 4 主题组合截图 |
| C 缺口首批 | `extra/`：Tree/Upload/Transfer/Carousel/QRCode | 交互测试 + 视觉核查 |
| D 图标 | `icon/`：44 个 antd outlined 常量 + 生成器 | 4 主题组合 + 颜色核查 |
| E 文档 | README、docs/08 迁移指南 | 人工核对 |
| F Gallery | 侧栏 + 卡片流 + 控制条（主题/密度/品牌色/亮暗对比）+ 代码视图 | 交互测试 + 截图基线 |
| 基线 | `samples/golden/`：亮标准、暗标准、亮紧凑，共 3 张 PNG | 二次运行比对稳定 |

### 11.2 v0.2（已推送、待打 tag）

- README 重写「快速开始」：CI/Release 徽章、Gallery 首图、可编译最小示例、`go run ./gallery`、Release 安装包入口。
- CI：Release 工作流（5 平台交叉编译安装包，release 后自动附加，支持 `workflow_dispatch`），已真实运行验证。
- `extra` 新增 **Cascader** 与 **TreeSelect**（基于 `ui.Popover` 组合），含示例与测试。
- 截图基线改为按平台生效（`samples/golden/manifest.json`），修复 Linux CI 上跨平台字体渲染导致的永久失败。

### 11.3 v0.3（已完成）：缺口控件全部 + 原生控件补全包装

| 控件 | 位置 | 示例 | 验证 |
| --- | --- | --- | --- |
| Empty / Result / Skeleton / Statistic / Descriptions / Timeline / Collapse | `widgets/data.go`、`widgets/collapse.go` | `samples/display`（新增 4 个示例） | 单测 + 亮/暗/紧凑渲染 + 人工图像核对 |
| Spin / Tooltip / Popconfirm | `widgets/feedback.go`、`widgets/overlay.go` | `samples/feedback`（新增 2 个示例） | 单测（含气泡确认/取消交互）+ 主题组合 |
| Form / FormItem / FormGroup | `widgets/form.go` | `samples/input`（新增 1 个示例） | 单测 + 主题组合 |
| List / Sidebar / Splitter / AutoComplete | `widgets/list.go`、`nav.go`、`layout.go`、`select.go` | display / nav / layout / input 各新增 1 个示例 | 单测（选中/键盘/拖拽/补全）+ 主题组合 |
| Drawer / Image / ColorPicker / 环形 Progress / Toggle / Affix | `widgets/drawer.go`、`media.go`、`feedback.go`、`select.go`、`layout.go` | feedback / display / input / selectdemo 新增示例 | 单测（抽屉开关、预览、进度环）+ 主题组合 |
| 原生控件补全包装 | `widgets/native.go`（Hyperlink/DecimalInput/SearchField/Stepper/Calendar/TokenField/EditableText/Meter/Collapsible/Accordion）、`overlay.go`（Popover/ContextMenu） | input / display / feedback 各新增 1 个示例 | 单测（渲染/气泡/右键菜单）+ 主题组合 |
| P3 缺口收尾 | `widgets/gaps2.go`：Watermark/Mentions/Tour/Notification/FloatButton/Anchor/BorderBeam/Masonry | feedback/display/layout/input 各新增 1 个示例 | 单测（水印/提及插入/引导步骤/通知/锚点/流光）+ 人工图像核对 |
| Gallery 增强 | `gallery/app.go` | 侧栏分类 + 搜索 + 每组数量、代码复制、亮暗/密度对比 | gallery 单测（搜索/密度对比/复制）+ 3 张截图基线 |
| IDE 风格示例 | `samples/ide` | `gallery/main.go` 注册 | gallery 集成测试 + `scripts/visualcheck` 截图核对 |
| 主题跟随系统 | `theme/theme.go`（`WithFollowSystem`） | — | 单测（高对比/文字缩放、默认不跟随） |
| 工具 | `scripts/visualcheck` | 把示例按主题渲染成 PNG，人工核对控件效果 | 无窗口 `ui.NewTester` |
| 质量门禁 | `samples/all_test.go` | 校验全部示例分组的代码片段提取与亮/暗渲染 | 新增测试 |

### 11.4 关键实现决策（实现期新增）

- 包名 `ant`（目录 `widgets/`）；示例目录 `samples/<分组>`，包名 `sample`。
- MyGo v0.10.0 无 `OnClick`：`ant.OnClick` 用 `ui.Local` 上升沿实现，一次点击只触发一次。
- `Wrap()` 无参（与早期设计文档 `Wrap(true)` 不同）。
- 品牌色与系统偏好冲突：`Apply` 默认不跟随系统强调色（ADR 0002）。
- `Changed()` / `Focused()` / `Submitted()` 必须查询内层原语元素（容器不可聚焦，Button 除外）。
- Breadcrumbs 语义：`chosen=-1` 未选择，最后一项是当前页不可点；菜单选择用 `ChooseMenuItem`。
- `ui.AlertDialog` 返回被点按钮下标（0 起）；`ant.Confirm` 只回调 `onOk`，不返回值。
- Tag 色系全部由主题语义色 `Mix` 派生（暗色自动成立）。
- **排坑**：触发器文本 `Grow(1).SingleLine()` 在 Row 布局下高度会塌陷（紧凑密度时为 0，命中区消失）——去掉该修饰即可；`ui.Popover` 的 anchor 必须传触发器元素本身，传空宿主容器会导致浮层定位错误。

## 12. 已知限制

- `ant.Table` 单元格为字符串比较排序（数字列请先格式化定宽文本）。
- `extra.Upload` 的文件对话框阻塞调用，Tester 无窗口环境不可点击触发。
- `ant.Message` / `Notify` 复用 MyGo 底部 toast viewport，无法做到 antd 的顶部居中 / 右上角独立浮层。
- 含加载动画或网络依赖的示例（Button 状态、反馈 加载中、RemoteImage 网络图片）不纳入逐字节比对；其余示例覆盖亮/暗 × 标准/紧凑 4 组合。
- `ant.Skeleton` 为静态骨架，暂未做呼吸动画；`ant.Descriptions` 的 `bordered` 是「外框 + 内边距」简化版，无内部网格线。
- `ant.Collapse` 用自建头部（无展开/收起动画）。
- `ant.ImageEl` 预览为模态大图，无缩放/旋转工具条。
- `ant.Drawer` / `ant.Modal` 已实现打开动效；关闭仍为即时（MyGo `DialogBase` 关闭即卸载面板）。

## 13. 变更记录

| 日期 | 变更 |
| --- | --- |
| 2026-11-26 | A1–A8、F1–F5、B1 完成；v0.1.0 进度表初版 |
| 2026-11-27 | v0.1.0 发布（9 组换肤 + 首批缺口 + 图标）；v0.2 Cascader/TreeSelect；本文档按「实现状态」重写，修正 MyGo 版本与 API 名、区分「设计/已实现」 |
| 2026-11-27 | v0.3 进行中：P1 缺口控件 Empty/Result/Skeleton/Statistic/Descriptions/Timeline/Spin 落地；P2 的 Tooltip/Popconfirm/Collapse/Form 落地（含示例、测试、visualcheck 工具）；截图基线改为按平台生效 |
