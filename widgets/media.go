package ant

import (
	"github.com/egoist/mygo/ui"
)

// ImageOption 定制 ImageEl。
type ImageOption func(*imageConfig)

type imageConfig struct {
	width    float32 // 主题单位数；0 = 原始尺寸
	preview  bool
	maxWidth float32
}

// ImageWidth 显示宽度（主题单位数，随密度缩放）。
func ImageWidth(units float32) ImageOption { return func(c *imageConfig) { c.width = units } }

// ImagePreview 点击放大预览（antd Image preview，默认开启）。
func ImagePreview(b bool) ImageOption { return func(c *imageConfig) { c.preview = b } }

// ImageMaxWidth 预览大图的最大宽度（主题单位数）。
func ImageMaxWidth(units float32) ImageOption {
	return func(c *imageConfig) { c.maxWidth = units }
}

// ImageEl 图片（antd Image）：可指定显示宽度，点击后在模态框中放大预览。
//
//	ant.ImageEl(c, img, ant.ImageWidth(40))
//	ant.ImageEl(c, svg, ant.ImageWidth(8), ant.ImagePreview(false))
func ImageEl(c *ui.Context, src ui.ImageSource, opts ...ImageOption) ui.Element {
	cfg := &imageConfig{preview: true, maxWidth: 200}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()

	host := ui.Column(c).AlignItems(ui.Start).Gap(t.Space(1))
	open := ui.Local(host, "ant.image.open", func() bool { return false })
	host.Children(func() {
		if cfg.preview {
			box := ui.ButtonBase(c).Padding(0).Radius(t.Radius)
			box.Children(func() {
				img := ui.Image(c, src)
				if cfg.width > 0 {
					img.Width(t.Space(cfg.width))
				}
				img.Radius(t.Radius)
			})
			if box.Clicked() {
				*open = true
			}
		} else {
			img := ui.Image(c, src)
			if cfg.width > 0 {
				img.Width(t.Space(cfg.width))
			}
			img.Radius(t.Radius)
		}
		ui.Modal(c, open, func() {
			ui.Image(c, src).MaxWidth(t.Space(cfg.maxWidth))
		})
	})
	return host
}

// ColorPicker 颜色选择（antd ColorPicker）：色块触发器 + 展开的取色面板，
// 复用 MyGo 的 ColorWell（点击/Esc/点外关闭，键盘可用）：
//
//	ant.ColorPicker(c, &app.accent)
func ColorPicker(c *ui.Context, color *ui.Color, opts ...ControlOption) ui.Element {
	cfg := &controlConfig{}
	for _, o := range opts {
		o(cfg)
	}
	host := ui.Row(c).AlignItems(ui.Center)
	host.Children(func() {
		el := ui.ColorWell(c, color)
		applyControl(el, cfg)
	})
	return host
}

// ColorPanel 常驻的取色面板（antd ColorPicker 展开后的面板，不收起）。
func ColorPanel(c *ui.Context, color *ui.Color, opts ...ControlOption) ui.Element {
	cfg := &controlConfig{}
	for _, o := range opts {
		o(cfg)
	}
	host := ui.Row(c).AlignItems(ui.Start)
	host.Children(func() {
		el := ui.ColorPicker(c, color)
		applyControl(el, cfg)
	})
	return host
}
