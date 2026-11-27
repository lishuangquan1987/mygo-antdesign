# 05 · 紧凑密度风格

> 对应需求 6：在 Ant Design 标准间距之外，再加一套紧凑风格（padding/margin 调小），用于有限空间显示更多内容——上位机、写代码 IDE、VS2026 等。
>
> 核对基线：MyGo v0.10.0、本库当前实现（含 v0.2 未发版改动）。

## 目标

在**完全复用同一套控件实现**的前提下，提供两档密度：

- **标准（Comfortable）**：Ant Design 默认观感，适合后台管理、设置面板、普通工具。
- **紧凑（Compact）**：padding、margin、行高、控件尺寸整体缩小，适合上位机监控面板、IDE、终端类工具、数据密集表格。

**关键**：密度不是第二套控件库，而是主题层的一个维度。切换密度 = 改 `Theme.Spacing` 等度量 + `SetTheme`，全部控件自动跟随。

## 为什么能「一处修改，全控件生效」

MyGo 官方对 `Spacing` 的定义（`docs/ui/styling.md`）：

> `Spacing` is the unit of the room widgets leave: their paddings and gaps, and the sizes of check boxes, switches, sliders and the rows of tables and trees, are multiples of it. It is 4 by default; **3 makes every widget compact**, 5 roomy.

并且提供 `t.Space(n)`（n 个 Spacing 单位）与 `t.Rem(n)`（n 倍字号），自定义元素用它们即可与内建控件同步缩放。

这与 Ant Design 的 `compactAlgorithm` 同构——antd 通过缩小 `sizeUnit`（4→2）、`controlHeight`（32→24）、各 margin token 实现紧凑。两者映射如下：

| Ant Design（default → compact） | 本库（标准 → 紧凑） |
| --- | --- |
| `sizeUnit: 4 → 2`（内部基础单位） | `Theme.Spacing: 4 → 3`（对外基本单位） |
| `controlHeight: 32 → 24` | 控件高度由 `internal.ControlHeight` 按 `t.Space` 倍数推出 |
| `marginXS/SM/M/L: 8/12/16/24 → 4/8/12/16` | `t.Space(n)` 的 n 不变，单位自动缩小 |
| `fontSize` 基本不变 | `Theme.FontSize: 14 → 13`（当前实现会调小 1px） |
| `borderRadiusLG/SM` 不变或略减 | `Theme.Radius: 6 → 4` |

## 两档密度数值表（当前实现）

| 维度 | 标准 Comfortable | 紧凑 Compact | 说明 |
| --- | --- | --- | --- |
| `Theme.Spacing` | 4 | 3 | 核心开关，影响所有内建控件 |
| `t.Space(1)` | 4px | 3px | 最小间隙 |
| `t.Space(2)` | 8px | 6px | 控件内边距（横向，small） |
| `t.Space(3)` | 12px | 9px | 卡片内边距、gap |
| `t.Space(4)` | 16px | 12px | 区块间距、控件横向 padding（middle/large） |
| `t.Space(6)` | 24px | 18px | 大区块间距 |
| 默认按钮高度（middle） | 32px（`t.Space(8)`） | 24px | `internal.ControlHeight` |
| 大按钮 / 小按钮 | 40px / 24px（`t.Space(10)` / `(6)`） | 30px / 18px | 同上 |
| 输入框高度 | 32px | 24px | 与按钮同源（`internal.ControlHeight`） |
| 表格/树行高 | 由 Spacing 倍数决定 | 约降 25% | MyGo 内建按 Spacing 计算 |
| 复选框/开关尺寸 | `t.Space(n)` | 自动缩小 | Base 控件随 Spacing |
| `Theme.FontSize` | 14 | 13 | 由 `WithDensity` 设置 |
| `Theme.Radius` | 6 | 4 | 由 `WithDensity` 设置 |
| 滚动条 | MyGo 默认 | 未单独调整 | 可用 `ScrollbarStyle` 做更窄的 IDE 风格（见下） |

