// Package native 提供「按 MyGo 原生命名补齐」的控件示例（docs/02 覆盖矩阵）。
package native

import (
	_ "embed"

	"github.com/egoist/mygo/ui"

	samples "github.com/lishuangquan1987/mygo-antdesign/samples"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

//go:embed demo.go
var demoSrc string

// Demo 返回原生命名补全控件的示例集合。
func Demo() samples.WidgetDemo {
	return samples.WidgetDemo{
		Name:     "原生补全",
		Package:  "widgets",
		Category: samples.CatNative,
		Summary:  "按 MyGo 命名补齐的控件：Toolbar、SplitButton、GridView、CheckboxGroup、ColorWell。",
		Samples: []samples.Sample{
			{
				Title:       "Toolbar 工具条",
				Description: "MyGo 命名的工具条容器：按钮等控件无独立面容，悬停才高亮。",
				Code:        samples.Extract(demoSrc, "toolbar"),
				View:        ToolbarDemo,
			},
			{
				Title:       "SplitButton 分裂按钮",
				Description: "主操作 + 右侧箭头菜单，对齐 antd Dropdown.Button。",
				Code:        samples.Extract(demoSrc, "splitbutton"),
				View:        SplitButtonDemo,
			},
			{
				Title:       "GridView 网格视图",
				Description: "等宽网格 + 虚拟滚动，适合图标/卡片列表（antd List grid 近似）。",
				Code:        samples.Extract(demoSrc, "gridview"),
				View:        GridViewDemo,
			},
			{
				Title:       "CheckboxGroup 复选组",
				Description: "一个总复选框控制下方一组复选框，MyGo 原生命名。",
				Code:        samples.Extract(demoSrc, "checkboxgroup"),
				View:        CheckboxGroupDemo,
			},
			{
				Title:       "ColorWell 取色块",
				Description: "点击色块展开系统取色面板，颜色随主题。",
				Code:        samples.Extract(demoSrc, "colorwell"),
				View:        ColorWellDemo,
			},
		},
	}
}

// sample:toolbar
func ToolbarDemo(c *ui.Context) {
	t := c.Theme()
	ui.Column(c).Gap(t.Space(2)).Children(func() {
		widgets.Toolbar(c, func() {
			widgets.Button(c, "新建", widgets.Text())
			widgets.Button(c, "打开", widgets.Text())
			widgets.Button(c, "保存", widgets.Text())
			ui.Spacer(c)
			widgets.Button(c, "设置", widgets.Text())
		})
		ui.Text(c, "工具条内的按钮以文字按钮呈现，悬停才显示底色。").TextColor(t.TextMuted).FontSize(t.Rem(0.86))
	})
}

// /sample

// sample:splitbutton
func SplitButtonDemo(c *ui.Context) {
	save := widgets.SplitButton(c, "保存")
	save.Button.OnClick(func() {})
	save.Arrow.Menu(func(m *ui.Menu) {
		m.Item("另存为")
		m.Separator()
		m.Item("导出为 PDF")
	})
	widgets.PrimarySplitButton(c, "发布").Arrow.Menu(func(m *ui.Menu) {
		m.Item("发布到测试环境")
		m.Item("发布到生产环境")
	})
}

// /sample

// sample:gridview
func GridViewDemo(c *ui.Context) {
	t := c.Theme()
	var st ui.GridState
	ui.Column(c).Gap(t.Space(2)).Children(func() {
		widgets.GridView(c, &st, 12, t.Space(24), t.Space(20), func(i int) {
			box := ui.Column(c).FillWidth().FillHeight().
				Background(t.SurfaceHover).
				Radius(t.Radius).
				Center()
			box.Children(func() {
				ui.Textf(c, "文件 %d", i+1).TextColor(t.TextMuted)
			})
		}).Height(t.Space(40))
	})
}

// /sample

// sample:checkboxgroup
func CheckboxGroupDemo(c *ui.Context) {
	mail, chat, sms := true, false, false
	widgets.CheckboxGroup(c, "通知", func() {
		widgets.Checkbox(c, &mail, "邮件")
		widgets.Checkbox(c, &chat, "聊天")
		widgets.Checkbox(c, &sms, "短信")
	})
}

// /sample

// sample:colorwell
func ColorWellDemo(c *ui.Context) {
	t := c.Theme()
	color := t.Accent
	ui.Row(c).Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
		widgets.ColorWell(c, &color)
		ui.Text(c, "点击色块打开取色器").TextColor(t.TextMuted).FontSize(t.Rem(0.86))
	})
}

// /sample
