// Package main 是 mygo-antdesign 的控件示例 Gallery（docs/07-samples-gallery.md）：
// 左侧控件目录、右侧示例卡片流（实时预览 + 真实代码），顶部可切换
// 主题（亮/暗）、密度（标准/紧凑）、品牌色，并支持亮暗并排对比。
//
// 本文件是应用骨架（A8 + F1/F2/F3）；示例注册在 main.go。
// 说明：gallery 是演示应用，允许少量展示用的 ui.Hex（不在静态检查范围）；
// 控件本体一律遵守样式铁律。
package main

import (
	"github.com/egoist/mygo/ui"

	samples "github.com/lishuangquan1987/mygo-antdesign/samples"
	"github.com/lishuangquan1987/mygo-antdesign/theme"
)

// 品牌色候选（antd 常用色系演示）。
var brandColors = []struct {
	name string
	hex  string
}{
	{"拂晓蓝", "#1677ff"},
	{"酱紫", "#7c3aed"},
	{"青柠", "#13c2c2"},
	{"洋红", "#eb2f96"},
}

// Gallery 持有展示应用的状态；View 是状态的函数，每次交互后下一帧重绘。
type Gallery struct {
	dark     bool
	density  theme.Density
	brand    string // 品牌色 hex，"" 表示 antd 默认
	compare  bool   // 亮暗并排对比
	selected int    // 当前展示的控件下标
}

// NewGallery 创建 Gallery。
func NewGallery() *Gallery {
	return &Gallery{density: theme.Comfortable, brand: brandColors[0].hex}
}

// options 汇总当前的品牌色与密度选项。
func (g *Gallery) options() []theme.Option {
	opts := []theme.Option{theme.WithDensity(g.density)}
	if g.brand != "" {
		opts = append(opts, theme.WithPrimary(ui.Hex(g.brand)))
	}
	return opts
}

// token 当前全局主题。
func (g *Gallery) token() *theme.AntToken {
	if g.dark {
		return theme.Dark(g.options()...)
	}
	return theme.Light(g.options()...)
}

// View 是 Gallery 的根视图。
func (g *Gallery) View(c *ui.Context) {
	theme.Apply(c, g.token())
	ui.Column(c).Fill().Children(func() {
		g.topBar(c)
		ui.Row(c).Fill().Grow(1).AlignItems(ui.Stretch).Children(func() {
			g.sidebar(c)
			g.content(c)
		})
	})
}

// toggler 控制条上的切换按钮：选中的实底高亮。返回值方式与 Button 一致
// （Clicked 查询 + 幂等赋值）。
func toggler(c *ui.Context, label string, on bool) ui.Element {
	t := c.Theme()
	el := ui.ButtonBase(c).
		Height(t.Space(7)).
		Padding(0, t.Space(3)).
		Radius(t.Radius).
		Center()
	if on {
		el.Background(t.Accent)
		el.Children(func() { ui.Text(c, label).TextColor(t.AccentText) })
	} else {
		el.Children(func() {
			ui.Text(c, label).TextColor(t.Text)
		})
	}
	return el
}

// topBar 顶部控制条：主题、密度、品牌色、并排对比。
func (g *Gallery) topBar(c *ui.Context) {
	t := c.Theme()
	bar := ui.Row(c).FillWidth().Padding(t.Space(2), t.Space(4)).Gap(t.Space(2)).
		Border(1, t.Border).BorderWidth(0, 0, 1, 0).
		AlignItems(ui.Center)
	bar.Children(func() {
		group := ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center)
		group.Children(func() {
			if toggler(c, "亮", !g.dark).Clicked() {
				g.dark = false
			}
			if toggler(c, "暗", g.dark).Clicked() {
				g.dark = true
			}
		})
		dens := ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center)
		dens.Children(func() {
			if toggler(c, "标准", g.density == theme.Comfortable).Clicked() {
				g.density = theme.Comfortable
			}
			if toggler(c, "紧凑", g.density == theme.Compact).Clicked() {
				g.density = theme.Compact
			}
		})
		// 品牌色色块。
		for _, bc := range brandColors {
			sw := ui.ButtonBase(c).Size(t.Space(4), t.Space(4)).Radius(t.Radius)
			sw.Background(ui.Hex(bc.hex))
			if g.brand == bc.hex {
				sw.Border(2, t.Text)
			}
			if sw.Clicked() {
				g.brand = bc.hex
			}
			sw.Children(func() {})
		}
		if toggler(c, "亮暗对比", g.compare).Clicked() {
			g.compare = !g.compare
		}
	})
}

