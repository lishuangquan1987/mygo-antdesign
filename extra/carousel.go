package extra

import (
	"github.com/egoist/mygo/ui"
)

// CarouselOption 定制 Carousel。
type CarouselOption func(*carouselConfig)

type carouselConfig struct {
	height float32 // 页面区高度；0 表示随内容
}

// CarouselHeight 页面区高度（DIP）。
func CarouselHeight(h float32) CarouselOption { return func(c *carouselConfig) { c.height = h } }

// Carousel 走马灯（antd Carousel）：一次显示一页，底部指示点切换。
// pages 为各页的构建函数（即时模式，每帧构建当前页）：
//
//	ant.Carousel(c, []func(*ui.Context){
//	    func(c *ui.Context) { widgets.Title(c, "第一页", widgets.H3) },
//	    func(c *ui.Context) { widgets.Title(c, "第二页", widgets.H3) },
//	})
func Carousel(c *ui.Context, pages []func(*ui.Context), opts ...CarouselOption) ui.Element {
	cfg := &carouselConfig{}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()

	host := ui.Column(c).Gap(t.Space(2)).FillWidth()
	index := ui.Local(host, "ant.extra.carousel.index", func() int { return 0 })

	host.Children(func() {
		// 页面区。
		page := ui.Box(c).FillWidth().Center()
		if cfg.height > 0 {
			page.Height(cfg.height)
		} else {
			page.MinHeight(t.Space(20))
		}
		page.Radius(t.Radius).Background(t.SurfaceHover)
		page.Children(func() {
			if len(pages) > 0 && *index >= 0 && *index < len(pages) {
				pages[*index](c)
			}
		})
		// 指示点：外层 12×12 命中区，内含 8×8 视觉圆点。
		if len(pages) > 1 {
			dots := ui.Row(c).Gap(t.Space(2)).Justify(ui.Center).FillWidth()
			dots.Children(func() {
				for i := range pages {
					i := i
					dot := ui.ButtonBase(c).Size(t.Space(3), t.Space(3)).Center()
					dot.Children(func() {
						inner := ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(t.Space(1))
						inner.Children(func() {})
						if i == *index {
							inner.Background(t.Accent)
						} else {
							inner.Background(t.Border)
						}
					})
					if dot.Clicked() {
						*index = i
					}
				}
			})
		}
	})
	return host
}
