package ant_test

import (
	"image"
	"testing"

	"github.com/egoist/mygo/ui"

	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

// TestWatermarkClipsToHost 守护水印的绘制范围：水印只在包裹内容的矩形内
// 平铺，不铺满整个窗口（覆盖层是窗口级的，早前用整窗矩形导致侧栏与工具条
// 也被铺上水印）。
func TestWatermarkClipsToHost(t *testing.T) {
	img := func(withWatermark bool) *image.RGBA {
		tt := ui.NewTester(func(c *ui.Context) {
			ui.Column(c).Fill().Padding(40).Children(func() {
				if withWatermark {
					widgets.Watermark(c, "内部资料", func() {
						ui.Column(c).Width(200).Height(80).Children(func() {})
					})
					return
				}
				ui.Column(c).Width(200).Height(80).Children(func() {})
			})
		}, 400, 300)
		tt.Frame()
		return tt.Image()
	}

	plain := img(false)
	marked := img(true)

	minX, minY, maxX, maxY := 1<<30, 1<<30, -1, -1
	b := marked.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if plain.RGBAAt(x, y) != marked.RGBAAt(x, y) {
				if x < minX {
					minX = x
				}
				if y < minY {
					minY = y
				}
				if x > maxX {
					maxX = x
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}
	if maxX < 0 {
		t.Fatal("水印未绘制：两帧完全一致")
	}
	// 水印宿主是 FillWidth 的列，宽度撑满内边距内的可用区，但高度只到内容
	// 底部（40+80）。关键约束是纵向不铺满整窗（此前的 bug 会画到窗口底边）。
	if minY < 38 || maxY > 122 {
		t.Errorf("水印纵向越出宿主：y 范围 %d~%d，期望落在 40~120 内", minY, maxY)
	}
	if minX < 38 {
		t.Errorf("水印横向越出宿主：x 起点 %d，期望不小于 40", minX)
	}
}