// sidebar 左侧控件目录（可滚动），列出注册的全部控件。
func (g *Gallery) sidebar(c *ui.Context) {
	t := c.Theme()
	ui.Scroll(c).Width(t.Space(40)).FillHeight().
		Border(1, t.Border).BorderWidth(0, 1, 0, 0).
		Children(func() {
			col := ui.Column(c).Padding(t.Space(2)).Gap(t.Space(1)).FillWidth()
			col.Children(func() {
				ui.Text(c, "控件").TextColor(t.TextMuted).FontSize(t.Rem(0.9))
				for i, d := range samples.All() {
					item := ui.ButtonBase(c).FillWidth().
						Height(t.Space(7)).
						Padding(0, t.Space(2)).
						Radius(t.Radius).
						Center()
					if i == g.selected {
						item.Background(t.Selection)
						item.Children(func() {
							ui.Text(c, d.Name).TextColor(t.Accent)
						})
					} else {
						item.Children(func() {
							ui.Text(c, d.Name).TextColor(t.Text)
						})
					}
					if item.Clicked() {
						g.selected = i
					}
				}
			})
		})
}

// content 右侧内容区：当前控件的标题、简介与示例卡片流（可滚动）。
func (g *Gallery) content(c *ui.Context) {
	all := samples.All()
	if len(all) == 0 {
		ui.Text(c, "尚未注册任何示例（在 gallery/main.go 中调用 samples.Register）")
		return
	}
	if g.selected >= len(all) {
		g.selected = 0
	}
	d := all[g.selected]
	t := c.Theme()

	ui.Scroll(c).Fill().Children(func() {
		col := ui.Column(c).FillWidth().Padding(t.Space(4)).Gap(t.Space(4))
		col.Children(func() {
			ui.Text(c, d.Name).FontSize(t.Rem(1.6))
			ui.Text(c, d.Summary).TextColor(t.TextMuted)
			for _, s := range d.Samples {
				g.sampleCard(c, d, s)
			}
		})
	})
}

// sampleCard 一张示例卡片：标题、描述、实时预览（或亮暗并排对比）、代码。
func (g *Gallery) sampleCard(c *ui.Context, d samples.WidgetDemo, s samples.Sample) {
	t := c.Theme()
	card := ui.Column(c).FillWidth().Gap(t.Space(2)).
		Padding(t.Space(3)).
		Radius(t.Radius).
		Background(t.Surface)
	card.Children(func() {
		ui.Text(c, s.Title).FontSize(t.Rem(1.15))
		ui.Text(c, s.Description).TextColor(t.TextMuted).FontSize(t.Rem(0.9))

		preview := ui.Column(c).FillWidth().
			Padding(t.Space(3)).
			Radius(t.Radius).
			Background(t.Background)
		preview.Children(func() {
			if g.compare {
				row := ui.Row(c).Gap(t.Space(3)).FillWidth().AlignItems(ui.Stretch)
				row.Children(func() {
					theme.Sub(c, g.light()).Grow(1).Padding(t.Space(2)).Children(func() { s.View(c) })
					theme.Sub(c, g.darkToken()).Grow(1).Padding(t.Space(2)).Children(func() { s.View(c) })
				})
			} else {
				s.View(c)
			}
		})

		// 代码视图：真实源码片段，等宽字体用系统默认。
		code := ui.Column(c).FillWidth().
			Padding(t.Space(2)).
			Radius(t.Radius).
			Background(t.Background).
			MaxHeight(t.Space(30))
		code.Children(func() {
			ui.Text(c, s.Code).TextColor(t.TextMuted).FontSize(t.Rem(0.86))
		})
	})
}

// light/darkToken 并排对比用的两套主题（携带密度与品牌色）。
func (g *Gallery) light() *theme.AntToken     { return theme.Light(g.options()...) }
func (g *Gallery) darkToken() *theme.AntToken { return theme.Dark(g.options()...) }
