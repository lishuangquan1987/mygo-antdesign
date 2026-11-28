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
	"github.com/egoist/mygo/ui"
	"github.com/lishuangquan1987/mygo-antdesign/icon"
	samples "github.com/lishuangquan1987/mygo-antdesign/samples"
	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
	"strings"
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
	brand    string  // 品牌色 hex，"" 表示 antd 默认
	compare  int     // compareNone / compareTheme / compareDensity
	fontSize float32 // 0 表示跟随密度（标准 14 / 紧凑 13）
	size     gallerySize
	selected string // 当前展示的控件名
	codeTab  int    // 代码区当前展示的示例下标
	query    string // 侧栏搜索词
}

// gallerySize Gallery 的全局尺寸档（docs/06 §10：控件尺寸统一注入）。
type gallerySize int

const (
	sizeMedium gallerySize = iota
	sizeSmall
	sizeLarge
)

// NewGallery 创建 Gallery。
func NewGallery() *Gallery {
	// brand 为空表示跟随 antd 预设主色（亮 #1677ff / 暗 #1668dc），
	// 不固定为拂晓蓝，否则暗色下主色无法切换到规范值。
	return &Gallery{density: theme.Comfortable, selected: homeName}
}

// homeName 首页在侧栏中的名字（g.selected == homeName 时展示首页）。
const homeName = "开始"

// fontScale 把尺寸、字号与品牌色并入主题选项（0 表示跟随密度）。
func (g *Gallery) fontScale(d theme.Density) []theme.Option {
	opts := []theme.Option{theme.WithDensity(d)}
	// 全局尺寸：在小/大档用主题度量整体缩放，弥补各控件尺寸选项无统一注入点。
	switch g.size {
	case sizeSmall:
		opts = append(opts, theme.WithSpacing(3), theme.WithRadius(4), theme.WithFontSize(12))
	case sizeLarge:
		opts = append(opts, theme.WithSpacing(5), theme.WithRadius(8), theme.WithFontSize(16))
	}
	if g.fontSize > 0 {
		opts = append(opts, theme.WithFontSize(g.fontSize))
	}
	if g.brand != "" {
		opts = append(opts, theme.WithPrimary(ui.Hex(g.brand)))
	}
	return opts
}

// options 汇总当前的品牌色、密度与字号选项。
func (g *Gallery) options() []theme.Option {
	return g.fontScale(g.density)
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
	opts := g.fontScale(d)
	if dark {
		return theme.Dark(opts...)
	}
	return theme.Light(opts...)
}

// filtered 返回符合搜索条件的控件页（保持侧栏顺序）。
func (g *Gallery) filtered() []samples.ControlPage {
	all := samples.Pages()
	out := make([]samples.ControlPage, 0, len(all))
	for _, p := range all {
		if samples.MatchPage(p, g.query) {
			out = append(out, p)
		}
	}
	return out
}

// current 返回当前展示的控件页（搜索后回退到第一项）。
func (g *Gallery) current() (samples.ControlPage, bool) {
	list := g.filtered()
	if len(list) == 0 {
		return samples.ControlPage{}, false
	}
	for _, p := range list {
		if p.Name == g.selected {
			return p, true
		}
	}
	g.selected = list[0].Name
	return list[0], true
}

