// Package nav 提供导航类控件的示例。
package nav

import (
	_ "embed"

	"github.com/egoist/mygo/ui"

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
