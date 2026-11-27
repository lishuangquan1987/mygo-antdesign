# 02 · MyGo 原生控件 → Ant Design 皮肤映射

> 对应需求 3：把 MyGo 自带控件加上 Ant Design 设计。
>
> **核对基线**：MyGo **v0.10.0**（`go.mod` 锁定）、Go 1.27.4、本库当前代码（v0.1.0 已发布，v0.3 开发中），核对日期 2026-11-27。
> 本表只把 `package ui` **实际存在的导出名**作为映射来源；早期草案里写错的 `ui.Dialog`、`ui.Input`、`ui.Tooltip`、`ui.Chart` 等名字已全部修正（见 §八 勘误表）。

## 映射图例

- **换肤方式**
  - `链式` — 直接在原生控件返回的 Element 上做样式调用
  - `Base` — 用对应的 `XxxBase` 重建外观（需要自定义内部结构）
  - `组合` — 用若干基础 Element / 控件拼装
  - `原生` — 不包装，直接使用 MyGo 控件，仅让主题色跟随
- **状态**：`✅` 已在本库实现（`widgets/` 包名 `ant`，或 `extra/`）；`◻` 已列入计划未实现；`—` 明确不移植/保留原生。

## 一、通用与基础

| MyGo v0.10.0 控件 | Ant Design 对应 | 换肤方式 | 本库 API | 状态 | 要点 |
| --- | --- | --- | --- | --- | --- |
| `ui.Button` / `ui.PrimaryButton` | Button | `ButtonBase` | `ant.Button` | ✅ | 五型（default/primary/dashed/text/link）× danger × 三档尺寸 × Loading/Block/OnClick |
| `ui.Text` / `ui.Textf` | Typography | 链式 | `ant.Title` / `Paragraph` / `Code` / `TextMuted` | ✅ | H1–H5、Secondary/Delete/Strong/Underline/Italic/Marked |
| `ui.Icon`（SVG） | Icon | 组合 | `icon.Icon` | ✅ | 44 个 antd outlined 常量，颜色跟随文本色 |
| `ui.Row` / `ui.Column` / `ui.Grid` | Flex / Grid / Space | 链式 | `ant.Space` / `Flex` / `Grid` | ✅ | `Gap` 走 `t.Space(n)`，实现 antd gutter |
| `ui.Divider` | Divider | 组合 | `ant.Divider` / `VerticalDivider` / `DashedVerticalDivider` | ✅ | 原生 `ui.Divider` 只有实线；虚线/带文字/对齐用 `Box` + `BorderWidth` 组合 |
| `ui.Split` / `SplitPercent` / `SplitVertical` | Splitter | 链式 | — | ◻ | 分隔条颜色对齐 `t.Border` 即可 |
| `ui.Breadcrumbs` | Breadcrumb | 链式 | `ant.Breadcrumb` | ✅ | `chosen=-1` 表示未选择，最后一项为当前页 |
| `ui.MenuButton` | Dropdown / Menu | 组合 | `ant.Dropdown` | ✅ | `MenuItemData` 声明菜单项、分隔线、危险项 |
| `el.ContextMenu(fn)` | Dropdown（右键触发） | 组合 | — | ◻ | 元素方法，不是 `ui.ContextMenu` 函数 |
| `ui.Sidebar` / `SidebarSection` / `SidebarItem` | Menu / Layout.Sider | 链式 | — | ◻ | 选中项 → `t.Selection` |
| `ui.Link` | Typography.Link | 链式 | — | ◻ | 外观即主色文字，可直接用原生 |
| `ui.Toolbar` | Space（紧凑工具条） | 原生 | — | — | MyGo 特有溢出逻辑，保留原生 |
| `ui.Toggle` / `ui.ToggleGroup` | Radio.Button / Segmented | 链式 | — | ◻ | 与 `ant.Segmented` 语义重叠，优先级低 |
| `ui.SplitButton` / `ui.PrimarySplitButton` | Dropdown.Button | 链式 | — | ◻ | 主按钮 + 下拉箭头两段 |
| `ui.PopoverMenuButton` | Dropdown | 组合 | — | ◻ | 可选，`ant.Dropdown` 已覆盖主场景 |

## 二、数据录入

