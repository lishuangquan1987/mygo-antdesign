# 10 · 控件改造清单（Checklist）

> 两个清单：**① MyGo 原生控件**（换肤/包装改造）、**② antd 有而 MyGo 没有的缺口控件**（自建）。
> 完成改造的打勾（`[x]`）。打勾标准：有实现 + 单元测试 + `samples/` 示例 + 截图基线（或属于 🧱 无需包装的底座/布局原语，已作出结论）。
>
> 核对基线：MyGo **v0.10.0**（`go.mod` 锁定）、本库 v0.5.0，核对日期 2026-11-27。
> 详细说明见 [02 映射表 §九](02-widget-mapping.md) 与 [03 缺口控件](03-gap-widgets.md)。

## 一、MyGo 原生控件清单（换肤/包装改造）

### 1.1 通用与基础

- [x] `ui.Button` → `ant.Button`（五型 × danger × 三档尺寸 × Loading/Block）
- [x] `ui.PrimaryButton` → `ant.PrimaryButton`
- [x] `ui.Text` / `ui.Textf` → `ant.Title` / `Paragraph` / `Code` / `TextMuted`
- [x] `ui.RichText` → `ant.RichText`
- [x] `ui.Icon` → `icon.Icon`（44 个 antd 图标）/ `ant.Icon`
- [x] `ui.Link` → `ant.Hyperlink`
- [x] `ui.Divider` → `ant.Divider` / `VerticalDivider` / `DashedVerticalDivider`
- [x] `ui.Row` / `ui.Column` / `ui.Grid` / `ui.Box` / `ui.Spacer` → `ant.Space` / `Flex` / `Grid`
- [x] `ui.Scroll` / `ui.ScrollBoth` / `ui.ScrollHorizontal` → `ant.Scroll` / `ScrollBoth` / `ScrollHorizontal`
- [x] `ui.Split` / `ui.SplitPercent` / `ui.SplitVertical` / `ui.SplitVerticalPercent` → `ant.Split*`（另有 antd 风格 `ant.Splitter` 系列）
- [x] `ui.Toolbar` → `ant.Toolbar`
- [x] `ui.Breadcrumbs` → `ant.Breadcrumbs`（别名 `Breadcrumb`）
- [x] `ui.MenuButton` → `ant.MenuButton`（另有声明式 `ant.Dropdown`）
- [x] `ui.PopoverMenuButton` → `ant.PopoverMenuButton`
- [x] `ui.SplitButton` / `ui.PrimarySplitButton` → `ant.SplitButton` / `PrimarySplitButton`
- [x] `ui.BackButton` / `ui.ForwardButton` → `ant.BackButton` / `ForwardButton`
- [x] `ui.AlertDialog` → `ant.AlertDialog`
- [x] `ui.FindBar` → `ant.FindBar`
- [x] `ui.Modal` → `ant.Modal`（头/体/脚 antd 内边距 + 打开缩放淡入动效）
- [x] `ui.Popover` → `ant.Popover`
- [x] `ui.Toggle` / `ui.ToggleGroup` / `ui.ToggleBase` → `ant.Toggle` / `ToggleGroup` / `ToggleBase`
- [x] `ui.Spinner` → `ant.Spinner`（别名 `Spin`）
- [x] `ui.Overlay` / `ui.Themed` — 🧱 渲染/主题底座，无需包装

### 1.2 数据录入

