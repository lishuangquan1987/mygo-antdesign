package ant_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

func TestSpaceLayouts(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		col := ui.Column(c).Gap(8)
		col.Children(func() {
			widgets.Space(c).Children(func() {
				widgets.Button(c, "取消")
				widgets.Button(c, "确定", widgets.Primary())
			})
			widgets.Space(c, widgets.Gap(widgets.SpaceSmall)).Children(func() {
				widgets.Button(c, "小间距一")
				widgets.Button(c, "小间距二")
			})
			widgets.Space(c, widgets.SpaceVertical(), widgets.SpaceFillWidth()).Children(func() {
				widgets.Paragraph(c, "纵向一")
				widgets.Paragraph(c, "纵向二")
			})
		})
	}
	tt := ui.NewTester(view, 360, 360)
	tt.Frame()
	for _, want := range []string{"取消", "确定", "小间距一", "纵向一", "纵向二"} {
		if !tt.HasText(want) {
			t.Errorf("Space 子元素 %q 未渲染", want)
		}
	}
}

func TestFlexLayouts(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		col := ui.Column(c).Gap(8).FillWidth()
		col.Children(func() {
			flex := widgets.Flex(c, widgets.FlexJustify(ui.SpaceBetween), widgets.FlexAlign(ui.Center), widgets.FlexFill())
			flex.Children(func() {
				widgets.Title(c, "左标题", widgets.H4)
				widgets.Button(c, "右按钮", widgets.Primary())
			})
			wrap := widgets.Flex(c, widgets.FlexWrap(), widgets.FlexGap(2), widgets.FlexFill())
			wrap.Children(func() {
				for _, label := range []string{"换行一", "换行二", "换行三", "换行四"} {
					widgets.Button(c, label, widgets.Text())
				}
			})
		})
	}
	tt := ui.NewTester(view, 320, 300)
	tt.Frame()
	if !tt.HasText("左标题") || !tt.HasText("右按钮") || !tt.HasText("换行一") {
		t.Error("Flex 子元素未渲染完整")
	}
}

func TestGridLayout(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		grid := widgets.Grid(c, 2)
		grid.Children(func() {
			for _, name := range []string{"格子一", "格子二", "格子三", "格子四"} {
				widgets.Paragraph(c, name)
			}
		})
	}
	tt := ui.NewTester(view, 320, 200)
	tt.Frame()
	for _, want := range []string{"格子一", "格子二", "格子三", "格子四"} {
		if !tt.HasText(want) {
			t.Errorf("Grid 子元素 %q 未渲染", want)
		}
	}
}

func TestLayoutFollowsThemeAndDensity(t *testing.T) {
	for _, tk := range []*theme.AntToken{
		theme.Light(),
		theme.Dark(),
		theme.Light(theme.WithDensity(theme.Compact)),
		theme.Dark(theme.WithDensity(theme.Compact)),
	} {
		token := tk
		view := func(c *ui.Context) {
			theme.Apply(c, token)
			widgets.Space(c, widgets.Gap(widgets.SpaceLarge)).Children(func() {
				widgets.Button(c, "间距")
			})
			widgets.Flex(c, widgets.FlexColumn(), widgets.FlexGap(1)).Children(func() {
				widgets.Paragraph(c, "纵向")
			})
			widgets.Grid(c, 2).Children(func() {
				widgets.Paragraph(c, "网格")
			})
		}
		tt := ui.NewTester(view, 320, 220)
		tt.Frame()
		if !tt.HasText("间距") || !tt.HasText("纵向") || !tt.HasText("网格") {
			t.Errorf("密度 %v 下布局控件未渲染完整", token.Density)
		}
	}
}