| MyGo v0.10.0 控件 | Ant Design 对应 | 换肤方式 | 本库 API | 状态 | 要点 |
| --- | --- | --- | --- | --- | --- |
| `ui.TextInput` / `ui.TextArea` | Input / Input.TextArea | `TextInputBase` | `ant.Input` / `TextArea` | ✅ | 描边 `t.Border`、聚焦环 `t.Focus`、占位符 `t.TextMuted`、前后缀、密码、字数 |
| `ui.TextInput` + 搜索图标 | Input.Search | 组合 | `ant.SearchInput` | ✅ | 回车提交，`OnSubmit` 回调 |
| `ui.NumberInput` | InputNumber | 组合 | `ant.NumberInputV` | ✅ | v0.10.0 原语收 `*float64` + 范围 + 步进 |
| `ui.DecimalInput` | InputNumber（金额） | 组合 | — | ◻ | 收 `*string`，适合金额/精度场景 |
| `ui.SearchField` | Input.Search | 原生组合 | — | ◻ | 原生自带放大镜/清除/Esc；当前 `ant.SearchInput` 用的是 `TextInput` 组合 |
| `ui.Select` | Select | 链式 | `ant.Select` | ✅ | 高度/字号/圆角对齐 antd |
| `ui.Combobox` | Select（showSearch） | 链式 | `ant.Combobox` | ✅ | 输入过滤 + 候选 |
| `ui.Autocomplete` | AutoComplete | `ComboboxBase` | — | ◻ | 与 Combobox 的差别是保留用户输入 |
| `ui.Checkbox` / `CheckboxGroup` | Checkbox | `CheckboxBase` | `ant.Checkbox` | ✅ | `CheckboxGroup` 容器未单独包装 |
| `ui.CheckboxState` | Checkbox（indeterminate） | `CheckboxStateBase` | `ant.CheckboxAll` | ✅ | 全选/半选三态 |
| `ui.Radio` / `ui.RadioGroup` | Radio / Radio.Group | `RadioBase` | `ant.Radio` / `ant.RadioGroup` | ✅ | 单项 `ant.Radio[T]` 与容器 `RadioGroup` 均已包装 |
| `ui.Switch` | Switch | 链式 | `ant.Switch` / `SwitchWithLabel` | ✅ | 轨道色 `t.Accent` |
| `ui.Slider` / `RangeSlider` | Slider | 链式 | `ant.Slider` / `RangeSliderDouble` | ✅ | 双端区间；`StepSlider` 未包装 |
| `ui.Stepper` | InputNumber（步进器） | 链式 | — | ◻ | ⚠️ 它是数字步进器，**不是** antd Steps，勿与 `ant.Steps` 混淆 |
| `ui.Rating` | Rate | 链式 | `ant.Rate` | ✅ | 选中色 `t.Warning`（antd 星色） |
| `ui.ColorWell` / `ui.ColorPicker` | ColorPicker | 链式 | — | ◻ | 面板圆角/阴影统一后即可用 |
| `ui.Calendar` | Calendar | 组合 | — | ◻ | 日期面板 + 触发器 + 快捷栏；`ant.DatePicker` 已覆盖字段场景 |
| `ui.DateInput` | DatePicker | 组合 | `ant.DatePicker` | ✅ | `*time.Time` 受控 |
| `ui.TimeInput` | TimePicker | 组合 | `ant.TimePicker` | ✅ | 分段输入，12/24h 随区域 |
| `ui.Form` / `Field` / `Fieldset` | Form | 链式 | — | ◻ | 标签对齐、必填星号 `t.Danger`、校验信息色 |
| `ui.TokenField` | Select（tags）/ Mentions 近似 | 组合 | — | ◻ | 标签圆角色块 + 关闭按钮 |
| `ui.EditableText` | Typography（editable） | 链式 | — | ◻ | 原位重命名 |
| `ui.Segmented` | Segmented | 链式 | `ant.Segmented` | ✅ | 支持 `ControlOption`（Disabled/OnChange） |
| `ui.FindBar` | 无直接对应 | 原生 | — | — | 编辑器内查找栏，antd 无同类，保留原生 |

## 三、数据展示

