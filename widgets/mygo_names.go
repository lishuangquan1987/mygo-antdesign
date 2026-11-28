package ant

import (
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/internal"
)

// 本文件集中提供与 MyGo 原生控件同名的入口（项目约定：mygo 与 antd 有
// 同类型控件时，命名以 mygo 为准；此前按 antd 命名的入口保留为别名，
// 不破坏既有代码）。对应关系见 docs/02-widget-mapping.md 的覆盖矩阵。

// Breadcrumbs 面包屑（MyGo `ui.Breadcrumbs` 命名；`Breadcrumb` 为别名）。
func Breadcrumbs(c *ui.Context, items []string, chosen *int) ui.Element {
	return Breadcrumb(c, items, chosen)
}

// Avatar 头像（MyGo `ui.Avatar` 命名；`AvatarEl` 为别名）。
func Avatar(c *ui.Context, name string, image ...*ui.Bitmap) ui.Element {
	return AvatarEl(c, name, image...)
}

// Autocomplete 自动补全（MyGo 命名；`AutoComplete` 为别名）。
func Autocomplete(c *ui.Context, value *string, suggestions []string, opts ...ControlOption) ui.Element {
	return AutoComplete(c, value, suggestions, opts...)
}

// Image 图片（MyGo 命名；`ImageEl` 为别名）。
func Image(c *ui.Context, src ui.ImageSource, opts ...ImageOption) ui.Element {
	return ImageEl(c, src, opts...)
}

// NumberInput 数字输入（MyGo 命名；`NumberInputV` 为别名）。
func NumberInput(c *ui.Context, value *float64, lo, hi, step float64, opts ...InputOption) ui.Element {
	return NumberInputV(c, value, lo, hi, step, opts...)
}

// Progress 进度条（MyGo 命名，线形；`ProgressLine` 为别名）。
func Progress(c *ui.Context, value float64, opts ...ProgressOption) ui.Element {
	return ProgressLine(c, value, opts...)
}

// RangeSlider 双端滑动条（MyGo 命名；`RangeSliderDouble` 为别名）。
func RangeSlider(c *ui.Context, low, high *float64, lo, hi, step float64, opts ...ControlOption) ui.Element {
	return RangeSliderDouble(c, low, high, lo, hi, step, opts...)
}

// StepSlider 带刻度滑动条（MyGo 命名）。
func StepSlider(c *ui.Context, value *float64, lo, hi, step float64, opts ...ControlOption) ui.Element {
	cfg := &controlConfig{}
	for _, o := range opts {
		o(cfg)
	}
	el := ui.StepSlider(c, value, lo, hi, step)
	applyControl(el, cfg)
	return el
}

// Rating 评分（MyGo 命名；`Rate` 为别名）。
func Rating(c *ui.Context, value *int, max int, opts ...ControlOption) ui.Element {
	return Rate(c, value, max, opts...)
}

// Spinner 加载指示器（沿用 MyGo 的 Spinner 命名；实现是自绘的 Spin，
// 不是对 ui.Spinner 的包装）。
func Spinner(c *ui.Context, opts ...SpinOption) ui.Element {
	return Spin(c, opts...)
}

// PrimaryButton 主要按钮（MyGo `ui.PrimaryButton` 命名）。
func PrimaryButton(c *ui.Context, label string, opts ...ButtonOption) ui.Element {
	all := make([]ButtonOption, 0, len(opts)+1)
	all = append(all, Primary())
	all = append(all, opts...)
	return Button(c, label, all...)
}

// CheckboxState 三态复选框（MyGo 命名；与 `CheckboxAll` 等价）。
func CheckboxState(c *ui.Context, state *ui.CheckState, label string, opts ...ControlOption) ui.Element {
	cfg := &controlConfig{}
	for _, o := range opts {
		o(cfg)
	}
	el := ui.CheckboxState(c, state, label)
	applyControl(el, cfg)
	return el
}

// CheckboxGroup 复选组：一个总复选框控制下方一组复选框（MyGo 命名）。
func CheckboxGroup(c *ui.Context, label string, fn func(), opts ...ControlOption) ui.Element {
	cfg := &controlConfig{}
	for _, o := range opts {
		o(cfg)
	}
	el := ui.CheckboxGroup(c, label, fn)
	applyControl(el, cfg)
	return el
}

// ColorWell 取色块：点击展开系统取色面板（MyGo 命名）。
func ColorWell(c *ui.Context, color *ui.Color, opts ...ControlOption) ui.Element {
	cfg := &controlConfig{}
	for _, o := range opts {
		o(cfg)
	}
	el := ui.ColorWell(c, color)
	applyControl(el, cfg)
	return el
}

