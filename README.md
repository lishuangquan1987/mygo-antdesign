# mygo-antdesign

> Ant Design 设计风格的 [MyGo](https://github.com/egoist/mygo) 原生 UI 控件库（**设计阶段 / RFC**）

本仓库目标：为 MyGo 的 `package ui` 原生控件提供一套完整的 **Ant Design 视觉与交互规范**，并补齐 Ant Design 有而 MyGo 暂缺的控件。所有控件共享一套主题，支持**运行时切换**；在 Ant Design 标准间距之外，额外提供一套**紧凑密度**风格，用于上位机、IDE 等信息密集型场景。

> 当前状态：**实现进行中**——主题层、Button 控件、示例框架与 Gallery 已落地（见 [docs/06](docs/06-development-plan.md) 进度章节）；其余控件按 WBS 推进。

## 一句话结论

**可行，且 MyGo 的底层设计天然契合**——它自带 headless Base、统一 Theme、`SetTheme` 运行时切换和 `Spacing` 密度单位，恰好对应 Ant Design 的「无样式组件 + Design Token + ConfigProvider + compactAlgorithm」四层能力。主要工作量在视觉规范翻译与缺口控件补齐，而非改造框架。详见 [可行性分析](docs/01-feasibility.md)。

## 设计原则

1. **不重复造轮子**：颜色、圆角、焦点环、禁用态、无障碍、键盘交互一律复用 MyGo 的 Theme 与 Base，只在「外观」层做 Ant Design 化。
2. **view = f(state)**：MyGo 是即时模式（immediate mode），控件是返回 Element 的 Go 函数，不是常驻对象。所有 API 按函数式设计，不引入命令式 `setText` 式调用。
3. **无 CSS 选择器**：MyGo 明确「State, not selectors」，因此 Ant Design 的 className / CSS-in-JS 动态样式要翻译成 Go 的条件分支与 Theme 字段。
4. **主题单一真源**：一套 Token 映射表驱动全部控件，杜绝每个控件各自硬编码颜色。
5. **密度是一等公民**：标准（Ant Design 默认）与紧凑两套密度共用同一组件实现，靠参数切换。
6. **示例是一等交付物**：每个控件都配可运行示例，展示不同主题、尺寸、密度与状态，代码与效果一一对应（见 [Samples & Gallery](docs/07-samples-gallery.md)）。

## 控件示例 Gallery

像 WPF 的 Material Design 示例程序一样，本库提供一个可运行的 **Gallery 应用**：左侧控件目录、中间实时预览、旁边可复制的真实代码，顶部可切换主题（亮/暗/品牌色）、密度（标准/紧凑）、尺寸（大/中/小），并支持**同屏并排对比**（借助 MyGo 的 `ui.Themed`）。示例同时用于视觉回归截图与文档插图。设计见 [07 Samples & Gallery](docs/07-samples-gallery.md)。

## 文档导航

| 文档 | 覆盖的需求 | 内容 |
| --- | --- | --- |
| [01 可行性分析](docs/01-feasibility.md) | 需求 2 | 技术契合度、架构分层、工作量、风险与里程碑 |
| [02 控件换肤映射](docs/02-widget-mapping.md) | 需求 3 | MyGo 全部原生控件 → Ant Design 皮肤的逐项改造方案 |
| [03 缺口控件实现](docs/03-gap-widgets.md) | 需求 4 | Ant Design 有而 MyGo 无的控件，基于 Base / Element 的实现路径 |
| [04 主题系统](docs/04-theme-system.md) | 需求 5 | Token 映射、运行时切换、全局/局部主题、预设 |
| [05 紧凑密度](docs/05-compact-density.md) | 需求 6 | 标准/紧凑两套密度规范、数值表、运行时切换 |
| [06 开发计划](docs/06-development-plan.md) | 整体 | WBS、里程碑、时间线、DoD、风险 |
| [07 Samples & Gallery](docs/07-samples-gallery.md) | 示例展示 | 每控件示例、四维展示、并排对比、真实代码、截图回归 |

## 目录约定（后续实现阶段）

```
mygo-antdesign/
├── theme/        # Token 定义、预设主题、与 MyGo Theme 的桥接
├── widgets/      # 换肤后的控件（与 MyGo 原生一一对应）
├── extra/        # 用 Base / Element 补齐的缺口控件
├── icon/         # Ant Design 图标的 SVG 路径数据（Go 常量）
├── samples/      # 每控件示例源码（真实可编译，go:embed 提取代码片段）
├── gallery/      # 示例展示应用：侧栏 + 实时预览 + 代码 + 主题/尺寸/密度切换 + 并排对比
└── docs/         # 本文档集（含自动生成的截图）
```

## License

待定（建议 MIT，与 MyGo 保持一致）。
