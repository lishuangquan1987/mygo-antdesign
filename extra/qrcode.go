package extra

import (
	"image/color"

	qrcode "github.com/skip2/go-qrcode"

	"github.com/egoist/mygo/ui"
)

// QRCodeOption 定制 QRCode。
type QRCodeOption func(*qrcodeConfig)

type qrcodeConfig struct {
	size float32 // 显示尺寸（DIP）；默认 96
}

// QRCodeSize 显示尺寸（DIP，正方形）。
func QRCodeSize(px float32) QRCodeOption { return func(c *qrcodeConfig) { c.size = px } }

// QRCode 二维码（antd QRCode）。码色与底色取自主题（亮色黑码白底、
// 暗色白码深底，对齐 antd dark 模式），内容为空时显示占位文字：
//
//	ant.QRCode(c, "https://example.com", ant.QRCodeSize(128))
func QRCode(c *ui.Context, content string, opts ...QRCodeOption) ui.Element {
	cfg := &qrcodeConfig{size: 96}
	for _, o := range opts {
		o(cfg)
	}
	t := c.Theme()

	pad := t.Space(3)
	panel := ui.Box(c).Padding(pad, pad).Radius(t.Radius).Background(t.Surface)
	panel.Children(func() {
		if content == "" {
			ui.Text(c, "无内容").TextColor(t.TextMuted)
			return
		}
		qr, err := qrcode.New(content, qrcode.Medium)
		if err != nil {
			ui.Text(c, "生成失败").TextColor(t.TextMuted)
			return
		}
		qr.DisableBorder = true
		qr.ForegroundColor = toGoColor(t.Text)
		qr.BackgroundColor = toGoColor(t.Surface)
		bmp := ui.NewBitmap(qr.Image(256))
		ui.Image(c, bmp).Width(cfg.size).Height(cfg.size)
	})
	return panel
}

// toGoColor 把 ui.Color 转为标准库 color.Color。
func toGoColor(c ui.Color) color.Color {
	return color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}
}
