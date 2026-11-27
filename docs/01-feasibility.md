# 01 · 可行性分析

> **核对基线**：MyGo **v0.10.0**（`go.mod` 锁定）、Go **1.27.4**、本库 v0.1.0 已发布（v0.3 开发中），核对日期 2026-11-27。
> 本文是立项时的可行性论证；早期的 `v0.7` 版本引用、以及「约 0.7 API 可能变动」等描述已按当前基线修正。落地进度见 [06 开发计划](06-development-plan.md)。

## 结论

**可行，且属于「顺势而为」而非「逆流改造」。** MyGo 暴露的扩展点与 Ant Design 的抽象层次高度同构：

| Ant Design 的能力层 | MyGo 的对应底座 | 契合度 |
| --- | --- | --- |
| ConfigProvider / Design Token | `ui.Theme` + `Context.SetTheme` + `ui.Themed` | ★★★★★ 原生支持 |
| 无样式组件（headless） | `ButtonBase` / `SelectBase` / `TreeBase` 等 20+ Base | ★★★★★ 原生支持 |
| compactAlgorithm | `Theme.Spacing` + `t.Space(n)` / `t.Rem(n)` | ★★★★★ 原生支持 |
| 组件 = 函数式 UI | `view(c *ui.Context)` 即时模式，每帧重建 | ★★★★☆ 需转换心智 |
| 动态样式（CSS-in-JS） | 「State, not selectors」——无样式表、无选择器 | ★★★☆☆ 需翻译为 Go 分支 |
| 组件事件 / 值绑定 | `Clicked()` / `Changed()` / `*T` 绑定 / `ui.Local` | ★★★★★ 原生支持 |
| 无障碍与键盘 | `Access` 系列、Focus Ring、Base 自带键盘处理 | ★★★★★ 原生支持 |

**核心判断**：这不是给 MyGo「套一层 CSS」，而是**用 MyGo 的 Base + Theme 重新实现一套 Ant Design 外观的控件函数**。MyGo 官方在 `docs/ui/custom-widgets.md` 中明确鼓励这种做法（"Build a design of your own on them, as headless component libraries do on the web"），本方案正是这条官方推荐路径的工程化。

## 为什么可行：四个关键技术论证

### 1. 换肤不需要 fork 框架

每个 MyGo 控件都返回其 Element，构建后的链式调用即可改外观：

```go
// 官方示例：Button 返回 element，后续调用改样式
ui.Button(c, "Save").Padding(10, 20).Radius(999)
```

配合 `docs/ui/styling.md` 提供的 `Background` / `Border` / `BorderWidth` / `Radius` / `Shadow` / `InsetShadow` / `Gradient` / `Opacity` / `Cursor` 等能力，Ant Design 的按钮、卡片、输入框的视觉要素（描边、投影、圆角、悬浮态）都能表达。对于需要彻底自定义外观的控件，Base 提供「无外观版本」，把指针、键盘、焦点、无障碍都保留，只把颜色尺寸形状留给我们。

> 结论：**换肤层是纯包装函数**，不改 MyGo 源码，不 fork。

### 2. 主题可以做到「单一真源 + 运行时切换」

MyGo v0.10.0 的 `ui.Theme` 结构体包含 Ant Design 所需的全部语义槽位（共 31 个字段）：

- 颜色：`Background` / `Surface` / `SurfaceHover` / `SurfacePressed` / `Border` / `Text` / `TextMuted` / `Accent`(+Hover/Pressed/Text) / `InactiveAccent` / `Danger` / `Warning` / `Success` / `Selection` / `Focus` / `Inverse` / `InverseText` / `Scrollbar` / `ScrollbarStyle`
- 度量：`Radius` / `Spacing` / `FontSize` / `Font` / `FocusRingWidth` / `FocusRingOffset` / `ScrollbarWidth` / `TooltipDelay`

切换机制现成：

