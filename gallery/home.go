package main

import (
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
	"github.com/lishuangquan1987/mygo-antdesign/icon"
	samples "github.com/lishuangquan1987/mygo-antdesign/samples"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

// homePage 首页：品牌主视觉 + 特性卡片 + 统计 + 快速开始。
func (g *Gallery) homePage(c *ui.Context) {
	t := c.Theme()
	controls := len(samples.Pages())
	fragments := samples.FragmentCount()

	ui.Scroll(c).Fill().Children(func() {
		col := ui.Column(c).FillWidth().Padding(t.Space(10), t.Space(14)).Gap(t.Space(8))

		// 主视觉：品牌色浅底 Hero 横幅（对齐 antd 首页 hero 形制）。
		col.Children(func() {
			hero := ui.Column(c).FillWidth().
				Padding(t.Space(10), t.Space(9)).
				Radius(t.Radius+t.Space(0.5)).
				Background(t.Accent.Alpha(0.05)).
				Border(1, t.Accent.Alpha(0.14)).
				Gap(t.Space(3))
			hero.Children(func() {
				meta := ui.Row(c)
				meta.Children(func() {
					widgets.Tag(c, "v"+version+" · MIT · 开箱即用", widgets.TagPrimary())
				})
				row := ui.Row(c).AlignItems(ui.Center).Gap(t.Space(2))
				row.Children(func() {
					ui.Text(c, "mygo").FontWeight(700).FontSize(t.Rem(2.86)).TextColor(t.Text)
					ui.Text(c, "antdesign").FontWeight(700).FontSize(t.Rem(2.86)).TextColor(t.Accent)
				})
				ui.Text(c, "Ant Design 风格的 MyGo 原生控件库").FontSize(t.Rem(1.29)).TextColor(t.TextMuted)
				ui.Text(c, "为 MyGo 原生 UI 提供 antd v5 视觉与交互，补齐 antd 有而 MyGo 没有的控件；亮暗双主题、标准/紧凑双密度，全部控件配真实代码示例与截图基线。").
					TextColor(t.TextMuted).FontSize(t.Rem(1.0))
				actions := ui.Row(c).AlignItems(ui.Center).Gap(t.Space(3)).Margin(t.Space(2), 0, 0, 0)
				actions.Children(func() {
					if widgets.Button(c, "浏览控件", widgets.Primary(), widgets.Size(widgets.Large)).Clicked() {
						g.selected = "Button"
					}
					widgets.Hyperlink(c, "GitHub 仓库", "https://github.com/lishuangquan1987/mygo-antdesign").
						FontSize(t.Rem(1.0))
				})
			})
		})

		// 统计数值：antd Statistic 左对齐形制。
		col.Children(func() {
			row := ui.Row(c).FillWidth().Gap(t.Space(4))
			row.Children(func() {
				for _, st := range []struct {
					title string
					value float64
				}{
					{"控件构造函数", 144},
					{"控件示例页", float64(controls)},
					{"示例片段", float64(fragments)},
					{"截图基线", 324},
				} {
					card := ui.Column(c).Grow(1).
						Padding(t.Space(6), t.Space(5)).
						Radius(t.Radius+t.Space(0.5)).
						Background(t.Surface).
						Border(1, t.Border).
						Gap(t.Space(1))
					card.Children(func() {
						widgets.Statistic(c, st.title, st.value)
					})
				}
			})
		})

		// 特性卡片：图标色块 + 2×2 等宽网格。
		col.Children(func() {
			row := ui.Row(c).FillWidth().Gap(t.Space(4)).Wrap()
			row.Children(func() {
				for _, f := range []struct {
					icon  string
					title string
					desc  string
				}{
					{icon.Eye, "主题切换", "亮/暗、品牌色、字号运行时切换，所有控件同步生效。"},
					{icon.Menu, "双密度", "标准 32px / 紧凑 28px 控件高度，适配上位机与 IDE 场景。"},
					{icon.Appstore, "缺口补齐", "antd 有而 MyGo 没有的 35 个控件全部自建实现。"},
					{icon.Copy, "视觉回归", "3 张 Gallery + 324 张逐示例截图基线，CI 逐字节比对。"},
				} {
					card := ui.Column(c).Grow(1).Width(t.Space(96)).
						Padding(t.Space(6)).
						Radius(t.Radius+t.Space(0.5)).
						Background(t.Surface).
						Border(1, t.Border).
						Gap(t.Space(2))
					card.Children(func() {
						head := ui.Row(c).AlignItems(ui.Center).Gap(t.Space(2))
						head.Children(func() {
							tile := ui.Column(c).Size(t.Space(7), t.Space(7)).
								Radius(t.Radius).
								Background(t.Accent.Alpha(0.1)).
								Center()
							tile.Children(func() {
								icon.Icon(c, f.icon, t.Space(4), icon.WithColor(t.Accent))
							})
							ui.Text(c, f.title).FontSize(t.Rem(1.14)).FontWeight(600).TextColor(t.Text)
						})
						ui.Text(c, f.desc).TextColor(t.TextMuted).FontSize(t.Rem(0.9))
					})
				}
			})
		})

		// 快速开始：语法高亮的只读代码视图（预览前 14 行，完整示例见 README）。
		col.Children(func() {
			code := ui.Column(c).FillWidth().
				Radius(t.Radius+t.Space(0.5)).
				Background(t.Surface).
				Border(1, t.Border).
				Shadow(0, 1, 2, 0, t.Text.Alpha(0.04))
			code.Children(func() {
				head := ui.Row(c).FillWidth().
					Padding(t.Space(3), t.Space(6)).
					Justify(ui.SpaceBetween).
					AlignItems(ui.Center).
					Border(1, t.Border).BorderWidth(0, 0, 1, 0)
				head.Children(func() {
					ui.Text(c, "快速开始").FontSize(t.Rem(1.14)).FontWeight(600).TextColor(t.Text)
					ui.Text(c, "Go").TextColor(t.TextMuted).FontSize(t.Rem(0.75)).FontWeight(600).
						Padding(t.Space(0.5), t.Space(2)).
						Radius(t.Space(2)).
						Background(t.Background)
				})
				body := ui.Column(c).FillWidth().Background(t.Background)
				body.Children(func() {
					st := ui.Local(code, "gallery.quickstart", func() editor.State {
						buf := editor.NewBuffer(truncateCode(quickStartCode, 14, "完整示例见 README 快速开始"))
						buf.SetLanguage(editor.LanguageNamed("go"))
						return editor.State{
							Buffer:       buf,
							ReadOnly:     true,
							NoActiveLine: true,
							FitContent:   true,
						}
					})
					widgets.CodeEditor(c, st).FillWidth()
				})
			})
		})
	})
}
