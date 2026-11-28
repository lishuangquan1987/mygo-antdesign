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
	"github.com/egoist/mygo/ui/editor"

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

// quickStartCode 首页「快速开始」展示的最小示例（与 README 快速开始同源）。
const quickStartCode = `go get github.com/lishuangquan1987/mygo-antdesign@v0.5.0

// 运行示例应用：go run ./gallery
package main

import (
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

func main() {
	remember := false
	mygo.App.WhenReady(func() {
		mygo.NewWindow(mygo.WindowOptions{Title: "登录", Width: 420, Height: 320,
			Content: ui.View(func(c *ui.Context) {
				theme.Apply(c, theme.Light())
				widgets.Title(c, "登录", widgets.H2)
				widgets.Checkbox(c, &remember, "记住我")
				widgets.Button(c, "登录", widgets.Primary())
			}),
		})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}`

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
func (g *Gallery) previews(c *ui.Context, d samples.ControlPage) {
	grid := ui.Column(c).Fill().Gap(c.Theme().Space(2))
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
		rows := (n + 1) / 2
		for i := 0; i < n; i += 2 {
			end := i + 2
			if end > n {
				end = n
			}
			row := ui.Row(c).FillWidth().Gap(c.Theme().Space(2)).AlignItems(ui.Stretch)
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
		Radius(t.Radius + t.Space(0.5)).
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
				// 两个子面板上下堆叠（横排每面板仅约 264px，宽内容必截断），
				// 整体超高由 body 层的滚动兜底（见 body 定义处）。
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

// codePanel 下方代码区：antd Tabs 切换该控件的全部示例，选中示例的
// 完整代码在编辑器中展示，编辑器撑满剩余高度（内部滚动）。
func (g *Gallery) codePanel(c *ui.Context, d samples.ControlPage) {
	t := c.Theme()
	panel := ui.Column(c).Fill().
		Background(t.Surface).
		Radius(t.Radius + t.Space(0.5)).
		Border(1, t.Border).
		Shadow(0, 1, 2, 0, t.Text.Alpha(0.04))
	panel.Children(func() {
		// 页头：示例标签页 + 复制按钮。
		if g.codeTab >= len(d.Samples) {
			g.codeTab = 0
		}
		head := ui.Row(c).FillWidth().
			Padding(t.Space(1), t.Space(4), 0, t.Space(4)).
			AlignItems(ui.Center).
			Justify(ui.SpaceBetween).
			Border(1, t.Border).BorderWidth(0, 0, 1, 0)
		head.Children(func() {
			titles := make([]string, len(d.Samples))
			for i, s := range d.Samples {
				titles[i] = s.Title
			}
			widgets.Tabs(c, &g.codeTab, titles...)
			if widgets.Button(c, "复制代码", widgets.Size(widgets.Small), widgets.Text()).Clicked() {
				c.WriteClipboard(d.Samples[g.codeTab].Code)
				c.Toast("代码已复制")
			}
		})

		// 代码正文：完整代码，编辑器撑满剩余空间（内部滚动）。
		// 编辑器原语固定绘制主题 Background（暗色为页面底 #000），
		// 这里对暗色用局部主题把编辑器底切到容器色 #141414，
		// 对齐 antd dark 的容器层级（代码区属于容器，不该比页面更暗）。
		body := ui.Column(c).Fill().Grow(1)
		body.Children(func() {
			s := d.Samples[g.codeTab]
			st := ui.Local(panel, "gallery.codeview."+d.Name+"."+s.Title, func() editor.State {
				return editor.State{
					Buffer:       editor.NewBuffer(s.Code),
					ReadOnly:     true,
					NoActiveLine: true,
				}
			})
			if t.Dark {
				// 编辑器原语绘制主题 Background（暗色=页面底 #000），
				// 代码区属于容器层级：用局部主题把 BgLayout 提到容器色
				// #141414（不能拷贝 ui.Theme 结构体——会使编辑器渲染
				// 异常，必须走 AntToken → theme.Sub 正规通道）。
				etk := g.token()
				etk.BgLayout = etk.BgContainer
				theme.Sub(c, etk).Fill().Children(func() {
					widgets.CodeEditor(c, st).Fill()
				})
			} else {
				widgets.CodeEditor(c, st).Fill()
			}
		})
	})
}

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
						return editor.State{
							Buffer:       editor.NewBuffer(truncateCode(quickStartCode, 14, "完整示例见 README 快速开始")),
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


// comparePane 并排对比中的一栏：小标题 + 局部主题容器。
func (g *Gallery) comparePane(c *ui.Context, label string, dark bool, density theme.Density, view func(*ui.Context)) {
	t := c.Theme()
	tk := g.tokenOf(dark, density)
	pane := ui.Column(c).FillWidth().Grow(1).Gap(t.Space(1))
	pane.Children(func() {
		ui.Text(c, label).TextColor(t.TextMuted).FontSize(t.Rem(0.8)).FontWeight(500)
		// 子主题面板：底色/描边跟随各自主题（暗色面板用 BgContainer，
		// 不再用外层 Surface，避免暗色内容衬白底）；纵向撑满面板剩余
		// 空间，内容水平居中。
		// 纵向内边距收紧（子面板高度有限，Space(4) 上下会让内容溢出
		// 底边），内容水平+垂直居中。
		sub := theme.Sub(c, tk).
			FillWidth().Grow(1).
			Padding(t.Space(0.5), t.Space(4)).
			Radius(t.Radius).
			Background(tk.BgContainer).
			Border(1, tk.Border)
		sub.AlignItems(ui.Center).Justify(ui.Center)
		sub.Children(func() {
			view(c)
		})
	})
}
