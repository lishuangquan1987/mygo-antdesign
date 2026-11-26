package ant

import (
	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/internal"
)

// InputOption 定制文本录入类控件（Input / TextArea / SearchInput / NumberInput）。
type InputOption func(*inputConfig)

type inputConfig struct {
	size     internal.SizeKind
	ph       string
	label    string
	password bool
	disabled bool
	readonly bool
	maxLen   int
	rows     int
	prefix   func(*ui.Context)
	suffix   func(*ui.Context)
	onChg    func()
	onSub    func()
}

// InputSize 尺寸档（antd size）。
func InputSize(s SizeKind) InputOption { return func(c *inputConfig) { c.size = s } }

// PlaceholderText 占位提示（antd placeholder）。
func PlaceholderText(s string) InputOption { return func(c *inputConfig) { c.ph = s } }

// InputLabel 无障碍标签（对齐 MyGo 原语 Label，屏幕阅读器与测试定位用）。
func InputLabel(s string) InputOption { return func(c *inputConfig) { c.label = s } }

// PasswordInput 密码模式（显示为圆点）。
func PasswordInput() InputOption { return func(c *inputConfig) { c.password = true } }

// InputDisabled 禁用（antd disabled）。
func InputDisabled(b bool) InputOption { return func(c *inputConfig) { c.disabled = b } }

// InputReadOnly 只读（antd readOnly）。
func InputReadOnly(b bool) InputOption { return func(c *inputConfig) { c.readonly = b } }

// MaxLengthText 最大输入长度。
func MaxLengthText(n int) InputOption { return func(c *inputConfig) { c.maxLen = n } }

// Rows TextArea 行数（antd rows，默认 3）。
func Rows(n int) InputOption { return func(c *inputConfig) { c.rows = n } }

// InputPrefix 前缀（antd prefix）。
func InputPrefix(fn func(*ui.Context)) InputOption { return func(c *inputConfig) { c.prefix = fn } }

// InputSuffix 后缀（antd suffix）。
func InputSuffix(fn func(*ui.Context)) InputOption { return func(c *inputConfig) { c.suffix = fn } }

// OnChanged 值变化回调（每帧 el.Changed() 边沿触发）。
func OnChanged(fn func()) InputOption { return func(c *inputConfig) { c.onChg = fn } }

// OnSubmit 单行输入框内按 Enter 回调（antd onPressEnter）。
func OnSubmit(fn func()) InputOption { return func(c *inputConfig) { c.onSub = fn } }

// inputFrame 统一的输入框容器：高度、底色、圆角、内边距与状态描边
// （默认 1px colorBorder，聚焦 2px 焦点色，禁用灰底——对齐 antd）。
func inputFrame(c *ui.Context, t *ui.Theme, cfg *inputConfig, focused bool) ui.Element {
	el := ui.Row(c).
		Height(internal.ControlHeight(t, cfg.size)).
		Padding(0, internal.ControlPaddingX(t, cfg.size)).
		Radius(t.Radius).
		Background(t.Surface).
		AlignItems(ui.Center)
	if cfg.disabled {
		el.Background(t.Background).Border(1, t.Border)
	} else if focused {
		el.Border(2, t.Focus)
	} else {
		el.Border(1, t.Border)
	}
	return el
}

// fireInputEvents Changed/Submitted 的边沿回调（Changed 本身是「本帧是否
// 变化」语义，天然是边沿）。
func fireInputEvents(el ui.Element, cfg *inputConfig) {
	if cfg.onChg != nil && el.Changed() {
		cfg.onChg()
	}
	if cfg.onSub != nil && el.Submitted() {
		cfg.onSub()
	}
}

