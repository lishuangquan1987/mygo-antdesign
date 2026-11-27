// Package main 是 mygo-antdesign 的控件示例 Gallery（docs/07-samples-gallery.md）：
// 左侧控件目录（按分类分组 + 搜索）、右侧示例卡片流（实时预览 + 真实代码），
// 顶部可切换主题（亮/暗）、密度（标准/紧凑）、品牌色，并支持亮暗 / 密度
// 两种并排对比模式。
//
// 本文件是应用骨架（A8 + F1/F2/F3）；示例注册在 main.go。
// 说明：gallery 是演示应用，允许少量展示用的 ui.Hex（不在静态检查范围）；
// 控件本体一律遵守样式铁律。
package main

import (
	"strings"

	"github.com/egoist/mygo/ui"

	samples "github.com/lishuangquan1987/mygo-antdesign/samples"
	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
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

// 对比模式。
const (
	compareNone = iota
	compareTheme
	compareDensity
)

// 侧栏分类顺序。
var categories = []samples.Category{
	samples.CatGeneral, samples.CatInput, samples.CatDisplay,
	samples.CatNav, samples.CatFeedback, samples.CatExtra,
}

// Gallery 持有展示应用的状态；View 是状态的函数，每次交互后下一帧重绘。
type Gallery struct {
	dark     bool
	density  theme.Density
	brand    string // 品牌色 hex，"" 表示 antd 默认
	compare  int    // compareNone / compareTheme / compareDensity
	selected string // 当前展示的控件名
	query    string // 侧栏搜索词
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

// tokenOf 指定密度与亮暗的主题（并排对比用）。
func (g *Gallery) tokenOf(dark bool, d theme.Density) *theme.AntToken {
	opts := []theme.Option{theme.WithDensity(d)}
	if g.brand != "" {
		opts = append(opts, theme.WithPrimary(ui.Hex(g.brand)))
	}
	if dark {
		return theme.Dark(opts...)
	}
	return theme.Light(opts...)
}

// matches 搜索过滤：控件名或简介包含关键词（大小写不敏感）。
func matches(d samples.WidgetDemo, query string) bool {
	q := strings.TrimSpace(strings.ToLower(query))
	if q == "" {
		return true
	}
	return strings.Contains(strings.ToLower(d.Name), q) ||
		strings.Contains(strings.ToLower(d.Summary), q) ||
		strings.Contains(strings.ToLower(string(d.Category)), q)
}

// filtered 返回符合搜索条件的示例分组（保持注册顺序）。
func (g *Gallery) filtered() []samples.WidgetDemo {
	all := samples.All()
	out := make([]samples.WidgetDemo, 0, len(all))
	for _, d := range all {
		if matches(d, g.query) {
			out = append(out, d)
		}
	}
	return out
}

// current 返回当前展示的示例分组（搜索后回退到第一项）。
func (g *Gallery) current() (samples.WidgetDemo, bool) {
	list := g.filtered()
	if len(list) == 0 {
		return samples.WidgetDemo{}, false
	}
	for _, d := range list {
		if d.Name == g.selected {
			return d, true
		}
	}
	g.selected = list[0].Name
	return list[0], true
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

// topBar 顶部控制条：主题、密度、品牌色、两种并排对比。
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
		if toggler(c, "亮暗对比", g.compare == compareTheme).Clicked() {
			if g.compare == compareTheme {
				g.compare = compareNone
			} else {
				g.compare = compareTheme
			}
		}
		if toggler(c, "密度对比", g.compare == compareDensity).Clicked() {
			if g.compare == compareDensity {
				g.compare = compareNone
			} else {
				g.compare = compareDensity
			}
		}
	})
}

