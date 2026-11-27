# 03 · 缺口控件与待补控件

> 对应需求 4：用 MyGo 实现 Ant Design 有而 MyGo 没有的控件。
>
> **核对基线**：MyGo v0.10.0（`go.mod` 锁定）、本库当前代码，核对日期 2026-11-27。
> 本文同时是**「还剩什么」的清单**：§一 已交付，§二 待实现，§三 明确不移植。

## 判定方法（已修正）

1. 以 Ant Design 组件总览（通用 / 布局 / 导航 / 数据录入 / 数据展示 / 反馈）为清单。
2. 与 MyGo v0.10.0 `package ui` 的**实际导出名**逐项比对（不是凭印象写控件名）。
3. 只有「MyGo 既没有同名控件、也没有等价组合原语」才归为缺口；有等价原语的（如 `ui.Divider`、`el.Sticky()`）不算缺口。
4. 缺口控件可以放在 `widgets/`（与原生控件语义同族，如 `ant.Pagination`、`ant.Steps`）或 `extra/`（跨原语的复杂组合，如 Tree/Upload/Transfer）。

> **早期草案的四处误判**（已在本版修正）：
> - `Divider`：v0.10.0 已有 `ui.Divider`，不是缺口；虚线/带文字由 `ant.Divider` 在原生基础上组合。
> - `Pagination`：已实现于 `widgets/nav.go`（`ant.Pagination`），不在 `extra/`。
> - `Timeline`：`ui.Timeline` 是**动效时间轴类型**，不是控件；antd Timeline 仍是待实现缺口。
> - `Spin` / `Affix`：MyGo 有 `ui.Spinner` 与元素方法 `el.Sticky().Top(...)`，只需薄包装，不属于「需要从零实现」。

## 一、已交付的缺口控件（v0.1.0 – v0.3 开发中）

