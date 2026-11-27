# 02 · MyGo 原生控件 → Ant Design 皮肤映射

> 对应需求 3：把 MyGo 自带控件加上 Ant Design 设计。
>
> **核对基线**：MyGo **v0.10.0**（`go.mod` 锁定）、Go 1.27.4、本库当前代码（v0.1.0 已发布，v0.3 开发中），核对日期 2026-11-27。
> 本表只把 `package ui` **实际存在的导出名**作为映射来源；早期草案里写错的 `ui.Dialog`、`ui.Input`、`ui.Tooltip`、`ui.Chart` 等名字已全部修正（见 §八 勘误表）。
>
> **命名规则**：MyGo 与 antd 有同类型控件时，**以 MyGo 命名为准**（如 `Breadcrumbs`、`Avatar`、`NumberInput`、`Progress`、`Rating`、`Spinner`、`RangeSlider`）；早期 antd 名保留为别名。**全量清单以 §九 覆盖审计为准**，§一~§八 为历史条目，若有冲突以 §九 为准。

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
| `ui.Split` / `SplitPercent` / `SplitVertical` / `SplitVerticalPercent` | Splitter | 链式 | `ant.Splitter` 系列 | ✅ | 分隔条可拖拽，尺寸/百分比两种模式 |
| `ui.Breadcrumbs` | Breadcrumb | 链式 | `ant.Breadcrumb` | ✅ | `chosen=-1` 表示未选择，最后一项为当前页 |
| `ui.MenuButton` | Dropdown / Menu | 组合 | `ant.Dropdown` | ✅ | `MenuItemData` 声明菜单项、分隔线、危险项 |
| `el.ContextMenu(fn)` | Dropdown（右键触发） | 组合 | `ant.ContextMenu` | ✅ | 右键/菜单键触发，菜单项可设危险色 |
| `ui.Sidebar` / `SidebarSection` / `SidebarItem` | Menu / Layout.Sider | 链式 | `ant.Sidebar` / `SidebarItem` / `SidebarItemWithIcon` / `SidebarSection` | ✅ | 选中项 `t.Selection`，键盘可导航 |
| `ui.Link` | Typography.Link | 链式 | `ant.Hyperlink` | ✅ | 主色文字（`Link()` 已被按钮类型占用，故命名 Hyperlink） |
| `ui.Toolbar` | Space（紧凑工具条） | 链式 | `ant.Toolbar` | ✅ | 保留 MyGo 溢出/键盘逻辑，套主题面与间距 |
| `ui.Toggle` / `ui.ToggleGroup` | Radio.Button / Segmented | 链式 | `ant.Toggle` / `ToggleGroup` | ✅ | 保持按下的按钮 |
| `ui.SplitButton` / `ui.PrimarySplitButton` | Dropdown.Button | 链式 | `ant.SplitButton` / `PrimarySplitButton` | ✅ | 返回 MyGo `SplitButtonParts`，箭头挂菜单 |
| `ui.PopoverMenuButton` | Dropdown | 链式 | `ant.PopoverMenuButton` | ✅ | MyGo 命名入口；`ant.Dropdown` 为声明式封装 |

## 二、数据录入

