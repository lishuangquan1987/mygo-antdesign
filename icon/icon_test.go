package icon_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/icon"
	"github.com/lishuangquan1987/mygo-antdesign/theme"
)

func TestIconRender(t *testing.T) {
	for _, tk := range []*theme.AntToken{
		theme.Light(),
		theme.Dark(),
		theme.Light(theme.WithDensity(theme.Compact)),
		theme.Dark(theme.WithDensity(theme.Compact)),
	} {
		token := tk
		view := func(c *ui.Context) {
			theme.Apply(c, token)
			row := ui.Row(c).Gap(8).AlignItems(ui.Center)
			row.Children(func() {
				icon.Icon(c, icon.Search, 16)
				icon.Icon(c, icon.Setting, 16, icon.WithColor(c.Theme().Accent))
				icon.Icon(c, icon.User, 24)
				icon.Icon(c, icon.Cloud, 16)
			})
		}
		tt := ui.NewTester(view, 200, 80)
		tt.Frame()
		if tt.FramePending() {
			t.Errorf("密度 %v 下图标持续请求帧", token.Density)
		}
	}
}

func TestIconInvalidSVG(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		icon.Icon(c, "<not-svg>", 16)
	}
	tt := ui.NewTester(view, 120, 60)
	tt.Frame()
	if !tt.HasText("?") {
		t.Error("非法 SVG 未显示占位符")
	}
}
