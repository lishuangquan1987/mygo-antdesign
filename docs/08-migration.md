# 08 · 迁移指南：从 MyGo 原生到 mygo-antdesign

本指南面向已使用 MyGo `package ui` 的应用，逐步迁移到 Ant Design 风格的
`mygo-antdesign`。控件是 `ui` 原语的换肤包装或组合，即时模式语义
（view = f(state)、`*T` 受控绑定、`ui.Local`）完全一致。

> 基线：MyGo **v0.10.0**、本库 **v0.1.0**（v0.2 开发中）、Go 1.27.4。
> 文中 API 名均已与代码核对；早期版本里写错的 `ant.Radio[T]`、`ui.Dialog` 等已修正。

## 1. 安装与主题接入

```go
import (
    "github.com/egoist/mygo/ui"

    "github.com/lishuangquan1987/mygo-antdesign/theme"
)

func (app *App) view(c *ui.Context) {
    theme.Apply(c, theme.Light())            // 或 theme.Dark()；token 可运行时替换
    // ... 之后构建的控件全部拿到 antd 主题
}
```

- **全局切换**：`theme.Apply(c, token)` 每帧调用安全，`token` 换成
  `theme.Dark(theme.WithPrimary(ui.Hex("#7c3aed")))` 即整窗变色。
- **局部切换**：`theme.Sub(c, token)` 返回 `ui.Themed` 容器，适合亮暗并排预览、独立主题面板。
- **紧凑密度（上位机 / IDE）**：`theme.WithDensity(theme.Compact)`。
- **主题 API 一览**：
  `theme.Light(opts...)` / `theme.Dark(opts...)`、`theme.Apply(c, *AntToken)`、`theme.Sub(c, *AntToken) ui.Element`；
  Option：`WithPrimary(ui.Color)`、`WithDensity(theme.Density)`、`WithFontSize(float32)`、`WithRadius(float32)`、`WithSpacing(float32)`。
  `WithDensity` 会重置度量，应放在细调 Option **之前**。

## 2. 控件对照表

### 通用组件

| MyGo 原生 | antd 风格 | 说明 |
| --- | --- | --- |
| `ui.Button` / `ui.PrimaryButton` / `ui.ButtonBase` | `ant.Button(c, label, opts...)` | `ant.Primary()` / `Danger()` / `Dashed()` / `Text()` / `Link()` / `Size(ant.Large)` / `Loading(true)` / `Disabled(true)` / `Block(true)` / `OnClick(fn)` |
| `ui.Text` / `ui.Textf` | `ant.Title(c, s, ant.H1)` / `ant.Paragraph(c, s, opts...)` / `ant.Code(c, s)` / `ant.TextMuted(c, s)` | `Paragraph` 选项：`Secondary()` / `Delete()` / `Strong()` / `Underline()` / `Italic()` / `Marked()` |
| `ui.Row` / `ui.Column` / `ui.Grid` | `ant.Space(c, opts...)` / `ant.Flex(c, opts...)` / `ant.Grid(c, cols)` | 间距经 `t.Space` 随密度缩放；`Space` 选项：`Gap` / `GapUnits` / `SpaceVertical` / `SpaceAlign` / `SpaceFillWidth`；`Flex`：`FlexColumn` / `FlexJustify` / `FlexAlign` / `FlexGap` / `FlexWrap` / `FlexGrow` / `FlexFill` |
| `ui.Divider` | `ant.Divider(c, opts...)` / `ant.VerticalDivider(c)` / `ant.DashedVerticalDivider(c)` | 选项：`DashedDivider()` / `WithText(s)` / `TextLeft()` / `TextCenter()` / `TextRight()` |
| `ui.Icon`（`ui.ParseSVG` + `ui.Image`） | `icon.Icon(c, icon.Search, 16, icon.WithColor(t.Accent))` | 44 个 outlined 常量；默认跟随文本色 |

### 数据录入

