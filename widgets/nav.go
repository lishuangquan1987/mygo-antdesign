package ant

import (
	"fmt"
	"sync"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/icon"
)

// Breadcrumb 面包屑（antd Breadcrumb）。chosen 为当前层级下标：
//
//	ant.Breadcrumb(c, &level, []string{"首页", "列表", "详情"})
func Breadcrumb(c *ui.Context, items []string, chosen *int) ui.Element {
	return ui.Breadcrumbs(c, items, chosen)
}

// TabsOption 定制 Tabs。
type TabsOption func(*tabsConfig)

type tabsConfig struct {
	onChange func(index int)
	card     bool
}

// TabsOnChange 页签切换回调（antd Tabs onChange）。
func TabsOnChange(fn func(index int)) TabsOption {
	return func(c *tabsConfig) { c.onChange = fn }
}

// TabsCard 卡片式页签（antd Tabs type="card"）：页签带描边与浅底，
// 选中项白底并盖住下边线。
func TabsCard() TabsOption { return func(c *tabsConfig) { c.card = true } }

// Tabs 页签（antd Tabs）。selected 为选中下标，Changed 报告切换：
//
//	ant.Tabs(c, &tab, []string{"常规", "外观", "高级"}, ant.TabsOnChange(app.load))
//
// 样式对齐 antd v5：页签间距 32px、选中项主色 + 2px 下划线、悬停变主色，
// 整行底部 1px 分隔线。
func Tabs(c *ui.Context, selected *int, labels []string, opts ...TabsOption) ui.Element {
	cfg := &tabsConfig{}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()
	row := ui.Row(c).FillWidth().
		Gap(t.Space(8)).
		Border(1, t.Border).
		BorderWidth(0, 0, 1, 0)
	if cfg.card {
		row.Gap(0).AlignItems(ui.End)
	}
	row.Children(func() {
		for i, label := range labels {
			i := i
			active := i == *selected
			var tab ui.Element
			if cfg.card {
				tab = ui.ButtonBase(c).
					Height(t.Space(8)).
					Padding(0, t.Space(4))
				tab.Border(1, t.Border)
				if active {
					tab.Background(t.Surface)
					tab.BorderWidth(0, 0, 0, 0)
					tab.Margin(0, 0, t.Space(-0.5), 0)
				} else {
					tab.Background(t.Background.Mix(t.Surface, 0.5))
				}
				tab.Children(func() {
					color := t.TextMuted
					if active {
						color = t.Text
					}
					ui.Text(c, label).TextColor(color).FontSize(t.FontSize).SingleLine()
				})
			} else {
				tab = ui.ButtonBase(c).Padding(t.Space(3), 0, 0, 0)
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
			}
			if tab.Clicked() {
				if *selected != i {
					*selected = i
					if cfg.onChange != nil {
						cfg.onChange(i)
					}
				}
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
	onChange  func(page int)
	disabled  bool
	showTotal bool
	pageSize  *int
}

// OnPageChange 页码变化回调。
func OnPageChange(fn func(page int)) PaginationOption {
	return func(c *paginationConfig) { c.onChange = fn }
}

// PaginationDisabled 禁用整个分页器（antd disabled）。
func PaginationDisabled(b bool) PaginationOption {
	return func(c *paginationConfig) { c.disabled = b }
}

// PaginationShowTotal 显示总数文案（antd showTotal「共 N 页」）。
func PaginationShowTotal() PaginationOption {
	return func(c *paginationConfig) { c.showTotal = true }
}

// PaginationPageSize 启用每页条数（antd showSizeChanger）：此时 total 视为
// 条目总数（而非页数），页数按 *pageSize 折算，并提供 10/20/50/100 的选择。
func PaginationPageSize(size *int) PaginationOption {
	return func(c *paginationConfig) { c.pageSize = size }
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

	// 页数：默认 total 即页数；启用每页条数时 total 视为条目总数。
	pageCount := total
	itemTotal := 0
	if cfg.pageSize != nil && *cfg.pageSize > 0 {
		itemTotal = total
		pageCount = (total + *cfg.pageSize - 1) / *cfg.pageSize
	}
	if pageCount < 1 {
		pageCount = 1
	}
	if *current < 1 {
		*current = 1
	}
	if *current > pageCount {
		*current = pageCount
	}

	fire := func() {
		if cfg.onChange != nil {
			cfg.onChange(*current)
		}
	}
	item := func(label string, active, disabled bool, onClick func()) ui.Element {
		disabled = disabled || cfg.disabled
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
			// antd disabled 分页整体置灰：文字变弱，选中项不再是主色实底
			//（否则看起来仍可点）。
			fg = t.TextMuted
			if active {
				btn.Background(t.Background.Mix(t.Surface, 0.5))
			}
		}
		btn.Children(func() {
			ui.Text(c, label).TextColor(fg).FontSize(t.FontSize).SingleLine()
		})
		if !disabled && btn.Clicked() && onClick != nil {
			onClick()
		}
		return btn
	}

	// 箭头页（‹ ›）：用 SVG 图标替代字符，保证在方形按钮内居中。
	itemIcon := func(svg string, disabled bool, onClick func()) {
		disabled = disabled || cfg.disabled
		btn := ui.ButtonBase(c).Size(t.Space(8), t.Space(8)).
			Padding(0, t.Space(2)).
			Radius(t.Radius).
			Center()
		fg := t.Text
		if btn.Hovered() && !disabled {
			btn.Border(1, t.Accent)
			fg = t.Accent
		}
		if disabled {
			fg = t.TextMuted
		}
		btn.Children(func() {
			icon.Icon(c, svg, t.Space(3), icon.WithColor(fg))
		})
		if !disabled && btn.Clicked() && onClick != nil {
			onClick()
		}
	}

	row := ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center)
	row.Children(func() {
		if cfg.showTotal {
			if itemTotal > 0 {
				ui.Textf(c, "共 %d 条", itemTotal).TextColor(t.TextMuted).
					FontSize(t.FontSize).SingleLine()
			} else {
				ui.Textf(c, "共 %d 页", pageCount).TextColor(t.TextMuted).
					FontSize(t.FontSize).SingleLine()
			}
		}
		row.Children(func() {
			itemIcon(icon.Left, *current <= 1, func() {
				*current--
				fire()
			})
		})
		pages := pageWindow(*current, pageCount)
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
		itemIcon(icon.Right, *current >= pageCount, func() {
			*current++
			fire()
		})
		// 每页条数选择（antd showSizeChanger）。
		if cfg.pageSize != nil && *cfg.pageSize > 0 {
			cur := fmt.Sprintf("%d 条/页", *cfg.pageSize)
			sel := Select(c, &cur, []string{"10 条/页", "20 条/页", "50 条/页", "100 条/页"})
			if cur != fmt.Sprintf("%d 条/页", *cfg.pageSize) {
				var n int
				if _, err := fmt.Sscanf(cur, "%d", &n); err == nil && n > 0 {
					*cfg.pageSize = n
					*current = 1
					fire()
				}
			}
			_ = sel
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

	// AlignItems(Start) + 连接线 MarginTop：让线落在圆心高度
	//（圆 Size(Space(8))，圆心距顶 Space(4)，线高 Space(0.5)）。
	row := ui.Row(c).FillWidth().Padding(0, t.Space(4)).AlignItems(ui.Start)
	row.Children(func() {
		for i, label := range labels {
			if i > 0 {
				// 连接线：完成的实心，未完成的描边。
				line := ui.Box(c).Grow(1).Height(t.Space(0.5)).
					Margin(t.Space(3.75), 0, 0, 0).
					Border(1, t.Border)
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
					circle.Label("已完成")
					circle.Children(func() {
						icon.Icon(c, icon.Check, t.Space(3.5), icon.WithColor(t.AccentText))
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

// sidebarSel 栈：Sidebar 构建期间记录当前 selected 指针，SidebarItem
// 据此实现 antd 选中态。MyGo 的 Sidebar 未选中焦点时用灰色底（AppKit
// 惯例），与 antd Menu 的 item-selected（主色浅底 + 主色文字）不符，
// 这里在包装层覆盖。
var sidebarSel struct {
	mu       sync.Mutex
	selected []*string
}

func sidebarSelPush(selected *string) {
	sidebarSel.mu.Lock()
	sidebarSel.selected = append(sidebarSel.selected, selected)
	sidebarSel.mu.Unlock()
}

func sidebarSelPop() {
	sidebarSel.mu.Lock()
	sidebarSel.selected = sidebarSel.selected[:len(sidebarSel.selected)-1]
	sidebarSel.mu.Unlock()
}

// sidebarSelChosen 返回当前构建中的侧栏是否选中了 id。
func sidebarSelChosen(id string) bool {
	sidebarSel.mu.Lock()
	defer sidebarSel.mu.Unlock()
	if n := len(sidebarSel.selected); n > 0 {
		return *sidebarSel.selected[n-1] == id
	}
	return false
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
	sidebarSelPush(selected)
	defer sidebarSelPop()
	return ui.Sidebar(c, selected, fn).
		Width(t.Space(cfg.width)).
		FillHeight()
}

// SidebarItem 侧栏条目（无图标）。选中态对齐 antd Menu：主色 12% 浅底 +
// 主色文字（亮色 ≈ #e6f4ff/#1677ff，暗色同理派生）。
func SidebarItem(c *ui.Context, id, label string) ui.Element {
	el := ui.SidebarItem(c, id, nil, label)
	return styleSidebarItem(c, el, id)
}

// SidebarItemWithIcon 带图标的侧栏条目。
func SidebarItemWithIcon(c *ui.Context, id, label string, svg *ui.SVG) ui.Element {
	el := ui.SidebarItem(c, id, svg, label)
	return styleSidebarItem(c, el, id)
}

// styleSidebarItem 覆盖 MyGo 侧栏未聚焦时的灰色选中底为 antd 主色浅底。
// 键盘聚焦时框架的焦点环仍会绘制，导航可用性不变。
func styleSidebarItem(c *ui.Context, el ui.Element, id string) ui.Element {
	if sidebarSelChosen(id) {
		t := c.Theme()
		el.Background(t.Accent.Mix(t.Surface, 0.88)).TextColor(t.Accent)
	}
	return el
}

// SidebarSection 侧栏分区；open 为 nil 时始终展开。
func SidebarSection(c *ui.Context, title string, open *bool, fn func()) ui.Element {
	return ui.SidebarSection(c, title, open, fn)
}
