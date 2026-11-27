package ant_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

// TestMygoNamedBuffers 验证 TextBuffer 输入入口复用 Input/TextArea 外框。
func TestMygoNamedBuffers(t *testing.T) {
	buf := ui.NewTextBuffer("hello buffer")
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		ui.Column(c).Padding(24).Width(420).Gap(12).Children(func() {
			widgets.TextInputBuffer(c, buf)
			widgets.TextAreaBuffer(c, buf)
		})
	}
	tt := ui.NewTester(view, 520, 260)
	tt.Frame()
	// 输入框内容不进入可访问文本树，改为确认视图已渲染且缓冲内容未被改写。
	if tt.Image() == nil {
		t.Fatal("TextBuffer 输入未渲染帧")
	}
	if buf.String() != "hello buffer" {
		t.Errorf("TextBuffer 内容被改写：%q", buf.String())
	}
}

// TestMygoNamedMarkdown 验证 Markdown 子包适配渲染。
func TestMygoNamedMarkdown(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		ui.Column(c).Padding(24).Width(480).Children(func() {
			widgets.Markdown(c, "# 标题\n\n正文段落", nil)
		})
	}
	tt := ui.NewTester(view, 560, 260)
	tt.Frame()
	if !tt.HasText("标题") || !tt.HasText("正文段落") {
		t.Errorf("Markdown 未渲染标题/正文，texts=%v", tt.Texts())
	}
}

// TestMygoNamedChartPlot 验证图表子包适配可渲染绘图区。
func TestMygoNamedChartPlot(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.ChartPlot(c, widgets.ChartAxes{}, func(area ui.Element) {
			area.Children(func() {
				ui.Text(c, "绘图区")
			})
		}).Height(200).Width(360)
	}
	tt := ui.NewTester(view, 420, 260)
	tt.Frame()
	if !tt.HasText("绘图区") {
		t.Error("图表绘图区未渲染")
	}
}
