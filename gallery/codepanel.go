package main

import (
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
	samples "github.com/lishuangquan1987/mygo-antdesign/samples"
	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

// codePanel 下方代码区：antd Tabs 切换该控件的全部示例，选中示例的
// 完整代码在编辑器中展示，编辑器撑满剩余高度（内部滚动）。
func (g *Gallery) codePanel(c *ui.Context, d samples.ControlPage) {
	t := c.Theme()
	panel := ui.Column(c).Fill().
		Background(t.Surface).
		Radius(t.Radius+t.Space(0.5)).
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
			widgets.Tabs(c, &g.codeTab, titles)
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
				buf := editor.NewBuffer(s.Code)
				// 示例代码均为 Go：设置语言启用语法高亮
				//（关键字/字符串/注释等按编辑器默认 GitHub 配色着色）。
				buf.SetLanguage(editor.LanguageNamed("go"))
				return editor.State{
					Buffer:       buf,
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