- [x] `ui.TextInput` / `ui.TextArea` → `ant.Input` / `TextArea`
- [x] `ui.TextInputBuffer` / `ui.TextAreaBuffer` → `ant.TextInputBuffer` / `TextAreaBuffer`
- [x] `ui.NumberInput` → `ant.NumberInput`（别名 `NumberInputV`）
- [x] `ui.DecimalInput` → `ant.DecimalInput`
- [x] `ui.SearchField` → `ant.SearchField`
- [x] `ui.Select` → `ant.Select`（`ui.SelectBase` 自绘）
- [x] `ui.Combobox` → `ant.Combobox`（`ui.ComboboxBase` 自绘）
- [x] `ui.Autocomplete` → `ant.Autocomplete`（别名 `AutoComplete`）
- [x] `ui.Checkbox` → `ant.Checkbox`
- [x] `ui.CheckboxGroup` → `ant.CheckboxGroup`
- [x] `ui.CheckboxState` → `ant.CheckboxState`（别名 `CheckboxAll`）
- [x] `ui.RadioGroup` → `ant.RadioGroup` / `ant.Radio[T]`
- [x] `ui.Switch` → `ant.Switch` / `SwitchWithLabel`（自绘 44×22）
- [x] `ui.Slider` → `ant.Slider`
- [x] `ui.StepSlider` → `ant.StepSlider`
- [x] `ui.RangeSlider` → `ant.RangeSlider`（别名 `RangeSliderDouble`）
- [x] `ui.Stepper` → `ant.Stepper`
- [x] `ui.Rating` → `ant.Rating`（别名 `Rate`）
- [x] `ui.ColorWell` / `ui.ColorPicker` → `ant.ColorWell` / `ColorPicker` / `ColorPanel`
- [x] `ui.Calendar` → `ant.Calendar`
- [x] `ui.DateInput` → `ant.DateInput`（别名 `DatePicker`）
- [x] `ui.TimeInput` → `ant.TimeInput`（别名 `TimePicker`）
- [x] `ui.Form` / `ui.Field` / `ui.Fieldset` → `ant.Form` / `FormItem` / `FormGroup`
- [x] `ui.TokenField` → `ant.TokenField`
- [x] `ui.EditableText` → `ant.EditableText`
- [x] `ui.Segmented` → `ant.Segmented`
- [x] `ui.Tabs` → `ant.Tabs`（自绘 antd 下划线页签）

### 1.3 数据展示与导航

- [x] `ui.Table` → `ant.Table`（自绘 antd 表格，数值列按数值排序）
- [x] `ui.List` → `ant.List`
- [x] `ui.Tree` / `ui.TreeItem` → `extra.Tree`
- [x] `ui.GridView` → `ant.GridView`
- [x] `ui.Avatar` → `ant.Avatar`（别名 `AvatarEl`）
- [x] `ui.Badge` → `ant.Badge` / `BadgeDot`（计数骑角为自建）
- [x] `ui.Progress` → `ant.Progress`（别名 `ProgressLine`）+ `ProgressCircle`
- [x] `ui.Meter` → `ant.Meter`
- [x] `ui.Image` → `ant.Image`（别名 `ImageEl`）
- [x] `ui.Video` → `ant.Video`
- [x] `ui.WebView` → `ant.WebView`
- [x] `ui.NativeView` → `ant.NativeView`
- [x] `ui.Sidebar` / `ui.SidebarItem` / `ui.SidebarSection` → `ant.Sidebar` 系列
- [x] `ui.Accordion` / `ui.AccordionItem` → `ant.Accordion` 系列
- [x] `ui.Collapsible` → `ant.Collapsible`

### 1.4 Base 底座（🧱 无外观，按设计直接使用 `ui.*Base`，不重复包装）

- [x] `CheckboxBase` / `CheckboxStateBase` / `SwitchBase` / `SliderBase` / `SegmentedBase` / `TabsBase` / `TreeBase` / `ButtonBase` / `TextInputBase` / `TextAreaBase` / `TextInputBufferBase` / `TextAreaBufferBase` / `ComboboxBase` / `AccordionBase` / `CollapsibleBase` / `SelectBase` / `SelectBaseOpen` / `MenuBase` / `ContextMenuBase` / `PopoverBase` / `TooltipBase` / `DialogBase` / `ToastBase` / `ToastViewportBase` — 已在自绘实现中直接使用（如 Button/Pagination/Select/Combobox/Drawer/Tooltip）

### 1.5 MyGo UI 子包适配

