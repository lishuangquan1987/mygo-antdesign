// Package ide 提供「IDE / 上位机」风格的紧凑密度示例（docs/06 E3）：
// 在同一个示例里用 theme.Sub 演示紧凑密度下的侧栏、标签页、代码区与状态栏。
package ide

import (
	_ "embed"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/icon"
	samples "github.com/lishuangquan1987/mygo-antdesign/samples"
	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

//go:embed demo.go
var demoSrc string

// Demo 返回 IDE 风格示例。
func Demo() samples.WidgetDemo {
	return samples.WidgetDemo{
		Name:     "IDE 风格",
		Package:  "samples",
		Category: samples.CatGeneral,
		Summary:  "紧凑密度落地上位机 / IDE 场景：侧栏 + 标签页 + 代码区 + 状态栏。",
		Samples: []samples.Sample{
			{
				Title:       "紧凑型 IDE 布局",
				Description: "theme.Sub + Compact 局部紧凑：文件树、编辑区、大纲、状态栏。",
				Code:        samples.Extract(demoSrc, "layout"),
				View:        Layout,
				Fill:        true,
			},
		},
	}
}

// sample:layout
func Layout(c *ui.Context) {
	// 局部紧凑密度：跟随当前亮暗，只覆盖密度与字号。
	tk := theme.Light(theme.WithDensity(theme.Compact))
	if c.Theme().Dark {
		tk = theme.Dark(theme.WithDensity(theme.Compact))
	}
	theme.Sub(c, tk).Fill().Children(func() {
		t := c.Theme()
		host := ui.Column(c).Fill()
		tab := ui.Local(host, "ant.sample.ide.tab", func() int { return 0 })
		host.Children(func() {
			ui.Column(c).Fill().Border(1, t.Border).Radius(t.Radius).Children(func() {
				// 主体：左文件树 + 中编辑区 + 右大纲（纵向撑满剩余空间）。
				ui.Row(c).FillWidth().Grow(1).AlignItems(ui.Stretch).Children(func() {
					// 左侧文件树。
					left := ui.Column(c).Width(t.Space(34)).
						Padding(0, t.Space(1)).
						Background(t.Surface).Border(1, t.Border)
					left.BorderWidth(0, 1, 0, 0)
					left.Children(func() {
						ui.Text(c, "资源管理器").TextColor(t.TextMuted).FontSize(t.Rem(0.75)).
							Padding(t.Space(1), t.Space(1))
						for _, f := range []struct {
							name string
							icon string
						}{
							{"main.go", icon.File},
							{"theme.go", icon.File},
							{"widgets.go", icon.File},
						} {
							row := ui.ButtonBase(c).FillWidth().Height(t.Space(6)).
								Padding(0, t.Space(1)).Radius(t.Radius * 0.5).
								Gap(t.Space(1.5)).AlignItems(ui.Center)
							if f.name == "theme.go" {
								row.Background(t.Selection)
							}
							row.Children(func() {
								icon.Icon(c, f.icon, t.Space(3), icon.WithColor(t.TextMuted))
								ui.Text(c, f.name).TextColor(t.Text).FontSize(t.Rem(0.8))
							})
						}
					})

					// 中间编辑区。
					main := ui.Column(c).Grow(1).Background(t.Surface)
					main.Children(func() {
						tabs := ui.Row(c).Gap(t.Space(1)).
							Padding(t.Space(1), t.Space(2)).
							Background(t.Background).Border(1, t.Border)
						tabs.BorderWidth(0, 0, 1, 0)
						tabs.Children(func() {
							widgets.Tabs(c, tab, "main.go", "theme.go", "widgets.go")
						})
						code := ui.Column(c).Fill().Padding(t.Space(2)).Gap(t.Space(1))
						code.Children(func() {
							for i, line := range []string{
								"package main",
								"",
								"func main() {",
								"    theme.Apply(c, theme.Light())",
								"    widgets.Button(c, \"保存\", widgets.Primary())",
								"}",
							} {
								row := ui.Row(c).Gap(t.Space(2))
								row.Children(func() {
									ui.Textf(c, "%d", i+1).TextColor(t.TextMuted).FontSize(t.Rem(0.7))
									ui.Text(c, line).TextColor(t.Text).FontSize(t.Rem(0.78))
								})
							}
						})
					})

					// 右侧大纲。
					right := ui.Column(c).Width(t.Space(28)).
						Padding(t.Space(1)).Gap(t.Space(1)).
						Background(t.Surface).Border(1, t.Border)
					right.BorderWidth(0, 0, 0, 1)
					right.Children(func() {
						ui.Text(c, "大纲").TextColor(t.TextMuted).FontSize(t.Rem(0.75))
						for _, s := range []string{"main()", "Layout()", "render()"} {
							ui.Text(c, s).TextColor(t.Text).FontSize(t.Rem(0.8))
						}
					})
				})

				// 底部状态栏。
				status := ui.Row(c).FillWidth().Height(t.Space(6)).
					Padding(0, t.Space(2)).
					Justify(ui.SpaceBetween).AlignItems(ui.Center).
					Background(t.Accent)
				status.Children(func() {
					ui.Text(c, "Go 1.27.4").TextColor(t.AccentText).FontSize(t.Rem(0.75))
					ui.Text(c, "紧凑密度 · Compact").TextColor(t.AccentText).FontSize(t.Rem(0.75))
				})
			})
		})
	})
}

// /sample
