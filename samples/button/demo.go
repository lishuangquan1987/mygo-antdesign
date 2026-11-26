// Package button 提供 Button 控件的示例。
// 本文件的示例函数是 Gallery 实时渲染的代码，也是用户复制粘贴的源码——
// 展示片段由 // sample: 标记与 go:embed 提取，二者同源（docs/07 第 6 节）。
package button

import (
	_ "embed"

	"github.com/egoist/mygo/ui"

	samples "github.com/lishuangquan1987/mygo-antdesign/samples"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

// 包名约定：目录 samples/ 的包名是 sample，本文件统一以别名 samples 引用。

//go:embed demo.go
var demoSrc string

// Demo 返回 Button 的示例集合（Gallery 启动时注册）。
func Demo() samples.WidgetDemo {
	return samples.WidgetDemo{
		Name:     "Button",
		Package:  "widgets",
		Category: samples.CatGeneral,
		Summary:  "按钮：五种 type × danger × 三档尺寸，基于 ButtonBase。",
		Samples: []samples.Sample{
			{
				Title:       "基础用法",
				Description: "点击「保存」后出现反馈，演示 OnClick 与 ui.Local 状态。",
				Code:        samples.Extract(demoSrc, "basic"),
				View:        Basic,
			},
			{
				Title:       "按钮类型",
				Description: "五种 type 对应 antd 的 primary / default / dashed / text / link。",
				Code:        samples.Extract(demoSrc, "variants"),
				View:        Variants,
			},
			{
				Title:       "尺寸",
				Description: "大 / 中 / 小三档，随密度（标准/紧凑）自动缩放。",
				Code:        samples.Extract(demoSrc, "sizes"),
				View:        Sizes,
			},
			{
				Title:       "状态",
				Description: "危险、禁用、加载中与危险文字按钮。",
				Code:        samples.Extract(demoSrc, "states"),
				View:        States,
			},
		},
	}
}

// sample:basic
func Basic(c *ui.Context) {
	row := ui.Row(c).Gap(8).AlignItems(ui.Center)
	saved := ui.Local(row, "ant.sample.button.saved", func() bool { return false })
	row.Children(func() {
		widgets.Button(c, "保存", widgets.Primary(), widgets.OnClick(func() { *saved = true }))
		if *saved {
			ui.Text(c, "已保存").TextColor(c.Theme().Success)
		}
	})
}

// /sample

// sample:variants
func Variants(c *ui.Context) {
	ui.Row(c).Gap(8).Wrap().Children(func() {
		widgets.Button(c, "默认按钮")
		widgets.Button(c, "主要按钮", widgets.Primary())
		widgets.Button(c, "虚线按钮", widgets.Dashed())
		widgets.Button(c, "文字按钮", widgets.Text())
		widgets.Button(c, "链接按钮", widgets.Link())
	})
}

// /sample

// sample:sizes
func Sizes(c *ui.Context) {
	ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
		widgets.Button(c, "大按钮", widgets.Size(widgets.Large))
		widgets.Button(c, "中按钮")
		widgets.Button(c, "小按钮", widgets.Size(widgets.Small))
	})
}

// /sample

// sample:states
func States(c *ui.Context) {
	ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
		widgets.Button(c, "危险", widgets.Primary(), widgets.Danger())
		widgets.Button(c, "禁用", widgets.Disabled(true))
		widgets.Button(c, "加载中", widgets.Loading(true))
		widgets.Button(c, "删除", widgets.Text(), widgets.Danger())
	})
}

// /sample