// View 是 Gallery 的根视图。
func (g *Gallery) View(c *ui.Context) {
	theme.Apply(c, g.token())
	ui.Row(c).Fill().Children(func() {
		g.sidebar(c)
		// 内容区显式用页面底色（antd colorBgLayout：亮 #f5f5f5 / 暗 #000），
		// 让面板容器（Surface：白 / #141414）形成层级对比。
		ui.Column(c).Fill().Grow(1).Background(c.Theme().Background).Children(func() {
			if g.selected == homeName {
				g.homePage(c)
				return
			}
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

// segmented 自绘分段控制器：浅灰承托底（无边框）、2px 内边距、
// 28px 项高，选中项容器底 + 1px 描边 + 轻阴影（对齐 antd Segmented）。
func segmented(c *ui.Context, items ...segItem) ui.Element {
	t := c.Theme()
	group := ui.Row(c).
		Padding(t.Space(0.5)).
		Radius(t.Radius + t.Space(0.5)).
		Background(t.Background).
		AlignItems(ui.Center)
	group.Children(func() {
		for _, it := range items {
			it := it
			btn := ui.ButtonBase(c).
				Height(t.Space(7)).
				Padding(0, t.Space(3)).
				Radius(t.Radius * 0.67).
				Center()
			if it.on {
				btn.Background(t.Surface)
				btn.Border(1, t.Border)
				btn.Shadow(0, 1, 2, 0, t.Text.Alpha(0.08))
				btn.Children(func() { ui.Text(c, it.label).TextColor(t.Text).FontSize(t.Rem(0.86)).FontWeight(500) })
			} else {
				btn.Children(func() { ui.Text(c, it.label).TextColor(t.TextMuted).FontSize(t.Rem(0.86)) })
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

// header 右侧标题区：面包屑、标题、简介与工具条。
func (g *Gallery) header(c *ui.Context) {
	d, ok := g.current()
	if !ok {
		d = samples.ControlPage{Name: "Gallery", Summary: "在 gallery/main.go 中调用 samples.Register 注册示例。"}
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
		g.tools(c)
	})
}

// truncateCode 把代码按整行截断到 maxLines 行，并追加一行提示注释；
// 配合 editor.State.FitContent 使用，高度永远与内容行数一致，
// 不出现半行截断，也不依赖具体行高像素。
func truncateCode(code string, maxLines int, hint string) string {
	ls := strings.Split(code, "\n")
	if len(ls) <= maxLines {
		return code
	}
	return strings.Join(ls[:maxLines], "\n") + "\n\n// … " + hint
}

// tools 主题/密度/尺寸/字号/品牌色/对比 的全局工具条（首页与控件页共用）。
// 固定两行：第一行主题/密度/尺寸/字号，第二行品牌色/对比——
// 不用单行 Wrap，窄窗口下换行行会被裁剪。
func (g *Gallery) tools(c *ui.Context) {
	t := c.Theme()
	tools := ui.Column(c).FillWidth().
		Gap(t.Space(2)).
		Margin(t.Space(1), 0, 0, 0)
	tools.Children(func() {
		row1 := ui.Row(c).AlignItems(ui.Center).Gap(t.Space(3))
		row1.Children(func() {
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
			toolGroup(c, "尺寸", func() {
				segmented(c,
					segItem{label: "小", on: g.size == sizeSmall, onClick: func() { g.size = sizeSmall }},
					segItem{label: "中", on: g.size == sizeMedium, onClick: func() { g.size = sizeMedium }},
					segItem{label: "大", on: g.size == sizeLarge, onClick: func() { g.size = sizeLarge }},
				)
			})
			widgets.VerticalDivider(c)
			toolGroup(c, "字号", func() {
				segmented(c,
					segItem{label: "12", on: g.fontSize == 12, onClick: func() { g.fontSize = 12 }},
					segItem{label: "14", on: g.fontSize == 0, onClick: func() { g.fontSize = 0 }},
					segItem{label: "16", on: g.fontSize == 16, onClick: func() { g.fontSize = 16 }},
				)
			})
		})
		row2 := ui.Row(c).AlignItems(ui.Center).Gap(t.Space(3))
		row2.Children(func() {
			toolGroup(c, "品牌色", func() {
				swatches := ui.Row(c).AlignItems(ui.Center).Gap(t.Space(1))
				swatches.Children(func() {
					for i, bc := range brandColors {
						bc := bc
						sw := ui.ButtonBase(c).
							Size(t.Space(5), t.Space(5)).
							Radius(t.Space(2.5)).
							Background(ui.Hex(bc.hex)).
							Center().
							Label("品牌色·" + bc.name)
						// 第一个色板（拂晓蓝）即 antd 默认：选中它等于清空
						// 覆盖，跟随主题预设（暗色自动切 #1668dc）。
						on := (i == 0 && g.brand == "") || g.brand == bc.hex
						if on {
							// 选中态对齐 antd ColorPicker：白色对勾 + 1px 描边环。
							sw.Shadow(0, 0, 0, 1, t.Border)
						}
						sw.Children(func() {
							if on {
								icon.Icon(c, icon.Check, t.Space(2.5), icon.WithColor(ui.Hex("#fff")))
							}
						})
						if sw.Clicked() {
							if i == 0 {
								g.brand = ""
							} else {
								g.brand = bc.hex
							}
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

// content 右侧内容区：单控件单页、上下撑满的分栏布局——
// 上方集中摆放该控件全部状态/示例的实时预览（两列网格、行高均分），
// 下方集中展示全部示例代码（antd Tabs 切换 + 代码编辑器撑满），
// 中间用可拖拽的垂直分隔面板分配比例，整页无空白。
func (g *Gallery) content(c *ui.Context) {
	d, ok := g.current()
	if !ok {
		ui.Column(c).Fill().Center().Children(func() {
			widgets.Empty(c, widgets.EmptyDescription("未找到匹配的控件，试试其他关键词"))
		})
		return
	}

	root := ui.Column(c).Fill()
	percent := ui.Local(root, "gallery.split", func() float32 { return 58 })
	root.Children(func() {
		widgets.SplitterVerticalPercent(c, percent, func() {
			g.previews(c, d)
		}, func() {
			g.codePanel(c, d)
		}).Fill()
	})
}

// previews 上方预览区：该控件全部示例（状态）的实时预览，
// 两列一行、行高均分，撑满整个预览面板，不留空白。
// 网格四周留 Padding、卡片之间留 Gap（antd Card gutter 16px），
// 避免卡片贴着面板边缘或相互粘连。
func (g *Gallery) previews(c *ui.Context, d samples.ControlPage) {
	t := c.Theme()
	grid := ui.Column(c).Fill().Padding(t.Space(4)).Gap(t.Space(4))
	n := len(d.Samples)
	if n == 0 {
		grid.Children(func() {
			ui.Column(c).Fill().Center().Children(func() {
				widgets.Empty(c, widgets.EmptyDescription("该控件暂无示例"))
			})
		})
		return
	}
	grid.Children(func() {
		// 对比模式：一次只展示当前示例（上下两栏各占一半面板），否则
		// 一格卡片高度不足以同时放下两种主题的预览。示例用下方页签切换。
		if g.compare != compareNone {
			if g.codeTab >= n {
				g.codeTab = 0
			}
			row := ui.Row(c).FillWidth().Fill()
			row.Children(func() {
				g.statePanel(c, g.codeTab, d.Samples[g.codeTab])
			})
			return
		}
		rows := (n + 1) / 2
		for i := 0; i < n; i += 2 {
			end := i + 2
			if end > n {
				end = n
			}
			row := ui.Row(c).FillWidth().Gap(t.Space(4)).AlignItems(ui.Stretch)
			if rows == 1 {
				// 仅一行时撑满整个预览区，避免上下留白。
				row.Fill()
			} else {
				row.Grow(1)
			}
			row.Children(func() {
				for j := i; j < end; j++ {
					g.statePanel(c, j, d.Samples[j])
				}
			})
		}
	})
}

// statePanel 预览网格中的一格：状态标题 + 居中的实时预览，
// 随网格行列拉伸撑满，无需滚动。
func (g *Gallery) statePanel(c *ui.Context, index int, s samples.Sample) {
	t := c.Theme()
	panel := ui.Column(c).Grow(1).
		Background(t.Surface).
		Radius(t.Radius+t.Space(0.5)).
		Border(1, t.Border).
		Shadow(0, 1, 2, 0, t.Text.Alpha(0.04))
	panel.Children(func() {
		head := ui.Row(c).FillWidth().
			Padding(t.Space(2), t.Space(4)).
			AlignItems(ui.Center).
			Gap(t.Space(2)).
			Border(1, t.Border).BorderWidth(0, 0, 1, 0)
		head.Children(func() {
			ui.Textf(c, "%02d", index+1).TextColor(t.Accent).FontSize(t.Rem(0.86)).FontWeight(600)
			ui.Text(c, s.Title).FontSize(t.Rem(1.0)).FontWeight(600).TextColor(t.Text).SingleLine()
			ui.Text(c, s.Description).TextColor(t.TextMuted).FontSize(t.Rem(0.78)).
				SingleLine().Grow(1)
		})
		body := ui.Column(c).Fill().Grow(1).
			Padding(t.Space(5))
		body.Children(func() {
			switch g.compare {
			case compareTheme:
				// 两个子面板上下平分高度（横排每面板仅约 264px，宽内容必截断）；
				// 各自内容超出时由 comparePane 内部的 Scroll 滚动，不再截断。
				col := ui.Column(c).Fill().Gap(t.Space(2))
				col.Children(func() {
					g.comparePane(c, "亮色", false, g.density, s.View)
					g.comparePane(c, "暗色", true, g.density, s.View)
				})
			case compareDensity:
				col := ui.Column(c).Fill().Gap(t.Space(2))
				col.Children(func() {
					g.comparePane(c, "标准密度", g.dark, theme.Comfortable, s.View)
					g.comparePane(c, "紧凑密度", g.dark, theme.Compact, s.View)
				})
			default:
				if s.Fill {
					host := ui.Column(c).Fill()
					host.Children(func() {
						s.View(c)
					})
				} else {
					host := ui.Column(c).Fill().Center()
					host.Children(func() {
						s.View(c)
					})
				}
			}
		})
	})
}
