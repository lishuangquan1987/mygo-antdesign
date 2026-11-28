package main

import (
	"github.com/egoist/mygo/ui"
	"github.com/lishuangquan1987/mygo-antdesign/icon"
	samples "github.com/lishuangquan1987/mygo-antdesign/samples"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

// sidebar 左侧：品牌、搜索、分类目录与统计。
func (g *Gallery) sidebar(c *ui.Context) {
	t := c.Theme()
	side := ui.Column(c).Width(t.Space(64)).FillHeight().
		Background(t.Surface).
		Border(1, t.Border).BorderWidth(0, 1, 0, 0)

	side.Children(func() {
		// 品牌区：Logo 标 + 双色字标 + 副标语。
		brand := ui.Column(c).FillWidth().
			Padding(t.Space(5), t.Space(4), t.Space(4), t.Space(4)).
			Gap(t.Space(2))
		brand.Children(func() {
			row := ui.Row(c).AlignItems(ui.Center).Gap(t.Space(2))
			row.Children(func() {
				logo := ui.Column(c).Size(t.Space(7), t.Space(7)).
					Radius(t.Radius + t.Space(0.5)).
					Background(t.Accent).
					Center()
				logo.Children(func() {
					icon.Icon(c, icon.Appstore, t.Space(4), icon.WithColor(ui.Hex("#fff")))
				})
				ui.Text(c, "mygo").FontWeight(700).FontSize(t.Rem(1.21)).TextColor(t.Text)
				ui.Text(c, "antdesign").FontWeight(700).FontSize(t.Rem(1.21)).TextColor(t.Accent)
			})
			ui.Text(c, "Ant Design · MyGo 原生 UI").TextColor(t.TextMuted).FontSize(t.Rem(0.72))
		})

		// 搜索框：复用库自身的 Input（antd Input 规范：容器底描边、
		// 聚焦主色描边 + 2px 焦点环），前缀搜索图标，后缀 allowClear。
		// 宽度用显式值（侧栏 Space(64) - 左右各 Space(4)）：
		// FillWidth 百分比按父容器全宽解析，Margin/Padding 都约束不住，
		// 会右侧溢出、边框压到分隔线。
		widgets.Input(c, &g.query,
			widgets.InputPrefix(func(c *ui.Context) {
				icon.Icon(c, icon.Search, t.Space(4), icon.WithColor(t.TextMuted))
				pad := ui.Box(c).Width(t.Space(1))
				pad.Children(func() {})
			}),
			widgets.InputSuffix(func(c *ui.Context) {
				if g.query == "" {
					return
				}
				clr := ui.ButtonBase(c).Size(t.Space(4), t.Space(4)).Center().Radius(t.Space(2)).Label("清空搜索")
				clr.Children(func() {
					icon.Icon(c, icon.Close, t.Space(2.5), icon.WithColor(t.TextMuted))
				})
				if clr.Clicked() {
					g.query = ""
				}
			}),
			widgets.PlaceholderText("搜索控件 / 分类"),
			widgets.InputLabel("搜索控件"),
		).Width(t.Space(56)).Margin(0, t.Space(4))

		// 目录。
		ui.Scroll(c).Fill().Grow(1).Children(func() {
			nav := ui.Column(c).FillWidth().
				Padding(0, t.Space(1)).
				Gap(t.Space(0.5))
			nav.Children(func() {
				filtered := g.filtered()

				// 「开始」入口固定在目录最上方。
				g.navItem(c, homeName, homeName == g.selected, 0)

				if len(filtered) == 0 {
					empty := ui.Column(c).FillWidth().Padding(t.Space(4), t.Space(2)).Gap(t.Space(1))
					empty.Children(func() {
						ui.Text(c, "未找到匹配的控件").TextColor(t.TextMuted).FontSize(t.Rem(0.86))
						ui.Text(c, "试试「Button」「表格」「反馈」").TextColor(t.TextMuted).FontSize(t.Rem(0.72))
					})
					return
				}
				for _, cat := range categories {
					var group []samples.ControlPage
					for _, p := range filtered {
						if p.Category == cat {
							group = append(group, p)
						}
					}
					if len(group) == 0 {
						continue
					}
					head := ui.Row(c).FillWidth().
						Justify(ui.SpaceBetween).
						Padding(t.Space(2), t.Space(3), t.Space(1), t.Space(3))
					head.Children(func() {
						ui.Text(c, string(cat)).TextColor(t.TextMuted).FontSize(t.Rem(0.75)).FontWeight(600)
						ui.Textf(c, "%d 控件", len(group)).TextColor(t.TextMuted).FontSize(t.Rem(0.72))
					})
					for _, p := range group {
						g.navItem(c, p.Name, p.Name == g.selected, len(p.Samples))
					}
				}
			})
		})

		// 底部统计。
		footer := ui.Column(c).FillWidth().
			Padding(t.Space(3), t.Space(5)).
			Gap(t.Space(0.5)).
			Border(1, t.Border).BorderWidth(1, 0, 0, 0)
		footer.Children(func() {
			ui.Textf(c, "%d 个控件 · %d 个示例", len(samples.Pages()), samples.FragmentCount()).
				TextColor(t.TextMuted).FontSize(t.Rem(0.7))
			ui.Text(c, "MIT License · "+version).TextColor(t.TextMuted).FontSize(t.Rem(0.7))
		})
	})
}

// navItem 侧栏目录里的一项（对齐 antd Menu：40px 项高、8px 圆角、
// 选中浅主色底 + 主色文字，hover 出浅填充；无左侧指示条）。
func (g *Gallery) navItem(c *ui.Context, name string, selected bool, samples int) {
	t := c.Theme()
	item := ui.ButtonBase(c).FillWidth().
		Height(t.Space(10)).
		Padding(0, t.Space(3)).
		Radius(t.Radius + t.Space(0.5)).
		AlignItems(ui.Center)
	if selected {
		item.Background(t.Selection)
	} else if item.Hovered() {
		item.Background(t.SurfaceHover)
	}
	item.Children(func() {
		ui.Text(c, name).
			TextColor(map[bool]ui.Color{true: t.Accent, false: t.Text}[selected]).
			FontSize(t.Rem(0.9)).
			SingleLine().
			Grow(1)
		if samples > 0 {
			ui.Textf(c, "%d 例", samples).
				TextColor(t.TextMuted).
				FontSize(t.Rem(0.7))
		}
	})
	if item.Clicked() {
		g.selected = name
		g.codeTab = 0
	}
}