| MyGo 原生 | antd 风格 | 说明 |
| --- | --- | --- |
| `ui.TextInput` / `ui.TextArea` | `ant.Input(c, *string, opts...)` / `ant.TextArea(c, *string, opts...)` | `PlaceholderText` / `InputLabel` / `PasswordInput` / `InputDisabled` / `InputReadOnly` / `MaxLengthText` / `Rows` / `InputPrefix` / `InputSuffix` / `InputSize` / `OnChanged` / `OnSubmit` |
| `ui.TextInput` 组合 | `ant.SearchInput(c, *string, opts...)` | 回车提交（antd Input.Search） |
| `ui.NumberInput` | `ant.NumberInputV(c, *float64, lo, hi, step, opts...)` | v0.10.0 原语要求 `*float64` |
| `ui.Select` | `ant.Select(c, *string, options []string, opts ...ant.ControlOption)` | 高度/字号/圆角对齐 antd |
| `ui.Combobox` | `ant.Combobox(c, *string, options []string, opts ...)` | 输入过滤 + 候选 |
| `ui.Autocomplete` | `ant.AutoComplete(c, *string, suggestions []string, opts ...)` | 保留键入内容 + 候选建议 |
| `ui.DateInput` / `ui.TimeInput` | `ant.DatePicker(c, *time.Time, opts...)` / `ant.TimePicker(c, *time.Time, opts...)` | `*time.Time` 受控 |
| `ui.Checkbox` / `ui.CheckboxState` | `ant.Checkbox(c, *bool, label, opts...)` / `ant.CheckboxAll(c, *ui.CheckState, label, opts...)` | 三态全选框 |
| `ui.RadioGroup` / `ui.Radio[T]` | `ant.RadioGroup(c, func())` / `ant.Radio(c, *T, value, label, opts...)` | 泛型受控值，选项共用 `ControlOption` |
| `ui.Switch` | `ant.Switch(c, *bool, opts...)` / `ant.SwitchWithLabel(c, *bool, label, opts...)` | — |
| `ui.Slider` / `ui.RangeSlider` | `ant.Slider(c, *float64, lo, hi, opts...)` / `ant.RangeSliderDouble(c, low, high *float64, lo, hi, step, opts...)` | — |
| `ui.Rating` | `ant.Rate(c, *int, max int, opts...)` | — |
| `ui.Segmented` | `ant.Segmented(c, *int, labels []string, opts ...ant.ControlOption)` | 支持 `ControlDisabled` / `OnChange` |
| `ui.Toggle` / `ui.ToggleGroup` | `ant.Toggle(c, *bool, label, opts...)` / `ant.ToggleGroup(c, fn)` | 保持按下的按钮，共用 `ControlOption` |
| `ui.ColorWell` / `ui.ColorPicker` | `ant.ColorPicker(c, *ui.Color, opts...)` / `ant.ColorPanel(c, *ui.Color, opts...)` | 色块展开面板 / 常驻面板 |
| `ui.Calendar` / `ui.Form` / `ui.TokenField` | — | 未包装，直接用原生 |

> `ant.ControlOption` 提供 `ControlDisabled(bool)` 与 `OnChange(func())`，供 Checkbox/Switch/Slider/Rate/Segmented/Select/Combobox 等共用。

### 数据展示

