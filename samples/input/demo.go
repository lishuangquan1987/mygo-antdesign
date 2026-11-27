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
			{
				Title:       "表单",
				Description: "Form 标签列对齐，FormItem 支持必填、帮助与错误文字。",
				Code:        samples.Extract(demoSrc, "form"),
				View:        FormDemo,
			},
			{
				Title:       "自动补全",
				Description: "AutoComplete 保留键入内容，同时给出候选建议。",
				Code:        samples.Extract(demoSrc, "autocomplete"),
				View:        AutoCompleteDemo,
			},
			{
				Title:       "颜色选择",
				Description: "ColorPicker 色块 + 展开面板，也可直接用常驻 ColorPanel。",
				Code:        samples.Extract(demoSrc, "color"),
				View:        ColorDemo,
			},
			{
				Title:       "更多录入控件",
				Description: "Hyperlink / DecimalInput / SearchField / Stepper / TokenField / EditableText / Calendar。",
				Code:        samples.Extract(demoSrc, "moreinput"),
				View:        MoreInputDemo,
			},
			{
				Title:       "提及",
				Description: "Mentions 输入 @ 弹出候选，选中后插入。",
				Code:        samples.Extract(demoSrc, "mentions"),
				View:        MentionsDemo,
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

// sample:form
func FormDemo(c *ui.Context) {
	name := ""
	email := ""
	widgets.Form(c, func() {
		widgets.FormItem(c, "用户名", func() {
			widgets.Input(c, &name, widgets.PlaceholderText("请输入用户名"))
		}, widgets.FormItemRequired(true))
		widgets.FormItem(c, "邮箱", func() {
			widgets.Input(c, &email, widgets.PlaceholderText("name@example.com"))
		}, widgets.FormItemHelp("用于接收通知"), widgets.FormItemError("邮箱格式不正确"))
		widgets.FormGroup(c, "其他", func() {
			widgets.FormItem(c, "备注", func() {
				widgets.TextArea(c, &name, widgets.Rows(2))
			})
		})
	})
}

// /sample

// sample:autocomplete
func AutoCompleteDemo(c *ui.Context) {
	host := ui.Column(c).FillWidth().Gap(c.Theme().Space(2))
	city := ui.Local(host, "ant.sample.input.city", func() string { return "" })
	host.Children(func() {
		widgets.AutoComplete(c, city, []string{"北京", "上海", "广州", "深圳"})
		widgets.Paragraph(c, "当前值："+*city, widgets.Secondary())
	})
}

// /sample

// sample:color
func ColorDemo(c *ui.Context) {
	host := ui.Column(c).FillWidth().Gap(c.Theme().Space(3))
	// 初始值取主题主色：亮色 #1677ff、暗色 #1668dc，避免暗色下演示
	// 硬编码亮色品牌色。
	col := ui.Local(host, "ant.sample.input.color", func() ui.Color { return c.Theme().Accent })
	host.Children(func() {
		widgets.ColorPicker(c, col)
		widgets.Paragraph(c, "色块可展开取色面板；也可直接使用常驻面板：", widgets.Secondary())
		widgets.ColorPanel(c, col)
	})
}

// /sample

// sample:moreinput
func MoreInputDemo(c *ui.Context) {
	host := ui.Column(c).FillWidth().Gap(c.Theme().Space(2))
	amount := ui.Local(host, "ant.sample.input.amount", func() string { return "12.30" })
	query := ui.Local(host, "ant.sample.input.searchfield", func() string { return "" })
	count := ui.Local(host, "ant.sample.input.count", func() float64 { return 3 })
	name := ui.Local(host, "ant.sample.input.editable", func() string { return "重命名我" })
	tags := ui.Local(host, "ant.sample.input.tags", func() []string { return []string{"Go"} })
	date := ui.Local(host, "ant.sample.input.calendar", func() time.Time {
		return time.Date(2026, 11, 27, 0, 0, 0, 0, time.Local)
	})
	host.Children(func() {
		widgets.Hyperlink(c, "查看文档", "https://example.com")
		widgets.DecimalInput(c, amount, "0", "1000000", "0.01")
		widgets.SearchField(c, query)
		widgets.Stepper(c, count, 0, 100, 1)
		widgets.TokenField(c, tags, []string{"Go", "GUI", "Ant Design"})
		widgets.EditableText(c, name)
		widgets.Calendar(c, date)
	})
}

// /sample

// sample:mentions
func MentionsDemo(c *ui.Context) {
	host := ui.Column(c).FillWidth().Gap(c.Theme().Space(2))
	reply := ui.Local(host, "ant.sample.input.mentions", func() string { return "" })
	host.Children(func() {
		widgets.Mentions(c, reply, []string{"张三", "李四", "王五"})
		widgets.Paragraph(c, "输入 @ 触发候选；当前值："+*reply, widgets.Secondary())
	})
}

// /sample