// sidebar 左侧控件目录：搜索 + 按分类分组，显示每个分组的示例数。
func (g *Gallery) sidebar(c *ui.Context) {
	t := c.Theme()
	ui.Scroll(c).Width(t.Space(46)).FillHeight().
		Border(1, t.Border).BorderWidth(0, 1, 0, 0).
		Children(func() {
			col := ui.Column(c).Padding(t.Space(2)).Gap(t.Space(2)).FillWidth()
			col.Children(func() {
				search := ui.Row(c).FillWidth().
					Height(t.Space(7)).
					Padding(0, t.Space(2)).
					Radius(t.Radius).
					Background(t.Surface).
					Border(1, t.Border).
					AlignItems(ui.Center)
				search.Children(func() {
					ui.TextInput(c, &g.query).
						Grow(1).
						Background(t.Surface).
						Placeholder("搜索控件…").
						PlaceholderColor(t.TextMuted).
						FontSize(t.Rem(0.86))
				})

				filtered := g.filtered()
				if len(filtered) == 0 {
					ui.Text(c, "未找到匹配的控件").TextColor(t.TextMuted).FontSize(t.Rem(0.86))
					return
				}
				for _, cat := range categories {
					var group []samples.WidgetDemo
					for _, d := range filtered {
						if d.Category == cat {
							group = append(group, d)
						}
					}
					if len(group) == 0 {
						continue
					}
					ui.Text(c, string(cat)).TextColor(t.TextMuted).FontSize(t.Rem(0.8)).
						Padding(t.Space(1), t.Space(2), 0, t.Space(2))
					for _, d := range group {
						d := d
						item := ui.ButtonBase(c).FillWidth().
							Height(t.Space(7)).
							Padding(0, t.Space(2)).
							Radius(t.Radius)
						if d.Name == g.selected {
							item.Background(t.Selection)
						}
						item.Children(func() {
							row := ui.Row(c).FillWidth().Justify(ui.SpaceBetween).AlignItems(ui.Center)
							row.Children(func() {
								color := t.Text
								if d.Name == g.selected {
									color = t.Accent
								}
								ui.Text(c, d.Name).TextColor(color).FontSize(t.Rem(0.9)).SingleLine()
								ui.Textf(c, "%d", len(d.Samples)).
									TextColor(t.TextMuted).FontSize(t.Rem(0.75))
							})
						})
						if item.Clicked() {
							g.selected = d.Name
						}
					}
				}
			})
		})
}

// content 右侧内容区：当前控件的标题、简介与示例卡片流（可滚动）。
func (g *Gallery) content(c *ui.Context) {
	d, ok := g.current()
	if !ok {
		ui.Text(c, "尚未注册任何示例（在 gallery/main.go 中调用 samples.Register）")
		return
	}
	t := c.Theme()

	ui.Scroll(c).Fill().Children(func() {
		col := ui.Column(c).FillWidth().Padding(t.Space(4)).Gap(t.Space(4))
		col.Children(func() {
			ui.Text(c, d.Name).FontSize(t.Rem(1.6))
			ui.Text(c, d.Summary).TextColor(t.TextMuted)
			for _, s := range d.Samples {
				g.sampleCard(c, s)
			}
		})
	})
}

// sampleCard 一张示例卡片：标题、描述、实时预览（或并排对比）、代码 + 复制。
func (g *Gallery) sampleCard(c *ui.Context, s samples.Sample) {
	t := c.Theme()
	card := ui.Column(c).FillWidth().Gap(t.Space(2)).
		Padding(t.Space(3)).
		Radius(t.Radius).
		Background(t.Surface)
	card.Children(func() {
		head := ui.Row(c).FillWidth().Justify(ui.SpaceBetween).AlignItems(ui.Center)
		head.Children(func() {
			ui.Text(c, s.Title).FontSize(t.Rem(1.15))
			if widgets.Button(c, "复制代码", widgets.Size(widgets.Small), widgets.Text()).Clicked() {
				c.WriteClipboard(s.Code)
				c.Toast("代码已复制")
			}
		})
		ui.Text(c, s.Description).TextColor(t.TextMuted).FontSize(t.Rem(0.9))

		preview := ui.Column(c).FillWidth().
			Padding(t.Space(3)).
			Radius(t.Radius).
			Background(t.Background)
		preview.Children(func() {
			switch g.compare {
			case compareTheme:
				row := ui.Row(c).Gap(t.Space(3)).FillWidth().AlignItems(ui.Stretch)
				row.Children(func() {
					theme.Sub(c, g.tokenOf(false, g.density)).Grow(1).Padding(t.Space(2)).
						Children(func() { s.View(c) })
					theme.Sub(c, g.tokenOf(true, g.density)).Grow(1).Padding(t.Space(2)).
						Children(func() { s.View(c) })
				})
			case compareDensity:
				row := ui.Row(c).Gap(t.Space(3)).FillWidth().AlignItems(ui.Stretch)
				row.Children(func() {
					theme.Sub(c, g.tokenOf(g.dark, theme.Comfortable)).Grow(1).Padding(t.Space(2)).
						Children(func() { s.View(c) })
					theme.Sub(c, g.tokenOf(g.dark, theme.Compact)).Grow(1).Padding(t.Space(2)).
						Children(func() { s.View(c) })
				})
			default:
				s.View(c)
			}
		})

		// 代码视图：真实源码片段。
		code := ui.Column(c).FillWidth().
			Padding(t.Space(2)).
			Radius(t.Radius).
			Background(t.Background).
			MaxHeight(t.Space(24))
		code.Children(func() {
			ui.Text(c, s.Code).TextColor(t.TextMuted).FontSize(t.Rem(0.8))
		})
	})
}
