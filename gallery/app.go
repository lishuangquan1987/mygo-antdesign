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
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
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

// Gallery 持有展示应用的状态；View 是状态的函数，每次交互后下一帧重绘。
type Gallery struct {
	dark      bool
	density   theme.Density
	brand     string  // 品牌色 hex，"" 表示 antd 默认
	compare   int     // compareNone / compareTheme / compareDensity
	fontSize  float32 // 0 表示跟随密度（标准 14 / 紧凑 13）
	size      gallerySize
	selected  string                    // 当前展示的控件名
	codeTab   int                       // 对比模式下当前展示的示例下标
	query     string                    // 侧栏搜索词
	collapsed map[samples.Category]bool // 侧栏分类折叠状态（nil 表示全部展开）
	expanded  map[string]bool           // demo 代码展开状态（key = 页名|示例标题）
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

// homeName 总览页在侧栏中的名字（g.selected == homeName 时展示总览）。
const homeName = "组件总览"

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
		ui.Textf(c, "%s / %s / %s", homeName, d.Category, d.Name).
			TextColor(t.TextMuted).FontSize(t.Rem(0.72))

		title := ui.Row(c).AlignItems(ui.Center).Gap(t.Space(2))
		title.Children(func() {
			ui.Text(c, d.Name).FontSize(t.Rem(1.71)).FontWeight(600).TextColor(t.Text)
			ui.Text(c, samples.PageZh(d.Name)).FontSize(t.Rem(1.14)).TextColor(t.TextMuted)
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

// content 右侧内容区：对齐 antd 组件文档页的结构——该控件全部示例
// 竖排为一张张 demo 卡片（标题 + 描述 + 居中实时预览 + 底部操作条），
// 代码默认收起，点「展开代码」就地展开（懒加载编辑器，不展开不构建）。
// 整页放进滚动容器，内容自然高、卡片永远完整显示。
func (g *Gallery) content(c *ui.Context) {
	d, ok := g.current()
	if !ok {
		ui.Column(c).Fill().Center().Children(func() {
			widgets.Empty(c, widgets.EmptyDescription("未找到匹配的控件，试试其他关键词"))
		})
		return
	}

	t := c.Theme()
	root := ui.Column(c).Fill()
	// 预览页的滚动位置按控件页保存，切页时互不干扰。
	scroll := ui.Local(root, "gallery.preview.scroll."+d.Name, func() ui.ScrollState {
		return ui.ScrollState{}
	})
	root.Children(func() {
		pane := ui.Scroll(c).Fill().Grow(1).TrackScroll(scroll)
		pane.Children(func() {
			col := ui.Column(c).FillWidth().Padding(t.Space(4)).Gap(t.Space(4))
			col.Children(func() {
				// 对比模式：一次只展示当前示例（亮暗/密度上下分栏），
				// 示例用页签切换。
				if g.compare != compareNone {
					n := len(d.Samples)
					if g.codeTab >= n {
						g.codeTab = 0
					}
					titles := make([]string, n)
					for i, s := range d.Samples {
						titles[i] = s.Title
					}
					widgets.Tabs(c, &g.codeTab, titles)
					g.demoCard(c, d.Name, g.codeTab, d.Samples[g.codeTab])
					return
				}
				for i := range d.Samples {
					g.demoCard(c, d.Name, i, d.Samples[i])
				}
			})
		})
	})
}

// demoCard 一张示例卡片（对齐 antd demo-block）：标题 + 描述 + 居中
// 实时预览 + 底部操作条（复制代码 / 展开代码）。代码区懒加载：展开时
// 才构建编辑器，收起时不保留任何状态（性能）。
func (g *Gallery) demoCard(c *ui.Context, page string, index int, s samples.Sample) {
	t := c.Theme()
	key := page + "|" + s.Title
	panel := ui.Column(c).FillWidth().
		Background(t.Surface).
		Radius(t.Radius+t.Space(0.5)).
		Border(1, t.Border).
		Shadow(0, 1, 2, 0, t.Text.Alpha(0.04))
	panel.Children(func() {
		// 标题 + 描述（antd demo 卡片头部，左对齐）。
		head := ui.Column(c).FillWidth().
			Padding(t.Space(5), t.Space(8), 0, t.Space(8)).
			Gap(t.Space(1))
		head.Children(func() {
			ui.Textf(c, "%02d. %s", index+1, s.Title).
				FontSize(t.Rem(1.0)).FontWeight(600).TextColor(t.Text)
			if s.Description != "" {
				ui.Text(c, s.Description).TextColor(t.TextMuted).FontSize(t.Rem(0.82))
			}
		})

		// 实时预览：居中（Fill 示例撑满），MinHeight 保证体量。
		body := ui.Column(c).FillWidth().
			Padding(t.Space(6), t.Space(8)).
			MinHeight(t.Space(34))
		body.Children(func() {
			switch g.compare {
			case compareTheme:
				// 两个子面板上下排列（各占一半面板高度）；
				// 各自内容超出时由 comparePane 内部的 Scroll 滚动，不截断。
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
					host := ui.Column(c).FillWidth().Grow(1)
					host.Children(func() {
						s.View(c)
					})
				} else {
					host := ui.Column(c).FillWidth().MinHeight(t.Space(30)).Center()
					host.Children(func() {
						s.View(c)
					})
				}
			}
		})

		// 展开的代码区（懒加载）：构建过才保留 ui.Local 状态。
		open := g.expanded[key]
		if open {
			g.codeView(c, panel, page, s)
		}

		// 底部操作条：右对齐「复制代码 / 展开代码」。
		foot := ui.Row(c).FillWidth().
			Padding(t.Space(2), t.Space(6)).
			Justify(ui.End).
			AlignItems(ui.Center).
			Gap(t.Space(3)).
			Border(1, t.Border).BorderWidth(1, 0, 0, 0)
		foot.Children(func() {
			if widgets.Button(c, "复制代码", widgets.Size(widgets.Small), widgets.Text()).Clicked() {
				c.WriteClipboard(s.Code)
				c.Toast("代码已复制")
			}
			label := "展开代码"
			arrow := icon.Down
			if open {
				label = "收起代码"
				arrow = icon.Up
			}
			if widgets.Button(c, label, widgets.Size(widgets.Small), widgets.Text()).Clicked() {
				if g.expanded == nil {
					g.expanded = make(map[string]bool)
				}
				g.expanded[key] = !open
			}
			icon.Icon(c, arrow, t.Space(3), icon.WithColor(t.TextMuted))
		})
	})
}

// codeView 展开的示例代码区：完整源码（FitContent 自适应行数），
// 编辑器状态按 页名+示例标题 缓存——收起时不构建、不保留（懒加载）。
func (g *Gallery) codeView(c *ui.Context, host ui.Element, page string, s samples.Sample) {
	t := c.Theme()
	area := ui.Column(c).FillWidth().
		Padding(t.Space(4), t.Space(6)).
		Background(t.Background)
	area.Children(func() {
		st := ui.Local(host, "gallery.code."+page+"."+s.Title, func() editor.State {
			buf := editor.NewBuffer(s.Code)
			// 示例代码均为 Go：设置语言启用语法高亮。
			buf.SetLanguage(editor.LanguageNamed("go"))
			return editor.State{
				Buffer:       buf,
				ReadOnly:     true,
				NoActiveLine: true,
				FitContent:   true,
			}
		})
		if t.Dark {
			// 编辑器原语绘制主题 Background（暗色=页面底 #000），
			// 代码区属于容器层级：用局部主题把 BgLayout 提到容器色
			// #141414（不能拷贝 ui.Theme 结构体——会使编辑器渲染
			// 异常，必须走 AntToken → theme.Sub 正规通道）。
			etk := g.token()
			etk.BgLayout = etk.BgContainer
			theme.Sub(c, etk).FillWidth().Children(func() {
				widgets.CodeEditor(c, st).FillWidth()
			})
		} else {
			widgets.CodeEditor(c, st).FillWidth()
		}
	})
}
