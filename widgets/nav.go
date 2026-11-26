package ant

import (
	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/internal"
)

// Breadcrumb 面包屑（antd Breadcrumb）。chosen 为当前层级下标：
//
//	ant.Breadcrumb(c, &level, []string{"首页", "列表", "详情"})
func Breadcrumb(c *ui.Context, items []string, chosen *int) ui.Element {
	return ui.Breadcrumbs(c, items, chosen)
}

// Tabs 页签（antd Tabs）。selected 为选中下标，Changed 报告切换：
//
//	ant.Tabs(c, &tab, "常规", "外观", "高级")
func Tabs(c *ui.Context, selected *int, labels ...string) ui.Element {
	return ui.Tabs(c, selected, labels...)
}

// DropdownOption 定制 Dropdown。
type DropdownOption func(*dropdownConfig)

type dropdownConfig struct {
	disabled bool
}

// DropdownDisabled 禁用触发按钮。
func DropdownDisabled(b bool) DropdownOption { return func(c *dropdownConfig) { c.disabled = b } }

// MenuItemData 下拉菜单项（antd Menu.Item）。
type MenuItemData struct {
	Label    string
	Danger   bool
	Disabled bool
	// Separator 为 true 时该项是一条分隔线（Label 忽略）。
	Separator bool
	OnSelect  func()
}

// Dropdown 下拉菜单（antd Dropdown + Menu）：
//
//	ant.Dropdown(c, "操作", []ant.MenuItemData{
//	    {Label: "编辑", OnSelect: app.edit},
//	    {Separator: true},
//	    {Label: "删除", Danger: true, OnSelect: app.remove},
//	})
func Dropdown(c *ui.Context, label string, items []MenuItemData, opts ...DropdownOption) ui.Element {
	cfg := &dropdownConfig{}
	for _, o := range opts {
		o(cfg)
	}
	el := ui.MenuButton(c, label, func(m *ui.Menu) {
		for _, it := range items {
			if it.Separator {
				m.Separator()
				continue
			}
			mi := m.Item(it.Label)
			if it.Disabled {
				mi.Disabled(true)
			}
			if it.Danger {
				mi.Danger()
			}
			if mi.Chosen() && it.OnSelect != nil {
				it.OnSelect()
			}
		}
	})
	if cfg.disabled {
		el.Disabled(true)
	}
	return el
}

// PaginationOption 定制 Pagination。
type PaginationOption func(*paginationConfig)

type paginationConfig struct {
	onChange func(page int)
}

// OnPageChange 页码变化回调。
func OnPageChange(fn func(page int)) PaginationOption {
	return func(c *paginationConfig) { c.onChange = fn }
}

// Pagination 分页器（antd Pagination）。current 由调用方持有（1 起）：
//
//	ant.Pagination(c, &page, 42, ant.OnPageChange(app.loadPage))
func Pagination(c *ui.Context, current *int, total int, opts ...PaginationOption) ui.Element {
	cfg := &paginationConfig{}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()

	fire := func() {
		if cfg.onChange != nil {
			cfg.onChange(*current)
		}
	}

	row := ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center)
	row.Children(func() {
		// 页码窗口：总数少则全展示，多则 1 … 窗口 … 总数。
		pages := pageWindow(*current, total)
		for _, p := range pages {
			if p == 0 { // 省略号
				ui.Text(c, "…").TextColor(t.TextMuted)
				continue
			}
			btn := ui.ButtonBase(c).
				Height(internal.ControlHeight(t, internal.SizeSmall)).
				MinWidth(internal.ControlHeight(t, internal.SizeSmall)).
				Padding(0, t.Space(1)).
				Radius(t.Radius).
				Center()
			if p == *current {
				btn.Background(t.Accent)
				btn.Children(func() {
					ui.Textf(c, "%d", p).TextColor(t.AccentText).SingleLine()
				})
			} else {
				btn.Children(func() {
					ui.Textf(c, "%d", p).TextColor(t.Text).SingleLine()
				})
			}
			if btn.Clicked() && p != *current {
				*current = p
				fire()
			}
		}
	})
	return row
}

// pageWindow 生成分页页码序列：0 表示省略号。
func pageWindow(cur, total int) []int {
	if total <= 7 {
		out := make([]int, 0, total)
		for i := 1; i <= total; i++ {
			out = append(out, i)
		}
		return out
	}
	var out []int
	add := func(p ...int) { out = append(out, p...) }
	switch {
	case cur <= 3:
		add(1, 2, 3, 4, 0, total)
	case cur >= total-2:
		add(1, 0, total-3, total-2, total-1, total)
	default:
		add(1, 0, cur-1, cur, cur+1, 0, total)
	}
	return out
}

// Steps 步骤条（antd Steps）。current 为当前步骤下标（0 起），之前的
// 步骤显示完成态：
//
//	ant.Steps(c, &step, []string{"填写信息", "确认订单", "支付"})
func Steps(c *ui.Context, current *int, labels []string) ui.Element {
	t := c.Theme()

	row := ui.Row(c).FillWidth().AlignItems(ui.Center)
	row.Children(func() {
		for i, label := range labels {
			if i > 0 {
				// 连接线：完成的实心，未完成的描边。
				line := ui.Box(c).Grow(1).Height(t.Space(1)).Border(1, t.Border)
				if i <= *current {
					line.BorderStyle(ui.BorderSolid)
					line.Border(1, t.Accent)
				}
				line.BorderWidth(0, 0, 1, 0)
				line.Children(func() {})
			}
			done := i < *current
			now := i == *current

			item := ui.Column(c).Gap(t.Space(1)).AlignItems(ui.Center)
			item.Children(func() {
				circle := ui.Box(c).Size(t.Space(6), t.Space(6)).Radius(t.Space(3)).Center()
				switch {
				case done:
					circle.Background(t.Accent)
					circle.Children(func() {
						ui.Text(c, "✓").TextColor(t.AccentText).FontSize(t.Rem(0.8))
					})
				case now:
					circle.Background(t.Accent.Mix(t.Background, 0.85))
					circle.Border(1, t.Accent)
					circle.Children(func() {
						ui.Textf(c, "%d", i+1).TextColor(t.Accent)
					})
				default:
					circle.Background(t.Background.Mix(t.Text, 0.05))
					circle.Children(func() {
						ui.Textf(c, "%d", i+1).TextColor(t.TextMuted)
					})
				}
				color := t.TextMuted
				if now || done {
					color = t.Text
				}
				ui.Text(c, label).TextColor(color).FontSize(t.Rem(0.86)).SingleLine()
			})
		}
	})
	return row
}