| MyGo v0.10.0 控件 | Ant Design 对应 | 换肤方式 | 本库 API | 状态 | 要点 |
| --- | --- | --- | --- | --- | --- |
| `ui.TextInput` / `ui.TextArea` | Input / Input.TextArea | `TextInputBase` | `ant.Input` / `TextArea` | ✅ | 描边 `t.Border`、聚焦环 `t.Focus`、占位符 `t.TextMuted`、前后缀、密码、字数 |
| `ui.TextInput` + 搜索图标 | Input.Search | 组合 | `ant.SearchInput` | ✅ | 回车提交，`OnSubmit` 回调 |
| `ui.NumberInput` | InputNumber | 组合 | `ant.NumberInput`（别名 `NumberInputV`） | ✅ | 收 `*float64` + 范围 + 步进；MyGo 命名优先 |
| `ui.DecimalInput` | InputNumber（金额） | 组合 | `ant.DecimalInput` | ✅ | 收 `*string`，金额/精度场景不被浮点舍入 |
| `ui.SearchField` | Input.Search | 原生组合 | `ant.SearchField` | ✅ | 原生放大镜/清除/Esc；`ant.SearchInput` 是换肤版 |
| `ui.Select` | Select | 链式 | `ant.Select` | ✅ | 高度/字号/圆角对齐 antd |
| `ui.Combobox` | Select（showSearch） | 链式 | `ant.Combobox` | ✅ | 输入过滤 + 候选 |
| `ui.Autocomplete` | AutoComplete | 链式 | `ant.AutoComplete` | ✅ | 保留用户输入 + 候选 |
| `ui.Checkbox` / `CheckboxGroup` | Checkbox | `CheckboxBase` | `ant.Checkbox` / `CheckboxGroup` | ✅ | 8px 标签间距；`CheckboxGroup` 已按 MyGo 命名补齐 |
| `ui.CheckboxState` | Checkbox（indeterminate） | `CheckboxStateBase` | `ant.CheckboxAll` | ✅ | 全选/半选三态 |
| `ui.Radio` / `ui.RadioGroup` | Radio / Radio.Group | `RadioBase` | `ant.Radio` / `ant.RadioGroup` | ✅ | 单项 `ant.Radio[T]` 与容器 `RadioGroup` 均已包装 |
| `ui.Switch` | Switch | 链式 | `ant.Switch` / `SwitchWithLabel` | ✅ | 轨道色 `t.Accent` |
| `ui.Slider` / `RangeSlider` | Slider | 链式 | `ant.Slider` / `RangeSlider`（别名 `RangeSliderDouble`）/ `StepSlider` | ✅ | 双端区间 + 带刻度；MyGo 命名优先 |
| `ui.Stepper` | InputNumber（步进器） | 链式 | `ant.Stepper` | ✅ | ⚠️ 只有上下箭头（AppKit 风格），数值请配 `ant.NumberInputV`；**不是** antd Steps |
| `ui.Rating` | Rate | 链式 | `ant.Rate` | ✅ | 选中色 `t.Warning`（antd 星色） |
| `ui.ColorWell` / `ui.ColorPicker` | ColorPicker | 链式 | `ant.ColorPicker` / `ColorPanel` | ✅ | 色块展开面板 / 常驻面板 |
| `ui.Calendar` | Calendar | 组合 | `ant.Calendar` | ✅ | 月历浏览/选择；`ant.DatePicker` 覆盖字段场景 |
| `ui.DateInput` | DatePicker | 组合 | `ant.DatePicker` | ✅ | `*time.Time` 受控 |
| `ui.TimeInput` | TimePicker | 组合 | `ant.TimePicker` | ✅ | 分段输入，12/24h 随区域 |
| `ui.Form` / `Field` / `Fieldset` | Form | 链式 | `ant.Form` / `FormItem` / `FormGroup` | ✅ | 标签对齐、必填标记、帮助/错误文字 |
| `ui.TokenField` | Select（tags）/ Mentions 近似 | 组合 | `ant.TokenField` | ✅ | 标签输入 + 候选建议 |
| `ui.EditableText` | Typography（editable） | 链式 | `ant.EditableText` | ✅ | 原位重命名 |
| `ui.Segmented` | Segmented | 链式 | `ant.Segmented` | ✅ | 支持 `ControlOption`（Disabled/OnChange） |
| `ui.FindBar` | 无直接对应 | 原生 | — | — | 编辑器内查找栏，antd 无同类，保留原生 |

## 三、数据展示

