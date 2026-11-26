// Package divider 提供 Divider 控件的示例。
package divider

import (
	_ "embed"

	"github.com/egoist/mygo/ui"

	samples "github.com/lishuangquan1987/mygo-antdesign/samples"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

//go:embed demo.go
var demoSrc string

// Demo 返回 Divider 的示例集合。
func Demo() samples.WidgetDemo {
	return samples.WidgetDemo{
		Name:     "Divider",
		Package:  "widgets",
		Category: samples.CatGeneral,
		Summary:  "分割线：水平/垂直、实线/虚线、带文字（左/中/右）。",
		Samples: []samples.Sample{
			{
				Title:       "水平分割线",
				Description: "实线与虚线，上下留白随密度缩放。",
				Code:        samples.Extract(demoSrc, "horizontal"),
				View:        Horizontal,
			},
			{
				Title:       "带文字的分割线",
				Description: "文字可左 / 中 / 右对齐。",
				Code:        samples.Extract(demoSrc, "text"),
				View:        WithTextDemo,
			},
			{
				Title:       "垂直分割线",
				Description: "用于行内分组。",
				Code:        samples.Extract(demoSrc, "vertical"),
				View:        Vertical,
			},
		},
	}
}

// sample:horizontal
func Horizontal(c *ui.Context) {
	col := ui.Column(c).FillWidth()
	col.Children(func() {
		widgets.Paragraph(c, "上方内容")
		widgets.Divider(c)
		widgets.Paragraph(c, "虚线分割")
		widgets.Divider(c, widgets.DashedDivider())
		widgets.Paragraph(c, "下方内容")
	})
}

// /sample

// sample:text
func WithTextDemo(c *ui.Context) {
	col := ui.Column(c).FillWidth()
	col.Children(func() {
		widgets.Divider(c, widgets.WithText("分组一"), widgets.TextLeft())
		widgets.Paragraph(c, "内容一")
		widgets.Divider(c, widgets.WithText("分组二"), widgets.TextCenter())
		widgets.Paragraph(c, "内容二")
		widgets.Divider(c, widgets.WithText("分组三"), widgets.TextRight())
		widgets.Paragraph(c, "内容三")
	})
}

// /sample

// sample:vertical
func Vertical(c *ui.Context) {
	row := ui.Row(c).Gap(8).AlignItems(ui.Center).Height(32)
	row.Children(func() {
		widgets.Button(c, "新增")
		widgets.VerticalDivider(c)
		widgets.Button(c, "编辑")
		widgets.VerticalDivider(c)
		widgets.Button(c, "删除", widgets.Danger())
	})
}

// /sample