| MyGo v0.10.0 控件 | Ant Design 对应 | 换肤方式 | 本库 API | 状态 | 要点 |
| --- | --- | --- | --- | --- | --- |
| `ui.Table` | Table | 链式 | `ant.Table` | ✅ | 列定义、表头排序、行多选；单元格为字符串比较排序 |
| `ui.List` | List（虚拟列表） | 链式 | — | ◻ | 虚拟滚动由 MyGo 原生提供，换肤价值主要是行 hover/选中色 |
| `ui.Tree` / `TreeItem` | Tree | 链式/组合 | `extra.Tree` | ✅ | 声明式节点树 + 展开状态 |
| `ui.Outline` / `OutlineTable` | Tree / Anchor（近似） | 链式 | — | ◻ | 大树的列表/表格形态 |
| `ui.GridView` | List（grid）/ Image.Preview | 链式 | — | ◻ | 网格视图 |
| `ui.Avatar` | Avatar | 组合 | `ant.AvatarEl` | ✅ | 圆形/姓名首字母；`image` 可选 |
| `ui.Badge` | Badge（计数/红点） | 组合 | `ant.Badge` / `BadgeDot` | ✅ | ⚠️ 原生 `ui.Badge` 是文本胶囊；计数与右上角骑角是自建（`ui.Overlay` + `AttachTo`） |
| `ui.Progress` | Progress（线形） | 链式 | `ant.ProgressLine` | ✅ | 百分比文字可隐藏；环形/仪表盘未做 |
| `ui.Meter` | Progress（仪表/电平） | 链式 | — | ◻ | 按阈值变色 |
| `ui.Image` | Image（圆角/预览） | 链式 | — | ◻ | 预览用 `DialogBase` |
| `ui.Collapsible` | Collapse（单项） | 链式 | — | ◻ | 头部背景 `t.Surface` |
| `ui.Accordion` / `AccordionItem` | Collapse | 链式 | — | ◻ | 手风琴分区 |
| `ui.Tabs` | Tabs | 链式 | `ant.Tabs` | ✅ | 下划线/卡片风格，选中线 `t.Accent` |
| `ui.Stepper`（数字） | — | — | — | — | 见 §二；antd Steps 对应的是 `ant.Steps`（自建缺口控件） |
| `ui.Timeline`（**动画类型**） | — | — | — | — | ⚠️ 这是动效时间轴类型，不是 antd Timeline；antd 时间轴尚未实现 |
| `ui.Popover` | Popover | 组合 | 内部使用 | ◻ | `extra.Cascader` / `extra.TreeSelect` 已用 `ui.Popover`，但未对外提供 `ant.Popover` |
| `ui.TooltipBase` | Tooltip | Base | — | ◻ | Tooltip 是元素方法 `el.Tooltip(...)`，**没有** `ui.Tooltip` 函数 |
| `ui.RichText` | Typography（富文本） | 原生 | — | — | span 级富文本，保留原生 |
| `ui/chart.Plot` 等 | 图表 | 原生 | — | — | 子包 `ui/chart`，保留原生 |
| `ui/editor.View` / `editor.DiffView` | 代码编辑器 / Diff | 原生 | — | — | 子包 `ui/editor`；Gallery 的代码视图后续可换用它 |
| `ui/markdown.View` / `Stream` | Markdown | 原生 | — | — | 子包 `ui/markdown` |
| `ui/richtext.Editor` | RichText | 原生 | — | — | 子包 `ui/richtext` |

## 四、反馈与浮层

