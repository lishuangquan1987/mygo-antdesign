# 05 · 紧凑密度风格

> 对应需求 6：在 Ant Design 标准间距之外，再加一套紧凑风格（padding/margin 调小），用于有限空间显示更多内容——上位机、写代码 IDE、VS2026 等。

## 目标

在**完全复用同一套控件实现**的前提下，提供两档（可扩为三档）密度：

- **标准（Comfortable）**：Ant Design 默认观感，适合后台管理、设置面板、普通工具。
- **紧凑（Compact）**：padding、margin、行高、控件尺寸整体缩小，适合上位机监控面板、IDE、终端类工具、数据密集表格。

**关键**：密度不是第二套控件库，而是主题层的一个维度。切换密度 = 改 `Theme.Spacing` 等度量 + 调 `SetTheme`，全部控件自动跟随。

## 为什么能「一处修改，全控件生效」

MyGo 官方对 `Spacing` 的定义（`styling.md`）：

> `Spacing` is the unit of the room widgets leave: their paddings and gaps, and the sizes of check boxes, switches, sliders and the rows of tables and trees, are multiples of it. It is 4 by default; **3 makes every widget compact**, 5 roomy.

并且提供 `t.Space(n)`（n 个 Spacing 单位）与 `t.Rem(n)`（n 倍字号），自定义元素用它们即可与内建控件同步缩放。

这与 Ant Design 的 `compactAlgorithm` 完全同构——antd 通过缩小 `sizeUnit`（4→2）、`controlHeight`（32→24）、各 margin token 实现紧凑。两者映射如下：

| Ant Design（default → compact）             | 本库（标准 → 紧凑）                              |
| ----------------------------------------- | ---------------------------------------- |
| `sizeUnit: 4 → 2`（内部基础单位）                 | `Theme.Spacing: 4 → 3`（对外基本单位）           |
| `controlHeight: 32 → 24`                  | 控件高度 `t.Space(n)` 的 n 下调 / 用 `Size()` 覆盖 |
| `marginXS/SM/M/L: 8/12/16/24 → 4/8/12/16` | `t.Space(n)` 的 n 下调                      |
| `fontSize` 基本不变                           | `Theme.FontSize` 可微调（IDE 常用 13px）        |
| `borderRadiusLG/SM` 不变或略减                 | `Theme.Radius` 可略减（如 6 → 4）              |

## 两档密度数值表（设计基线）

| 维度               | 标准 Comfortable      | 紧凑 Compact        | 说明                         |
| ---------------- | ------------------- | ----------------- | -------------------------- |
| `Theme.Spacing`  | 4                   | 3                 | 核心开关，影响所有内建控件              |
| `t.Space(1)`     | 4px                 | 3px               | 最小间隙                       |
| `t.Space(2)`     | 8px                 | 6px               | 控件内边距（横向）                  |
| `t.Space(3)`     | 12px                | 9px               | 卡片内边距、gap                  |
| `t.Space(4)`     | 16px                | 12px              | 区块间距                       |
| `t.Space(6)`     | 24px                | 18px              | 大区块间距                      |
| 默认按钮高度           | ~32px               | ~24–26px          | 由 padding + 字号推得           |
| 输入框高度            | 32px                | 24–26px           | `TextInputBase` padding 下调 |
| 表格/树行高           | 由 Spacing 倍数        | 约降 25%            | MyGo 内建按 Spacing 计算        |
| 复选框/开关尺寸         | `t.Space(n)`        | 自动缩小              | Base 控件随 Spacing           |
| `Theme.FontSize` | 系统默认（14/13）         | 可设 12–13          | IDE 场景常用                   |
| `Theme.Radius`   | 6                   | 4（可选）             | 紧凑界面更利落                    |
| 滚动条              | `ScrollbarStyle` 默认 | 可设更窄 `ThumbWidth` | MyGo 原生支持「VS Code 方形滚动条」   |

> 说明：紧凑档同时提供 `Compact`（3）与可选的 `UltraCompact`（2.5）用于极端信息密度；但 2 以下控件会过于拥挤，不建议常规使用。

## 控件层落地规则

1. **禁止字面间距**：换肤控件内部所有 `Padding`/`Gap`/`Margin` 必须写成 `t.Space(n)` 或 `t.Rem(n)`，绝不能写死像素。这是密度能生效的前提（见 02 文档统一范式）。
2. **自定义尺寸同步**：缺口控件（`extra/`）自己的 padding、图标尺寸、行高同样用 `t.Space` / `t.Rem`。
3. **Size 参数**：为常用控件提供 `Size("large"|"middle"|"small")` 选项，映射为不同的 `t.Space` 倍数，兼容 antd 的组件级 `size`。
4. **浮层与表格**：弹层 padding、表格行高、树缩进是最容易被遗漏的地方，需专门用紧凑档截图核对。
5. **固定像素的例外**：1px 描边、焦点环宽度等保持不变，避免视觉破碎。

## 运行时切换

与主题共用同一机制（见 [04 主题系统](04-theme-system.md)）：

```go
// 概念示意：用户在设置里切换密度，下一帧全部控件同步
theme.Apply(c, theme.Light(theme.WithDensity(theme.Compact)))
```

- 切换发生在 `Apply` 内部调整 `Spacing`/`FontSize`/`Radius` 后 `SetTheme`；
- 无需遍历控件、无需手动重排；
- 可与亮/暗、品牌色**自由组合**（如「暗色 + 紧凑 + 品牌紫」）；
- 支持局部：在 IDE 的某个面板内用 `theme.Sub(c, compactToken)` 单独紧凑，其余保持标准。

## 面向 IDE / 上位机的额外约定

针对 VS2026、上位机等场景，紧凑档默认再启用以下优化（可单独关闭）：

- **更窄的滚动条**：`ScrollbarStyle{TrackWidth: 14, ThumbWidth: 14, Square: true, Flush: true, MinLength: 20}`（MyGo 官方给出的 VS Code 风格配置）。
- **表格斑马纹减弱**：用 `t.Surface.Alpha(...)` 轻量区分行，减少视觉噪声。
- **控件边框收敛**：输入框、按钮边框用 `t.Border` 细线，聚焦时才用 `t.Focus` 强调。
- **字号优先可读性**：紧凑的是间距不是字号；字号下限建议 12px，避免长时间编码疲劳。
- **信息密度开关持久化**：密度选择存入 app state / 配置，重启后恢复（MyGo 有 `Window.StateKey` 与配置机制可复用）。

## 验收清单

- [ ] 标准与紧凑两档下，全部控件无硬编码像素导致的「切了没变」。
- [ ] 表格、树、弹层、表单在紧凑档下信息密度显著提升且可辨识。
- [ ] 亮/暗 × 标准/紧凑 四种组合均无视觉错位。
- [ ] 运行时切换密度无闪烁、无需重启。
- [ ] IDE 风格示例中，滚动条、标签页、侧栏workspace 符合紧凑基线。
