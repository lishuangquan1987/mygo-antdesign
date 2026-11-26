# 02 · MyGo 原生控件 → Ant Design 皮肤映射

> 对应需求 3：把 MyGo 所有自带控件加上 Ant Design 设计。

## 映射图例

- **换肤方式**
  - `链式` — 直接在控件返回的 Element 上做样式调用即可
  - `Base` — 用对应的 `XxxBase` 重建外观（需要自定义内部结构）
  - `组合` — 用若干基础 Element / 控件拼装
- **复杂度**：低 / 中 / 高

## 一、通用与基础

| MyGo 控件 | Ant Design 对应 | 换肤方式 | 复杂度 | 要点 |
| --- | --- | --- | --- | --- |
| `ui.Button` / `ui.PrimaryButton` | Button | 链式 / `ButtonBase` | 低 | 主色→`t.Accent`；危险按钮→`t.Danger`；文字/链接按钮用 Base；实现 `type`=primary/default/dashed/text/link |
| `ui.Text` | Typography | 链式 | 低 | 标题/正文/次要文本映射 `t.Text` / `t.TextMuted`，字号用 `t.Rem(n)` |
| `ui.Icon` (SVG) | Icon | 组合 | 中 | 建立 antd 图标路径表，尺寸随 `t.Rem` |
| `ui.Row` / `ui.Column` / `ui.Grid` | Flex / Grid / Space | 链式 | 低 | `Gap` 用 `t.Space(n)`，实现 antd 的 gutter |
| `ui.Box` + `Border` | Divider | 组合 | 低 | `BorderWidth(0,0,1,0)` + 虚线样式 |
| `ui.Split` | Splitter | 链式 | 低 | 分隔条颜色→`t.Border` |
| `ui.Breadcrumbs` | Breadcrumb | 链式 | 低 | 分隔符与链接色→`t.TextMuted` / `t.Accent` |
| `ui.MenuButton` / `ui.ContextMenu` | Dropdown / Menu | `MenuBase` / 链式 | 中 | 弹层圆角、阴影、项高亮→主题 |
| `ui.Sidebar` / `ui.Navigation` | Menu / Layout | 链式 | 中 | 选中项→`t.Selection`，图标色跟随 |

## 二、数据录入

| MyGo 控件 | Ant Design 对应 | 换肤方式 | 复杂度 | 要点 |
| --- | --- | --- | --- | --- |
| `ui.Input` / `ui.TextInput` / `ui.TextArea` | Input / Input.TextArea | `TextInputBase` | 中 | 描边 `t.Border`、聚焦环 `t.Focus`、占位符 `t.TextMuted`；`TextAreaBase` 同理 |
| `ui.NumberInput` / `ui.DecimalInput` | InputNumber | `TextInputBase` + 组合 | 中 | 步进按钮复用 ButtonBase |
| `ui.SearchField` | Input.Search | `TextInputBase` + 组合 | 低 | 尾部搜索图标按钮 |
| `ui.Select` / `ui.Combobox` | Select | `SelectBase` / `ComboboxBase` | 高 | 触发器 + 弹层 + 选项高亮，重点复刻 antd 下拉 |
| `ui.Autocomplete` | AutoComplete | `ComboboxBase` | 中 | 同上，去掉多选 |
| `ui.Checkbox` | Checkbox | `CheckboxBase` / `CheckboxStateBase` | 低 | 勾选框尺寸用 `t.Space`，半选状态复用 Base |
| `ui.Radio` | Radio | `RadioBase` | 低 | 同上 |
| `ui.Switch` / `ui.Toggle` | Switch | `SwitchBase` / `ToggleBase` | 低 | 开关轨道色→`t.Accent` |
| `ui.Slider` / `ui.RangeSlider` | Slider | `SliderBase` | 中 | 轨道/滑块/刻度配色，双端区间 |
| `ui.Rating` | Rate | 组合 | 中 | 星星 SVG，选中色→`t.Warning`（antd 星色） |
| `ui.ColorPicker` | ColorPicker | 链式 | 中 | 面板圆角/阴影统一 |
| `ui.Calendar` / `ui.DateInput` | DatePicker / Calendar | `PopoverBase` + 组合 | 高 | 日期面板 + 触发器 + 快捷栏 |
| `ui.TimeInput` | TimePicker | 组合 | 中 | 时间列选择器 |
| `ui.Form` | Form | 链式 | 低 | 标签对齐、必填星号 `t.Danger`、校验信息色 |
| `ui.TokenField` | Tag（输入态）/ Mentions 近似 | 组合 | 中 | 标签用圆角色块 + 关闭按钮 |

## 三、数据展示

