// Package typography 提供 Typography 控件的示例。
package typography

import (
	_ "embed"

	"github.com/egoist/mygo/ui"

	samples "github.com/lishuangquan1987/mygo-antdesign/samples"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

//go:embed demo.go
var demoSrc string

// Demo 返回 Typography 的示例集合。
func Demo() samples.WidgetDemo {
	return samples.WidgetDemo{
		Name:     "Typography",
		Package:  "widgets",
		Category: samples.CatGeneral,
		Summary:  "排版：五级标题、正文变体（次要/加粗/删除线）、行内代码。",
		Samples: []samples.Sample{
			{
				Title:       "标题",
				Description: "H1–H5 五级标题，字号对齐 antd Typography.Title，随基准字号缩放。",
				Code:        samples.Extract(demoSrc, "titles"),
				View:        Titles,
			},
			{
				Title:       "文本与文本组件",
				Description: "次要、加粗、删除线、下划线、斜体与高亮等文本变体。",
				Code:        samples.Extract(demoSrc, "texts"),
				View:        Texts,
			},
			{
				Title:       "行内代码",
				Description: "Code 展示行内代码样式（容器底色 + 描边）。",
				Code:        samples.Extract(demoSrc, "code"),
				View:        Code,
			},
		},
	}
}

// sample:titles
func Titles(c *ui.Context) {
	col := ui.Column(c).Gap(4)
	col.Children(func() {
		widgets.Title(c, "一级标题", widgets.H1)
		widgets.Title(c, "二级标题", widgets.H2)
		widgets.Title(c, "三级标题", widgets.H3)
		widgets.Title(c, "四级标题", widgets.H4)
	})
}

// /sample

// sample:texts
func Texts(c *ui.Context) {
	col := ui.Column(c).Gap(6)
	col.Children(func() {
		widgets.Paragraph(c, "默认正文，颜色跟随主题主文本色。")
		widgets.Paragraph(c, "次要文本，用于辅助说明。", widgets.Secondary())
		widgets.Paragraph(c, "加粗文本", widgets.Strong())
		widgets.Paragraph(c, "删除文本", widgets.Delete())
		widgets.Paragraph(c, "下划线文本", widgets.Underline())
		widgets.Paragraph(c, "高亮文本", widgets.Marked())
	})
}

// /sample

// sample:code
func Code(c *ui.Context) {
	row := ui.Row(c).Gap(4).AlignItems(ui.Center)
	row.Children(func() {
		widgets.Paragraph(c, "运行")
		widgets.Code(c, "go run ./cmd/server")
		widgets.Paragraph(c, "启动服务。")
	})
}

// /sample
