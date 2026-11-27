// Package nav 提供导航类控件的示例。
package nav

import (
	_ "embed"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/icon"
	samples "github.com/lishuangquan1987/mygo-antdesign/samples"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

//go:embed demo.go
var demoSrc string

// Demo 返回导航控件的示例集合。
func Demo() samples.WidgetDemo {
	return samples.WidgetDemo{
		Name:     "导航",
		Package:  "widgets",
		Category: samples.CatNav,
		Summary:  "Breadcrumb / Tabs / Dropdown / Pagination / Steps。",
		Samples: []samples.Sample{
			{
				Title:       "面包屑与页签",
				Description: "Breadcrumb 层级导航（最后一项为当前页）；Tabs 页签切换。",
				Code:        samples.Extract(demoSrc, "breadcrumb"),
				View:        BreadcrumbDemo,
			},
			{
				Title:       "下拉菜单",
				Description: "Dropdown 菜单项、分隔线与危险项。",
				Code:        samples.Extract(demoSrc, "dropdown"),
				View:        DropdownDemo,
			},
			{
				Title:       "分页器",
				Description: "超过 7 页自动折叠为「1 2 3 4 … N」窗口。",
				Code:        samples.Extract(demoSrc, "pagination"),
				View:        PaginationDemo,
			},
			{
				Title:       "步骤条",
				Description: "已完成（✓）/ 进行中 / 等待三种状态。",
				Code:        samples.Extract(demoSrc, "steps"),
				View:        StepsDemo,
			},
			{
				Title:       "侧栏导航",
				Description: "Sidebar 分区 + 条目 + 选中态，键盘可导航。",
				Code:        samples.Extract(demoSrc, "sidebar"),
				View:        SidebarDemo,
			},
			{
				Title:       "带图标的侧栏",
				Description: "SidebarItemWithIcon 在条目左侧显示图标。",
				Code:        samples.Extract(demoSrc, "sidebaricons"),
				View:        SidebarIconsDemo,
			},
		},
	}
}

// sample:breadcrumb
func BreadcrumbDemo(c *ui.Context) {
	col := ui.Column(c).FillWidth().Gap(12)
	level := ui.Local(col, "ant.sample.nav.level", func() int { return -1 })
	tab := ui.Local(col, "ant.sample.nav.tab", func() int { return 0 })
	col.Children(func() {
		widgets.Breadcrumb(c, []string{"首页", "列表", "详情"}, level)
		widgets.Tabs(c, tab, "常规", "外观", "高级")
	})
}

// /sample

// sample:dropdown
func DropdownDemo(c *ui.Context) {
	row := ui.Row(c).Gap(8)
	picked := ui.Local(row, "ant.sample.nav.picked", func() string { return "" })
	row.Children(func() {
		widgets.Dropdown(c, "更多操作", []widgets.MenuItemData{
			{Label: "编辑", OnSelect: func() { *picked = "已编辑" }},
			{Label: "复制", OnSelect: func() { *picked = "已复制" }},
			{Separator: true},
			{Label: "删除", Danger: true, OnSelect: func() { *picked = "已删除" }},
		})
		if *picked != "" {
			widgets.Tag(c, *picked, widgets.TagSuccess())
		}
	})
}

// /sample

// sample:pagination
func PaginationDemo(c *ui.Context) {
	col := ui.Column(c).FillWidth()
	page := ui.Local(col, "ant.sample.nav.page", func() int { return 1 })
	col.Children(func() {
		widgets.Pagination(c, page, 42, widgets.OnPageChange(func(int) {}))
	})
}

// /sample

// sample:steps
func StepsDemo(c *ui.Context) {
	col := ui.Column(c).FillWidth().Gap(12)
	step := ui.Local(col, "ant.sample.nav.step", func() int { return 1 })
	col.Children(func() {
		widgets.Steps(c, step, []string{"填写信息", "确认订单", "支付"})
		widgets.Space(c, widgets.GapUnits(1)).Children(func() {
			widgets.Button(c, "上一步", widgets.OnClick(func() {
				if *step > 0 {
					*step--
				}
			}))
			widgets.Button(c, "下一步", widgets.Primary(), widgets.OnClick(func() {
				if *step < 2 {
					*step++
				}
			}))
		})
	})
}

// /sample

// sample:sidebar
func SidebarDemo(c *ui.Context) {
	host := ui.Column(c).FillWidth().Height(c.Theme().Space(44))
	chosen := ui.Local(host, "ant.sample.nav.sidebar", func() string { return "inbox" })
	host.Children(func() {
		widgets.Sidebar(c, chosen, func() {
			widgets.SidebarSection(c, "邮件", nil, func() {
				widgets.SidebarItem(c, "inbox", "收件箱")
				widgets.SidebarItem(c, "sent", "已发送")
				widgets.SidebarItem(c, "draft", "草稿箱")
			})
		})
	})
	widgets.Paragraph(c, "当前选择："+*chosen, widgets.Secondary())
}

// /sample

// sample:sidebaricons
func SidebarIconsDemo(c *ui.Context) {
	host := ui.Column(c).FillWidth().Height(c.Theme().Space(40))
	selected := ui.Local(host, "ant.sample.nav.sidebaricon", func() string { return "home" })
	openSection := ui.Local(host, "ant.sample.nav.sidebarsection", func() bool { return true })
	home, errHome := ui.ParseSVG([]byte(icon.Home))
	setting, errSetting := ui.ParseSVG([]byte(icon.Setting))
	if errHome != nil || errSetting != nil {
		return
	}
	host.Children(func() {
		widgets.Sidebar(c, selected, func() {
			widgets.SidebarSection(c, "导航", openSection, func() {
				widgets.SidebarItemWithIcon(c, "home", "首页", home)
				widgets.SidebarItemWithIcon(c, "settings", "设置", setting)
			})
		})
	})
}

// /sample
