package ant_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

func TestDividerVariants(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		col := ui.Column(c).Gap(8)
		col.Children(func() {
			widgets.Divider(c)
			widgets.Divider(c, widgets.DashedDivider())
			widgets.Divider(c, widgets.WithText("居中"), widgets.TextCenter())
			widgets.Divider(c, widgets.WithText("靠左"))
			widgets.Divider(c, widgets.WithText("靠右"), widgets.TextRight())
		})
	}
	tt := ui.NewTester(view, 360, 400)
	tt.Frame()
	for _, want := range []string{"居中", "靠左", "靠右"} {
		if !tt.HasText(want) {
			t.Errorf("分割线文字 %q 未渲染", want)
		}
	}
}

func TestVerticalDivider(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		row := ui.Row(c).Gap(8).AlignItems(ui.Center).Height(48)
		row.Children(func() {
			widgets.Paragraph(c, "左侧")
			widgets.VerticalDivider(c)
			widgets.Paragraph(c, "右侧")
			widgets.DashedVerticalDivider(c)
			widgets.Paragraph(c, "末尾")
		})
	}
	tt := ui.NewTester(view, 320, 120)
	tt.Frame()
	if !tt.HasText("左侧") || !tt.HasText("右侧") || !tt.HasText("末尾") {
		t.Error("垂直分割线两侧内容未渲染完整")
	}
}

func TestDividerFollowsThemeAndDensity(t *testing.T) {
	for _, tk := range []*theme.AntToken{
		theme.Light(),
		theme.Dark(),
		theme.Light(theme.WithDensity(theme.Compact)),
		theme.Dark(theme.WithDensity(theme.Compact)),
	} {
		token := tk
		view := func(c *ui.Context) {
			theme.Apply(c, token)
			widgets.Divider(c, widgets.WithText("分组"), widgets.TextLeft())
			widgets.Divider(c, widgets.DashedDivider())
		}
		tt := ui.NewTester(view, 280, 160)
		tt.Frame()
		if !tt.HasText("分组") {
			t.Errorf("密度 %v 下分割线未渲染完整", token.Density)
		}
	}
}
