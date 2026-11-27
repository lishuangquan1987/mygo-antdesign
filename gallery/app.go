// Package main 是 mygo-antdesign 的控件示例 Gallery（docs/07-samples-gallery.md）：
// 左侧品牌区 + 搜索 + 分类目录，右侧标题区 + 示例卡片流（实时预览 + 真实代码），
// 顶部可切换主题（亮/暗）、密度（标准/紧凑）、品牌色，并支持亮暗 / 密度
// 两种并排对比模式。
//
// 本文件是应用骨架（A8 + F1/F2/F3）；示例注册在 main.go。
// 说明：gallery 是演示应用，允许少量展示用的 ui.Hex（不在静态检查范围）；
// 控件本体一律遵守样式铁律。
package main

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/icon"
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

// sampleCount 返回全部示例片段数量（侧栏底部统计）。
func sampleCount() int {
	n := 0
	for _, d := range samples.All() {
		n += len(d.Samples)
	}
	return n
}

// View 是 Gallery 的根视图。
func (g *Gallery) View(c *ui.Context) {
	theme.Apply(c, g.token())
	ui.Row(c).Fill().Children(func() {
		g.sidebar(c)
		ui.Column(c).Fill().Grow(1).Children(func() {
			g.header(c)
			g.content(c)
		})
	})
}

// segItem 是自绘分段控制器的一项（对齐 antd Segmented 的外观）。
type segItem struct {
	label   string
	on      bool
	onClick func()
}

// segmented 自绘分段控制器：浅灰底、选中项白底 + 轻阴影（antd Segmented）。
func segmented(c *ui.Context, items ...segItem) ui.Element {
	t := c.Theme()
	group := ui.Row(c).
		Padding(t.Space(0.5)).
		Radius(t.Radius).
		Background(t.Background).
		Border(1, t.Border).
		AlignItems(ui.Center)
	group.Children(func() {
		for _, it := range items {
			it := it
			btn := ui.ButtonBase(c).
				Height(t.Space(6)).
				Padding(0, t.Space(2)).
				Radius(t.Radius * 0.67).
				Center()
			if it.on {
				btn.Background(t.Surface)
				btn.Shadow(0, 1, 2, 0, t.Text.Alpha(0.08))
				btn.Children(func() { ui.Text(c, it.label).TextColor(t.Text).FontSize(t.Rem(0.82)) })
			} else {
				btn.Children(func() { ui.Text(c, it.label).TextColor(t.TextMuted).FontSize(t.Rem(0.82)) })
			}
			if btn.Clicked() && it.onClick != nil {
				it.onClick()
			}
		}
	})
	return group
}

// toolGroup 控制条上的一组「标签 + 控件」。
func toolGroup(c *ui.Context, label string, body func()) ui.Element {
	t := c.Theme()
	group := ui.Row(c).AlignItems(ui.Center).Gap(t.Space(1))
	group.Children(func() {
		ui.Text(c, label).TextColor(t.TextMuted).FontSize(t.Rem(0.75)).SingleLine()
		body()
	})
	return group
}