| MyGo v0.10.0 控件 | Ant Design 对应 | 换肤方式 | 本库 API | 状态 | 要点 |
| --- | --- | --- | --- | --- | --- |
| `ui.Modal` | Modal | `DialogBase` | `ant.Modal` | ✅ | 遮罩、居中卡片、页脚按钮顺序 |
| `ui.AlertDialog` | Modal.confirm / Popconfirm | `DialogBase` | `ant.Confirm` | ✅ | 返回被点按钮下标；`Confirm` 无返回值，只回调 `onOk` |
| `c.Toast` / `c.AddToast` / `Toast{Type}` | Message / Notification | `ToastViewportBase` | `ant.Message` / `MessageTypeTyped` / `Notify` | ✅ | 复用 MyGo 底部 toast viewport；**不是**顶部居中/右上角独立浮层 |
| `ui.Spinner` | Spin | 链式 | — | ◻ | `ant.Button` 的 Loading 内部已用 `ui.Spinner` |
| `ui.Overlay` + 绝对定位 | Drawer | 组合 | — | ◻ | 从边缘滑入的面板，需 `Modal`/`DialogBase` + `Motion` |
| `el.Sticky()` / `Top()`（元素方法） | Affix | 原生 | — | ◻ | 没有 `ui.Sticky` 函数；粘性由元素方法提供 |
| `ui.DialogBase` | Modal / Drawer 底座 | Base | `ant.Modal` 内部 | ✅ | `DialogBase(c, open, fn(backdrop, panel))` |
| `ui.Alert`（不存在） | Alert | 组合 | `ant.Alert` | ✅ | 四级别 + 可关闭 + 描述，纯组合，不是原生控件 |
| `ui.FindBar` | 无 | 原生 | — | — | 保留 |
| `ui.Drop` / `DataDragOver` / `DragOver` | Upload | 组合 | `extra.Upload` | ✅ | 桌面语义：系统文件对话框 + 原生拖放 + 列表 |
| `Inspector`（调试浮层） | 无 | 原生 | — | — | 非公开控件，不包装 |

## 五、MyGo 特有、Ant Design 无对应（保留原生外观）

以下能力没有 antd 对应物，不做「antd 化」，只保证主题色/圆角/焦点环一致：

- **绘图与动效**：`ui.Canvas`（`NewCanvas`）、`ui.Painter`、`ui.Path`、`el.AnimateTo(...)`、`ui.Motion`、`ui.ElementTransition`、`ui.Keyframe`、`ui.Placement`、`ui.Transition`、`ui.Timeline`（时间轴类型）。
- **编辑器与文档**：`ui/editor`（`View` / `DiffView`）、`ui/markdown`（`View` / `Stream`）、`ui/richtext`（`Editor`）、`ui/chart`（`Plot`）、`ui/shader`。
- **原生视图**：`ui.WebView`、`ui.NativeView`、`ui.Video`。
- **窗口与历史**：`ui.BackButton` / `ui.ForwardButton`（历史按钮）、`ui.TitleBar`、窗口/标题栏 API。
- **工具栏与输入法**：`ui.Toolbar`、`ui.Key` / `ui.KeyCode` / `ui.Modifiers`、IME、拖放、触摸。
- **无障碍**：`ui.AccessInfo`、`ui.Role`、`el.Label(...)`、`el.AccessRole(...)`。
- **测试与记录**：`ui.NewTester`、`ui.NewRecorder`、`ui.Render`。

> 策略：**保持 MyGo 原生实现 + 主题配色同步**，不伪造不存在的 antd 组件。

## 六、换肤统一范式

所有换肤控件遵循同一模板（以 Button 为例，与 `widgets/button.go` 一致）：

```go
// 真实实现的结构示意：状态色集中在 internal.ButtonColors，
// 颜色/间距全部来自 c.Theme()。
func Button(c *ui.Context, label string, opts ...ButtonOption) ui.Element {
	cfg := applyButton(opts)
	t := c.Theme()                       // 唯一颜色来源
	p := internal.ButtonColors(t, cfg.typ, cfg.danger) // 状态集中映射
	el := ui.ButtonBase(c)               // 复用交互底座（焦点/键盘/无障碍）
	el.Height(internal.ControlHeight(t, cfg.size)).
		Padding(0, internal.ControlPaddingX(t, cfg.size)).
		Radius(t.Radius).
		Background(p.Bg).
		Border(p.BorderWidth, p.Border).
		TextColor(p.Fg)
	// 内部监听 Hovered()/Pressed()/Disabled()，按 p 的 hover/pressed 槽位切换
	_ = label
	return el
}
```

要点（与 `AGENTS.md` 2.1 一致）：

1. **禁止字面颜色**：所有颜色经 `t` 取，不写 `ui.Hex("#...")`（`internal/check` 会失败）。
2. **禁止字面间距**：所有 `Padding` / `Gap` / `Margin` / `Size` / `Width` / `Height` 经 `t.Space(n)` / `t.Rem(n)`（0 与 1px 描边豁免）。
3. **状态集中映射**：hover / pressed / disabled 在 `internal/variant.go` 统一派生，不在每个控件里各写分支。
4. **签名贴近原生**：`ant.Button(c, "保存", ant.Primary())`，降低迁移成本；有状态值仍用 `*T`。

