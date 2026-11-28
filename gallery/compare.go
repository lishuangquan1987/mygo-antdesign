package main

import (
	"github.com/egoist/mygo/ui"
	"github.com/lishuangquan1987/mygo-antdesign/theme"
)

// comparePane 并排对比中的一栏：小标题 + 局部主题容器。
func (g *Gallery) comparePane(c *ui.Context, label string, dark bool, density theme.Density, view func(*ui.Context)) {
	t := c.Theme()
	tk := g.tokenOf(dark, density)
	pane := ui.Column(c).FillWidth().Grow(1).Gap(t.Space(1))
	pane.Children(func() {
		ui.Text(c, label).TextColor(t.TextMuted).FontSize(t.Rem(0.8)).FontWeight(500)
		// 子主题面板：底色/描边跟随各自主题（暗色面板用 BgContainer）。
		// 面板半高显示；内容超出时在面板内滚动（此前注释称 body 滚动兜底
		// 但 body 并无 Scroll，超高内容会被截断）。
		sub := theme.Sub(c, tk).
			FillWidth().Grow(1).
			Radius(t.Radius).
			Background(tk.BgContainer).
			Border(1, tk.Border)
		sub.Children(func() {
			ui.Scroll(c).Fill().Children(func() {
				inner := ui.Column(c).FillWidth().Padding(t.Space(3), t.Space(4))
				inner.Children(func() {
					view(c)
				})
			})
		})
	})
}
