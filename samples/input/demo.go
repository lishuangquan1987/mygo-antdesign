// Package input 提供文本录入类控件的示例。
package input

import (
	_ "embed"
	"time"

	"github.com/egoist/mygo/ui"

	samples "github.com/lishuangquan1987/mygo-antdesign/samples"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

//go:embed demo.go
var demoSrc string

// Demo 返回文本录入控件的示例集合。
func Demo() samples.WidgetDemo {
	return samples.WidgetDemo{
		Name:     "输入框",
		Package:  "widgets",
		Category: samples.CatInput,
		Summary:  "Input / TextArea / SearchInput / NumberInput / DatePicker / TimePicker：聚焦描边、禁用、前缀后缀。",
		Samples: []samples.Sample{
			{
				Title:       "基础输入框",
				Description: "单行输入，占位提示；聚焦时描边切换为焦点色。",
				Code:        samples.Extract(demoSrc, "basic"),
				View:        Basic,
			},
			{
				Title:       "多行文本",
				Description: "TextArea 指定行数，带占位提示。",
				Code:        samples.Extract(demoSrc, "textarea"),
				View:        TextAreaDemo,
			},
			{
				Title:       "搜索框与数字输入",
				Description: "SearchInput 回车提交；NumberInput 限定范围与步长。",
				Code:        samples.Extract(demoSrc, "search"),
				View:        SearchDemo,
			},
			{
				Title:       "尺寸与状态",
				Description: "大 / 中 / 小三档；密码、禁用、带前缀。",
				Code:        samples.Extract(demoSrc, "states"),
				View:        States,
			},
			{
				Title:       "日期与时间",
				Description: "DatePicker 日期选择与 TimePicker 时间选择。",
				Code:        samples.Extract(demoSrc, "picker"),
				View:        PickerDemo,
			},
		},
	}
}

// sample:basic
func Basic(c *ui.Context) {
	col := ui.Column(c).FillWidth().Gap(8)
	name := ui.Local(col, "ant.sample.input.name", func() string { return "" })
	col.Children(func() {
		widgets.Input(c, name,
			widgets.PlaceholderText("请输入姓名"),
			widgets.InputLabel("姓名"))
		widgets.Input(c, name, widgets.InputDisabled(true), widgets.PlaceholderText("禁用状态"))
	})
}

// /sample

// sample:textarea
func TextAreaDemo(c *ui.Context) {
	col := ui.Column(c).FillWidth()
	bio := ui.Local(col, "ant.sample.input.bio", func() string { return "" })
	col.Children(func() {
		widgets.TextArea(c, bio, widgets.Rows(3), widgets.PlaceholderText("介绍一下自己"))
	})
}

// /sample

// sample:search
func SearchDemo(c *ui.Context) {
	col := ui.Column(c).FillWidth().Gap(8)
	q := ui.Local(col, "ant.sample.input.q", func() string { return "" })
	port := ui.Local(col, "ant.sample.input.port", func() float64 { return 8080 })
	col.Children(func() {
		widgets.SearchInput(c, q, widgets.PlaceholderText("搜索…"), widgets.InputLabel("搜索"), widgets.OnSubmit(func() {}))
		widgets.NumberInputV(c, port, 1, 65535, 1, widgets.InputLabel("端口"))
	})
}

// /sample

// sample:states
func States(c *ui.Context) {
	col := ui.Column(c).FillWidth().Gap(8)
	v := ui.Local(col, "ant.sample.input.states", func() string { return "" })
	col.Children(func() {
		widgets.Input(c, v, widgets.InputSize(widgets.Large), widgets.PlaceholderText("大号"))
		widgets.Input(c, v, widgets.PlaceholderText("默认"))
		widgets.Input(c, v, widgets.InputSize(widgets.Small), widgets.PasswordInput(), widgets.PlaceholderText("密码"))
		widgets.Input(c, v, widgets.PlaceholderText("金额"), widgets.InputPrefix(func(c *ui.Context) {
			ui.Text(c, "¥").TextColor(c.Theme().TextMuted)
		}))
	})
}

// /sample

// sample:picker
func PickerDemo(c *ui.Context) {
	col := ui.Column(c).FillWidth().Gap(8)
	date := ui.Local(col, "ant.sample.input.date", func() time.Time {
		return time.Date(2026, 11, 27, 0, 0, 0, 0, time.Local)
	})
	start := ui.Local(col, "ant.sample.input.time", func() time.Time {
		return time.Date(2026, 11, 27, 14, 30, 0, 0, time.Local)
	})
	col.Children(func() {
		widgets.DatePicker(c, date, widgets.InputLabel("日期"))
		widgets.TimePicker(c, start, widgets.InputLabel("开始时间"))
	})
}

// /sample