- [x] `ui/markdown` → `ant.Markdown`
- [x] `ui/editor` → `ant.CodeEditor` / `ant.DiffView`（Gallery 代码视图已使用）
- [x] `ui/chart` → `ant.ChartPlot`
- [x] `ui/richtext` → `ant.RichTextEditor`
- [x] `ui/netimage` → `ant.RemoteImage`
- [x] `ui/video` → `ant.OpenVideo`（配合 `ant.Video`）

## 二、antd 有而 MyGo 没有的控件清单（缺口控件自建）

### 2.1 通用与数据展示

- [x] **Alert** 提示条（四级别 + 描述 + 可关闭）
- [x] **Card** 卡片（头部 extra / 底部 actions）
- [x] **Tag** 标签（7 色系 + 可关闭）
- [x] **Badge** 计数/红点（右上角骑角，原生 `ui.Badge` 只是文本胶囊）
- [x] **Empty** 空状态
- [x] **Result** 结果页（四态）
- [x] **Skeleton** 骨架屏
- [x] **Statistic** 统计数值（趋势色）
- [x] **Descriptions** 描述列表
- [x] **Timeline** 时间轴（节点色）
- [x] **Collapse** 折叠面板（含手风琴）
- [x] **QRCode** 二维码（码色随主题）
- [x] **Carousel** 走马灯
- [x] **Image 预览**（点击放大模态）
- [x] **Progress 环形**（`ant.ProgressCircle`，状态色）
- [x] **Masonry** 瀑布流
- [x] **BorderBeam** 边框流光

### 2.2 导航

- [x] **Pagination** 分页器（窗口折叠 + 翻页箭头）
- [x] **Steps** 步骤条（三态）
- [x] **Anchor** 锚点目录
- [x] **Tour** 引导（简化版步骤卡片）
- [x] **FloatButton** 悬浮按钮

### 2.3 数据录入

- [x] **Cascader** 级联选择
- [x] **TreeSelect** 树形下拉
- [x] **Transfer** 穿梭框
- [x] **Upload** 上传（拖放区 + 文件列表）
- [x] **Mentions** 提及（@ 触发候选）
- [x] **Popconfirm** 气泡确认

### 2.4 反馈

- [x] **Tooltip** 文字提示（反色气泡）
- [x] **Modal.confirm** → `ant.Confirm`
- [x] **Message** 轻提示（toast viewport）
- [x] **Notification** 通知（右上角独立浮层）
- [x] **Watermark** 水印
- [x] **Spin** 容器加载（`ui.Spinner` + 提示文字）
- [x] **Drawer** 抽屉（打开滑入动效）

## 三、明确不移植（保持 ☐，不纳入完成统计）

- [ ] **App / ConfigProvider / Util** — 纯 Web 概念，桌面原生 UI 无对应物（AGENTS.md §6 禁止伪造）
- [ ] **Pro 组件**（ProTable / ProForm 等）— 超出 v0.x 范围（docs/06 §10）
- [ ] **`ui/mobile`** — 移动端专用控件（BottomBar/Sheet/PageView 等），桌面库不包装
- [ ] **`ui/shader`** — GPU 着色器绘制能力，非控件
- [ ] **系统亮暗自动跟随** — MyGo `ui` 未暴露系统外观 API（ADR 0002 上游限制）

## 统计

| 清单 | 总数 | 已完成 | 完成率 |
| --- | --- | --- | --- |
| MyGo 原生控件（根包 111 个构造函数 = 89 个换肤/包装 + 22 个 Base 底座结论） | 111 | 111 | **100%** |
| MyGo UI 子包适配（markdown / editor / chart / richtext / netimage / video） | 6 包 7 入口 | 6 包 7 入口 | **100%** |
| antd 缺口控件 | 35 | 35 | **100%** |

> 勾选规则重申：`[x]` = 改造完成（实现 + 单测 + 示例 + 基线，或 🧱 已决策不包装/直接使用）；`[ ]` = 待办或不移植。
> 新增控件时：先在本清单加行 `[ ]`，完成改造后改 `[x]`，并同步 [02 §九](02-widget-mapping.md) / [03](03-gap-widgets.md)。