| antd 组件 | 本库 API | 位置 | 组合方式 | 复杂度 | 备注 |
| --- | --- | --- | --- | --- | --- |
| Alert | `ant.Alert` | widgets | `Row` + `Box` + `Text` + SVG 图标 | 低 | 四级别（success/info/warning/error）+ 描述 + 可关闭 |
| Card | `ant.Card` | widgets | `Column` + `Radius` + `Surface` + 头/体/操作区 | 低 | 头部 extra、底部 actions |
| Tag | `ant.Tag` | widgets | `Box` + 小圆角 + 语义色派生（7 色系） | 低 | 可关闭；暗色由 `Mix` 自动成立 |
| Badge（计数/红点） | `ant.Badge` / `BadgeDot` | widgets | `ui.Overlay` + `AttachTo` 骑角 | 低 | 原生 `ui.Badge` 只是文本胶囊 |
| Pagination | `ant.Pagination` | widgets | `Row` + `ButtonBase` 页码 + 窗口折叠 | 中 | `OnPageChange`；当前固定每页 10 条 |
| Steps | `ant.Steps` | widgets | `Row` + 圆点 + 连接线 | 中 | 完成/进行/等待三态 |
| Message / Notification | `ant.Message` / `MessageTypeTyped` / `Notify` | widgets | `c.Toast` / `c.AddToast` | 低 | 复用 MyGo 底部 toast viewport，不是独立浮层 |
| Modal.confirm / Popconfirm | `ant.Confirm` | widgets | `ui.AlertDialog` | 低 | 无返回值，只回调 `onOk`；Popconfirm 仍未实现 |
| Divider（虚线/带文字/竖线） | `ant.Divider` / `VerticalDivider` / `DashedVerticalDivider` | widgets | `Box` + `BorderWidth` | 低 | 实线可直接用原生 `ui.Divider` |
| Tree | `extra.Tree` | extra | `ui.Tree` + `TreeItem` 递归 | 中 | 声明式节点 + 展开状态 + `OnTreeNodeSelect` |
| Upload | `extra.Upload` | extra | 系统文件对话框 + `ui.Drop` / `DragOver` + 列表 | 高 | 桌面语义；Tester 无窗口环境不能点 |
| Transfer | `extra.Transfer` | extra | 两列 `Column` + 行 `ButtonBase` + 中间移动按钮 | 中 | 右列选中 Key 由调用方持有 |
| Carousel | `extra.Carousel` | extra | 当前页容器 + `ButtonBase` 指示点 | 中 | 直接切页，暂无滑入动画；`el.AnimateTo` 可后补 |
| QRCode | `extra.QRCode` | extra | `skip2/go-qrcode` → `ui.NewBitmap` → `ui.Image` | 中 | 码色随主题（前景=Text，背景=Surface） |
| Cascader | `extra.Cascader` | extra | `ui.Popover` + 多列 `Column` | 高 | 逐级下钻，叶子写回路径 |
| TreeSelect | `extra.TreeSelect` | extra | `ui.Popover` + `extra.Tree` | 高 | 选中后收起 |
| Empty | `ant.Empty` | widgets | `Column` + 虚线框 + 图标 + 描述 + 操作区 | 低 | 默认插画由图标 + 主题色生成，支持自定义 image/footer |
| Result | `ant.Result` | widgets | 状态图标 + 标题 + 副标题 + 操作区 | 低 | success/info/warning/error 四态 |
| Skeleton | `ant.Skeleton` | widgets | `Box` + `SurfaceHover` + `Radius` | 低 | 静态骨架（头像/标题/段落行数可配），暂无呼吸动画 |
| Statistic | `ant.Statistic` | widgets | `Text` 大字号 + 趋势色 | 低 | 前缀/后缀/精度/涨跌（红涨绿跌） |
| Descriptions | `ant.Descriptions` | widgets | `Grid` 键值成对布局 | 低 | 列数/标题/外框（简化：外框 + 内边距） |
| Timeline（antd） | `ant.Timeline` | widgets | `Column` + 圆点 + 连接线 | 中 | 节点色（蓝/绿/红/灰）、标题 + 描述 |
| Spin | `ant.Spin` | widgets | `ui.Spinner` + 提示文字 | 低 | 可容器内居中；复用 MyGo 原生指示器 |
| Tooltip | `ant.Tooltip` | widgets | `ui.TooltipBase` 反色气泡 | 低 | 悬停 0.6s / 键盘聚焦立即显示，Esc 关闭 |
| Popconfirm | `ant.Popconfirm` | widgets | `ui.Popover` + `ButtonBase` | 中 | 触发器点击弹出，确认/取消回调 |
| Collapse | `ant.Collapse` | widgets | `ButtonBase` 头部 + 图标箭头 + 内容区 | 中 | 手风琴模式、默认展开、按 Key 记忆 |
| Form | `ant.Form` / `FormItem` / `FormGroup` | widgets | `ui.Form` / `Field` / `Fieldset` | 低 | 标签列对齐、必填标记、帮助/错误文字 |
| List | `ant.List` | widgets | `ui.List`（虚拟滚动）+ 行内边距与分隔线 | 低 | 行选中/多选由 `ListState` 提供 |
| Sidebar / Menu | `ant.Sidebar` / `SidebarItem` / `SidebarItemWithIcon` / `SidebarSection` | widgets | `ui.Sidebar` 系列 | 低 | 分区 + 条目 + 选中态，键盘可导航 |
| Splitter | `ant.Splitter` / `SplitterVertical` / `SplitterPercent` / `SplitterVerticalPercent` | widgets | `ui.Split` 系列 | 低 | 分隔条可拖拽，尺寸/百分比两种模式 |
| AutoComplete | `ant.AutoComplete` | widgets | `ui.Autocomplete` | 低 | 保留键入内容 + 候选建议 |
| Drawer | `ant.Drawer` | widgets | `ui.DialogBase` + 右侧锚定 + `Scroll` 内容区 | 中 | 遮罩/标题栏/关闭按钮，点遮罩或 Esc 关闭 |
| Image | `ant.ImageEl` | widgets | `ui.Image` + `ButtonBase` + `ui.Modal` | 中 | 可指定宽度，点击放大预览 |
| ColorPicker | `ant.ColorPicker` / `ColorPanel` | widgets | `ui.ColorWell` / `ui.ColorPicker` | 低 | 色块 + 展开面板 / 常驻面板 |
| Progress（环形） | `ant.ProgressCircle` | widgets | `ui.Box.Draw` + `ui.Path` 圆弧 | 中 | 中心百分比，成功/错误状态色 |
| Toggle | `ant.Toggle` / `ToggleGroup` | widgets | `ui.Toggle` / `ui.ToggleGroup` | 低 | 保持按下的按钮 |
| Affix | `ant.Affix` | widgets | `el.Sticky().Top(...)` | 低 | 吸顶，无需额外状态 |

> 说明：`ant.AvatarEl`、`ant.Breadcrumb`、`ant.Tabs`、`ant.Dropdown`、`ant.Space` / `Flex` / `Grid` 等是**原生控件换肤**，不是缺口控件，见 [02 映射](02-widget-mapping.md)。

## 二、待实现清单（按优先级）

> 复杂度是 1 名熟悉 Go 的开发者的粗估，不含示例与截图回归。

### P1 · 纯组合、低风险（✅ 已完成，见 §一）

Empty / Result / Skeleton / Statistic / Descriptions / Timeline / Spin 七项
已在 `widgets/data.go` 与 `widgets/feedback.go` 实现，并配单元测试与
`samples/` 示例（`samples/display`、`samples/feedback`）。