- 全局：`c.SetTheme(&t)`，改变整个窗口，可每帧按用户选择调用；
- 局部：`ui.Themed(c, t)`，只给某个子树（如「亮/暗主题并排预览」）换主题，且弹出层、菜单、Tooltip 自动跟随；
- 品牌色：`Theme.WithAccent(color)` 自动生成 hover / pressed / selection / focus / AccentText 的成套色阶；
- 跟随系统：`Theme.WithPreferences(c.Preferences())` 存在，但**本库默认不使用**（品牌色优先，见 [ADR 0002](decisions/0002-apply-does-not-follow-system-accent-by-default.md)）。

> 结论：需求 5 的「共用 + 运行时切换」**基本无需自研机制**，只需在 MyGo Theme 之上定义 Ant Design 语义的包装层。

### 3. 紧凑密度有现成的度量单位

`Theme.Spacing` 是 MyGo 官方定义的「控件内外留白的基本单位」，官方文档原话：

> `Spacing` is the unit of the room widgets leave… It is 4 by default; 3 makes every widget compact, 5 roomy.

控件的 padding、gap、复选框/开关/滑块尺寸、表格与树行高，全部是它的倍数；自定义元素用 `t.Space(3)`、`t.Rem(2)` 即可同步缩放。这与 Ant Design 的 `compactAlgorithm` + `sizeUnit`/`controlHeight` 思路一致。

> 结论：需求 6 的第二套紧凑风格，**只需定义一组 Spacing/字号/行高数值**，不必逐控件改 padding。

### 4. 缺口控件可以用 Base 组合实现

Ant Design 有而 MyGo 无的控件（分页、级联选择、穿梭框、上传、抽屉、骨架屏、二维码等），绝大多数可由 MyGo 的 Base + Element + Overlay 组合而成。MyGo 还提供 `ToastViewportBase` / `DialogBase` / `MenuBase` / `TooltipBase` 等专门支撑「自建浮层类组件」的底座。少数依赖外部能力的（如二维码编码、上传）需要引入纯 Go 库，但都在可行范围内。详见 [03 缺口控件实现](03-gap-widgets.md)。

## 推荐架构（分层）

```
┌──────────────────────────────────────────────────────────┐
│  应用层   用户代码：ui.View(app.view) + 主题选择器          │
├──────────────────────────────────────────────────────────┤
│  控件层   widgets/ 换肤控件   ·   extra/ 缺口控件           │
│           （函数签名对齐 MyGo，外观对齐 Ant Design）        │
├──────────────────────────────────────────────────────────┤
│  主题层   theme/ AntToken → ui.Theme 映射 · 预设 · 密度档   │
├──────────────────────────────────────────────────────────┤
│  底座层   MyGo package ui：Theme / Base / Element / 颜色    │
└──────────────────────────────────────────────────────────┘
```

关键约束：**控件层只能依赖主题层和 MyGo 底座**，禁止控件之间互相硬编码颜色；所有视觉常量必须来自主题层的 Token。

## 当前实现状态（2026-11-27）

| 维度 | 现状 |
| --- | --- |
| 主题层 | ✅ `AntToken` + `Light()` / `Dark()` + `Apply` / `Sub`，单测覆盖映射与密度 |
| 换肤控件 | ✅ `widgets/`（包 `ant`）78 个构造函数，覆盖 35 个原生原语 |
| 缺口控件 | ✅ `extra/` 7 个（Tree / Upload / Transfer / Carousel / QRCode / Cascader / TreeSelect）+ `widgets/` 内自建（Pagination / Steps / Alert / Card / Tag / Badge / Empty / Result / Skeleton / Statistic / Descriptions / Timeline / Spin 等） |
| 图标 | ✅ 44 个 antd outlined 常量 + `icon.Icon` |
| 示例 | ◐ 11 个示例分组 / 58 个示例片段；未逐控件覆盖 |
| Gallery | ✅ 侧栏 + 卡片流 + 主题/密度/品牌色切换 + 亮暗并排对比；搜索/分类/全局尺寸为待办 |
| 视觉回归 | ◐ 4 张 Gallery 级截图基线（`samples/golden/`），尚无逐示例基线 |
| 静态检查 | ✅ `go run ./internal/check` 禁裸颜色/裸间距 |
| 版本 | v0.1.0 已打 tag；v0.2 已推送待打 tag；v0.3 开发中 |