// sidebar 左侧：品牌、搜索、分类目录与统计。
func (g *Gallery) sidebar(c *ui.Context) {
	t := c.Theme()
	side := ui.Column(c).Width(t.Space(64)).FillHeight().
		Background(t.Surface).
		Border(1, t.Border).BorderWidth(0, 1, 0, 0)

	side.Children(func() {
		// 品牌区。
		brand := ui.Column(c).FillWidth().
			Padding(t.Space(5), t.Space(5), t.Space(4), t.Space(5)).
			Gap(t.Space(1))
		brand.Children(func() {
			row := ui.Row(c).AlignItems(ui.Center).Gap(t.Space(1))
			row.Children(func() {
				ui.Text(c, "mygo").FontWeight(700).FontSize(t.Rem(1.29)).TextColor(t.Text)
				ui.Text(c, "antdesign").FontWeight(700).FontSize(t.Rem(1.29)).TextColor(t.Accent)
			})
			ui.Text(c, "Ant Design · MyGo 原生 UI").TextColor(t.TextMuted).FontSize(t.Rem(0.72))
		})

		// 搜索框。
		search := ui.Row(c).FillWidth().
			Height(t.Space(8)).
			Margin(0, t.Space(4)).
			Padding(0, t.Space(2)).
			Radius(t.Radius).
			Background(t.Background).
			Border(1, t.Border).
			AlignItems(ui.Center).
			Gap(t.Space(1))
		search.Children(func() {
			icon.Icon(c, icon.Search, t.Space(4), icon.WithColor(t.TextMuted))
			ui.TextInput(c, &g.query).
				Grow(1).
				Background(t.Background).
				Placeholder("搜索控件 / 分类").
				PlaceholderColor(t.TextMuted).
				FontSize(t.Rem(0.86))
			if g.query != "" {
				clear := ui.ButtonBase(c).Size(t.Space(5), t.Space(5)).Center()
				clear.Children(func() {
					ui.Text(c, "✕").TextColor(t.TextMuted).FontSize(t.Rem(0.7))
				})
				if clear.Clicked() {
					g.query = ""
				}
			}
		})

		// 目录。
		ui.Scroll(c).Fill().Grow(1).Children(func() {
			nav := ui.Column(c).FillWidth().
				Padding(0, t.Space(2)).
				Gap(t.Space(0.5))
			nav.Children(func() {
				filtered := g.filtered()
				if len(filtered) == 0 {
					empty := ui.Column(c).FillWidth().Padding(t.Space(4), t.Space(2)).Gap(t.Space(1))
					empty.Children(func() {
						ui.Text(c, "未找到匹配的控件").TextColor(t.TextMuted).FontSize(t.Rem(0.86))
						ui.Text(c, "试试「Button」「表格」「反馈」").TextColor(t.TextMuted).FontSize(t.Rem(0.72))
					})
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
					head := ui.Row(c).FillWidth().
						Justify(ui.SpaceBetween).
						Padding(t.Space(2), t.Space(2), t.Space(1), t.Space(2))
					head.Children(func() {
						ui.Text(c, string(cat)).TextColor(t.TextMuted).FontSize(t.Rem(0.7)).FontWeight(600)
						ui.Textf(c, "%d 控件", len(group)).TextColor(t.TextMuted).FontSize(t.Rem(0.68))
					})
					for _, d := range group {
						d := d
						selected := d.Name == g.selected
						item := ui.ButtonBase(c).FillWidth().
							Height(t.Space(8)).
							Padding(0, t.Space(2)).
							Radius(t.Radius * 0.67).
							AlignItems(ui.Center).
							Gap(t.Space(2))
						if selected {
							item.Background(t.Selection)
						}
						item.Children(func() {
							bar := ui.Box(c).Width(t.Space(0.5)).Height(t.Space(4)).
								Radius(t.Space(0.25))
							if selected {
								bar.Background(t.Accent)
							} else {
								bar.Background(ui.Transparent)
							}
							bar.Children(func() {})
							ui.Text(c, d.Name).
								TextColor(map[bool]ui.Color{true: t.Accent, false: t.Text}[selected]).
								FontSize(t.Rem(0.9)).
								SingleLine().
								Grow(1)
							ui.Textf(c, "%d 例", len(d.Samples)).
								TextColor(t.TextMuted).
								FontSize(t.Rem(0.7))
						})
						if item.Clicked() {
							g.selected = d.Name
						}
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
			ui.Textf(c, "%d 个分组 · %d 个示例", len(samples.All()), sampleCount()).
				TextColor(t.TextMuted).FontSize(t.Rem(0.7))
			ui.Text(c, "MIT License · "+version).TextColor(t.TextMuted).FontSize(t.Rem(0.7))
		})
	})
}

// header 右侧标题区：面包屑、标题、简介与工具条。
func (g *Gallery) header(c *ui.Context) {
	d, ok := g.current()
	if !ok {
		d = samples.WidgetDemo{Name: "Gallery", Summary: "在 gallery/main.go 中调用 samples.Register 注册示例。"}
	}
	t := c.Theme()

	head := ui.Column(c).FillWidth().
		Padding(t.Space(4), t.Space(6)).
		Gap(t.Space(2)).
		Background(t.Surface).
		Border(1, t.Border).BorderWidth(0, 0, 1, 0)
	head.Children(func() {
		ui.Textf(c, "控件 / %s", d.Category).TextColor(t.TextMuted).FontSize(t.Rem(0.72))

		title := ui.Row(c).AlignItems(ui.Center).Gap(t.Space(2))
		title.Children(func() {
			ui.Text(c, d.Name).FontSize(t.Rem(1.71)).FontWeight(600).TextColor(t.Text)
			if len(d.Samples) > 0 {
				widgets.Tag(c, fmt.Sprintf("%d 个示例", len(d.Samples)), widgets.TagPrimary())
			}
		})

		ui.Text(c, d.Summary).TextColor(t.TextMuted).FontSize(t.Rem(0.9))

		// 工具条独占一行，避免与标题争宽后折行。
		tools := ui.Row(c).FillWidth().
			AlignItems(ui.Center).
			Gap(t.Space(3)).
			Margin(t.Space(1), 0, 0, 0)
		tools.Children(func() {
			toolGroup(c, "主题", func() {
				segmented(c,
					segItem{label: "亮", on: !g.dark, onClick: func() { g.dark = false }},
					segItem{label: "暗", on: g.dark, onClick: func() { g.dark = true }},
				)
			})
			widgets.VerticalDivider(c)
			toolGroup(c, "密度", func() {
				segmented(c,
					segItem{label: "标准", on: g.density == theme.Comfortable, onClick: func() { g.density = theme.Comfortable }},
					segItem{label: "紧凑", on: g.density == theme.Compact, onClick: func() { g.density = theme.Compact }},
				)
			})
			widgets.VerticalDivider(c)
			toolGroup(c, "品牌色", func() {
				swatches := ui.Row(c).AlignItems(ui.Center).Gap(t.Space(1))
				swatches.Children(func() {
					for _, bc := range brandColors {
						bc := bc
						sw := ui.ButtonBase(c).
							Size(t.Space(5), t.Space(5)).
							Radius(t.Space(2.5)).
							Background(ui.Hex(bc.hex))
						if g.brand == bc.hex {
							sw.Border(2, t.Text)
						}
						sw.Children(func() {})
						if sw.Clicked() {
							g.brand = bc.hex
						}
					}
				})
			})
			widgets.VerticalDivider(c)
			toolGroup(c, "对比", func() {
				segmented(c,
					segItem{label: "无", on: g.compare == compareNone, onClick: func() { g.compare = compareNone }},
					segItem{label: "亮暗对比", on: g.compare == compareTheme, onClick: func() {
						if g.compare == compareTheme {
							g.compare = compareNone
						} else {
							g.compare = compareTheme
						}
					}},
					segItem{label: "密度对比", on: g.compare == compareDensity, onClick: func() {
						if g.compare == compareDensity {
							g.compare = compareNone
						} else {
							g.compare = compareDensity
						}
					}},
				)
			})
		})
	})
}

// content 右侧内容区：示例卡片流（可滚动）。
func (g *Gallery) content(c *ui.Context) {
	d, ok := g.current()
	if !ok {
		ui.Column(c).Fill().Center().Children(func() {
			widgets.Empty(c, widgets.EmptyDescription("未找到匹配的控件，试试其他关键词"))
		})
		return
	}
	t := c.Theme()

	ui.Scroll(c).Fill().Children(func() {
		col := ui.Column(c).FillWidth().Padding(t.Space(6)).Gap(t.Space(5))
		col.Children(func() {
			for i, s := range d.Samples {
				g.sampleCard(c, i, s)
			}
		})
	})
}

// sampleCard 一张示例卡片：标题/描述 + 实时预览（或并排对比）+ 可复制代码。
func (g *Gallery) sampleCard(c *ui.Context, index int, s samples.Sample) {
	t := c.Theme()
	card := ui.Column(c).FillWidth().
		Background(t.Surface).
		Radius(t.Radius+t.Space(0.5)).
		Border(1, t.Border).
		Shadow(0, 1, 2, 0, t.Text.Alpha(0.04))

	expanded := ui.Local(card, "gallery.code."+s.Title, func() bool { return false })

	card.Children(func() {
		// 卡片头：序号 + 标题 + 描述 + 复制。
		head := ui.Row(c).FillWidth().
			Padding(t.Space(4), t.Space(6)).
			Justify(ui.SpaceBetween).
			AlignItems(ui.Center).
			Gap(t.Space(4)).
			Border(1, t.Border).BorderWidth(0, 0, 1, 0)
		head.Children(func() {
			left := ui.Row(c).AlignItems(ui.Center).Gap(t.Space(2)).Grow(1)
			left.Children(func() {
				ui.Textf(c, "%02d", index+1).TextColor(t.TextMuted).FontSize(t.Rem(0.72))
				title := ui.Column(c).Gap(t.Space(0.25))
				title.Children(func() {
					ui.Text(c, s.Title).FontSize(t.Rem(1.14)).FontWeight(600).TextColor(t.Text)
					ui.Text(c, s.Description).TextColor(t.TextMuted).FontSize(t.Rem(0.84))
				})
			})
			actions := ui.Row(c).AlignItems(ui.Center).Gap(t.Space(1))
			actions.Children(func() {
				if widgets.Button(c, "复制代码", widgets.Size(widgets.Small), widgets.Text()).Clicked() {
					c.WriteClipboard(s.Code)
					c.Toast("代码已复制")
				}
			})
		})

		// 预览区：浅灰底承托，四周 24px 留白（对齐 antd demo 区块）。
		preview := ui.Column(c).FillWidth().
			Padding(t.Space(6)).
			Background(t.Background)
		preview.Children(func() {
			switch g.compare {
			case compareTheme:
				row := ui.Row(c).FillWidth().Gap(t.Space(4)).AlignItems(ui.Stretch)
				row.Children(func() {
					g.comparePane(c, "亮色", false, g.density, s.View)
					g.comparePane(c, "暗色", true, g.density, s.View)
				})
			case compareDensity:
				row := ui.Row(c).FillWidth().Gap(t.Space(4)).AlignItems(ui.Stretch)
				row.Children(func() {
					g.comparePane(c, "标准密度", g.dark, theme.Comfortable, s.View)
					g.comparePane(c, "紧凑密度", g.dark, theme.Compact, s.View)
				})
			default:
				s.View(c)
			}
		})

		// 代码区：可展开/收起。
		code := ui.Column(c).FillWidth().
			Border(1, t.Border).BorderWidth(1, 0, 0, 0)
		code.Children(func() {
			codeHead := ui.Row(c).FillWidth().
				Padding(t.Space(1), t.Space(4)).
				Justify(ui.SpaceBetween).
				AlignItems(ui.Center)
			codeHead.Children(func() {
				ui.Text(c, "Go").TextColor(t.TextMuted).FontSize(t.Rem(0.7)).FontWeight(600)
				toggle := widgets.Button(c, map[bool]string{true: "收起代码", false: "展开代码"}[*expanded],
					widgets.Size(widgets.Small), widgets.Text())
				if toggle.Clicked() {
					*expanded = !*expanded
				}
			})
			// 代码正文：用 Scroll 承载，超出部分滚动而不是溢出卡片。
			body := ui.Scroll(c).FillWidth()
			if !*expanded {
				body.MaxHeight(t.Space(30))
			}
			body.Children(func() {
				ui.Text(c, s.Code).
					TextColor(t.TextMuted).
					FontSize(t.Rem(0.78)).
					Padding(0, t.Space(4), t.Space(3), t.Space(4))
			})
		})
	})
}

// comparePane 并排对比中的一栏：小标题 + 局部主题容器。
func (g *Gallery) comparePane(c *ui.Context, label string, dark bool, density theme.Density, view func(*ui.Context)) {
	t := c.Theme()
	pane := ui.Column(c).Grow(1).Gap(t.Space(1))
	pane.Children(func() {
		ui.Text(c, label).TextColor(t.TextMuted).FontSize(t.Rem(0.72))
		theme.Sub(c, g.tokenOf(dark, density)).
			FillWidth().
			Padding(t.Space(4)).
			Radius(t.Radius).
			Background(t.Surface).
			Border(1, t.Border).
			Children(func() { view(c) })
	})
}
