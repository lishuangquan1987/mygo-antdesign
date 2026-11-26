# 03 · 缺口控件实现方案

> 对应需求 4：用 MyGo 实现 Ant Design 有而 MyGo 没有的控件。

## 判定方法

以 Ant Design 官方组件总览（通用 4 / 布局 7 / 导航 7 / 数据录入 18 / 数据展示 21 / 反馈 11）为清单，逐项与 MyGo `docs/ui` 的 70+ 组件比对。下表只列 **MyGo 无直接对应** 的控件；有对应的见 [02 映射](02-widget-mapping.md)。

MyGo 可复用的关键底座（来自 `custom-widgets.md`）：`ButtonBase`、`CheckboxBase/StateBase`、`RadioBase`、`SliderBase`、`TabsBase`、`SegmentedBase`、`CollapsibleBase`、`AccordionBase`、`TreeBase`、`SelectBase`、`ComboboxBase`、`PopoverBase`、`DialogBase`、`MenuBase/ContextMenuBase`、`TooltipBase`、`ToastViewportBase/ToastBase`、`TextInputBase/TextAreaBase`。浮层可用 `AttachTo` / `PressedOutside` / `OverlayShortcut` / `Modal` 定位与关闭。

## 一、优先级 P1 — 高频、易实现（建议 M1–M2 完成）

| 缺口控件 | antd 组件 | 实现路径（MyGo 原语） | 复杂度 |
| --- | --- | --- | --- |
| **Pagination 分页** | Pagination | `Row` + `ButtonBase` 页码 + `SelectBase` 每页条数；页码状态存 `ui.Local`/app state | 低 |
| **Card 卡片** | Card | `ui.Box` + `Radius` + `Background(t.Surface)` + `Shadow` + 标题/内容分区。MyGo 无 Card 控件，纯组合 | 低 |
| **Tag 标签** | Tag | `Box` + 小 `Radius` + `t.Surface` 背景 + 可选 `ButtonBase` 关闭图标；可关闭/可选中态 | 低 |
| **Empty 空状态** | Empty | 居中 `Column` + SVG 插画 + `ui.Text(t.TextMuted)` | 低 |
| **Result 结果** | Result | 状态图标（success/error/warning/info）+ 标题 + 描述 + 操作区，纯组合 | 低 |
| **Skeleton 骨架屏** | Skeleton | `Box` + `Background(t.SurfaceHover)` + `Radius`，用 `ui.Motion`/`AnimateTo` 做呼吸动画 | 低 |
| **Divider 分割线** | Divider | `Box` + `BorderWidth(0,0,1,0)` + `BorderStyle(ui.BorderDashed)`，可带文字 | 低 |
| **Statistic 统计数值** | Statistic | `Text` 大字号 `t.Rem` + 趋势色（涨红跌绿，遵循中文习惯） | 低 |
| **Descriptions 描述列表** | Descriptions | `Grid`/`Row` + 标签(`t.TextMuted`)/值成对布局 | 低 |
| **Timeline 时间轴** | Timeline | `Column` + 节点圆点 `Box` + 竖线 `Border`，纯组合 | 中 |

## 二、优先级 P2 — 中频、需浮层或状态（M2–M3）

| 缺口控件 | antd 组件 | 实现路径（MyGo 原语） | 复杂度 |
| --- | --- | --- | --- |
| **Cascader 级联选择** | Cascader | `SelectBase` 嵌套多级 `Popup`，每级一列 `List`，选项 `*T` 绑定；复用 `TreeBase` 的层级思路 | 高 |
| **TreeSelect 树选择** | TreeSelect | `SelectBase` 触发器 + `TreeBase` 作为 Popup 内容，勾选态复用 `CheckboxStateBase` | 高 |
| **Transfer 穿梭框** | Transfer | 两列 `List/Table` + 中间 `ButtonBase` 移动按钮，左右数据源为 app state | 中 |
| **Mentions 提及** | Mentions | `TextInputBase` + 检测 `@` 触发 `PopoverBase` 候选列表，插入逻辑自管 | 高 |
| **Drawer 抽屉** | Drawer | `DialogBase` 改为边缘锚定：`AttachTo` + `Modal` + 滑入 transition；关闭同 Dialog | 中 |
| **Popconfirm 气泡确认** | Popconfirm | `PopoverBase` + 确认/取消 `ButtonBase`，anchored 模式 | 中 |
| **Watermark 水印** | Watermark | `ui.Draw` / `Painter` 铺满平铺文字，或 `ui.Overlay` + 旋转文本；`PassThrough()` 不挡事件 | 中 |
| **Tour 漫游式引导** | Tour | `PopoverBase` 依次锚定目标元素 + 上一步/下一步；步骤状态存 app state | 高 |
| **Carousel 走马灯** | Carousel | `ui.Motion` + 横排 `Row` 位移，指示点用 `ButtonBase` | 中 |
| **Affix 固钉** | Affix | MyGo 有 `ui.Sticky`/`sticky.go`，语义等价，直接换肤即可 | 低 |