> 详细任务状态见 [06 开发计划 §5](06-development-plan.md#5-wbs-详细任务分解与状态)。

## 工作量估算（立项估算，未随实现重估）

| 阶段 | 内容 | 人日（1 名熟悉 Go 的开发者） |
| --- | --- | --- |
| P0 主题层 | Token 映射、亮/暗/紧凑预设、桥接函数 | 5–8 |
| P1 基础控件 | 按钮/输入/选择/复选/单选/开关/标签等 ~20 个 | 10–15 |
| P2 数据展示 | 表格/树/列表/卡片/徽标/进度等 ~20 个 | 15–20 |
| P3 反馈与浮层 | 对话框/抽屉/提示/气泡/通知等 ~15 个 | 12–18 |
| P4 缺口控件 | 分页/级联/穿梭/上传/二维码等 ~12 个 | 15–25 |
| P5 图标与 Gallery | 图标数据、预览应用、视觉回归 | 6–10 |
| **合计** | | **约 63–96 人日** |

> 说明：MyGo 现版本 v0.10.0、生态仍早期，API 可能有变动；建议按周迭代并与 MyGo 版本对齐。实际交付按里程碑分批，见 06 文档。

## 风险与缓解

| 风险 | 说明 | 缓解措施 |
| --- | --- | --- |
| **API 稳定性** | MyGo 仍在 0.x，官方声明 API 可能变化 | 本库锁定 `v0.10.0`；换肤层做薄，隔离变化；升级时跑全量回归 |
| **Token 非 1:1** | Ant Design 有数百个 token，MyGo Theme 约 31 个字段 | 定义语义收敛的 `AntToken` 子集，多余的忽略或局部用 Element 样式表达 |
| **无 CSS 选择器** | 无法像 Web 那样用类名批量改样式 | 统一封装「状态→样式」函数，集中管理 hover/active/disabled 分支 |
| **图标体系不同** | antd 图标是 React 组件，MyGo 用 SVG 路径 | 建立 Go 常量图标表，从官方 SVG 资源生成，MIT 兼容 |
| **字体栈差异** | antd 用系统字体栈，MyGo `Theme.Font` 为单一族名 | 默认留空走系统字体，文档说明跨平台差异 |
| **交互语义差异** | Web 组件（Upload/Affix）在桌面场景语义不同 | 重新定义为桌面语义（本地文件选择、元素级 Sticky），见 03 文档 |
| **工作量较大** | 80+ 控件不是短期可全部完成 | 按 P0→P5 里程碑推进，优先高频控件，社区可认领 |
| **文档与实现漂移** | 设计文档容易写成「已实现」 | 每篇文档标注核对日期与状态列，任务收尾同步更新 |

## 里程碑建议（含状态）

- **M1**：主题层 + 基础控件 20 个，Gallery 可预览 → 验证「换肤 + 切换主题」闭环。✅ 已完成
- **M2**：数据展示与反馈控件，覆盖常见后台/工具类应用 80% 场景。◐ 主体完成，Form/Collapse/Tooltip/Sidebar/Drawer 未做
- **M3**：紧凑密度全控件适配，出上位机 / IDE 风格示例。◐ 密度可切换、4 组合基线已有；逐控件审计与 IDE 示例待做
- **M4**：缺口控件补齐，发布 v0.1.0。◐ v0.1.0 已发布；P2/P3 缺口控件仍在推进

## 建议

**立项并继续推进。** 主题与换肤闭环已被 v0.1.0 验证；剩余工作的风险集中在「缺口控件的组合复杂度」与「示例/截图覆盖」，按 03 文档的优先级分批推进即可。
