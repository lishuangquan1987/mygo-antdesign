// Package icon 提供 Ant Design 图标（SVG 源码常量，outlined 子集，
// 生成自 scripts/genicons）。图标 fill=currentColor，颜色跟随所在
// 文本环境（主题色系自动生效）。
//
// 导入：import "github.com/lishuangquan1987/mygo-antdesign/icon"
package icon

import (
	"sync"

	"github.com/egoist/mygo/ui"
)

// svgCache 缓存解析后的 SVG。ui.ParseSVG 每次返回新指针，而 MyGo 的
// 图片缓存以 *SVG 指针为键——不缓存则同一图标每帧都要重新解析并重新
// 栅格化（Gallery 一帧数十次）。key 为 SVG 源码常量，进程内不可变。
var svgCache sync.Map // map[string]*ui.SVG（nil 表示解析失败）

// parseSVG 返回缓存中的解析结果；同一源码只解析一次。
func parseSVG(svgData string) (*ui.SVG, bool) {
	if v, ok := svgCache.Load(svgData); ok {
		s, _ := v.(*ui.SVG)
		return s, s != nil
	}
	s, err := ui.ParseSVG([]byte(svgData))
	if err != nil {
		svgCache.Store(svgData, (*ui.SVG)(nil))
		return nil, false
	}
	svgCache.Store(svgData, s)
	return s, true
}

// Icon 渲染一个图标（默认 16 DIP）。svgData 为包内常量之一。颜色
// 跟随文本色（ui.Image 对 SVG 以其 TextColor 作为 currentColor），
// 即主题色自动生效；也可显式指定：
//
//	icon.Icon(c, icon.Search, 16)                       // 跟随文本色
//	icon.Icon(c, icon.Search, 16, icon.WithColor(t.Accent)) // 指定色
func Icon(c *ui.Context, svgData string, size float32, opts ...IconOption) ui.Element {
	cfg := &iconConfig{}
	for _, o := range opts {
		o(cfg)
	}
	svg, ok := parseSVG(svgData)
	if !ok {
		return ui.Text(c, "?").FontSize(10)
	}
	el := ui.Image(c, svg).Width(size).Height(size)
	if cfg.color != nil {
		el.TextColor(*cfg.color)
	}
	return el
}

// IconOption 定制图标。
type IconOption func(*iconConfig)

type iconConfig struct {
	color *ui.Color
}

// WithColor 显式指定图标颜色（默认跟随文本色）。
func WithColor(c ui.Color) IconOption { return func(cfg *iconConfig) { cfg.color = &c } }
