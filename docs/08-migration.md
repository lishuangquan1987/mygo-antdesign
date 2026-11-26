# 08 迁移指南：从 MyGo 原生到 mygo-antdesign

本指南面向已使用 MyGo `package ui` 的应用，逐步迁移到 Ant Design 风格的
`mygo-antdesign`。控件是 `ui` 原语的换肤包装或组合，即时模式语义
（view = f(state)、`*T` 受控绑定、`ui.Local`）完全一致。

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

- 全局切换：`theme.Apply(c, token)` 每帧调用安全，`token` 换成
  `theme.Dark(theme.WithPrimary(ui.Hex("#7c3aed")))` 即整窗变色。
- 局部切换：`theme.Sub(c, token)` 返回 `ui.Themed` 容器，适合亮暗并排
  预览、独立主题面板。
- 紧凑密度（上位机 / IDE）：`theme.WithDensity(theme.Compact)`。

## 2. 控件对照表

### 通用组件

| MyGo 原生 | antd 风格 | 说明 |
| --- | --- | --- |
| `ui.Button` / `ButtonBase` | `ant.Button(c, label, opts...)` | `ant.Primary()` / `Danger()` / `Size(ant.Large)` / `Loading` / `OnClick` |
| `ui.Text` 组合 | `ant.Title` / `ant.Paragraph` / `ant.Code` | H1–H5 字级、Secondary/Delete/Strong 等变体 |
| `ui.Row` + 手画线 | `ant.Divider` / `ant.VerticalDivider` | 虚线、带文字、对齐方向 |

### 数据录入

| MyGo 原生 | antd 风格 | 说明 |
| --- | --- | --- |
| `ui.TextInput` / `TextArea` | `ant.Input` / `ant.TextArea` | 聚焦描边、Placeholder、InputLabel、OnChanged/OnSubmit |
| `ui.TextInput` 组合 | `ant.SearchInput` | 回车提交（antd Input.Search） |
| `ui.NumberInput` | `ant.NumberInputV(c, *float64, lo, hi, step)` | v0.10.0 原语要求 `*float64` |
| `ui.Select` | `ant.Select(c, *string, options)` | 高度/字号/圆角对齐 antd |
| `ui.Combobox` | `ant.Combobox` | 输入过滤 + 候选 |
| `ui.DateInput` / `TimeInput` | `ant.DatePicker` / `ant.TimePicker` | `*time.Time` 受控 |
| `ui.Checkbox` / `CheckboxState` | `ant.Checkbox` / `ant.CheckboxAll` | 三态全选框 |
| `ui.RadioGroup` / `Radio` | `ant.RadioGroup` / `ant.Radio[T]` | 泛型受控值 |
| `ui.Switch` | `ant.Switch` / `ant.SwitchWithLabel` | |
| `ui.Slider` / `RangeSlider` | `ant.Slider` / `ant.RangeSliderDouble` | |
| `ui.Rating` | `ant.Rate` | |
| `ui.Segmented` | `ant.Segmented(c, *int, labels)` | 支持 `ControlOption`（Disabled/OnChange） |

### 数据展示

| MyGo 原生 | antd 风格 | 说明 |
| --- | --- | --- |
| `ui.Table` | `ant.Table(c, columns, rows)` | 列定义、表头排序、`TableSelection` 行多选 |
| `ui.Tree` / `TreeItem` | `extra.Tree(c, nodes)` | 节点树声明式构建，`OnTreeNodeSelect` |
| `ui.Badge` 组合 | `ant.Badge` / `ant.BadgeDot` | 右上角骑角定位（Overlay + AttachTo） |
| `ui.Avatar` | `ant.AvatarEl` | |
| Row/Column 组合 | `ant.Card` | 头部/内容/操作区 |
| `ui.Image` + 自绘 | `extra.QRCode` | 码色随主题 |

### 反馈

| MyGo 原生 | antd 风格 | 说明 |
| --- | --- | --- |
| `ui.Row` 组合 | `ant.Alert` | 四级别、可关闭 |
| `ui.Modal` / `AlertDialog` | `ant.Modal` / `ant.Confirm` | OnOk/OnCancel；Confirm 返回确定与否 |
| `c.Toast` / `AddToast` | `ant.Message` / `ant.Notify` | 轻提示与通知 |
| `ui.Progress` | `ant.ProgressLine` | `ProgressHideText` 隐藏百分比 |
| — | `extra.Upload` | 系统文件对话框 + 原生拖放 + 列表 |

### 导航与布局

| MyGo 原生 | antd 风格 | 说明 |
| --- | --- | --- |
| `ui.Breadcrumbs` | `ant.Breadcrumb` | |
| `ui.Tabs` | `ant.Tabs` | |
| `ui.MenuButton` / `Menu` | `ant.Dropdown(c, label, items)` | MenuItemData 声明菜单 |
| Row + ButtonBase | `ant.Pagination` | 页码窗口折叠，OnPageChange |
| Row 组合 | `ant.Steps` | 完成/进行/等待三态 |
| `ui.Row` / `Column` / `Grid` | `ant.Space` / `ant.Flex` / `ant.Grid` | 间距经 `t.Space` 随密度缩放 |

### 图标

| MyGo 原生 | antd 风格 | 说明 |
| --- | --- | --- |
| `ui.ParseSVG` + `ui.Image` | `icon.Icon(c, icon.Search, 16)` | 44 个 outlined 常量，颜色跟随文本色 |

## 3. 事件回调差异

MyGo 原语是帧内查询（`el.Clicked()` / `el.Changed()`）；antd 包装层补齐了
React 风格回调：

- `ant.Button(..., ant.OnClick(fn))`：内部用 `ui.Local` 上升沿实现，
  一次点击只触发一次。
- `ant.Checkbox(..., ant.OnChange(fn))` 等：`Changed()` 本帧变化语义。
- 有状态值仍由调用方持有（`*string` / `*bool` / `*time.Time`），
  与 MyGo 的受控绑定一致。

## 4. 已知差异与限制

- `ant.Table` 单元格为字符串比较排序（数字列请先格式化定宽文本）。
- `extra.Upload` 的文件对话框调用 `mygo.Dialog.Open`，会阻塞至对话框关闭；
  无窗口（Tester）环境下不要点击拖放区。
- `Confirm` 基于 `ui.AlertDialog`，返回被点按钮语义（确定=true）。
- 颜色/间距铁律：包装层所有取值来自 `c.Theme()` 与 `t.Space/Rem`，
  运行时换肤自动生效；应用侧不要对 antd 控件硬编码颜色。

## 5. 迁移步骤建议

1. `go get github.com/lishuangquan1987/mygo-antdesign@v0.1.0`。
2. 在根视图首行加 `theme.Apply(c, theme.Light())`。
3. 逐控件替换（对照第 2 节）：先静态展示类（Title/Paragraph/Card），
   再交互类（Button/Input），最后容器类（Space/Flex/Grid）。
4. 需要紧凑密度时给 token 加 `theme.WithDensity(theme.Compact)`，
   其余代码零改动。
5. 用 Gallery（`go run ./gallery`）核对目标控件的视觉与交互。
