package main

import (
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/egoist/mygo/ui"

	samples "github.com/lishuangquan1987/mygo-antdesign/samples"
	"github.com/lishuangquan1987/mygo-antdesign/samples/button"
	"github.com/lishuangquan1987/mygo-antdesign/samples/display"
	"github.com/lishuangquan1987/mygo-antdesign/samples/divider"
	extrademo "github.com/lishuangquan1987/mygo-antdesign/samples/extrademo"
	"github.com/lishuangquan1987/mygo-antdesign/samples/feedback"
	"github.com/lishuangquan1987/mygo-antdesign/samples/ide"
	"github.com/lishuangquan1987/mygo-antdesign/samples/input"
	"github.com/lishuangquan1987/mygo-antdesign/samples/layout"
	"github.com/lishuangquan1987/mygo-antdesign/samples/nav"
	"github.com/lishuangquan1987/mygo-antdesign/samples/selectdemo"
	"github.com/lishuangquan1987/mygo-antdesign/samples/table"
	"github.com/lishuangquan1987/mygo-antdesign/samples/typography"
	"github.com/lishuangquan1987/mygo-antdesign/theme"
	widgets "github.com/lishuangquan1987/mygo-antdesign/widgets"
)

// registerAll 注册全部示例（与 main.go 保持一致）。
func registerAll() {
	if len(samples.All()) > 0 {
		return
	}
	samples.Register(button.Demo())
	samples.Register(typography.Demo())
	samples.Register(divider.Demo())
	samples.Register(layout.Demo())
	samples.Register(input.Demo())
	samples.Register(selectdemo.Demo())
	samples.Register(display.Demo())
	samples.Register(table.Demo())
	samples.Register(feedback.Demo())
	samples.Register(nav.Demo())
	samples.Register(extrademo.Demo())
	samples.Register(ide.Demo())
}

// savePNG 用无窗口 Tester 渲染一帧并写 PNG。
func savePNG(t *testing.T, path string, view func(*ui.Context), w, h int) {
	t.Helper()
	tt := ui.NewTester(view, w, h)
	tt.Frame()
	img := tt.Image()
	if img == nil {
		t.Fatalf("%s: 无渲染帧", path)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(f, img); err != nil {
		f.Close()
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
}

// TestWriteDocsScreenshots 生成 README 用的截图（仓库内提交，便于审阅视觉回归）。
// 默认跳过；需要更新时运行：
//
//	WRITE_DOCS_SCREENSHOTS=1 go test ./gallery -run TestWriteDocsScreenshots
func TestWriteDocsScreenshots(t *testing.T) {
	if os.Getenv("WRITE_DOCS_SCREENSHOTS") != "1" {
		t.Skip("设置 WRITE_DOCS_SCREENSHOTS=1 生成 docs/images 截图")
	}
	registerAll()

	dir := filepath.Join("..", "docs", "images")
	g := NewGallery()
	g.selected = "Button"

	// 亮色 Gallery。
	g.dark = false
	savePNG(t, filepath.Join(dir, "gallery.png"), g.View, 1280, 800)

	// 暗色 Gallery（放在同一张卡片流上，选择「反馈」分组以展示更多控件）。
	g.dark = true
	g.selected = "反馈"
	savePNG(t, filepath.Join(dir, "gallery-dark.png"), g.View, 1280, 800)

	// 最小示例（README 快速开始）。
	name := ""
	remember := false
	quickstart := func(c *ui.Context) {
		theme.Apply(c, theme.Light(theme.WithPrimary(ui.RGB(0x16, 0x77, 0xff))))
		t := c.Theme()
		col := ui.Column(c).Gap(t.Space(3)).Padding(t.Space(6)).FillWidth()
		col.Children(func() {
			widgets.Title(c, "登录", widgets.H2)
			widgets.Input(c, &name, widgets.PlaceholderText("用户名"))
			widgets.Checkbox(c, &remember, "记住我")
			widgets.Button(c, "登录", widgets.Primary())
		})
	}
	savePNG(t, filepath.Join(dir, "quickstart.png"), quickstart, 640, 360)
}