| MyGo v0.10.0 控件 | Ant Design 对应 | 换肤方式 | 本库 API | 状态 | 要点 |
| --- | --- | --- | --- | --- | --- |
| `ui.Table` | Table | 链式 | `ant.Table` | ✅ | 列定义、表头排序、行多选；单元格为字符串比较排序 |
| `ui.List` | List（虚拟列表） | 链式 | `ant.List` | ✅ | 虚拟滚动 + 行分隔线 + 选中/多选 |
| `ui.Tree` / `TreeItem` | Tree | 链式/组合 | `extra.Tree` | ✅ | 声明式节点树 + 展开状态 |
| `ui.Outline` / `OutlineTable` | Tree / Anchor（近似） | 原生 | — | — | 高级原生用法，不包装；ant.Anchor 覆盖目录场景 |
| `ui.GridView` | List（grid）/ Image.Preview | 链式 | `ant.GridView` | ✅ | 等宽网格 + 虚拟滚动，间距随主题 |
| `ui.Avatar` | Avatar | 组合 | `ant.Avatar`（别名 `AvatarEl`） | ✅ | 圆形/姓名首字母；`image` 可选；MyGo 命名优先 |
| `ui.Badge` | Badge（计数/红点） | 组合 | `ant.Badge` / `BadgeDot` | ✅ | ⚠️ 原生 `ui.Badge` 是文本胶囊；计数与右上角骑角是自建（`ui.Overlay` + `AttachTo`） |
| `ui.Progress` | Progress（线形） | 链式 | `ant.Progress`（别名 `ProgressLine`）+ `ProgressCircle` | ✅ | 百分比文字可隐藏；环形/仪表盘为 antd 补充 |
| `ui.Meter` | Progress（仪表/电平） | 链式 | `ant.Meter` | ✅ | 按阈值变色（`ui.MeterLevels`） |
| `ui.Image` | Image（圆角/预览） | 链式 | `ant.Image`（别名 `ImageEl`） | ✅ | 宽度可调 + 点击模态预览；MyGo 命名优先 |
| `ui.Collapsible` | Collapse（单项） | 链式 | `ant.Collapsible` | ✅ | 原生展开/收起动画 |
| `ui.Accordion` / `AccordionItem` | Collapse | 链式 | `ant.Accordion` / `AccordionItem` | ✅ | 手风琴分区 |
| `ui.Tabs` | Tabs | 链式 | `ant.Tabs` | ✅ | 下划线/卡片风格，选中线 `t.Accent` |
| `ui.Stepper`（数字） | — | — | — | — | 见 §二；antd Steps 对应的是 `ant.Steps`（自建缺口控件） |
| `ui.Timeline`（**动画类型**） | — | — | — | — | ⚠️ 这是动效时间轴类型，不是 antd Timeline；antd 时间轴尚未实现 |
| `ui.Popover` | Popover | 组合 | `ant.Popover` / `ant.Popconfirm` | ✅ | 通用锚定气泡 + 气泡确认 |

| `ui.TooltipBase`（元素方法 `el.Tooltip(...)`） | Tooltip | Base | `ant.Tooltip` | ✅ | 反色气泡 + 悬停/聚焦触发 |
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
| `ui.Spinner` | Spin | 链式 | `ant.Spin` | ✅ | 指示器 + 提示文字 |
| `ui.Overlay` + 绝对定位 | Drawer | 组合 | `ant.Drawer` | ✅ | `DialogBase` 右侧锚定 + 遮罩/标题栏 |
| `el.Sticky()` / `Top()`（元素方法） | Affix | 原生 | `ant.Affix` | ✅ | 吸顶包装 |
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
- 本库当前**已包装绝大多数原生控件与 Base 原语**，提供 **144 个控件构造函数**：
  - `widgets/`（包 `ant`）99 个：Button / Typography（4）/ Divider（3）/ Space·Flex·Grid / Input·TextArea·SearchInput·NumberInputV / Checkbox·CheckboxAll·RadioGroup·Switch·SwitchWithLabel·Slider·RangeSliderDouble·Rate·Segmented / Select·Combobox / DatePicker·TimePicker / Table / Tag·Badge·BadgeDot·AvatarEl·Card / Empty·Result·Skeleton·Statistic·Descriptions·Timeline·List·Drawer·ImageEl·ColorPicker·ColorPanel / Breadcrumb·Tabs·Dropdown·Pagination·Steps·Sidebar·SidebarItem·SidebarItemWithIcon·SidebarSection / Alert·Modal·Message·MessageTypeTyped·Notify·ProgressLine·Confirm·Spin·Tooltip·Popconfirm·Collapse·Form·FormItem·FormGroup·Splitter·SplitterVertical·SplitterPercent·SplitterVerticalPercent·AutoComplete·ProgressCircle·Toggle·ToggleGroup·Affix·Hyperlink·DecimalInput·SearchField·Stepper·Calendar·TokenField·EditableText·Meter·Collapsible·Accordion·AccordionItem·Popover·ContextMenu / Watermark·Mentions·Tour·Notification·FloatButton·Anchor·BorderBeam·Masonry；
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

## 九、全量覆盖审计（MyGo v0.10.0 · `package ui`）

> 起因：review 发现「文档只覆盖常用控件，MyGo `package ui` 实际导出 111 个控件构造函数（含 Base 底座与布局原语），存在遗漏」。
> 本节以 `go doc`/源码导出名为准逐项核对，作为后续开发的**唯一覆盖清单**；新增/修改控件时必须同步本表状态。

### 9.1 命名约定（本项目强制）

