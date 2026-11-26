package ant_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

func TestTitleRendersAllLevels(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		col := ui.Column(c).Gap(8)
		col.Children(func() {
			widgets.Title(c, "一级标题", widgets.H1)
			widgets.Title(c, "二级标题", widgets.H2)
			widgets.Title(c, "三级标题", widgets.H3)
			widgets.Title(c, "四级标题", widgets.H4)
			widgets.Title(c, "五级标题", widgets.H5)
		})
	}
	tt := ui.NewTester(view, 360, 480)
	tt.Frame()
	for _, want := range []string{"一级标题", "二级标题", "三级标题", "四级标题", "五级标题"} {
		if !tt.HasText(want) {
			t.Errorf("标题 %q 未渲染", want)
		}
	}
}

func TestParagraphVariants(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		col := ui.Column(c).Gap(4)
		col.Children(func() {
			widgets.Paragraph(c, "默认文本")
			widgets.Paragraph(c, "次要文本", widgets.Secondary())
			widgets.Paragraph(c, "加粗文本", widgets.Strong())
			widgets.Paragraph(c, "删除文本", widgets.Delete())
			widgets.Paragraph(c, "下划线", widgets.Underline())
			widgets.Code(c, "go run .")
		})
	}
	tt := ui.NewTester(view, 320, 300)
	tt.Frame()
	for _, want := range []string{"默认文本", "次要文本", "加粗文本", "删除文本", "下划线", "go run ."} {
		if !tt.HasText(want) {
			t.Errorf("文本 %q 未渲染", want)
		}
	}
}

func TestTypographyThemeAndDensity(t *testing.T) {
	// 4 种主题组合下排版控件渲染不 panic。
	for _, tk := range []*theme.AntToken{
		theme.Light(),
		theme.Dark(),
		theme.Light(theme.WithDensity(theme.Compact)),
		theme.Dark(theme.WithDensity(theme.Compact)),
	} {
		token := tk
		view := func(c *ui.Context) {
			theme.Apply(c, token)
			widgets.Title(c, "标题", widgets.H2)
			widgets.Paragraph(c, "正文", widgets.Secondary())
			widgets.Code(c, "code")
		}
		tt := ui.NewTester(view, 280, 160)
		tt.Frame()
		if !tt.HasText("标题") || !tt.HasText("正文") || !tt.HasText("code") {
			t.Errorf("密度 %v 下排版控件未渲染完整", token.Density)
		}
	}
}