## 七、覆盖度小结（2026-11-27）

- MyGo v0.10.0 的 `ui` 包共 **102 个返回 `Element` 的构造函数**（含 `Base` 与布局原语），`docs/ui/` 下 **88 篇** 组件/主题文档。
- 本库当前**已包装 35 个原生原语**，提供 **77 个控件构造函数**：
  - `widgets/`（包 `ant`）70 个：Button / Typography（4）/ Divider（3）/ Space·Flex·Grid / Input·TextArea·SearchInput·NumberInputV / Checkbox·CheckboxAll·RadioGroup·Switch·SwitchWithLabel·Slider·RangeSliderDouble·Rate·Segmented / Select·Combobox / DatePicker·TimePicker / Table / Tag·Badge·BadgeDot·AvatarEl·Card / Empty·Result·Skeleton·Statistic·Descriptions·Timeline·List / Breadcrumb·Tabs·Dropdown·Pagination·Steps·Sidebar·SidebarItem·SidebarItemWithIcon·SidebarSection / Alert·Modal·Message·MessageTypeTyped·Notify·ProgressLine·Confirm·Spin·Tooltip·Popconfirm·Collapse·Form·FormItem·FormGroup·Splitter·SplitterVertical·SplitterPercent·SplitterVerticalPercent·AutoComplete；
  - `extra/` 7 个：Tree / Upload / Transfer / Carousel / QRCode / Cascader / TreeSelect；
  - 另有 `icon.Icon` + 44 个图标常量。
- 仍**未覆盖**的原生控件（约 60 个）与 antd 缺口控件（约 25 个）统一列在 [03 缺口控件](03-gap-widgets.md)，按优先级推进。
- 早期草案中的「55 个可直接换肤、12 个需 Base、15 个特有」是**估算**，不是当前完成度；以本节计数为准。

## 八、勘误表（早期文档 → 正确 API）

| 早期写法 | 正确写法（MyGo v0.10.0） | 说明 |
| --- | --- | --- |
| `ui.Dialog` | `ui.Modal` / `ui.DialogBase` | 对话框原语叫 Modal |
| `ui.Input` | `ui.TextInput` | 单行输入 |
| `ui.Tooltip` | `el.Tooltip(...)` / `ui.TooltipBase` | Tooltip 是元素方法 + Base |
| `ui.ContextMenu` | `el.ContextMenu(fn)` / `ui.ContextMenuBase` | 元素方法 + Base |
| `ui.Sticky` | `el.Sticky().Top(0)` | 元素方法 |
| `ui.Inspector` | —（非公开控件） | 调试浮层，不包装 |
| `ui.Navigation` | `ui.Router` / `ui.Sidebar` | 导航由路由 + 侧栏组成 |
| `ui.DragAndDrop` | `ui.Drop` / `ui.DragOver` / `ui.DataDragOver` | 拖放 API |
| `ui.Chart` | `ui/chart.Plot` 等 | 子包 |
| `ui.CodeEditor` | `ui/editor.View` | 子包 |
| `ui.DiffViewer` | `ui/editor.DiffView` | 子包 |
| `ui.RichTextEditor` | `ui/richtext.Editor` | 子包 |
| `ui.Markdown` | `ui/markdown.View` / `Stream` | 子包 |
| `ui.Drawing` / `ui.Painter` | `ui.Canvas`（`ui.NewCanvas`）/ `ui.Painter` | 绘图原语 |
| `ui.Shader` | `ui/shader` + `ui.ShaderInput` | 子包 + 类型 |
| `ui.Windows` | 应用层 `mygo.Window` / `ui.TitleBar` | 窗口不是控件 |
| `ui.HistoryButtons` | `ui.BackButton` / `ui.ForwardButton` | 历史按钮 |
| `ui.Keyboard` | `ui.Key` / `ui.KeyCode` / `ui.Modifiers` | 键盘类型 |
| `ui.Accessibility` | `ui.AccessInfo` / `ui.Role` / `el.Label(...)` | 无障碍 API |
| `extra.Pagination` | `ant.Pagination` | 已实现于 widgets |
| `ant.Radio[T]` | `ant.Radio(c, *T, value, label, opts...)` | 泛型包装已实现（此前误记为未暴露） |
