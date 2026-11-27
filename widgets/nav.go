package ant

import (
	"fmt"

	"github.com/egoist/mygo/ui"
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
//
// 样式对齐 antd v5：页签间距 32px、选中项主色 + 2px 下划线、悬停变主色，
// 整行底部 1px 分隔线。
func Tabs(c *ui.Context, selected *int, labels ...string) ui.Element {
	t := c.Theme()
	row := ui.Row(c).FillWidth().
		Gap(t.Space(8)).
		Border(1, t.Border).
		BorderWidth(0, 0, 1, 0)
	row.Children(func() {
		for i, label := range labels {
			i := i
			active := i == *selected
			tab := ui.ButtonBase(c).Padding(t.Space(3), 0, 0, 0)
			tab.Children(func() {
				col := ui.Column(c).Gap(t.Space(1))
				col.Children(func() {
					color := t.Text
					if active {
						color = t.Accent
					} else if tab.Hovered() {
						color = t.AccentHover
					}
					ui.Text(c, label).TextColor(color).FontSize(t.FontSize).SingleLine()
					bar := ui.Box(c).FillWidth().Height(t.Space(0.5))
					if active {
						bar.Background(t.Accent)
					} else {
						bar.Background(ui.Transparent)
					}
					bar.Children(func() {})
				})
			})
			if tab.Clicked() {
				*selected = i
			}
		}
	})
	return row
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
//
// 样式对齐 antd v5：32×32 页码、6px 圆角、选中主色底白字、非选中悬停
// 显示主色描边，首尾带上一页/下一页箭头。
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
	item := func(label string, active, disabled bool, onClick func()) ui.Element {
		btn := ui.ButtonBase(c).
			Height(t.Space(8)).
			MinWidth(t.Space(8)).
			Padding(0, t.Space(2)).
			Radius(t.Radius).
			Center()
		fg := t.Text
		if active {
			btn.Background(t.Accent)
			fg = t.AccentText
		} else {
			btn.Border(1, ui.Transparent)
			if btn.Hovered() && !disabled {
				btn.Border(1, t.Accent)
				fg = t.Accent
			}
		}
		if disabled {
			fg = t.TextMuted
		}
		btn.Children(func() {
			ui.Text(c, label).TextColor(fg).FontSize(t.FontSize).SingleLine()
		})
		if !disabled && btn.Clicked() && onClick != nil {
			onClick()
		}
		return btn
	}

	row := ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center)
	row.Children(func() {
		row.Children(func() {
			item("‹", false, *current <= 1, func() {
				*current--
				fire()
			})
		})
		pages := pageWindow(*current, total)
		for _, p := range pages {
			p := p
			if p == 0 { // 省略号
				box := ui.Box(c).Size(t.Space(8), t.Space(8)).Center()
				box.Children(func() {
					ui.Text(c, "…").TextColor(t.TextMuted).FontSize(t.FontSize)
				})
				continue
			}
			item(fmt.Sprintf("%d", p), p == *current, false, func() {
				*current = p
				fire()
			})
		}
		item("›", false, *current >= total, func() {
			*current++
			fire()
		})
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

	row := ui.Row(c).FillWidth().Padding(0, t.Space(4)).AlignItems(ui.Center)
	row.Children(func() {
		for i, label := range labels {
			if i > 0 {
				// 连接线：完成的实心，未完成的描边。
				line := ui.Box(c).Grow(1).Height(t.Space(0.5)).Border(1, t.Border)
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
				circle := ui.Box(c).Size(t.Space(8), t.Space(8)).Radius(t.Space(4)).Center()
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
						ui.Textf(c, "%d", i+1).TextColor(t.Accent).FontSize(t.Rem(1.14))
					})
				default:
					circle.Background(t.Background.Mix(t.Text, 0.05))
					circle.Children(func() {
						ui.Textf(c, "%d", i+1).TextColor(t.TextMuted).FontSize(t.Rem(1.14))
					})
				}
				color := t.TextMuted
				if now || done {
					color = t.Text
				}
				ui.Text(c, label).TextColor(color).FontSize(t.FontSize).SingleLine()
			})
		}
	})
	return row
}

// SidebarOption 定制 Sidebar。
type SidebarOption func(*sidebarConfig)

type sidebarConfig struct {
	width float32 // 主题单位数
}

// SidebarWidth 侧栏宽度（主题单位数，默认 48；随密度缩放）。
func SidebarWidth(units float32) SidebarOption {
	return func(c *sidebarConfig) { c.width = units }
}

// Sidebar 侧栏导航（antd Menu / Layout.Sider）：分区 + 条目 + 选中态，
// 复用 MyGo 的 Sidebar（键盘导航、滚动、选中项跟随）：
//
//	var chosen string
//	ant.Sidebar(c, &chosen, func() {
//	    ant.SidebarSection(c, "收件箱", nil, func() {
//	        ant.SidebarItem(c, "inbox", "收件箱")
//	        ant.SidebarItem(c, "sent", "已发送")
//	    })
//	})
func Sidebar(c *ui.Context, selected *string, fn func(), opts ...SidebarOption) ui.Element {
	cfg := &sidebarConfig{width: 48}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()
	return ui.Sidebar(c, selected, fn).
		Width(t.Space(cfg.width)).
		FillHeight()
}

// SidebarItem 侧栏条目（无图标）。
func SidebarItem(c *ui.Context, id, label string) ui.Element {
	return ui.SidebarItem(c, id, nil, label)
}

// SidebarItemWithIcon 带图标的侧栏条目。
func SidebarItemWithIcon(c *ui.Context, id, label string, svg *ui.SVG) ui.Element {
	return ui.SidebarItem(c, id, svg, label)
}

// SidebarSection 侧栏分区；open 为 nil 时始终展开。
func SidebarSection(c *ui.Context, title string, open *bool, fn func()) ui.Element {
	return ui.SidebarSection(c, title, open, fn)
}