### P2 · 需要浮层 / 状态（部分完成）

已完成：Popconfirm、Collapse、Form、Tooltip、List、Sidebar、Splitter、AutoComplete、Drawer、Image、ColorPicker、环形 Progress、Toggle、Affix（见 §一）。剩余：

| antd 组件 | MyGo 原语 | 实现要点 | 复杂度 |
| --- | --- | --- | --- |
| **Watermark 水印** | `ui.NewCanvas` + `Painter` + `ui.Overlay` | `el.PassThrough()` 不挡事件 | 中 |
| **Tour 漫游引导** | `ui.PopoverBase` | 逐步锚定 + 上一步/下一步，步骤状态由调用方持有 | 高 |
| **Mentions 提及** | `ui.TextInput` + `ui.PopoverBase` | 检测 `@` 触发候选，插入逻辑自管 | 高 |

### P3 · 依赖额外能力或强 Web 语义

| antd 组件 | MyGo 原语 | 实现要点 | 复杂度 |
| --- | --- | --- | --- |
| **FloatButton 悬浮按钮** | `Box` + `Absolute()` + `ButtonBase` | 角落锚定；BackTop 用 `ScrollState` | 低 |
| **Anchor 锚点** | `ui.Outline` + `ScrollState` | 滚动位置联动高亮 | 中 |
| **BorderBeam 边框流光** | `el.BorderConicGradient(...)` + `el.AnimateTo(...)` | 旋转角度动画；MyGo 原生渐变边框 | 中 |
| **Notification（独立浮层）** | `ui.ToastViewportBase` | 自定义右上角 viewport + 关闭按钮 | 中 |
| **Masonry 瀑布流** | `Scroll` + 分桶 `Column` | 按高度分配到最短列 | 高 |

> QRCode / Upload / Transfer / Carousel / Cascader / TreeSelect 已在 §一 交付，无需再排期。

## 三、明确不移植

- **App / ConfigProvider / Util**：Web 框架层概念。ConfigProvider 的职责由 `theme/` 承担，App 包裹由 `ui.View` 承担，不伪造同名控件。
- **Pro 组件（ProTable / ProForm / ProLayout 等）**：商业组件，超出 v0.x 范围。
- **移动端专有组件**：MyGo 已原生支持移动端，后续版本再评估。
- **强 Web 语义的原生能力**：如 SSR、CSS-in-JS、路由级 Code Splitting，不映射到桌面控件。
- **SplitButton**：antd 的「主按钮 + 下拉箭头」在桌面端可由 `ant.Dropdown` + `ant.Button` 组合表达，不再单独提供同名控件。

## 四、缺口控件统一开发规范

1. **位置**：与原生控件同族的放 `widgets/`（如 Pagination、Steps、Alert）；跨多个原语的复杂组合放 `extra/`（如 Tree、Upload、Transfer）。
2. **签名**：函数式选项 + 贴近原生，例如 `ant.Pagination(c, &page, total, ant.OnPageChange(fn))`、`extra.Cascader(c, &path, nodes, extra.CascaderOnChange(fn))`。
3. **主题**：与换肤控件共用同一套 Token，禁止硬编码颜色；`internal/check` 会扫描 `widgets/` 与 `extra/`。
4. **密度**：必须同时支持标准/紧凑（间距走 `t.Space` / `t.Rem`，见 [05](05-compact-density.md)）。
5. **无障碍**：基于 Base 的实现天然带键盘与朗读；纯组合控件需补 `Role` / `Label`。
6. **示例与回归**：在 `samples/<分类>/` 写可编译示例并注册到 `gallery/main.go`，跑 Gallery 截图基线（见 [07](07-samples-gallery.md)）。

## 五、实现顺序（更新）

```
已完成  Card → Tag → Badge → Pagination → Steps → Alert → Message/Notify → Confirm
        → Tree → Upload → Transfer → Carousel → QRCode → Cascader → TreeSelect
        → Empty → Result → Skeleton → Statistic → Descriptions → Timeline → Spin
        → Popconfirm → Collapse → Tooltip → Form → List → Sidebar → Splitter → AutoComplete
        → Drawer → Image → ColorPicker → 环形 Progress → Toggle → Affix
下一步  Watermark → Mentions → Tour
再后    Anchor → Notification
按需    FloatButton → BorderBeam → Masonry
不提供  SplitButton（由 ant.Dropdown 覆盖，见 §三）
```

每完成一个控件，同步在 `samples/` 添加示例、在 `gallery/main.go` 注册，并生成/更新截图基线。
可用 `go run ./scripts/visualcheck -group 数据展示 -theme light` 把某个示例分组的
每个示例渲染成 PNG（`ui.NewTester` 无窗口渲染），人工核对效果。