// Input 单行文本输入框（antd Input）。value 由调用方持有（即时模式）：
//
//	var name string
//	ant.Input(c, &name, ant.PlaceholderText("请输入姓名"), ant.OnChanged(app.save))
func Input(c *ui.Context, value *string, opts ...InputOption) ui.Element {
	cfg := &inputConfig{size: internal.SizeMiddle}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()

	el := inputFrame(c, t, cfg, false)
	var inp ui.Element
	el.Children(func() {
		if cfg.prefix != nil {
			cfg.prefix(c)
		}
		inp = ui.TextInput(c, value).Grow(1).Background(t.Surface).FontSize(internal.ControlFont(t, cfg.size))
		if cfg.label != "" {
			inp.Label(cfg.label)
		}
		if cfg.ph != "" {
			inp.Placeholder(cfg.ph).PlaceholderColor(t.TextMuted)
		}
		if cfg.password {
			inp.Password()
		}
		if cfg.maxLen > 0 {
			inp.MaxLength(cfg.maxLen)
		}
		if cfg.suffix != nil {
			cfg.suffix(c)
		}
	})
	el.Disabled(cfg.disabled)
	if cfg.readonly {
		el.ReadOnly(true)
	}

	// 聚焦描边：查询内层输入元素的焦点状态（容器本身不可聚焦）。
	if !cfg.disabled && inp.Focused() {
		el.Border(2, t.Focus)
	}
	fireInputEvents(inp, cfg)
	return el
}

// TextArea 多行文本输入框（antd Input.TextArea）：
//
//	var bio string
//	ant.TextArea(c, &bio, ant.Rows(4), ant.PlaceholderText("介绍一下自己"))
func TextArea(c *ui.Context, value *string, opts ...InputOption) ui.Element {
	cfg := &inputConfig{size: internal.SizeMiddle, rows: 3}
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
		inp = ui.TextArea(c, value).Grow(1).FillWidth().Background(t.Surface).
			FontSize(internal.ControlFont(t, cfg.size))
		if cfg.label != "" {
			inp.Label(cfg.label)
		}
		if cfg.ph != "" {
			inp.Placeholder(cfg.ph).PlaceholderColor(t.TextMuted)
		}
		if cfg.maxLen > 0 {
			inp.MaxLength(cfg.maxLen)
		}
	})
	el.Disabled(cfg.disabled)
	if cfg.readonly {
		el.ReadOnly(true)
	}

	// 聚焦描边：查询内层输入元素的焦点状态。
	if !cfg.disabled && inp.Focused() {
		el.Border(2, t.Focus)
	}
	fireInputEvents(inp, cfg)
	return el
}

// SearchInput 搜索框（antd Input.Search）：Enter 提交触发 OnSubmit，
// 默认带搜索图标后缀。
//
//	var q string
//	ant.SearchInput(c, &q, ant.PlaceholderText("搜索…"), ant.OnSubmit(app.search))
func SearchInput(c *ui.Context, query *string, opts ...InputOption) ui.Element {
	return Input(c, query, append(opts, InputOption(func(cfg *inputConfig) {
		if cfg.suffix == nil {
			cfg.suffix = func(c *ui.Context) {
				ui.Text(c, "⌕").TextColor(c.Theme().TextMuted)
			}
		}
	}))...)
}

// NumberInputV 数字输入（antd InputNumber）。注意 v0.10.0 的 NumberInput
// 原语要求 *float64：
//
//	var port float64 = 8080
//	ant.NumberInputV(c, &port, 1, 65535, 1, ant.OnChanged(app.restart))
func NumberInputV(c *ui.Context, value *float64, lo, hi, step float64, opts ...InputOption) ui.Element {
	cfg := &inputConfig{size: internal.SizeMiddle}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()

	el := inputFrame(c, t, cfg, false)
	var inp ui.Element
	el.Children(func() {
		inp = ui.NumberInput(c, value, lo, hi, step).Grow(1).
			Background(t.Surface).FontSize(internal.ControlFont(t, cfg.size))
		if cfg.label != "" {
			inp.Label(cfg.label)
		}
		if cfg.ph != "" {
			inp.Placeholder(cfg.ph).PlaceholderColor(t.TextMuted)
		}
	})
	el.Disabled(cfg.disabled)
	fireInputEvents(inp, cfg)
	return el
}