1. **同类型控件以 MyGo 命名为准**：如 `ui.Breadcrumbs` → `ant.Breadcrumbs`、`ui.Avatar` → `ant.Avatar`、`ui.NumberInput` → `ant.NumberInput`、`ui.Progress` → `ant.Progress`、`ui.Rating` → `ant.Rating`、`ui.Spinner` → `ant.Spinner`、`ui.RangeSlider` → `ant.RangeSlider`。
2. 早期按 antd 命名的入口**保留为别名**（`Breadcrumb`、`AvatarEl`、`NumberInputV`、`ProgressLine`、`Rate`、`Spin`、`RangeSliderDouble`…），不破坏既有代码；新代码请用 MyGo 名。
3. 仅 antd 有、MyGo 没有的缺口控件（`Tabs` 以外的 `Collapse`、`Carousel`、`QRCode`、`Transfer`…）保持 antd 命名（见 `extra/`）。
4. `XxxBase` 是「无外观底座」，本项目一般不重复包装；需要自绘时直接使用 `ui.XxxBase` 并套主题（见 🧱 行）。

### 9.2 覆盖矩阵

图例：✅ 已包装（同名或别名）｜🧱 无外观底座/布局原语，不单独包装｜⏳ 已列入计划待补｜➖ MyGo/平台专属，仅保留原生。

#### 通用与基础

| MyGo 构造函数 | 本库入口 | 状态 | 说明 |
| --- | --- | --- | --- |
| `Button` | `ant.Button` | ✅ | 五型 × danger × 三档尺寸 |
| `PrimaryButton` | `ant.PrimaryButton` | ✅ | 本轮补齐，等价 `Button(..., Primary())` |
| `Text` / `Textf` | `ant.Title` / `Paragraph` / `Code` / `TextMuted` | ✅ | 排版族 |
| `RichText` | `ant.RichText` | ✅ | 本轮补齐，统一字号/正文色 |
| `Icon` | `icon.Icon` / `ant.Icon` | ✅ | 44 个 antd 图标 + 原生 SVG |
| `Link` | `ant.Hyperlink` | ✅ | `Link()` 被按钮类型占用，故别名 Hyperlink |
| `Divider` | `ant.Divider` / `VerticalDivider` | ✅ | 含虚线/带文字变体 |
| `Row` / `Column` / `Grid` / `Box` / `Spacer` | `ant.Space` / `Flex` / `Grid` | ✅ | 布局原语，按 antd 语义包装 |
| `Scroll` / `ScrollBoth` / `ScrollHorizontal` | `ant.Scroll` / `ScrollBoth` / `ScrollHorizontal` | ✅ | 本轮补齐，直接透传 |
| `Split` / `SplitPercent` / `SplitVertical` / `SplitVerticalPercent` | `ant.Split` / `SplitPercent` / `SplitVertical` / `SplitVerticalPercent` | ✅ | 本轮补齐；另有 antd 风格 `Splitter` 系列 |
| `Toolbar` | `ant.Toolbar` | ✅ | 本轮补齐：主题面 + 间距 |
| `Overlay` / `Themed` | — | 🧱 | 渲染/主题底座 |
| `Breadcrumbs` | `ant.Breadcrumbs`（别名 `Breadcrumb`） | ✅ | 本轮补齐 MyGo 命名 |
| `MenuButton` | `ant.MenuButton`（另有 `ant.Dropdown`） | ✅ | 本轮补齐；`Dropdown` 为声明式封装 |
| `PopoverMenuButton` | `ant.PopoverMenuButton` | ✅ | 本轮补齐 |
| `SplitButton` / `PrimarySplitButton` | `ant.SplitButton` / `PrimarySplitButton` | ✅ | 本轮补齐（返回 `SplitButtonParts`，对齐 Dropdown.Button） |
| `ContextMenuBase` / `MenuBase` / `PopoverBase` / `TooltipBase` / `DialogBase` | — | 🧱 | 自绘底座，直接用 `ui.*` + 主题 |
| `BackButton` / `ForwardButton` | `ant.BackButton` / `ForwardButton` | ✅ | 本轮补齐，配合 `ui.Router` |
| `AlertDialog` | `ant.AlertDialog` | ✅ | 本轮补齐；返回按钮下标，antd 场景可改用 `ant.Modal` |
| `FindBar` | `ant.FindBar` | ✅ | 本轮补齐，编辑器查找栏 |
| `Toggle` / `ToggleGroup` / `ToggleBase` | `ant.Toggle` / `ToggleGroup` / `ToggleBase` | ✅ | 保持按下的按钮 |

#### 数据录入

