package ant_test

import (
	"image"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/lishuangquan1987/mygo-antdesign/icon"
	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

func TestDrawerOpensAndCloses(t *testing.T) {
	open := false
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.Button(c, "打开抽屉", widgets.Primary(),
			widgets.OnClick(func() { open = true }))
		widgets.Drawer(c, &open, "设置", func() {
			widgets.Paragraph(c, "抽屉内容")
		})
	}
	tt := ui.NewTester(view, 640, 400)
	tt.Frame()
	if tt.HasText("抽屉内容") {
		t.Fatal("未打开时抽屉内容不应渲染")
	}
	if err := tt.Click("打开抽屉"); err != nil {
		t.Fatalf("点击打开失败: %v", err)
	}
	tt.Frame()
	if !tt.HasText("设置") || !tt.HasText("抽屉内容") {
		t.Fatal("抽屉未打开")
	}
	if err := tt.Click("✕"); err != nil {
		t.Fatalf("点击关闭失败: %v", err)
	}
	tt.Frame()
	if open {
		t.Error("点击关闭后 open 仍为 true")
	}
}

func TestProgressCircle(t *testing.T) {
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		row := ui.Row(c).Gap(16).AlignItems(ui.Center)
		row.Children(func() {
			widgets.ProgressCircle(c, 68)
			widgets.ProgressCircle(c, 100, widgets.ProgressStatusOpt(widgets.ProgressSuccess))
			widgets.ProgressCircle(c, 30, widgets.ProgressStatusOpt(widgets.ProgressError))
		})
	}
	tt := ui.NewTester(view, 420, 180)
	tt.Frame()
	for _, want := range []string{"68%", "100%", "30%"} {
		if !tt.HasText(want) {
			t.Errorf("环形进度未渲染 %q", want)
		}
	}
}

func TestImageAndColorPickerRender(t *testing.T) {
	img := ui.NewBitmap(image.NewRGBA(image.Rect(0, 0, 8, 8)))
	color := ui.Hex("#1677ff")
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		row := ui.Row(c).Gap(16).AlignItems(ui.Center)
		row.Children(func() {
			widgets.ImageEl(c, img, widgets.ImageWidth(8))
			widgets.ColorPicker(c, &color)
			widgets.ColorPanel(c, &color)
			widgets.Toggle(c, new(bool), "B")
			widgets.Affix(c, 4, func() {
				widgets.Paragraph(c, "固钉内容")
			})
		})
	}
	tt := ui.NewTester(view, 640, 320)
	tt.Frame()
	if !tt.HasText("固钉内容") {
		t.Error("Affix 内容未渲染")
	}
	if tt.FramePending() {
		t.Error("图片/取色器持续请求帧")
	}
}

func TestIconImageSourcePreview(t *testing.T) {
	svg, err := ui.ParseSVG([]byte(icon.Search))
	if err != nil {
		t.Fatal(err)
	}
	view := func(c *ui.Context) {
		theme.Apply(c, theme.Light())
		widgets.ImageEl(c, svg, widgets.ImageWidth(6))
	}
	tt := ui.NewTester(view, 320, 200)
	tt.Frame()
	if tt.FramePending() {
		t.Error("SVG 图片持续请求帧")
	}
}