// Toolbar 工具条容器（MyGo `ui.Toolbar` 命名）：间距与圆角随主题。
func Toolbar(c *ui.Context, fn func()) ui.Element {
	t := c.Theme()
	return ui.Toolbar(c, fn).
		Gap(t.Space(1)).
		Padding(t.Space(1)).
		Radius(t.Radius).
		Background(t.Surface).
		Border(1, t.Border)
}

// GridView 网格视图（MyGo 命名）：等宽网格 + 主题间距。
func GridView(c *ui.Context, s *ui.GridState, n int, minWidth, height float32, item func(i int)) ui.Element {
	return ui.GridView(c, s, n, minWidth, height, item).Gap(c.Theme().Space(2))
}

// FindBar 查找栏（MyGo 命名）：与 antd 无对应，保留原生交互并套主题面。
func FindBar(c *ui.Context, open *bool, query *string, matches int, current *int) ui.Element {
	t := c.Theme()
	return ui.FindBar(c, open, query, matches, current).
		Padding(t.Space(1)).
		Radius(t.Radius).
		Background(t.Surface).
		Border(1, t.Border)
}

// BackButton 历史后退按钮（MyGo 命名）。
func BackButton(c *ui.Context, r *ui.Router) ui.Element {
	return ui.BackButton(c, r)
}

// ForwardButton 历史前进按钮（MyGo 命名）。
func ForwardButton(c *ui.Context, r *ui.Router) ui.Element {
	return ui.ForwardButton(c, r)
}

// RichText 富文本片段（MyGo 命名）：统一字号与正文色。
func RichText(c *ui.Context, spans ...ui.Span) ui.Element {
	t := c.Theme()
	return ui.RichText(c, spans...).FontSize(t.FontSize).TextColor(t.Text)
}

// Icon 图标（MyGo `ui.Icon` 命名）：颜色跟随文本色，尺寸由 FontSize 或显式尺寸决定。
func Icon(c *ui.Context, svg *ui.SVG) ui.Element {
	return ui.Icon(c, svg)
}

// TextInputBuffer 大文本单行输入（MyGo 命名）：复用 Input 的外框与焦点环，
// 内部使用索引化 TextBuffer，适合百万行文档的局部编辑。
func TextInputBuffer(c *ui.Context, value *ui.TextBuffer, opts ...InputOption) ui.Element {
	cfg := &inputConfig{size: internal.SizeMiddle}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()

	el := inputFrame(c, t, cfg)
	var inp ui.Element
	el.Children(func() {
		inp = ui.TextInputBuffer(c, value).Grow(1).Background(t.Surface).
			FontSize(internal.ControlFont(t, cfg.size)).
			Border(0, ui.Transparent).Padding(0).Radius(0)
		if cfg.label != "" {
			inp.Label(cfg.label)
		}
	})
	el.Disabled(cfg.disabled)
	if cfg.readonly {
		el.ReadOnly(true)
	}
	if !cfg.disabled && inp.Focused() {
		el.Border(2, t.Focus)
	}
	fireInputEvents(inp, cfg)
	return el
}

// TextAreaBuffer 大文本多行输入（MyGo 命名）：复用 TextArea 的外框与焦点环。
func TextAreaBuffer(c *ui.Context, value *ui.TextBuffer, opts ...InputOption) ui.Element {
	cfg := &inputConfig{size: internal.SizeMiddle}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()

	rows := cfg.rows
	if rows < 2 {
		rows = 2
	}
	el := ui.Row(c).
		MinHeight(t.Space(2)+t.Rem(1.7)*float32(rows)).
		Padding(t.Space(1), internal.ControlPaddingX(t, cfg.size)).
		Radius(t.Radius).
		Background(t.Surface).
		AlignItems(ui.Start)
	if cfg.disabled {
		el.Background(t.Background).Border(1, t.Border)
	} else {
		el.Border(1, t.Border)
	}
	var inp ui.Element
	el.Children(func() {
		inp = ui.TextAreaBuffer(c, value).Grow(1).FillWidth().Background(t.Surface).
			Border(0, ui.Transparent).Padding(0).Radius(0).
			FontSize(internal.ControlFont(t, cfg.size))
		if cfg.label != "" {
			inp.Label(cfg.label)
		}
	})
	el.Disabled(cfg.disabled)
	if cfg.readonly {
		el.ReadOnly(true)
	}
	if !cfg.disabled && inp.Focused() {
		el.Border(2, t.Focus)
	}
	fireInputEvents(inp, cfg)
	return el
}