## 三、优先级 P3 — 依赖外部能力或平台语义（M3–M4）

| 缺口控件 | antd 组件 | 实现路径（MyGo 原语） | 复杂度 |
| --- | --- | --- | --- |
| **QRCode 二维码** | QRCode | 引入纯 Go 编码库（如 `skip2/go-qrcode`）生成位图，用 `ui.Image`/`ui.Bitmap` 或 `Painter` 逐格绘制 | 中 |
| **Upload 上传** | Upload | **桌面语义**：本地文件选择走 `mygo` 的 dialog API + `ReadFile`；拖放走 `ui.Drop`/`DragAndDrop`；进度条复用换肤后的 Progress | 高 |
| **FloatButton 悬浮按钮** | FloatButton | `Box` `Absolute()` 锚定角落 + `ButtonBase`，可选 BackTop 滚动监听 | 低 |
| **Masonry 瀑布流** | Masonry | `Column` 分桶 + 按高度分配，虚拟滚动复用 `ui.List` 思路 | 高 |
| **Anchor 锚点** | Anchor | `ui.Outline`/侧栏 + 滚动位置联动高亮 | 中 |
| **BorderBeam 边框流光** | BorderBeam | `BorderConicGradient` + `AnimateTo` 旋转角度，MyGo 原生支持渐变边框 | 中 |
| **Notification 通知提醒** | Notification | `ToastViewportBase` 自定义右上角 viewport + 关闭按钮，与 Message 分区 | 中 |
| **Spin 加载中** | Spin | 用 `ui.Motion` + SVG 旋转；loading 插件已有 29 种指示器可复用 | 低 |

## 四、明确不移植的控件

- **App / ConfigProvider / Util**：这些是 Web 框架层概念。本库中 ConfigProvider 的职责由 `theme/` 包承担，App 包裹由 `ui.View` 承担，无需伪造同名控件。
- **List（已废弃）/ Grid 语义重复**：MyGo 的 `ui.List`/`ui.Grid` 已覆盖。
- **强 Web 语义的原生能力**：如 `Affix` 在桌面窗口中意义有限，仅保留等价实现。

## 五、缺口控件的统一开发规范

1. **位置**：全部放在 `extra/` 包，不污染 `widgets/`（一一对应换肤控件）。
2. **签名**：与新 Ant Design 风格一致的函数式 API，例如 `extra.Pagination(c, &page, total)`。
3. **主题**：与换肤控件共用同一套 Token，禁止自带硬编码颜色。
4. **密度**：必须支持标准/紧凑两档（见 [05](05-compact-density.md)）。
5. **无障碍**：基于 Base 的实现天然带键盘与屏幕阅读器支持；纯组合控件需补 `Role`/`Label`。
6. **自测**：用 Gallery 的无窗口渲染出 PNG 做视觉回归。

## 六、实现顺序建议

```
P1（低风险、快速见效）
  Card → Tag → Empty → Result → Divider → Statistic → Skeleton → Timeline
P2（浮层与状态）
  Pagination → Drawer → Popconfirm → Descriptions → Transfer → Cascader/TreeSelect → Tour
P3（依赖能力）
  Spin → FloatButton → BorderBeam → Notification → QRCode → Upload → Anchor → Carousel → Masonry
```

每完成一个控件，同步在 Gallery 中添加预览与紧凑密度对照。