| MyGo 控件 | Ant Design 对应 | 换肤方式 | 复杂度 | 要点 |
| --- | --- | --- | --- | --- |
| `ui.Table` | Table | 链式 | 高 | 表头背景 `t.Surface`、斑马纹、行 hover、边框、尺寸档 |
| `ui.Tree` | Tree | `TreeBase` | 高 | 缩进线、展开图标、选中/悬停行 |
| `ui.List` / `ui.Grid` | List / Listy（虚拟列表） | 链式 | 中 | 虚拟滚动由 MyGo 原生提供 |
| `ui.Image` | Image | 链式 | 低 | 圆角、预览用 `DialogBase` |
| `ui.Avatar` | Avatar | 组合 | 低 | 圆形/方形、尺寸档 |
| `ui.Badge` | Badge | 组合 | 低 | 计数点→`t.Danger`，尺寸随密度 |
| `ui.Progress` / `ui.Meter` | Progress | 链式 | 中 | 线/环/仪表盘，状态色 success/warning/error |
| `ui.Collapsible` / `ui.Accordion` | Collapse | `CollapsibleBase` / `AccordionBase` | 中 | 手风琴头部背景 `t.Surface` |
| `ui.Tabs` | Tabs | `TabsBase` | 中 | 下划线/卡片风格，选中线→`t.Accent` |
| `ui.Segmented` | Segmented | `SegmentedBase` | 低 | 官方示例已给出完整换肤写法 |
| `ui.Stepper` | Steps | 组合 | 中 | 圆点/序号连接线，完成/进行中/等待三态 |
| `ui.Popover` | Popover / Tooltip | `PopoverBase` / `TooltipBase` | 中 | 箭头、圆角、阴影、延迟统一 |
| `ui.Tooltip` | Tooltip | `TooltipBase` | 低 | 用 `t.Inverse` / `t.InverseText` 反色 |
| `ui.Outline` | Anchor（近似） | 组合 | 中 | 目录树高亮当前项 |
| `ui.Timeline`（源码） | Timeline | 组合 | 中 | 时间轴节点与连接线 |

## 四、反馈与浮层

| MyGo 控件 | Ant Design 对应 | 换肤方式 | 复杂度 | 要点 |
| --- | --- | --- | --- | --- |
| `ui.Dialog` | Modal | `DialogBase` | 中 | 遮罩、居中卡片、页脚按钮顺序 |
| `ui.AlertDialog` | Popconfirm / Modal.confirm | `DialogBase` | 低 | 图标 + 文案 + 确认取消 |
| `ui.Toast` | Message / Notification | `ToastViewportBase` / `ToastBase` | 中 | 成功/警告/错误图标与状态色；`Message` 顶部居中、`Notification` 右上角 |
| `ui.Spinner`（loading 插件） | Spin | 组合 | 低 | antd 旋转图标 |
| `ui.Overlay` / `ui.Sidebar` | Drawer | `DialogBase` + 组合 | 中 | 从边缘滑入的面板 |
| `ui.Sticky` | Affix | 链式 | 低 | 滚动固钉 |
| `ui.Inspector` | Descriptions（近似） | 组合 | 中 | 键值对布局 |
| `ui.FindBar` | 无直接对应（可保留） | 链式 | 低 | 编辑器内查找栏，antd 无同类 |
| `ui.DragAndDrop` / `ui.Drop` | Upload（部分，见 03） | 组合 | 高 | 桌面拖放语义更强，见缺口章节 |

## 五、MyGo 特有、Ant Design 无对应（保留原生外观即可）

这些控件无需强行 antd 化，按需套用主题的颜色/圆角即可：`ui.Chart`、`ui.CodeEditor`、`ui.DiffViewer`、`ui.RichTextEditor`、`ui.Markdown`、`ui.WebView`、`ui.Video`、`ui.Drawing` / `ui.Painter`、`ui.Shader`、`ui.Motion` / `ui.Transition`、`ui.Windows`、`ui.HistoryButtons`、`ui.Toolbar`、`ui.Keyboard` 相关能力、`ui.Accessibility` 相关能力。

> 策略：**保持 MyGo 原生实现 + 主题配色同步**，不伪造不存在的 antd 组件。

## 六、换肤统一范式（设计草案，非实现代码）

所有换肤控件遵循同一模板，确保主题一致性：

```go
// 概念示意：每个控件从 c.Theme() 取色与度量，状态走条件分支
func Button(c *ui.Context, label string, opts ...Option) ui.Element {
    t := c.Theme()           // 唯一颜色来源
    o := apply(opts)
    el := ui.ButtonBase(c.Key(...), ...)   // 复用交互底座
    el.Padding(t.Space(2), t.Space(4))      // 度量来自主题
      .Radius(t.Radius)
      .Background(variantBg(t, o))          // 状态→颜色集中映射
      .Border(1, variantBorder(t, o))
    // hover / pressed / disabled / focus 由 variantXxx 统一处理
    return el
}
```

要点：
1. **禁止字面颜色**：所有颜色经 `t` 取，不写 `ui.Hex("#...")`。
2. **禁止字面间距**：所有间距经 `t.Space(n)` / `t.Rem(n)`。
3. **状态集中映射**：hover / pressed / disabled 交给共享函数，避免每个控件各写一套。
4. **签名贴近原生**：参数风格与 `ui.Button(c, "Save")` 一致，降低迁移成本。

## 七、覆盖度小结

- MyGo 组件文档共 **70+ 个**（官方口径约 87 个组件，含子构件）。
- 其中约 **55 个**可直接换肤或轻度组合；约 **12 个**需要基于 Base 重建；约 **15 个**属 MyGo 特有，保持原生配色同步即可。
- 缺口控件（antd 有、MyGo 无）另见 [03 缺口控件实现](03-gap-widgets.md)。