// Scroll 纵向滚动容器（MyGo 命名）。
func Scroll(c *ui.Context) ui.Element { return ui.Scroll(c) }

// ScrollBoth 双向滚动容器（MyGo 命名）。
func ScrollBoth(c *ui.Context) ui.Element { return ui.ScrollBoth(c) }

// ScrollHorizontal 横向滚动容器（MyGo 命名）。
func ScrollHorizontal(c *ui.Context) ui.Element { return ui.ScrollHorizontal(c) }

// Split 左右分栏 + 可拖拽分隔条（MyGo 命名）。
func Split(c *ui.Context, size *float32, first, second func()) ui.Element {
	return ui.Split(c, size, first, second)
}

// SplitPercent 按百分比分栏（MyGo 命名）。
func SplitPercent(c *ui.Context, p *float32, first, second func()) ui.Element {
	return ui.SplitPercent(c, p, first, second)
}

// SplitVertical 上下分栏（MyGo 命名）。
func SplitVertical(c *ui.Context, size *float32, first, second func()) ui.Element {
	return ui.SplitVertical(c, size, first, second)
}

// SplitVerticalPercent 按百分比上下分栏（MyGo 命名）。
func SplitVerticalPercent(c *ui.Context, p *float32, first, second func()) ui.Element {
	return ui.SplitVerticalPercent(c, p, first, second)
}

// Video 视频播放器画面（MyGo 命名）：按主题圆角裁切。
func Video(c *ui.Context, p *ui.VideoPlayer) ui.Element {
	return ui.Video(c, p).Radius(c.Theme().Radius)
}

// WebView 网页视图（MyGo 命名）：按主题圆角裁切。
func WebView(c *ui.Context, v ui.NativeWebView) ui.Element {
	return ui.WebView(c, v).Radius(c.Theme().Radius)
}

// NativeView 平台原生视图（MyGo 命名）：按主题圆角裁切。
func NativeView(c *ui.Context, v ui.PlatformView) ui.Element {
	return ui.NativeView(c, v).Radius(c.Theme().Radius)
}

// SplitButton 分裂按钮（MyGo 命名；antd 的 Dropdown.Button）。
func SplitButton(c *ui.Context, label string) ui.SplitButtonParts {
	parts := ui.SplitButton(c, label)
	parts.Group.Radius(c.Theme().Radius)
	return parts
}

// PrimarySplitButton 主色分裂按钮（MyGo 命名）。
func PrimarySplitButton(c *ui.Context, label string) ui.SplitButtonParts {
	parts := ui.PrimarySplitButton(c, label)
	parts.Group.Radius(c.Theme().Radius)
	return parts
}

// MenuButton 菜单按钮（MyGo 命名）：原生下拉菜单。
func MenuButton(c *ui.Context, label string, build func(m *ui.Menu)) ui.Element {
	t := c.Theme()
	return ui.MenuButton(c, label, build).
		Radius(t.Radius).
		FontSize(t.FontSize)
}

// PopoverMenuButton 浮层菜单按钮（MyGo 命名）。
func PopoverMenuButton(c *ui.Context, label string, build func(m *ui.Menu), options ...ui.PopoverMenuButtonOptions) ui.Element {
	t := c.Theme()
	return ui.PopoverMenuButton(c, label, build, options...).
		Radius(t.Radius).
		FontSize(t.FontSize)
}

// AlertDialog 系统警示对话框（MyGo 命名）：返回按钮下标 / AlertDismissed / -1。
func AlertDialog(c *ui.Context, open *bool, title, message string, buttons ...string) int {
	return ui.AlertDialog(c, open, title, message, buttons...)
}

// TimeInput 时间输入（MyGo 命名；`TimePicker` 为 antd 风格别名）。
func TimeInput(c *ui.Context, tm *time.Time, opts ...InputOption) ui.Element {
	return TimePicker(c, tm, opts...)
}

// DateInput 日期输入（MyGo 命名；`DatePicker` 为 antd 风格别名）。
func DateInput(c *ui.Context, date *time.Time, opts ...InputOption) ui.Element {
	return DatePicker(c, date, opts...)
}

// ToggleBase 保持按下的按钮底座（MyGo 命名）。
func ToggleBase(c *ui.Context, on *bool) ui.Element {
	t := c.Theme()
	return ui.ToggleBase(c, on).Radius(t.Radius).FontSize(t.FontSize)
}