| MyGo 原生 | antd 风格 | 说明 |
| --- | --- | --- |
| `ui.Table` | `ant.Table(c, columns []ant.TableColumn, rows [][]string, opts...)` | `TableSelection(*ui.Selection[int])` 行多选、`TableOnSort(func())` 表头排序 |
| `ui.Tree` / `TreeItem` | `extra.Tree(c, nodes []extra.TreeNode, opts...)` | `extra.OnTreeNodeSelect(func(key string))` |
| `ui.Badge` 组合 | `ant.Badge(c, count int, fn func())` / `ant.BadgeDot(c, fn func())` | 右上角骑角定位（`ui.Overlay` + `AttachTo`） |
| `ui.Avatar` | `ant.AvatarEl(c, name string, image ...*ui.Bitmap)` | 未传图片时显示姓名首字母 |
| `Row` / `Column` 组合 | `ant.Card(c, title string, body func(), opts...)` | 选项：`CardExtra(func(*ui.Context))` / `CardActions(func(*ui.Context))` |
| `ui.Image` + 自绘 | `extra.QRCode(c, content string, opts...)` | 码色随主题；`extra.QRCodeSize(px)` |
| `ui.Image` | `ant.ImageEl(c, src, ant.ImageWidth(units), ant.ImagePreview(b))` | 圆角显示 + 点击模态预览 |
| `ui.List` | `ant.List(c, *ui.ListState, n, row, opts...)` | `ListSelected` / `ListSelection` 绑定选中态 |
| `ui.Collapsible` / `ui.Accordion` | `ant.Collapsible(c, label, *bool, fn)` / `ant.Accordion` / `ant.AccordionItem` | 原生折叠动画；antd 风格面板用 `ant.Collapse` |
| `ui.Meter` | `ant.Meter(c, value, lo, hi, *ui.MeterLevels)` | 按阈值变色 |
| 自建 | `ant.Masonry(c, cols, items)` | 多列瀑布流（按列轮转分配） |
| 自建 | `ant.BorderBeam(c, content)` | 主题色渐变边框（静态） |

### 反馈

| MyGo 原生 | antd 风格 | 说明 |
| --- | --- | --- |
| `Row` 组合 | `ant.Alert(c, title string, opts...)` | `AlertTypeOpt(typ)` / `AlertDescription(s)` / `AlertClosable(fn)`；四级别 |
| `ui.Modal` / `ui.AlertDialog` | `ant.Modal(c, open *bool, title string, body func(), opts...)` / `ant.Confirm(c, open *bool, title, message string, onOk func())` | Modal 选项：`OkText` / `CancelText` / `OnOk` / `OnCancel` / `HideCancel`；**Confirm 无返回值**，只回调 `onOk` |
| `c.Toast` / `c.AddToast` | `ant.Message(c, text)` / `ant.MessageTypeTyped(c, typ, text)` / `ant.Notify(c, title, description, action, onAction)` | 复用 MyGo 底部 toast viewport |
| `ui.Progress` | `ant.ProgressLine(c, value float64, opts...)` / `ant.ProgressCircle(c, value, opts...)` | 线形/环形；`ProgressHideText()`、`ProgressStatusOpt(ProgressSuccess/Error)` |
| `ui.DialogBase` + 右侧锚定 | `ant.Drawer(c, *bool, title, body, opts...)` | `DrawerWidth(units)` / `DrawerClosable(b)` |
| `ui.Spinner` | `ant.Spin(c, opts...)` | 指示器 + 提示文字；Button 的 Loading 内部也使用它 |
| `ui.Popover` | `ant.Popover(c, anchor, *bool, body)` | 通用锚定气泡（调用方控制开关） |
| `el.ContextMenu(fn)` | `ant.ContextMenu(c, target, items, onSelect)` | 右键/菜单键触发；`ContextMenuItem{Danger}` 危险项 |
| — | `extra.Upload(c, opts...)` | 系统文件对话框 + 原生拖放 + 列表；选项：`UploadMultiple` / `UploadAccept` / `UploadHint` / `UploadOnChange` |

### 导航与布局

| MyGo 原生 | antd 风格 | 说明 |
| --- | --- | --- |
| `ui.Breadcrumbs` | `ant.Breadcrumb(c, items []string, chosen *int)` | `chosen=-1` 表示未选择 |
| `ui.Tabs` | `ant.Tabs(c, selected *int, labels ...string)` | — |
| `ui.MenuButton` / `Menu` | `ant.Dropdown(c, label string, items []ant.MenuItemData, opts...)` | `MenuItemData` 声明菜单；`DropdownDisabled(bool)` |
| 自建 | `ant.Pagination(c, current *int, total int, opts...)` | 页码窗口折叠，`OnPageChange(func(page int))` |
| 自建 | `ant.Steps(c, current *int, labels []string)` | 完成/进行/等待三态 |
| `ui.Split` / `SplitPercent` / `SplitVertical*` | `ant.Splitter` / `SplitterVertical` / `SplitterPercent` / `SplitterVerticalPercent` | 分隔条可拖拽 |
| `ui.Sidebar` / `SidebarItem` / `SidebarSection` | `ant.Sidebar(c, *string, fn, opts...)` / `ant.SidebarItem` / `ant.SidebarItemWithIcon` / `ant.SidebarSection` | `SidebarWidth(units)` 控宽 |

