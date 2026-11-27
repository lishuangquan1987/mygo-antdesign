// Package native 提供「按 MyGo 原生命名补齐」的控件示例（docs/02 覆盖矩阵）。
package native

import (
	_ "embed"

	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/chart"
	"github.com/egoist/mygo/ui/editor"
	"github.com/egoist/mygo/ui/richtext"

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
		Summary:  "按 MyGo 命名补齐的控件：Toolbar、SplitButton、GridView、CheckboxGroup、ColorWell，以及 Markdown/CodeEditor/ChartPlot/RichTextEditor 子包适配。",
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
			{
				Title:       "Markdown 文档",
				Description: "MyGo ui/markdown 子包适配：标题、列表、行内代码与链接。",
				Code:        samples.Extract(demoSrc, "markdown"),
				View:        MarkdownDemo,
			},
			{
				Title:       "CodeEditor 代码编辑器",
				Description: "MyGo ui/editor 子包适配：语法高亮、行号、只读浏览。",
				Code:        samples.Extract(demoSrc, "codeeditor"),
				View:        CodeEditorDemo,
			},
			{
				Title:       "ChartPlot 坐标图",
				Description: "MyGo ui/chart 子包适配：轴刻度 + 绘图区，绘制折线。",
				Code:        samples.Extract(demoSrc, "chartplot"),
				View:        ChartPlotDemo,
			},
			{
				Title:       "RichTextEditor 富文本",
				Description: "MyGo ui/richtext 子包适配：段落、加粗与列表。",
				Code:        samples.Extract(demoSrc, "richtexteditor"),
				View:        RichTextEditorDemo,
			},
			{
				Title:       "RemoteImage 网络图片",
				Description: "MyGo ui/netimage 子包适配：占位图加载完成后淡入（依赖网络，不纳入逐字节基线）。",
				Code:        samples.Extract(demoSrc, "remoteimage"),
				View:        RemoteImageDemo,
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
	row := ui.Row(c).Gap(8)
	row.Children(func() {
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

// sample:markdown
func MarkdownDemo(c *ui.Context) {
	widgets.Markdown(c, "## 版本说明\n\n- 新增 `Toolbar` / `GridView`\n- 修正 **Select** 内边距\n\n> 详见 [MyGo 文档](https://github.com/egoist/mygo)\n", nil)
}

// /sample

// sample:codeeditor
func CodeEditorDemo(c *ui.Context) {
	t := c.Theme()
	host := ui.Column(c).FillWidth()
	st := ui.Local(host, "sample.codeeditor.state", func() editor.State {
		return editor.State{
			Buffer:   editor.NewBuffer("package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"hello mygo\")\n}\n"),
			ReadOnly: true,
			SoftWrap: true,
		}
	})
	host.Children(func() {
		widgets.CodeEditor(c, st).Height(t.Space(30)).FillWidth()
	})
}

// /sample

// sample:chartplot
func ChartPlotDemo(c *ui.Context) {
	t := c.Theme()
	axes := widgets.ChartAxes{
		Left: []widgets.ChartTick{
			{At: 0.0, Value: 0, Label: "0"},
			{At: 0.5, Value: 50, Label: "50"},
			{At: 1.0, Value: 100, Label: "100"},
		},
		Bottom: []widgets.ChartTick{
			{At: 0.0, Value: 0, Label: "周一"},
			{At: 0.5, Value: 3, Label: "周三"},
			{At: 1.0, Value: 6, Label: "周六"},
		},
	}
	widgets.ChartPlot(c, axes, func(area ui.Element) {
		area.Children(func() {
			area.Draw(func(p *ui.Painter, r ui.Rect) {
				p.Line(chart.X(r, 0.1), chart.Y(r, 0.15), chart.X(r, 0.9), chart.Y(r, 0.85), 2, t.Accent)
			})
		})
	}).Height(t.Space(40)).Width(t.Space(80))
}

// /sample

// sample:richtexteditor
func RichTextEditorDemo(c *ui.Context) {
	t := c.Theme()
	host := ui.Column(c).FillWidth()
	doc := ui.Local(host, "sample.richtext.doc", func() richtext.Document {
		text := richtext.Text{}
		text = text.Append("富文本示例\n", richtext.Attributes{Weight: 700})
		text = text.Append("支持段落、加粗与列表。\n", richtext.Attributes{})
		return *richtext.NewDocument(text)
	})
	host.Children(func() {
		widgets.RichTextEditor(c, doc, widgets.RichTextEditorOptions{ReadOnly: true}).
			Height(t.Space(24)).FillWidth()
	})
}

// /sample

// sample:remoteimage
func RemoteImageDemo(c *ui.Context) {
	t := c.Theme()
	ui.Row(c).Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
		widgets.RemoteImage(c, "https://avatars.githubusercontent.com/u/15365329?v=4",
			widgets.RemoteImageOptions{MaxSide: 128}).
			Width(t.Space(16)).Height(t.Space(16)).Radius(t.Space(8))
		ui.Text(c, "网络图片：占位图加载完成后淡入").TextColor(t.TextMuted).FontSize(t.Rem(0.86))
	})
}

// /sample