| MyGo 构造函数 | 本库入口 | 状态 | 说明 |
| --- | --- | --- | --- |
| `TextInput` / `TextArea` | `ant.Input` / `TextArea` | ✅ | 自绘外框，聚焦环 1px 主色 + 2px 阴影 |
| `TextInputBase` / `TextAreaBase` / `TextAreaBufferBase` / `TextInputBufferBase` | — | 🧱 | 无外观底座 |
| `TextInputBuffer` / `TextAreaBuffer` | `ant.TextInputBuffer` / `TextAreaBuffer` | ✅ | 本轮补齐：复用 Input/TextArea 外框与焦点环，内部索引化编辑 |
| `NumberInput` | `ant.NumberInput`（别名 `NumberInputV`） | ✅ | 本轮补齐 MyGo 命名 |
| `DecimalInput` | `ant.DecimalInput` | ✅ | 金额/精度场景 |
| `SearchField` | `ant.SearchField` | ✅ | 原生放大镜/清除/Esc |
| `Select` | `ant.Select` | ✅ | `SelectBase` 自绘：32px/11px 内边距/箭头/32px 选项 |
| `Combobox` | `ant.Combobox` | ✅ | `ComboboxBase` 自绘输入框 + 候选面板 |
| `Autocomplete` | `ant.Autocomplete`（别名 `AutoComplete`） | ✅ | 本轮补齐 MyGo 命名；透传原生交互 |
| `Checkbox` | `ant.Checkbox` | ✅ | 8px 间距 + 标签可点 |
| `CheckboxGroup` | `ant.CheckboxGroup` | ✅ | 本轮补齐 |
| `CheckboxState` | `ant.CheckboxState`（别名 `CheckboxAll`） | ✅ | 本轮补齐三态复选框 |
| `CheckboxBase` / `CheckboxStateBase` | — | 🧱 | 无外观底座 |
| `RadioGroup` / `Radio` | `ant.RadioGroup` / `Radio[T]` | ✅ | 泛型受控值 |
| `Switch` / `SwitchBase` | `ant.Switch` / `SwitchBase` | ✅/🧱 | 自绘 44×22 轨道 |
| `Slider` / `SliderBase` | `ant.Slider` / `SliderBase` | ✅/🧱 | 主题化滑块 |
| `StepSlider` | `ant.StepSlider` | ✅ | 本轮补齐（带刻度） |
| `RangeSlider` | `ant.RangeSlider`（别名 `RangeSliderDouble`） | ✅ | 本轮补齐 MyGo 命名 |
| `Stepper` | `ant.Stepper` | ✅ | ⚠️ 上下箭头步进器，不是 antd Steps |
| `Rating` | `ant.Rating`（别名 `Rate`） | ✅ | 本轮补齐 MyGo 命名 |
| `ColorWell` / `ColorPicker` | `ant.ColorWell` / `ColorPicker` | ✅ | 本轮补齐 ColorWell |
| `Calendar` | `ant.Calendar` | ✅ | 月历浏览/选择 |
| `DateInput` / `TimeInput` | `ant.DateInput` / `TimeInput`（别名 `DatePicker` / `TimePicker`） | ✅ | 本轮补齐 MyGo 命名 |
| `Form` / `Field` / `Fieldset` | `ant.Form` / `FormItem` / `FormGroup` | ✅ | 纵向标签 + 必填星号 |
| `TokenField` | `ant.TokenField` | ✅ | 标签输入 + 候选 |
| `EditableText` | `ant.EditableText` | ✅ | 原位重命名 |
| `Segmented` / `SegmentedBase` | `ant.Segmented` / `SegmentedBase` | ✅/🧱 | 主题化分段 |
| `Tabs` / `TabsBase` | `ant.Tabs` / `TabsBase` | ✅/🧱 | 自绘 antd 下划线页签 |

#### 数据展示