> **未实现**：早期草案提到的 `UltraCompact`（Spacing 2.5）目前没有 `Density` 常量；如需极端密度，可直接 `theme.WithSpacing(2.5)` 试运行，但不保证所有控件视觉完整。

## 控件层落地规则

1. **禁止字面间距**：换肤控件内部所有 `Padding` / `Gap` / `Margin` / `Size` / `Width` / `Height` 必须写成 `t.Space(n)` 或 `t.Rem(n)`，绝不能写死像素。这是密度能生效的前提（见 [02 统一范式](02-widget-mapping.md#六换肤统一范式)）。
2. **自定义尺寸同步**：缺口控件（`extra/`）自己的 padding、图标尺寸、行高同样用 `t.Space` / `t.Rem`。
3. **Size 参数**：控件级尺寸用强类型枚举，而不是字符串——
   `ant.Size(ant.Small)` / `ant.Middle` / `ant.Large`（`ant.SizeKind = internal.SizeKind`）。
   它决定 `ControlHeight` / `ControlPaddingX` / `ControlFont`，与密度正交：`紧凑 + Large` 合法。
4. **浮层与表格**：弹层 padding、表格行高、树缩进是最容易被遗漏的地方，需专门用紧凑档截图核对。
5. **固定像素的例外**：1px 描边、焦点环宽度等保持不变，避免视觉破碎；`internal/check` 对 0 与 1px 有豁免。

## 运行时切换

与主题共用同一机制（见 [04 主题系统](04-theme-system.md)）：

```go
// 用户在设置里切换密度，下一帧全部控件同步。
theme.Apply(c, theme.Light(theme.WithDensity(theme.Compact)))
```

- 切换发生在 `Apply` 内部调整 `Spacing` / `Radius` / `FontSize` 后 `SetTheme`；
- 无需遍历控件、无需手动重排；
- 可与亮/暗、品牌色**自由组合**（如「暗色 + 紧凑 + 品牌紫」）；
- 支持局部：在某个面板内用 `theme.Sub(c, compactToken)` 单独紧凑，其余保持标准。

> **Option 顺序**：`WithDensity` 会重置度量；需要细调字号时把 `WithFontSize` / `WithRadius` / `WithSpacing` 放在它**之后**。

## 面向 IDE / 上位机的额外约定（部分未实现）

| 优化 | 状态 | 做法 |
| --- | --- | --- |
| 更窄的滚动条 | ◻ 未实现 | `ScrollbarStyle{TrackWidth/ThumbWidth/Square/Flush/MinLength}`，MyGo 原生支持 VS Code 风格配置 |
| 表格斑马纹减弱 | ◻ 未实现 | 用 `t.Surface.Alpha(...)` 轻量区分行 |
| 控件边框收敛 | ✅ 已贯彻 | 输入框/按钮边框用 `t.Border` 细线，聚焦时才用 `t.Focus` |
| 字号优先可读性 | ✅ 已贯彻 | 紧凑只缩间距；字号下限建议 12px，`WithDensity` 当前为 13px |
| 密度选择持久化 | ◻ 未实现 | 存入 app state / 配置，重启恢复；可配合 `Window.StateKey` |

## 验收清单

当前可用以下检查；完整清单是后续目标。

- [x] 标准与紧凑两档下，全部已实现控件切换密度后外观变化（Gallery 控制条可切换）。
- [x] 亮/暗 × 标准/紧凑 4 种组合的 Gallery 截图基线通过（`samples/golden/`）。
- [ ] 表格、树、弹层、表单在紧凑档下逐控件人工核对（当前只有 Gallery 级基线，缺逐控件截图）。
- [ ] 运行时切换密度无闪烁、无需重启（即时模式已满足，仍缺自动化断言）。
- [ ] IDE 风格示例应用（侧栏 / 标签页 / 窄滚动条）落地（E3 任务，未开始）。