### 缺口控件（extra/）

| antd 组件 | 本库 API |
| --- | --- |
| Tree | `extra.Tree(c, nodes, extra.OnTreeNodeSelect(fn))` |
| TreeSelect | `extra.TreeSelect(c, *string, nodes, extra.TreeSelectPlaceholder(s), extra.TreeSelectDisabled(b), extra.TreeSelectOnSelect(fn))` |
| Cascader | `extra.Cascader(c, *[]string, nodes, extra.CascaderPlaceholder(s), extra.CascaderDisabled(b), extra.CascaderOnChange(fn))` |
| Transfer | `extra.Transfer(c, items []extra.TransferItem, right *[]string, extra.TransferTitles(l, r), extra.TransferOnChange(fn))` |
| Carousel | `extra.Carousel(c, pages []func(*ui.Context), extra.CarouselHeight(h))` |
| Upload | `extra.Upload(c, opts...)` |
| QRCode | `extra.QRCode(c, content, extra.QRCodeSize(px))` |

## 3. 事件回调差异

MyGo 原语是帧内查询（`el.Clicked()` / `el.Changed()`）；antd 包装层补齐了 React 风格回调：

- `ant.Button(..., ant.OnClick(fn))`：内部用 `ui.Local` 上升沿实现，一次点击只触发一次。
- `ant.Input(..., ant.OnChanged(fn))` / `ant.OnSubmit(fn)`、`ant.ControlOption` 的 `OnChange(fn)`：语义是「本帧发生变化」。
- 有状态值仍由调用方持有（`*string` / `*bool` / `*time.Time` / `*[]string`），与 MyGo 的受控绑定一致。
- **没有** 通用 `onClick` 属性：`ui.Modal` 的确定/取消用 `ModalOption`，`ant.Pagination` / `ant.Steps` 用各自的回调。

## 4. 已知差异与限制

- `ant.Table` 以 `[][]string` 持有全部行数据，但内部走 `ui.Table` 的可见行回调，虚拟滚动仍由 MyGo 提供；排序对数字列（含千分位逗号与 `%` 后缀）按数值比较，其余按字符串比较。
- `extra.Upload` 的文件对话框调用 `mygo.Dialog.Open`，会阻塞至对话框关闭；无窗口（Tester）环境下不要点击拖放区。
- `ant.Confirm` 基于 `ui.AlertDialog`，**返回被点按钮下标（0 起），`ant.Confirm` 本身不返回值**，只回调 `onOk`。
- `ant.Message` / `Notify` 复用 MyGo 底部 toast viewport，无法做到 antd 的顶部居中 / 右上角独立浮层。
- 颜色/间距铁律：包装层所有取值来自 `c.Theme()` 与 `t.Space/Rem`，运行时换肤自动生效；应用侧不要对 antd 控件硬编码颜色。
- 缺口控件已全部落地（见 03）；`Outline` / `GridView` / `StepSlider` / `SplitButton` 维持原生或由 Dropdown 组合。语义简化：`Tour` 不高亮目标、`Watermark` 为窗口级、`Masonry` 无高度测量、`BorderBeam` 无旋转动画。

## 5. 迁移步骤建议

1. `go get github.com/lishuangquan1987/mygo-antdesign@v0.1.0`（v0.2 发布后可用 `@v0.2.0`）。
2. 在根视图首行加 `theme.Apply(c, theme.Light())`。
3. 逐控件替换（对照第 2 节）：先静态展示类（Title/Paragraph/Card），
   再交互类（Button/Input），最后容器类（Space/Flex/Grid）。
4. 需要紧凑密度时给 token 加 `theme.WithDensity(theme.Compact)`，其余代码零改动。
5. 用 Gallery（`go run ./gallery`）核对目标控件的视觉与交互，并参考 `samples/` 下的可编译示例。