| MyGo 构造函数 | 本库入口 | 状态 | 说明 |
| --- | --- | --- | --- |
| `Table` | `ant.Table` | ✅ | 自绘 antd 表格 |
| `List` | `ant.List` | ✅ | 虚拟滚动 + 行状态 |
| `Tree` / `TreeItem` / `TreeBase` | `extra.Tree` / 🧱 | ✅ | 声明式节点树 |
| `Accordion` / `AccordionItem` / `AccordionBase` | `ant.Accordion` 系列 | ✅/🧱 | 折叠面板 |
| `Collapsible` / `CollapsibleBase` | `ant.Collapsible` 系列 | ✅/🧱 | 单项折叠 |
| `GridView` | `ant.GridView` | ✅ | 本轮补齐：等宽网格 + 虚拟滚动 |
| `Avatar` | `ant.Avatar`（别名 `AvatarEl`） | ✅ | 本轮补齐 MyGo 命名 |
| `Badge` | `ant.Badge` / `BadgeDot`（另有 `ant.Tag`） | ✅ | ⚠️ 原生 `ui.Badge` 是文本胶囊；计数骑角为自建 |
| `Progress` | `ant.Progress`（别名 `ProgressLine`）+ `ProgressCircle` | ✅ | 本轮补齐 MyGo 命名；环形为 antd 补充 |
| `Meter` | `ant.Meter` | ✅ | 阈值变色电平条 |
| `Image` | `ant.Image`（别名 `ImageEl`） | ✅ | 本轮补齐 MyGo 命名 |
| `RichText` | `ant.RichText` | ✅ | 见通用 |
| `Icon` | `ant.Icon` / `icon.Icon` | ✅ | 见通用 |
| `Video` / `WebView` / `NativeView` | `ant.Video` / `WebView` / `NativeView` | ✅ | 本轮补齐；按主题圆角裁切 |
| `Sidebar` / `SidebarItem` / `SidebarSection` | `ant.Sidebar` 系列 | ✅ | 菜单导航 |
| `ToastBase` / `ToastViewportBase` | — | 🧱 | 由 `ant.Message` / `Notify` 覆盖 |
| `Spinner` | `ant.Spinner`（别名 `Spin`） | ✅ | 本轮补齐 MyGo 命名 |

#### 平台/子包范围

| 项 | 状态 | 说明 |
| --- | --- | --- |
| `ui/markdown` | ✅ | `ant.Markdown`（`MarkdownOptions` 别名） |
| `ui/editor` | ✅ | `ant.CodeEditor` / `DiffView`（`CodeEditorState` / `EditorDiffState` 别名） |
| `ui/chart` | ✅ | `ant.ChartPlot`（`ChartAxes` / `ChartTick` 别名） |
| `ui/richtext` | ✅ | `ant.RichTextEditor`（`RichTextDocument` / `RichTextEditorOptions` 别名） |
| `ui/netimage` | ✅ | `ant.RemoteImage`（`RemoteImageOptions` 别名） |
| `ui/video` | ✅ | `ant.OpenVideo` + `ant.Video` |
| `ui/mobile` | ➖ | 移动端专用（BottomBar/Drawer/Sheet/PageView 等），桌面库不包装 |
| `ui/shader` | ➖ | 自定义 GPU 着色器绘制，非控件；需要时直接用 `ui/shader` |
| `ui.Inspector`、`ui.Outline`、`ui.Drawing`、`ui.Diff` 等（文档主题/平台能力） | ➖ | MyGo 平台专属或调试工具，antd 无对应，保留原生 |

### 9.3 已落地的补齐批次

- 第九轮：新增 `widgets/mygo_names.go`：`Breadcrumbs`、`Avatar`、`Autocomplete`、`Image`、`NumberInput`、`Progress`、`RangeSlider`、`Rating`、`Spinner`、`PrimaryButton`、`CheckboxState`、`CheckboxGroup`、`ColorWell`、`Toolbar`、`GridView`、`FindBar`、`BackButton`、`ForwardButton`、`RichText`、`Icon`、`Scroll*`、`Split*`、`Video`/`WebView`/`NativeView`、`SplitButton`、`PrimarySplitButton`、`MenuButton`、`PopoverMenuButton`、`AlertDialog`、`TimeInput`/`DateInput`、`ToggleBase`、`StepSlider`。
- 第十轮：`TextInputBuffer` / `TextAreaBuffer`（复用 Input/TextArea 外框）；新增 `widgets/native_packages.go` 适配 `ui/markdown`、`ui/editor`、`ui/chart`、`ui/richtext`、`ui/netimage`、`ui/video`；`samples/native` 增加 Markdown 示例（共 6 个）。
- 测试：`widgets/mygo_names_test.go`、`widgets/native_packages_test.go` 覆盖别名渲染、SplitButton 点击、GridView、CheckboxGroup、TextBuffer、Markdown、ChartPlot。
- 根包 `package ui` 已无 ⏳ 项；`XxxBase` 按设计以 🧱 直接使用 `ui.*Base`；子包覆盖状态见 9.2 平台/子包表。
